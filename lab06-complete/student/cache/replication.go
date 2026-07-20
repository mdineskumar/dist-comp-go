package main

// ============================================================
// Lab 06 — Caching and Replication
// File: replication.go
// Role: Primary-replica replication
//
// TASKS IN THIS FILE:
//   Task 8  — RegisterReplica()
//   Task 9  — Replicate()
//   Task 10 — ApplyReplica()
// ============================================================

import (
	"fmt"
	"sync"
)

// ReplicationManager tracks all replica addresses and
// handles propagating writes from primary to replicas
type ReplicationManager struct {
	mu       sync.RWMutex
	replicas []string // replica addresses e.g. ["cache-replica1:9200"]
}

// NewReplicationManager creates an empty replication manager
func NewReplicationManager() *ReplicationManager {
	return &ReplicationManager{
		replicas: []string{},
	}
}

// ============================================================
// TASK 8 — RegisterReplica
// ============================================================
// Register a new replica address with the primary.
//
// Steps:
//  1. Lock: rm.mu.Lock() / defer rm.mu.Unlock()
//  2. Append addr to rm.replicas
//  3. Print: [REPLICATION] Registered replica: addr
//
// TODO: implement this function
func (rm *ReplicationManager) RegisterReplica(addr string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.replicas = append(rm.replicas, addr)
	fmt.Printf("[REPLICATION] Registered replica: %v\n", addr)
}

// ============================================================
// TASK 9 — Replicate
// ============================================================
// Push a write (key, value) to ALL registered replicas.
// Called by the primary after every successful write.
//
// Steps:
//  1. Read replica list (use read lock: rm.mu.RLock())
//  2. For each replica, send concurrently using goroutines:
//     go func(addr string) {
//     callRPC(addr, "CacheRPC.ApplyReplica",
//     &ApplyReplicaArgs{Key: key, Value: value}, &ApplyReplicaReply{})
//     }(replica)
//  3. Print: [REPLICATION] Replicated key="..." to N replicas
//
// NOTE: use goroutines so a slow replica does not block the primary
// NOTE: do NOT wait for replicas — fire and forget (eventual consistency)
//
// TODO: implement this function
func (rm *ReplicationManager) Replicate(key, value string) {
	rm.mu.RLock()
	for _, replica := range rm.replicas {
		go func(addr string) {
			callRPC(addr, "CacheRPC.ApplyReplica", &ApplyReplicaArgs{Key: key, Value: value}, &ApplyReplicaReply{})
		}(replica)
	}

	fmt.Printf("[REPLICATION] replicated key=%v to %v replicas\n", key, len(rm.replicas))
}

// ============================================================
// TASK 10 — ApplyReplica
// ============================================================
// Apply an incoming replicated write to THIS replica's cache.
// Called on the REPLICA side when it receives a Replicate RPC.
//
// Steps:
//  1. Call c.Set(key, value, false)
//     dirty=false because this came from the primary — already in sync
//  2. Print: [REPLICA] Applied key="..." value="..."
//
// TODO: implement this function
func ApplyReplica(c *Cache, key, value string) {
	c.Set(key, value, false)
	fmt.Printf("[REPLICATION] Applied key=%v value=%v\n", key, value)
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// ReplicaCount returns how many replicas are registered
func (rm *ReplicationManager) ReplicaCount() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.replicas)
}
