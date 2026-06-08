package main

// ============================================================
// Lab 03 — CAP Theorem
// File: rpc.go  (AP System)
// Role: RPC server and handlers
//
// TASKS IN THIS FILE:
//   Task 8 — Complete RPC handlers: Put, Get, Sync
// ============================================================

import (
	"fmt"
	"net"
	"net/rpc"
)

// ── RPC argument / reply types ────────────────────────────

type PutArgs  struct{ Key, Value string }
type PutReply struct{}

type GetArgs  struct{ Key string }
type GetReply struct {
	Value string
	Found bool
}

type SyncArgs  struct{ Store map[string]Entry }
type SyncReply struct{ Store map[string]Entry }

type PingArgs  struct{}
type PingReply struct{}

// APRPC is the RPC handler — all remote-callable methods go here
type APRPC struct{ node *Node }

// ============================================================
// TASK 8 — RPC Handlers
// ============================================================
// Each handler below is called by a remote node or the CLI.
//
// ── Put ───────────────────────────────────────────────────
// Called when another node or CLI wants to store a key.
// Call n.Put(args.Key, args.Value)
//
// TODO: implement
func (r *APRPC) Put(args *PutArgs, reply *PutReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Get ───────────────────────────────────────────────────
// Called when another node or CLI wants to retrieve a key.
// Call n.Get(args.Key)
// Set reply.Value and reply.Found
//
// TODO: implement
func (r *APRPC) Get(args *GetArgs, reply *GetReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Sync ──────────────────────────────────────────────────
// Called during periodic sync by a peer.
// Steps:
//   1. Merge the incoming store: n.merge(args.Store)
//   2. Take a snapshot of our updated store: n.snapshot()
//   3. Set reply.Store to the snapshot
//
// This bidirectional sync means both nodes benefit from
// each sync call — not just the caller.
//
// TODO: implement
func (r *APRPC) Sync(args *SyncArgs, reply *SyncReply) error {
	// YOUR CODE HERE
	return nil
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// Ping is a health check
func (r *APRPC) Ping(args *PingArgs, reply *PingReply) error {
	return nil
}

// startRPCServer starts the RPC listener on the given port
func startRPCServer(n *Node, port string) error {
	handler := &APRPC{node: n}
	server := rpc.NewServer()
	if err := server.Register(handler); err != nil {
		return fmt.Errorf("register failed: %v", err)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen on %s failed: %v", port, err)
	}
	fmt.Printf("[RPC] AP node listening on port %s\n", port)
	go server.Accept(ln)
	return nil
}

// callRPC makes a synchronous RPC call to addr
func callRPC(addr, method string, args, reply interface{}) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial %s: %v", addr, err)
	}
	defer client.Close()
	return client.Call(method, args, reply)
}
