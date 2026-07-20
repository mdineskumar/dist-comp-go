package main

// ============================================================
// Lab 07 — Consensus (Simplified Raft)
// File: state.go
// Role: Key-value state machine and node status
//
// TASKS IN THIS FILE:
//   Task 12 — applyToStateMachine()
//   Task 13 — GetValue()
//   Task 14 — PrintStatus()
// ============================================================

// ── THE STATE MACHINE ─────────────────────────────────────
//
// The state machine is what Raft is actually agreeing ON.
// In our case it is a simple key-value store.
//
// Commands have the format:
//   "SET <key> <value>"  → store the key-value pair
//   "DEL <key>"          → delete the key
//
// The magic of Raft: every node applies the SAME sequence of
// committed log entries to its state machine. Because the log
// is identical on all nodes (guaranteed by Raft), all state
// machines reach the same state. This is CONSENSUS.
// ─────────────────────────────────────────────────────────

import (
	"fmt"
	"strings"
)

// ============================================================
// TASK 12 — applyToStateMachine
// ============================================================
// Apply one committed log command to the key-value state machine.
//
// Parse the command string:
//   "SET city London" → n.stateMachine["city"] = "London"
//   "DEL city"        → delete(n.stateMachine, "city")
//   anything else     → print warning and ignore
//
// Steps:
//   1. Split command by spaces: parts := strings.Fields(command)
//   2. If len(parts) == 0: return
//   3. Switch on parts[0] (uppercase):
//        "SET": if len(parts) >= 3: n.stateMachine[parts[1]] = parts[2]
//        "DEL": if len(parts) >= 2: delete(n.stateMachine, parts[1])
//        default: print warning
//   4. Print: [STATE id] Applied "command" → stateMachine now has N keys
//
// NOTE: caller must hold n.mu when calling this function
//
// TODO: implement this function
func (n *Node) applyToStateMachine(command string) {
	// YOUR CODE HERE
	_ = strings.Fields // remove when implementing
	_ = fmt.Sprintf    // remove when implementing
}

// ============================================================
// TASK 13 — GetValue
// ============================================================
// Read a value from the state machine.
//
// Steps:
//   1. Lock: n.mu.Lock() / defer n.mu.Unlock()
//   2. Look up key in n.stateMachine
//   3. Return value, true if found
//   4. Return "", false if not found
//
// TODO: implement this function
func (n *Node) GetValue(key string) (string, bool) {
	// YOUR CODE HERE
	return "", false
}

// ============================================================
// TASK 14 — PrintStatus
// ============================================================
// Print a formatted summary of this node's current state.
//
// Example output:
//   ── Node Status ──────────────────────────────────────
//     ID:           raft-node1
//     State:        leader
//     Term:         3
//     Log length:   5
//     Committed:    5
//     Applied:      5
//     State machine: 3 keys
//       city = London
//       country = UK
//
// TODO: implement this function
func (n *Node) PrintStatus() {
	// YOUR CODE HERE
}
