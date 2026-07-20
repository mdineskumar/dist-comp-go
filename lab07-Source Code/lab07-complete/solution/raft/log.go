package main

import (
	"fmt"
	"time"
)

func (n *Node) AppendEntry(command string) (bool, error) {
	n.mu.Lock()
	if n.state != Leader { n.mu.Unlock(); return false, fmt.Errorf("not the leader") }
	entry := LogEntry{Term: n.currentTerm, Index: n.lastLogIndex() + 1, Command: command}
	n.log = append(n.log, entry)
	fmt.Printf("[LEADER %s] Appended entry index=%d command=%q\n", n.id, entry.Index, command)
	n.mu.Unlock()

	n.replicateLog()

	for i := 0; i < 50; i++ {
		n.mu.Lock()
		if n.commitIndex >= entry.Index { n.mu.Unlock(); return true, nil }
		n.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	return false, fmt.Errorf("commit timeout")
}

func (n *Node) replicateLog() {
	for _, peer := range n.peers {
		go n.sendAppendEntries(peer)
	}
}

func (n *Node) commitEntries() {
	for N := n.commitIndex + 1; N <= len(n.log); N++ {
		if n.log[N-1].Term != n.currentTerm { continue }
		count := 1
		for _, peer := range n.peers {
			if n.matchIndex[peer] >= N { count++ }
		}
		if count > len(n.peers)/2 {
			n.commitIndex = N
			fmt.Printf("[LEADER %s] Committed up to index %d\n", n.id, N)
			go n.applyCommitted()
		}
	}
}

func (n *Node) applyCommitted() {
	n.mu.Lock()
	defer n.mu.Unlock()
	for n.lastApplied < n.commitIndex {
		n.lastApplied++
		if n.lastApplied <= len(n.log) {
			entry := n.log[n.lastApplied-1]
			n.applyToStateMachine(entry.Command)
		}
	}
}
