package main

// ============================================================
// Lab 03 — CAP Theorem
// File: sync.go  (AP System)
// Role: Periodic background synchronisation with peers
//
// TASKS IN THIS FILE:
//   Task 6 — syncWith()
//   Task 7 — startSync()
// ============================================================

import (
	"fmt"
	"log"
	"time"
)

// ============================================================
// TASK 6 — syncWith
// ============================================================
// Exchange data with a single peer — send our store, receive
// theirs, merge both.
//
// Steps:
//  1. Take a snapshot of our local store (call n.snapshot())
//  2. Call the Sync RPC on peerAddr:
//     args  = SyncArgs{Store: snapshot}
//     reply = SyncReply{}
//     method = "APRPC.Sync"
//  3. If RPC fails → log the error and return
//     (AP property: never fail because a peer is unreachable)
//  4. If RPC succeeds → call n.merge(reply.Store)
//
// ── AP PROPERTY ───────────────────────────────────────────
// Notice we do NOT return an error here. If a peer is down
// we simply skip it and continue. This is what makes the
// system Available — it never refuses to operate just
// because some nodes are unreachable.
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func (n *Node) syncWith(peerAddr string) {
	snapshot := n.snapshot()
	args := SyncArgs{Store: snapshot}
	reply := SyncReply{}
	method := "APRPC.Sync"
	err := callRPC(peerAddr, method, &args, &reply)
	if err != nil {
		//log.Fatal("error: ", err)
		log.Println("Failed to reach peer:", peerAddr)
		return
	}
	n.merge(reply.Store)
}

// ============================================================
// TASK 7 — startSync
// ============================================================
// Start a background goroutine that periodically syncs with
// ALL peers.
//
// Steps:
//  1. Print: [SYNC] Background sync started (interval: ...)
//  2. Launch a goroutine that:
//     a. Creates a ticker: time.NewTicker(n.syncInterval)
//     b. On each tick: loops through n.peers and calls
//     syncWith(peer) for each one
//
// HINT:
//
//	go func() {
//	    ticker := time.NewTicker(n.syncInterval)
//	    for range ticker.C {
//	        for _, peer := range n.peers {
//	            ...
//	        }
//	    }
//	}()
//
// ── SIMPLIFIED METHOD ──────────────────────────────────────
// We broadcast to ALL peers every second. This is simple
// but uses O(N) bandwidth per sync cycle.
//
// COMING LATER (Week 6 — Message-oriented communication):
// Gossip protocols randomly select a subset of peers each
// cycle, achieving eventual consistency with O(log N)
// messages and much better scalability.
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func (n *Node) startSync() {
	fmt.Println("[SYNC] Background sync started (interval:", n.syncInterval, ")")
	go func() {
		ticker := time.NewTicker(n.syncInterval)
		for range ticker.C {
			for _, peer := range n.peers {
				n.syncWith(peer)
			}
		}
	}()
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// waitForConvergence blocks until all peers report the same
// value for the given key, or timeout is reached
func (n *Node) waitForConvergence(key string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
		localVal, ok := n.Get(key)
		if !ok {
			continue
		}
		allMatch := true
		for _, peer := range n.peers {
			var reply GetReply
			err := callRPC(peer, "APRPC.Get", &GetArgs{Key: key}, &reply)
			if err != nil || !reply.Found || reply.Value != localVal {
				allMatch = false
				break
			}
		}
		if allMatch {
			return true
		}
	}
	log.Printf("[SYNC] Convergence timeout for key %q", key)
	return false
}
