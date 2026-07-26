package main

// ============================================================
// Lab 07 — Consensus (Simplified Raft)
// File: node.go
// Role: Node state, data structures, term management
//
// TASKS IN THIS FILE:
//   Task 1 — NewNode()
//   Task 2 — becomeFollower()
//   Task 3 — becomeCandidate()
//
// ── BACKGROUND: RAFT BASICS ──────────────────────────────
// Raft nodes are always in one of three states:
//
//   Follower  → default state. Waits for heartbeats from leader.
//               If no heartbeat within election timeout → become Candidate.
//
//   Candidate → trying to become leader. Requests votes from all peers.
//               If majority vote → become Leader.
//               If hears from Leader → back to Follower.
//
//   Leader    → sends heartbeats to prevent new elections.
//               All client writes go through the leader.
//               Replicates log entries to followers.
//
// TERMS: Every election starts a new "term" (monotonically increasing
// integer). Terms act like logical clocks — a higher term always wins.
// If a node sees a message with a higher term, it immediately becomes
// a Follower, regardless of its current state.
// ─────────────────────────────────────────────────────────
// ============================================================

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Node states
const (
	Follower  = "follower"
	Candidate = "candidate"
	Leader    = "leader"
)

// LogEntry is one entry in the replicated log
type LogEntry struct {
	Term    int    // term when entry was received by leader
	Index   int    // position in the log (1-based)
	Command string // e.g. "SET city London"
}

// Node is the full Raft node structure
type Node struct {
	mu sync.Mutex

	// Identity
	id    string   // e.g. "raft-node1"
	addr  string   // e.g. "raft-node1:9400"
	peers []string // addresses of all other nodes

	// Persistent Raft state (in real Raft these survive crashes)
	currentTerm int    // latest term this node has seen
	votedFor    string // candidateID voted for in currentTerm ("" = none)
	log         []LogEntry

	// Volatile Raft state
	state       string // Follower, Candidate, or Leader
	commitIndex int    // highest log index known to be committed
	lastApplied int    // highest log index applied to state machine

	// Leader-only state (valid only when state == Leader)
	nextIndex  map[string]int // for each peer: next log index to send
	matchIndex map[string]int // for each peer: highest log index replicated

	// Election timer
	electionTimer *time.Timer

	// Key-value state machine — what Raft is actually agreeing ON
	stateMachine map[string]string
}

// ============================================================
// TASK 1 — NewNode
// ============================================================
// Create and return a new Raft node starting as a Follower.
//
// Steps:
//  1. Create a Node with:
//     id:           id
//     addr:         addr
//     peers:        peers
//     state:        Follower
//     currentTerm:  0
//     votedFor:     ""
//     log:          []LogEntry{}
//     commitIndex:  0
//     lastApplied:  0
//     nextIndex:    make(map[string]int)
//     matchIndex:   make(map[string]int)
//     stateMachine: make(map[string]string)
//  2. Print: [NODE id] Starting as Follower (term 0)
//  3. Start the election timer: n.resetElectionTimer()
//  4. Return the node
//
// TODO: implement this function
func NewNode(id, addr string, peers []string) *Node {
	node := Node{
		id:           id,
		addr:         addr,
		peers:        peers,
		state:        Follower,
		currentTerm:  0,
		votedFor:     "",
		log:          []LogEntry{},
		commitIndex:  0,
		lastApplied:  0,
		nextIndex:    make(map[string]int),
		matchIndex:   make(map[string]int),
		stateMachine: make(map[string]string),
	}
	fmt.Printf("[NODE %v] Starting as Follower (term %v)\n", id, node.currentTerm)
	node.resetElectionTimer()
	return &node
}

// ============================================================
// TASK 2 — becomeFollower
// ============================================================
// Transition this node to Follower state.
// Called when: we discover a higher term, we lose an election,
// or we receive a valid heartbeat.
//
// Steps:
//  1. Set n.state = Follower
//  2. Update n.currentTerm = term
//  3. Reset n.votedFor = ""   (new term = no vote cast yet)
//  4. Print: [NODE id] Became Follower (term N)
//  5. Reset election timer: n.resetElectionTimer()
//
// NOTE: caller must hold n.mu when calling this function
//
// TODO: implement this function
func (n *Node) becomeFollower(term int) {
	n.state = Follower
	n.currentTerm = term
	n.votedFor = ""
	fmt.Printf("[NODE %v] Became Follower (term %v)\n", n.id, n.currentTerm)
	n.resetElectionTimer()
}

// ============================================================
// TASK 3 — becomeCandidate
// ============================================================
// Transition to Candidate and start an election.
// Called when the election timer fires (no heartbeat from leader).
//
// Steps:
//  1. Lock: n.mu.Lock()
//  2. Increment term: n.currentTerm++
//  3. Set n.state = Candidate
//  4. Vote for self: n.votedFor = n.id
//  5. Print: [NODE id] Became Candidate (term N) — starting election
//  6. Reset election timer (in case this election fails)
//  7. Unlock: n.mu.Unlock()
//  8. Start election: go n.startElection()
//
// TODO: implement this function
func (n *Node) becomeCandidate() {
	// YOUR CODE HERE
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// resetElectionTimer sets a random election timeout (150-300ms).
// If no heartbeat arrives before the timer fires, becomeCandidate()
// is called. The randomness prevents all nodes timing out simultaneously.
func (n *Node) resetElectionTimer() {
	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}
	timeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
	n.electionTimer = time.AfterFunc(timeout, func() {
		n.mu.Lock()
		isFollowerOrCandidate := n.state != Leader
		n.mu.Unlock()
		if isFollowerOrCandidate {
			n.becomeCandidate()
		}
	})
}

// lastLogIndex returns the index of the last log entry (0 if empty)
func (n *Node) lastLogIndex() int {
	if len(n.log) == 0 {
		return 0
	}
	return n.log[len(n.log)-1].Index
}

// lastLogTerm returns the term of the last log entry (0 if empty)
func (n *Node) lastLogTerm() int {
	if len(n.log) == 0 {
		return 0
	}
	return n.log[len(n.log)-1].Term
}

// String returns a one-line summary of this node's state
func (n *Node) String() string {
	return fmt.Sprintf("Node{id:%s state:%s term:%d log:%d committed:%d}",
		n.id, n.state, n.currentTerm, len(n.log), n.commitIndex)
}
