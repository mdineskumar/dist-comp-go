package main

// ============================================================
// Lab 03 — CAP Theorem
// File: node.go  (CP System)
// Role: Node data structures and peer health check
//
// TASKS IN THIS FILE:
//   Task 9  — NewNode()
//   Task 10 — isAlive()
// ============================================================

import (
	"fmt"
	"sync"
)

// Entry stores a value with a timestamp.
// Same structure as AP — but used differently:
// In CP, timestamps help us return the most recent value
// when reading from multiple nodes.
type Entry struct {
	Value     string
	Timestamp int64
}

// Node represents a CP (Consistent + Partition-tolerant) store.
// Key property: writes require majority acknowledgement (quorum).
// If not enough nodes are reachable, writes FAIL rather than
// risk inconsistency. Reads also contact a majority.
//
// ── SIMPLIFIED METHOD ──────────────────────────────────────
// We use simple majority voting: quorum = len(peers)/2 + 1
// This means 3 out of 5 nodes must agree for any write/read.
//
// COMING LATER (Week 9 — Consensus):
// Raft and Paxos add leader election and log replication
// so the system remains consistent even across failures,
// not just partitions. Our simplified version works for
// the lab but wouldn't survive a node crash mid-write.
// ──────────────────────────────────────────────────────────
type Node struct {
	mu     sync.RWMutex
	addr   string           // this node's "host:port"
	peers  []string         // addresses of all other nodes
	store  map[string]Entry // local key-value storage
	quorum int              // minimum nodes needed (len(peers)/2 + 1)
}

// ============================================================
// TASK 9 — NewNode
// ============================================================
// Create and return a new CP node.
//
// Steps:
//   1. Create a Node with addr, peers, empty store
//   2. Calculate quorum: len(peers)/2 + 1
//      Example: 4 peers → quorum = 4/2 + 1 = 3
//      (need 3 out of 5 total nodes including self)
//   3. Return a pointer to the node
//
// TODO: implement this function
func NewNode(addr string, peers []string) *Node {
	// YOUR CODE HERE
	return nil
}

// ============================================================
// TASK 10 — isAlive
// ============================================================
// Check if a peer node is reachable by sending a Ping RPC.
// Return true if the ping succeeds, false if it fails.
//
// Use: callRPC(peerAddr, "CPRPC.Ping", &PingArgs{}, &PingReply{})
//
// TODO: implement this function
func (n *Node) isAlive(peerAddr string) bool {
	// YOUR CODE HERE
	return false
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// String returns a human-readable description of this node
func (n *Node) String() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return fmt.Sprintf("CPNode{addr:%s peers:%d quorum:%d keys:%d}",
		n.addr, len(n.peers), n.quorum, len(n.store))
}

// printStore shows all locally stored key-value pairs
func (n *Node) printStore() {
	n.mu.RLock()
	defer n.mu.RUnlock()
	fmt.Printf("\n── Store on %s ──────────────────────────────\n", n.addr)
	if len(n.store) == 0 {
		fmt.Println("  (empty)")
	} else {
		fmt.Printf("  %d keys:\n", len(n.store))
		for k, e := range n.store {
			fmt.Printf("  key=%-20q  value=%-20q  ts=%d\n", k, e.Value, e.Timestamp)
		}
	}
	fmt.Println()
}

// countAlive returns how many peers are currently reachable
func (n *Node) countAlive() int {
	alive := 1 // count self
	for _, peer := range n.peers {
		if n.isAlive(peer) {
			alive++
		}
	}
	return alive
}
