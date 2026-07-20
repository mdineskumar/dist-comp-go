package main

// ============================================================
// Lab 06 — Caching and Replication
// File: rpc.go
// Role: RPC server and handlers
//
// TASK IN THIS FILE:
//   Task 11 — Implement RPC handlers: Get, Set, ApplyReplica
// ============================================================

import (
	"fmt"
	"net"
	"net/rpc"
)

// ── RPC argument / reply types ────────────────────────────

type GetArgs   struct{ Key string }
type GetReply  struct{ Value string; Found bool; FromOrigin bool }

type SetArgs   struct{ Key, Value, Strategy string }
type SetReply  struct{ Success bool }

type RegisterReplicaArgs  struct{ Addr string }
type RegisterReplicaReply struct{}

type ApplyReplicaArgs  struct{ Key, Value string }
type ApplyReplicaReply struct{}

type StatsArgs  struct{}
type StatsReply struct{ Hits, Misses, Evictions int64; HitRate float64 }

// CacheRPC is the RPC handler
type CacheRPC struct {
	cache         *Cache
	originClient  *OriginClient
	replication   *ReplicationManager
}

// ============================================================
// TASK 11 — RPC Handlers
// ============================================================
//
// ── Get ───────────────────────────────────────────────────
// Try to get from cache first (c.cache.Get(args.Key))
// If cache HIT: set reply.Value, reply.Found=true, reply.FromOrigin=false
// If cache MISS:
//   1. Fetch from origin: originClient.Get(args.Key)
//   2. If found in origin: store in cache (c.cache.Set(key, value, false))
//      set reply.Value, reply.Found=true, reply.FromOrigin=true
//   3. If not in origin: reply.Found=false
//
// TODO: implement
func (r *CacheRPC) Get(args *GetArgs, reply *GetReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Set ───────────────────────────────────────────────────
// Write using the strategy specified in args.Strategy:
//   "write-through" → call WriteThrough(cache, originClient, key, value)
//   "write-back"    → call WriteBack(cache, key, value)
//   anything else   → return error "unknown strategy"
// Then replicate to all replicas: r.replication.Replicate(key, value)
// Set reply.Success = true
//
// TODO: implement
func (r *CacheRPC) Set(args *SetArgs, reply *SetReply) error {
	// YOUR CODE HERE
	return fmt.Errorf("not implemented")
}

// ── RegisterReplica ───────────────────────────────────────
// Call r.replication.RegisterReplica(args.Addr)
// Already implemented — do not change
func (r *CacheRPC) RegisterReplica(args *RegisterReplicaArgs, reply *RegisterReplicaReply) error {
	r.replication.RegisterReplica(args.Addr)
	return nil
}

// ── ApplyReplica ──────────────────────────────────────────
// Call ApplyReplica(r.cache, args.Key, args.Value)
// Used by replicas to receive writes from the primary
//
// TODO: implement
func (r *CacheRPC) ApplyReplica(args *ApplyReplicaArgs, reply *ApplyReplicaReply) error {
	// YOUR CODE HERE
	return nil
}

// ── GetStats ──────────────────────────────────────────────
// Return current cache statistics
// Already implemented — do not change
func (r *CacheRPC) GetStats(args *StatsArgs, reply *StatsReply) error {
	s := r.cache.stats
	reply.Hits = s.hits
	reply.Misses = s.misses
	reply.Evictions = s.evictions
	reply.HitRate = s.HitRate()
	return nil
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

func startRPCServer(cache *Cache, originClient *OriginClient,
	replication *ReplicationManager, port string) error {

	handler := &CacheRPC{
		cache:        cache,
		originClient: originClient,
		replication:  replication,
	}
	server := rpc.NewServer()
	if err := server.Register(handler); err != nil {
		return fmt.Errorf("register failed: %v", err)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen on %s failed: %v", port, err)
	}
	fmt.Printf("[RPC] Cache server listening on port %s\n", port)
	go server.Accept(ln)
	return nil
}

func callRPC(addr, method string, args, reply interface{}) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial %s: %v", addr, err)
	}
	defer client.Close()
	return client.Call(method, args, reply)
}
