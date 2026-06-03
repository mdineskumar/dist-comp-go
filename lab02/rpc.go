package main

// ============================================================
// Lab 02 — Chord DHT
// File: rpc.go
// Role: RPC server, handlers, and client helper
//
// TASKS IN THIS FILE:
//   Task 8 — Complete the RPC handlers:
//     FindSuccessor, GetPredecessor, Notify, Put, Get, Delete
// ============================================================

import (
	"fmt"
	"net"
	"net/rpc"
)

// ── RPC argument / reply types ────────────────────────────────

type FindSuccessorArgs struct{ ID uint8 }
type FindSuccessorReply struct{ Node NodeInfo }

type GetPredecessorArgs struct{}
type GetPredecessorReply struct{ Node *NodeInfo }

type NotifyArgs struct{ Node NodeInfo }
type NotifyReply struct{}

type PutArgs struct{ Key, Value string }
type PutReply struct{}

type GetArgs struct{ Key string }
type GetReply struct {
	Value string
	Found bool
}

type DeleteArgs struct{ Key string }
type DeleteReply struct{ Deleted bool }

type PingArgs struct{}
type PingReply struct{}

// ChordRPC is the RPC handler — all remote-callable methods are on this type
type ChordRPC struct{ node *Node }

// ============================================================
// TASK 8 — RPC Handlers
// ============================================================
// Each handler below is called by a REMOTE node over the network.
// Your job is to call the correct local function and fill in the reply.
//
// ── FindSuccessor ─────────────────────────────────────────────
// Called when another node wants to find the successor of an ID.
// Call n.findSuccessor(args.ID) and put the result in reply.Node
//
// TODO: implement
func (r *ChordRPC) FindSuccessor(args *FindSuccessorArgs, reply *FindSuccessorReply) error {
	node, err := r.node.findSuccessor(args.ID)
	if err != nil {
		return nil
	}
	reply.Node = node
	return err
}

// ── GetPredecessor ────────────────────────────────────────────
// Called by stabilize — returns this node's current predecessor.
// Read n.predecessor and assign it to reply.Node
// (reply.Node can be nil if we don't know our predecessor yet)
//
// TODO: implement
func (r *ChordRPC) GetPredecessor(args *GetPredecessorArgs, reply *GetPredecessorReply) error {
	reply.Node = r.node.predecessor
	return nil
}

// ── Notify ────────────────────────────────────────────────────
// Called when another node thinks it might be our predecessor.
// Update n.predecessor if:
//   - We have no predecessor (n.predecessor == nil), OR
//   - args.Node.ID falls in (n.predecessor.ID, n.info.ID)
//     meaning args.Node is closer to us than our current predecessor
//
// HINT: use inRange(args.Node.ID, n.predecessor.ID, n.info.ID)
//
//	use n.mu.Lock() / n.mu.Unlock()
//
// TODO: implement
func (r *ChordRPC) Notify(args *NotifyArgs, reply *NotifyReply) error {
	r.node.mu.Lock()
	defer r.node.mu.Unlock()
	if r.node.predecessor == nil || inRange(args.Node.ID, r.node.predecessor.ID, r.node.info.ID) {
		r.node.predecessor = &args.Node
	}
	return nil
}

// ── Put ───────────────────────────────────────────────────────
// Called when another node is forwarding a Put operation here.
// Call n.localPut(args.Key, args.Value)
//
// TODO: implement
func (r *ChordRPC) Put(args *PutArgs, reply *PutReply) error {
	r.node.localPut(args.Key, args.Value)
	return nil
}

// ── Get ───────────────────────────────────────────────────────
// Called when another node is forwarding a Get operation here.
// Call n.localGet(args.Key)
// Set reply.Value and reply.Found
//
// TODO: implement
func (r *ChordRPC) Get(args *GetArgs, reply *GetReply) error {
	reply.Value, reply.Found = r.node.localGet(args.Key)
	return nil
}

// ── Delete ────────────────────────────────────────────────────
// Called when another node is forwarding a Delete operation here.
// Call n.localDelete(args.Key)
// Set reply.Deleted
//
// TODO: implement
func (r *ChordRPC) Delete(args *DeleteArgs, reply *DeleteReply) error {
	reply.Deleted = r.node.localDelete(args.Key)
	return nil
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// Ping is a health check — always returns nil
func (r *ChordRPC) Ping(args *PingArgs, reply *PingReply) error {
	return nil
}

// startRPCServer starts the RPC listener on the given port
func startRPCServer(n *Node, port string) error {
	handler := &ChordRPC{node: n}
	server := rpc.NewServer()
	if err := server.Register(handler); err != nil {
		return fmt.Errorf("RPC register failed: %v", err)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen on port %s failed: %v", port, err)
	}
	fmt.Printf("[RPC] Server listening on port %s\n", port)
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
