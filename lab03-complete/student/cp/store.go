package main

// ============================================================
// Lab 03 — CAP Theorem
// File: store.go  (CP System)
// Role: Quorum-based reads and writes
//
// TASKS IN THIS FILE:
//   Task 11 — Put()  (quorum write)
//   Task 12 — Get()  (quorum read)
// ============================================================

import (
	"fmt"
	"time"
)

// ============================================================
// TASK 11 — Put (Quorum Write)
// ============================================================
// Store a key-value pair — but only if a majority of nodes
// acknowledge the write. This guarantees consistency.
//
// Steps:
//  1. Create entry: Entry{Value: value, Timestamp: time.Now().UnixNano()}
//  2. Call broadcastWrite(key, entry) — sends write to all peers
//     returns the number of peers that acknowledged
//  3. Count total acks = peer acks + 1 (self always counts)
//  4. If total acks < n.quorum:
//     return error: "quorum not reached: got X need Y"
//     DO NOT store locally (CP property — reject if no quorum)
//  5. If quorum reached:
//     store locally: n.store[key] = entry (use write lock)
//     print: [CP] Committed key="..." value="..." (acks=X/Y)
//     return nil
//
// ── CP PROPERTY ───────────────────────────────────────────
// Notice this function RETURNS AN ERROR if quorum is not met.
// This is what makes the system Consistent — it refuses to
// write rather than risk different nodes having different values.
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func (n *Node) Put(key, value string) error {
	entry := Entry{Value: value, Timestamp: time.Now().UnixNano()}
	acks := n.broadcastWrite(key, entry)

	total_acks := acks + 1

	if total_acks < n.quorum {
		return fmt.Errorf("quorum not reached: got %d need %d", total_acks, n.quorum)
	}
	n.mu.Lock()
	n.store[key] = entry
	n.mu.Unlock()
	fmt.Printf("[CP] Committed key=%s value=%s (acks=%d/%d)", total_acks, n.quorum)
	return nil
}

// ============================================================
// TASK 12 — Get (Quorum Read)
// ============================================================
// Retrieve a value — contact a majority of nodes and return
// the most recently written value.
//
// Steps:
//  1. Call broadcastRead(key) — asks all peers for the value
//     returns a slice of Entry from peers that responded
//  2. Add local value if it exists:
//     n.mu.RLock()
//     if local, ok := n.store[key]; ok { append to results }
//     n.mu.RUnlock()
//  3. If len(results) < n.quorum:
//     return "", error: "quorum not reached for read: got X need Y"
//  4. Find the entry with the highest Timestamp in results
//  5. If no entry found: return "", error: "key not found"
//  6. Return best.Value, nil
//
// ── WHY QUORUM READ? ──────────────────────────────────────
// A quorum read ensures we always see the latest written value.
// Example: if we wrote to nodes 1,2,3 (quorum) and read from
// nodes 3,4,5 — node 3 has the latest value and will be
// included in the read result.
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func (n *Node) Get(key string) (string, error) {
	results := n.broadcastRead(key)
	n.mu.RLock()
	if local, ok := n.store[key]; ok {
		results = append(results, local)
	}
	n.mu.RUnlock()
	if len(results) < n.quorum {
		return "", fmt.Errorf("quorum not reached for read: got %d need %d", len(results), n.quorum)
	}

	return "", fmt.Errorf("not implemented")
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// timestamp returns current time as int64 nanoseconds
func timestamp() int64 {
	return time.Now().UnixNano()
}
