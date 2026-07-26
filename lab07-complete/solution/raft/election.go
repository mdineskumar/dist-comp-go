package main

import (
	"fmt"
	"time"
)

func (n *Node) becomeLeader() {
	n.mu.Lock()
	if n.state != Candidate { n.mu.Unlock(); return }
	n.state = Leader
	for _, peer := range n.peers {
		n.nextIndex[peer] = n.lastLogIndex() + 1
		n.matchIndex[peer] = 0
	}
	fmt.Printf("[NODE %s] Became LEADER (term %d) ← %d peers\n", n.id, n.currentTerm, len(n.peers))
	n.mu.Unlock()
	go n.sendHeartbeats()
}

func (n *Node) startElection() {
	n.mu.Lock()
	term := n.currentTerm
	lastIdx := n.lastLogIndex()
	lastTerm := n.lastLogTerm()
	n.mu.Unlock()

	votes := make(chan bool, len(n.peers))
	for _, peer := range n.peers {
		go func(p string) {
			args := &RequestVoteArgs{Term: term, CandidateID: n.id, LastLogIndex: lastIdx, LastLogTerm: lastTerm}
			var reply RequestVoteReply
			if err := callRPC(p, "RaftRPC.RequestVote", args, &reply); err != nil {
				votes <- false; return
			}
			if reply.Term > term {
				n.mu.Lock(); n.becomeFollower(reply.Term); n.mu.Unlock()
			}
			votes <- reply.VoteGranted
		}(peer)
	}

	voteCount := 1
	for i := 0; i < len(n.peers); i++ {
		if <-votes { voteCount++ }
		if voteCount > len(n.peers)/2 { n.becomeLeader(); return }
	}
}

func (n *Node) sendHeartbeats() {
	ticker := time.NewTicker(50 * time.Millisecond)
	for range ticker.C {
		n.mu.Lock()
		isLeader := n.state == Leader
		n.mu.Unlock()
		if !isLeader { ticker.Stop(); return }
		for _, peer := range n.peers {
			go n.sendAppendEntries(peer)
		}
	}
}

func (n *Node) sendAppendEntries(peerAddr string) {
	n.mu.Lock()
	if n.state != Leader { n.mu.Unlock(); return }
	nextIdx := n.nextIndex[peerAddr]
	prevLogIndex := nextIdx - 1
	prevLogTerm := 0
	if prevLogIndex > 0 && prevLogIndex <= len(n.log) {
		prevLogTerm = n.log[prevLogIndex-1].Term
	}
	var entries []LogEntry
	if nextIdx <= len(n.log) { entries = n.log[nextIdx-1:] }
	args := &AppendEntriesArgs{Term: n.currentTerm, LeaderID: n.id,
		PrevLogIndex: prevLogIndex, PrevLogTerm: prevLogTerm,
		Entries: entries, LeaderCommit: n.commitIndex}
	term := n.currentTerm
	n.mu.Unlock()

	var reply AppendEntriesReply
	if err := callRPC(peerAddr, "RaftRPC.AppendEntries", args, &reply); err != nil { return }

	n.mu.Lock()
	defer n.mu.Unlock()
	if reply.Term > n.currentTerm { n.becomeFollower(reply.Term); return }
	if n.state != Leader || n.currentTerm != term { return }
	if reply.Success {
		if len(entries) > 0 {
			n.matchIndex[peerAddr] = args.PrevLogIndex + len(entries)
			n.nextIndex[peerAddr] = n.matchIndex[peerAddr] + 1
			fmt.Printf("[LEADER %s] Replicated to %s (matchIndex=%d)\n", n.id, peerAddr, n.matchIndex[peerAddr])
		}
		n.commitEntries()
	} else {
		if n.nextIndex[peerAddr] > 1 { n.nextIndex[peerAddr]-- }
	}
}
