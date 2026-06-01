package main

// ============================================================
// Lab 02 — Chord DHT
// File: node.go
// Role: Core data structures, hash function, ring arithmetic
//
// TASKS IN THIS FILE:
//   Task 1 — hashKey()
//   Task 2 — inRange()
// ============================================================

import (
	"crypto/sha1"
	"fmt"
	"sync"
)

const M = 8              // Number of bits in node ID
const RING_SIZE = 1 << M // Ring size = 2^8 = 256 positions (0-255)

// NodeInfo identifies a single node on the ring
type NodeInfo struct {
	ID   uint8  // Position on the ring (0–255)
	Addr string // "ip:port" e.g. "10.8.100.22:7000"
}

// Node is the full Chord node structure
type Node struct {
	mu          sync.RWMutex
	info        NodeInfo          // This node's own identity
	predecessor *NodeInfo         // Predecessor on the ring (nil if unknown)
	fingers     [M]NodeInfo       // Finger table — shortcuts across the ring
	next        int               // Finger index for fixFingers rotation
	storage     map[string]string // Local key-value data stored on this node
}

// ============================================================
// TASK 1 — hashKey
// ============================================================
// Hash a string to a position on the ring (0–255).
//
// Use SHA-1: sha1.Sum([]byte(key)) returns a [20]byte array.
// Take the first byte (index 0) of that array as the ring ID.
//
// Example:
//
//	hashKey("hello") → some value 0–255
//	hashKey("hello") always returns the same value (deterministic)
//
// HINT: import "crypto/sha1"
//
//	h := sha1.Sum([]byte(key))   → h is [20]byte
//	return h[0]                  → first byte = ring position
//
// TODO: implement this function
func hashKey(key string) uint8 {
	// _ = sha1.New() // remove this line when you implement
	h := sha1.Sum([]byte(key))
	return h[0]
}

// ============================================================
// TASK 2 — inRange
// ============================================================
// Check if id falls in the half-open arc (start, end] on the ring.
// The ring wraps around from 255 back to 0.
//
// Three cases to handle:
//  1. start < end  (normal):    id > start && id <= end
//  2. start > end  (wraps):     id > start || id <= end
//  3. start == end (full ring): return true (entire ring)
//
// Examples:
//
//	inRange(5,   0,  10)  → true   (5 is between 0 and 10)
//	inRange(15,  0,  10)  → false  (15 is outside 0-10)
//	inRange(5,  250, 10)  → true   (wraps: 5 is in 250→255→0→10)
//	inRange(200, 250, 10) → false  (200 is not in 250→10)
//	inRange(5,   5,   5)  → true   (full ring)
//
// TODO: implement this function
func inRange(id, start, end uint8) bool {
	if start == end {
		return true // full ring
	}
	if start < end {
		return id > start && id <= end
	}
	// wraparound case (start > end)
	return id > start || id <= end
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// newNode creates a Chord node at the given address.
func newNode(addr string) *Node {
	id := hashKey(addr)
	n := &Node{
		info:    NodeInfo{ID: id, Addr: addr},
		storage: make(map[string]string),
	}
	// Start as a single-node ring: point all fingers to self
	for i := range n.fingers {
		n.fingers[i] = n.info
	}
	return n
}

// successor returns this node's immediate successor (fingers[0])
func (n *Node) successor() NodeInfo {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.fingers[0]
}

// String returns a human-readable summary of this node
func (n *Node) String() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	pred := "nil"
	if n.predecessor != nil {
		pred = fmt.Sprintf("ID:%d", n.predecessor.ID)
	}
	return fmt.Sprintf("Node{ID:%d Addr:%s Successor:ID:%d Predecessor:%s}",
		n.info.ID, n.info.Addr, n.fingers[0].ID, pred)
}
