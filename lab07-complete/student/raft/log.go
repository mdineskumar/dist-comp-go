package main

// ============================================================
// Lab 07 — Consensus (Simplified Raft)
// File: log.go
// Role: Log replication — appending, replicating, committing
//
// TASKS IN THIS FILE:
//   Task 8  — AppendEntry()
//   Task 9  — replicateLog()
//   Task 11 — commitEntries()
// ============================================================

// ── THE RAFT LOG ──────────────────────────────────────────
//
// The log is the heart of Raft. Every command (e.g. "SET city London")
// is appended as a LogEntry before being applied to the state machine.
//
// Key log invariants (guaranteed by Raft):
//   1. If two entries in different logs have the same index and term,
//      they store the same command.
//   2. If two entries in different logs have the same index and term,
//      all preceding entries are identical.
//
// Commit process:
//   1. Client sends command to Leader
//   2. Leader appends to its own log (uncommitted)
//   3. Leader sends to all followers (AppendEntries)
//   4. Once MAJORITY ack → entry is committed
//   5. Leader applies to state machine, responds to client
//   6. Leader tells followers the new commitIndex on next heartbeat
//   7. Followers apply committed entries to their state machines
// ─────────────────────────────────────────────────────────

import (
	"fmt"
	"time"
)

// ============================================================
// TASK 8 — AppendEntry
// ============================================================
// Leader receives a command from a client and starts replication.
//
// Steps:
//  1. Lock: n.mu.Lock()
//  2. Verify this node is still Leader — if not, unlock and return false, error
//  3. Create new log entry:
//     entry := LogEntry{
//     Term:    n.currentTerm,
//     Index:   n.lastLogIndex() + 1,
//     Command: command,
//     }
//  4. Append to leader's own log: n.log = append(n.log, entry)
//  5. Print: [LEADER id] Appended entry index=N command="..."
//  6. Unlock: n.mu.Unlock()
//  7. Replicate to all peers: n.replicateLog()
//  8. Wait for commit (poll commitIndex):
//     for i := 0; i < 50; i++ {   (timeout after 50 * 10ms = 500ms)
//     n.mu.Lock()
//     if n.commitIndex >= entry.Index {
//     n.mu.Unlock()
//     return true, nil
//     }
//     n.mu.Unlock()
//     time.Sleep(10 * time.Millisecond)
//     }
//  9. Return false, error "commit timeout"
//
// TODO: implement this function
func (n *Node) AppendEntry(command string) (bool, error) {
	n.mu.Lock()
	if n.state != Leader {
		n.mu.Unlock()
		return false, fmt.Errorf("node is not leader\n")
	}
	entry := LogEntry{
		Term:    n.currentTerm,
		Index:   n.lastLogIndex() + 1,
		Command: command,
	}
	n.log = append(n.log, entry)
	fmt.Printf("[LEADER %v] Appended entry index=%v command=%v\n", n.id, entry.Index, command)

	n.mu.Unlock()
	n.replicateLog()

	for i := 0; i < 50; i++ {
		n.mu.Lock()
		if n.commitIndex >= entry.Index {
			n.mu.Unlock()
			return true, nil
		}
		n.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	return false, fmt.Errorf("commit timeout\n")
}

// ============================================================
// TASK 9 — replicateLog
// ============================================================
// Send AppendEntries to ALL peers concurrently to replicate
// the latest log entry.
//
// Steps:
//  1. For each peer, launch a goroutine:
//     go n.sendAppendEntries(peer)
//     (sendAppendEntries is already implemented in election.go)
//
// That's it — just fan out to all peers concurrently.
// sendAppendEntries handles the retry logic (backing up nextIndex)
// and calls commitEntries() when enough peers have replicated.
//
// TODO: implement this function
func (n *Node) replicateLog() {
	for _, peer := range n.peers {
		go n.sendAppendEntries(peer)
	}
}

// ============================================================
// TASK 11 — commitEntries
// ============================================================
// Advance commitIndex if a majority of nodes have replicated
// an entry from the current term.
//
// This is called after each successful AppendEntries reply.
//
// Steps:
//
//	Find the highest index N such that:
//	  a. N > n.commitIndex
//	  b. n.log[N-1].Term == n.currentTerm  (only commit current term entries)
//	  c. A majority of peers have n.matchIndex[peer] >= N
//	     (i.e. count peers where matchIndex[peer] >= N, add 1 for self)
//	     if count > len(n.peers)/2 → majority
//
//	If such N exists:
//	  n.commitIndex = N
//	  Print: [LEADER id] Committed up to index N
//	  go n.applyCommitted()
//
// NOTE: caller must hold n.mu when calling this function
//
// TODO: implement this function
func (n *Node) commitEntries() {
	for N := n.commitIndex + 1; N <= len(n.log); N++ {
		if n.log[N-1].Term != n.currentTerm {
			continue
		}

		count := 1
		for _, peer := range n.peers {
			if n.matchIndex[peer] >= N {
				count++
			}
		}

		if count > len(n.peers)/2 {
			n.commitIndex = N
			fmt.Printf("[LEADER %v] Committed up to index %v\n", n.id, N)
		}
	}
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// applyCommitted applies all committed but not yet applied log
// entries to the state machine. Called after commitIndex advances.
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
