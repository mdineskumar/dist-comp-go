package main

// ============================================================
// Lab 03 — CAP Theorem
// File: node.go  (AP System)
// Role: Node data structures and helper functions
//
// TASKS IN THIS FILE:
//   Task 1 — NewNode()
//   Task 2 — timestamp()
// ============================================================

import (
	"fmt"
	"sync"
	"time"
)

// Entry stores a value with a timestamp for conflict resolution.
//
// ── SIMPLIFIED METHOD ──────────────────────────────────────
// We use time.Now().UnixNano() as a version number.
// This works well for our lab but has one limitation:
// if two nodes write the same key at exactly the same
// nanosecond on different machines, the result is arbitrary.
//
// COMING LATER (Week 10 — Time and Synchronisation):
// Vector clocks give each node its own counter per peer.
// Concurrent writes are always detected and resolved correctly
// regardless of clock differences between machines.
// ──────────────────────────────────────────────────────────
type Entry struct {
	Value     string
	Timestamp int64 // time.Now().UnixNano() — higher = more recent
}

// Node represents an AP (Available + Partition-tolerant) store.
// Key property: every node ALWAYS accepts reads and writes,
// even during a network partition. Consistency is eventual.
type Node struct {
	mu           sync.RWMutex
	addr         string           // this node's "host:port"
	peers        []string         // addresses of all other nodes
	store        map[string]Entry // local key-value storage
	syncInterval time.Duration    // how often to sync with peers
}

// ============================================================
// TASK 1 — NewNode
// ============================================================
// Create and return a new AP node.
//
// Steps:
//   1. Create a Node with the given addr and peers
//   2. Initialise store as empty map[string]Entry
//   3. Set syncInterval to 1 * time.Second
//   4. Return a pointer to the node
//
// TODO: implement this function
func NewNode(addr string, peers []string) *Node {
	// YOUR CODE HERE
	return nil
}

// ============================================================
// TASK 2 — timestamp
// ============================================================
// Return the current time as a version number (int64).
//
// Use: time.Now().UnixNano()
//
// This is called every time we store a value so we know
// which version is newer when syncing with peers.
//
// TODO: implement this function
func timestamp() int64 {
	// YOUR CODE HERE
	return 0
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// String returns a human-readable description of this node
func (n *Node) String() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return fmt.Sprintf("APNode{addr:%s peers:%d keys:%d}",
		n.addr, len(n.peers), len(n.store))
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
