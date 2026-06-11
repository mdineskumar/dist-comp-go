package main

// ============================================================
// Lab 03 — CAP Theorem
// File: quorum.go  (CP System)
// Role: Broadcast operations across all peers
//
// TASKS IN THIS FILE:
//   Task 13 — broadcastWrite()
//   Task 14 — broadcastRead()
// ============================================================

// ============================================================
// TASK 13 — broadcastWrite
// ============================================================
// Send a write to ALL peers concurrently using goroutines.
// Return the number of peers that acknowledged successfully.
//
// Steps:
//  1. Create a buffered channel: ackCh := make(chan bool, len(n.peers))
//  2. For each peer, launch a goroutine that:
//     - Calls Write RPC: callRPC(peer, "CPRPC.Write", &WriteArgs{Key:key, Entry:entry}, &WriteReply{})
//     - Sends true to ackCh if RPC succeeded, false if failed
//  3. Collect results: loop len(n.peers) times, read from ackCh
//     Count how many were true
//  4. Return the count of successful acks
//
// ── SIMPLIFIED METHOD ──────────────────────────────────────
// We use parallel goroutines + a channel to collect results.
// This works well for our lab.
//
// COMING LATER (Week 9 — Consensus):
// Raft uses a two-phase commit with a persistent log:
// Phase 1 (prepare): leader appends to log, peers ack
// Phase 2 (commit): leader commits, peers commit
// This survives crashes mid-write, which our version cannot.
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func (n *Node) broadcastWrite(key string, entry Entry) int {
	ackCh := make(chan bool)

	for _, peer := range n.peers {
		go func() {
			err := callRPC(peer, "CPRPC.Write", &WriteArgs{Key: key, Entry: entry}, &WriteReply{})
			if err != nil {
				ackCh <- false
			} else {
				ackCh <- true
			}
		}()
	}
	count := 0

	for i := 0; i < len(n.peers); i++ {
		v := <-ackCh
		if v {
			count++
		}
	}

	return count
}

// ============================================================
// TASK 14 — broadcastRead
// ============================================================
// Ask ALL peers for the value of a key concurrently.
// Return a slice of Entry from peers that responded with a value.
//
// Steps:
//  1. Create a buffered channel: resultCh := make(chan Entry, len(n.peers))
//  2. For each peer, launch a goroutine that:
//     - Calls Read RPC: callRPC(peer, "CPRPC.Read", &ReadArgs{Key:key}, &reply)
//     - If RPC succeeded AND reply.Found: send reply.Entry to resultCh
//     - Otherwise: send empty Entry{} to resultCh
//  3. Collect results: loop len(n.peers) times, read from resultCh
//     Only add to results slice if entry.Timestamp > 0 (valid entry)
//  4. Return the results slice
//
// TODO: implement this function
func (n *Node) broadcastRead(key string) []Entry {
	resultCh := make(chan Entry)

	for _, peer := range n.peers {
		go func() {
			reply := ReadReply{}
			err := callRPC(peer, "CPRPC.Read", &ReadArgs{Key: key}, &reply)
			if err != nil {
				resultCh <- Entry{}
			} else {
				resultCh <- reply.Entry
			}
		}()
	}

	results := make([]Entry, len(n.peers))

	for i := 0; i < len(n.peers); i++ {
		v := <-resultCh
		if v.Timestamp > 0 {
			results = append(results, v)
		}
	}
	return results
}
