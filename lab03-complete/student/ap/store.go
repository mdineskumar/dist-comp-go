package main

// ============================================================
// Lab 03 — CAP Theorem
// File: store.go  (AP System)
// Role: Local key-value storage operations
//
// TASKS IN THIS FILE:
//   Task 3 — Put()
//   Task 4 — Get()
//   Task 5 — merge()
// ============================================================

import "fmt"

// ============================================================
// TASK 3 — Put
// ============================================================
// Store a key-value pair on THIS node immediately.
//
// This is the AP approach — always accept writes without
// waiting for other nodes. Other nodes learn about this
// write during the next sync cycle.
//
// Steps:
//   1. Create an Entry{Value: value, Timestamp: timestamp()}
//   2. Acquire write lock (n.mu.Lock())
//   3. Store entry: n.store[key] = entry
//   4. Release lock (n.mu.Unlock())
//   5. Print: [AP] Stored key="..." value="..." ts=...
//
// TODO: implement this function
func (n *Node) Put(key, value string) {
	// YOUR CODE HERE
}

// ============================================================
// TASK 4 — Get
// ============================================================
// Retrieve a value from THIS node's local store.
//
// This is the AP approach — always return immediately from
// local state, even if it might be slightly stale.
//
// Steps:
//   1. Acquire read lock (n.mu.RLock())
//   2. Look up key in n.store
//   3. Release lock (n.mu.RUnlock())
//   4. Return value and true if found
//   5. Return "" and false if not found
//
// TODO: implement this function
func (n *Node) Get(key string) (string, bool) {
	// YOUR CODE HERE
	return "", false
}

// ============================================================
// TASK 5 — merge
// ============================================================
// Update local store with entries received from a peer.
// Keep the entry with the HIGHER timestamp (last-write-wins).
//
// For each key in incoming:
//   - If key does NOT exist locally → store it
//   - If key exists AND incoming.Timestamp > local.Timestamp
//     → update to incoming (it is more recent)
//   - If key exists AND local.Timestamp >= incoming.Timestamp
//     → keep local (our version is newer or equal)
//
// ── SIMPLIFIED METHOD ──────────────────────────────────────
// "Last-write-wins" using system timestamps is simple and
// works well when clock drift between nodes is small.
//
// COMING LATER (Week 10 — Time and Synchronisation):
// Vector clocks can detect truly concurrent writes (where
// neither node's write happened "after" the other) and
// allow application-level conflict resolution instead of
// silently discarding one write.
// ──────────────────────────────────────────────────────────
//
// HINT: acquire write lock before modifying n.store
//
// TODO: implement this function
func (n *Node) merge(incoming map[string]Entry) {
	// YOUR CODE HERE
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// snapshot returns a copy of the local store (thread-safe)
func (n *Node) snapshot() map[string]Entry {
	n.mu.RLock()
	defer n.mu.RUnlock()
	copy := make(map[string]Entry, len(n.store))
	for k, v := range n.store {
		copy[k] = v
	}
	return copy
}

// printDiff shows which keys differ between this node and another
func (n *Node) printDiff(other map[string]Entry) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	fmt.Println("\n── Consistency Check ────────────────────────────")
	allKeys := make(map[string]bool)
	for k := range n.store { allKeys[k] = true }
	for k := range other  { allKeys[k] = true }

	same, diff := 0, 0
	for k := range allKeys {
		local, lok := n.store[k]
		remote, rok := other[k]
		if lok && rok && local.Value == remote.Value {
			same++
		} else {
			diff++
			lv := "(missing)"
			rv := "(missing)"
			if lok { lv = local.Value }
			if rok { rv = remote.Value }
			fmt.Printf("  DIFF key=%-15q  local=%-15q  remote=%q\n", k, lv, rv)
		}
	}
	fmt.Printf("  Same: %d  Different: %d  Total: %d\n", same, diff, same+diff)
	fmt.Println()
}
