package main

// ============================================================
// Lab 07 — Consensus (Simplified Raft)
// File: election.go
// Role: Leader election and heartbeat sending
//
// TASKS IN THIS FILE:
//   Task 4 — becomeLeader()
//   Task 5 — startElection()
//   Task 7 — sendHeartbeats()
// ============================================================

import (
	"fmt"
	"time"
)

// ============================================================
// TASK 4 — becomeLeader
// ============================================================
// Transition to Leader state after winning an election.
//
// Steps:
//   1. Lock: n.mu.Lock()
//   2. Verify we are still a Candidate (we might have lost the
//      election while goroutines were running) — if n.state != Candidate, unlock and return
//   3. Set n.state = Leader
//   4. Initialise nextIndex for each peer:
//        n.nextIndex[peer] = n.lastLogIndex() + 1
//   5. Initialise matchIndex for each peer:
//        n.matchIndex[peer] = 0
//   6. Print: [NODE id] Became LEADER (term N) ← N peers
//   7. Unlock: n.mu.Unlock()
//   8. Send immediate heartbeats: go n.sendHeartbeats()
//
// TODO: implement this function
func (n *Node) becomeLeader() {
	// YOUR CODE HERE
}

// ============================================================
// TASK 5 — startElection
// ============================================================
// Send RequestVote RPCs to all peers concurrently and count votes.
// This is the core of leader election.
//
// Steps:
//   1. Lock, snapshot current state (term, lastLogIndex, lastLogTerm), unlock
//   2. Create a vote channel: votes := make(chan bool, len(n.peers))
//   3. For each peer, launch a goroutine:
//        a. Call RequestVote RPC
//        b. If reply.Term > our term → becomeFollower(reply.Term), return
//        c. Send reply.VoteGranted to votes channel
//   4. Count votes: start with 1 (our own vote)
//      Read from votes channel len(n.peers) times
//      If granted: voteCount++
//      If voteCount > len(n.peers)/2 → we have majority → becomeLeader(), return
//   5. If election ends without majority → stay as candidate (timer will retry)
//
// ── KEY INSIGHT ───────────────────────────────────────────
// "First to majority wins" — as soon as we have N/2+1 votes
// we call becomeLeader() immediately without waiting for all
// peers to respond. This makes elections fast.
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func (n *Node) startElection() {
	// YOUR CODE HERE
}

// ============================================================
// TASK 7 — sendHeartbeats
// ============================================================
// Leader sends periodic AppendEntries RPCs to all followers.
// Heartbeats prevent followers from timing out and starting
// unnecessary elections.
//
// Steps:
//   1. Create a ticker: ticker := time.NewTicker(50 * time.Millisecond)
//   2. On each tick:
//        a. Lock, check if still Leader, unlock
//        b. If NOT leader → ticker.Stop() and return
//        c. If leader → for each peer, go sendAppendEntries(peer)
//
// ── WHY 50ms? ─────────────────────────────────────────────
// Election timeout is 150-300ms random.
// Heartbeat interval (50ms) must be much less than election timeout
// so followers always receive a heartbeat before they time out.
// In real Raft: heartbeat << election timeout << MTBF (mean time between failures)
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func (n *Node) sendHeartbeats() {
	// YOUR CODE HERE
	_ = time.NewTicker // remove when implementing
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// sendAppendEntries sends one AppendEntries RPC to a single peer.
// Used for both heartbeats (empty entries) and log replication.
func (n *Node) sendAppendEntries(peerAddr string) {
	n.mu.Lock()
	if n.state != Leader {
		n.mu.Unlock()
		return
	}

	// Figure out what entries to send to this peer
	nextIdx := n.nextIndex[peerAddr]
	prevLogIndex := nextIdx - 1
	prevLogTerm := 0
	if prevLogIndex > 0 && prevLogIndex <= len(n.log) {
		prevLogTerm = n.log[prevLogIndex-1].Term
	}

	// Entries to send (may be empty = heartbeat only)
	var entries []LogEntry
	if nextIdx <= len(n.log) {
		entries = n.log[nextIdx-1:]
	}

	args := &AppendEntriesArgs{
		Term:         n.currentTerm,
		LeaderID:     n.id,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
		LeaderCommit: n.commitIndex,
	}
	term := n.currentTerm
	n.mu.Unlock()

	var reply AppendEntriesReply
	err := callRPC(peerAddr, "RaftRPC.AppendEntries", args, &reply)
	if err != nil {
		return
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if reply.Term > n.currentTerm {
		n.becomeFollower(reply.Term)
		return
	}

	if n.state != Leader || n.currentTerm != term {
		return
	}

	if reply.Success {
		// Update nextIndex and matchIndex for this peer
		if len(entries) > 0 {
			n.matchIndex[peerAddr] = args.PrevLogIndex + len(entries)
			n.nextIndex[peerAddr] = n.matchIndex[peerAddr] + 1
			fmt.Printf("[LEADER %s] Replicated to %s (matchIndex=%d)\n",
				n.id, peerAddr, n.matchIndex[peerAddr])
		}
		// Try to advance commitIndex
		n.commitEntries()
	} else {
		// Follower rejected — back up nextIndex and retry
		if n.nextIndex[peerAddr] > 1 {
			n.nextIndex[peerAddr]--
		}
	}
}
