package main

// ============================================================
// Lab 05 — Event Queues and Pub/Sub
// File: rpc.go  (Event Queue system)
// Role: RPC server and handlers
//
// TASK IN THIS FILE:
//   Task 6 — Implement RPC handlers
// ============================================================

import (
	"fmt"
	"net"
	"net/rpc"
)

// ── RPC argument / reply types ────────────────────────────

type EnqueueArgs  struct{ QueueName, Payload string }
type EnqueueReply struct{ TaskID string }

type DequeueArgs  struct{ QueueName, WorkerID string }
type DequeueReply struct{ Task Task }

type AckArgs  struct{ TaskID string }
type AckReply struct{ Success bool }

type NackArgs  struct{ TaskID string }
type NackReply struct{ Success bool }

// QueueRPC is the RPC handler
type QueueRPC struct{ qm *QueueManager }

var taskCounter int

// ============================================================
// TASK 6 — RPC Handlers
// ============================================================
//
// ── Enqueue ───────────────────────────────────────────────
// Generate a unique task ID (use generateTaskID() — provided below)
// Create a Task{ID: id, QueueName: args.QueueName, Payload: args.Payload}
// Call qm.Enqueue(args.QueueName, task)
// Set reply.TaskID = id
//
// TODO: implement
func (r *QueueRPC) Enqueue(args *EnqueueArgs, reply *EnqueueReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Dequeue ───────────────────────────────────────────────
// Call qm.Dequeue(args.QueueName, args.WorkerID) — this BLOCKS
// until a task is available
// Set reply.Task to the result
//
// TODO: implement
func (r *QueueRPC) Dequeue(args *DequeueArgs, reply *DequeueReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Ack ───────────────────────────────────────────────────
// Call qm.Ack(args.TaskID)
// Set reply.Success to the result
//
// TODO: implement
func (r *QueueRPC) Ack(args *AckArgs, reply *AckReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Nack ──────────────────────────────────────────────────
// Call qm.Nack(args.TaskID)
// Set reply.Success to the result
//
// TODO: implement
func (r *QueueRPC) Nack(args *NackArgs, reply *NackReply) error {
	// YOUR CODE HERE
	return nil
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// generateTaskID creates a unique task identifier
func generateTaskID() string {
	taskCounter++
	return fmt.Sprintf("task-%d", taskCounter)
}

// startRPCServer starts the RPC listener
func startRPCServer(qm *QueueManager, port string) error {
	handler := &QueueRPC{qm: qm}
	server := rpc.NewServer()
	if err := server.Register(handler); err != nil {
		return fmt.Errorf("register failed: %v", err)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen failed: %v", err)
	}
	fmt.Printf("[RPC] Queue server listening on port %s\n", port)
	go server.Accept(ln)
	return nil
}

// callRPC makes a synchronous RPC call
func callRPC(addr, method string, args, reply interface{}) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial %s: %v", addr, err)
	}
	defer client.Close()
	return client.Call(method, args, reply)
}
