package main

// ============================================================
// Lab 03 — CAP Theorem
// File: rpc.go  (CP System)
// Role: RPC server and handlers
//
// TASKS IN THIS FILE:
//   Task 15 — Complete RPC handlers: Write, Read, Put, Get
// ============================================================

import (
	"fmt"
	"net"
	"net/rpc"
)

// ── RPC argument / reply types ────────────────────────────

type WriteArgs  struct{ Key string; Entry Entry }
type WriteReply struct{}

type ReadArgs  struct{ Key string }
type ReadReply struct {
	Entry Entry
	Found bool
}

type PutArgs  struct{ Key, Value string }
type PutReply struct{ Err string }

type GetArgs  struct{ Key string }
type GetReply struct {
	Value string
	Found bool
	Err   string
}

type PingArgs  struct{}
type PingReply struct{}

// CPRPC is the RPC handler
type CPRPC struct{ node *Node }

// ============================================================
// TASK 15 — RPC Handlers
// ============================================================
//
// ── Write ─────────────────────────────────────────────────
// Called by broadcastWrite — stores an entry locally.
// Acquire write lock, store args.Entry in n.store[args.Key],
// release lock. Return nil.
//
// NOTE: This bypasses quorum — it is called BY the quorum
// coordinator, not by the client directly.
//
// TODO: implement
func (r *CPRPC) Write(args *WriteArgs, reply *WriteReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Read ──────────────────────────────────────────────────
// Called by broadcastRead — returns local entry for a key.
// Acquire read lock, look up args.Key in n.store.
// Set reply.Entry and reply.Found. Return nil.
//
// TODO: implement
func (r *CPRPC) Read(args *ReadArgs, reply *ReadReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Put ───────────────────────────────────────────────────
// Called by CLI — runs the full quorum Put.
// Call n.Put(args.Key, args.Value)
// If error: set reply.Err = err.Error()
//
// TODO: implement
func (r *CPRPC) Put(args *PutArgs, reply *PutReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Get ───────────────────────────────────────────────────
// Called by CLI — runs the full quorum Get.
// Call n.Get(args.Key)
// If found: set reply.Value and reply.Found = true
// If error: set reply.Err = err.Error()
//
// TODO: implement
func (r *CPRPC) Get(args *GetArgs, reply *GetReply) error {
	// YOUR CODE HERE
	return nil
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// Ping is a health check
func (r *CPRPC) Ping(args *PingArgs, reply *PingReply) error {
	return nil
}

// startRPCServer starts the RPC listener on the given port
func startRPCServer(n *Node, port string) error {
	handler := &CPRPC{node: n}
	server := rpc.NewServer()
	if err := server.Register(handler); err != nil {
		return fmt.Errorf("register failed: %v", err)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen on %s failed: %v", port, err)
	}
	fmt.Printf("[RPC] CP node listening on port %s\n", port)
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
