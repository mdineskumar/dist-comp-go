package main

// ============================================================
// Lab 04 — RPC and Web Services
// File: rpc_client.go  (net/rpc client)
// Role: Connect to net/rpc server and call methods
//
// TASK IN THIS FILE:
//   Task 4 — Implement RPCClient
// ============================================================

// ── HOW net/rpc CLIENT WORKS ──────────────────────────────
//
// net/rpc client workflow:
//   1. rpc.Dial("tcp", "server:7000")  → connect
//   2. client.Call("Handler.Method", args, &reply)  → call
//   3. client.Close()  → disconnect
//
// "KVHandler.Put" means:
//   - Type registered with rpc.Register: KVHandler
//   - Method name: Put
//
// ──────────────────────────────────────────────────────────

import (
	"fmt"
	"net/rpc"
)

// RPCClient wraps a net/rpc connection
type RPCClient struct {
	addr   string
	client *rpc.Client
}

// ============================================================
// TASK 4 — Implement RPCClient
// ============================================================
//
// ── Connect ───────────────────────────────────────────────
// Dial the server and store the connection in r.client
// Use: rpc.Dial("tcp", r.addr)
// Print: [net/rpc] Connected to addr
//
// TODO: implement
func (r *RPCClient) Connect() error {
	client, err := rpc.Dial("tcp", r.addr)
	r.client = client

	if err != nil {
		return fmt.Errorf("client connection %s: %v", r.addr, err.Error())
	}
	fmt.Printf("[net/rpc] Connected to addr: ", r.addr)
	return nil
}

// ── Put ───────────────────────────────────────────────────
// Call "KVHandler.Put" with PutArgs{Key: key, Value: value}
// Return error if call fails or Success is false
//
// TODO: implement
func (r *RPCClient) Put(key, value string) error {
	reply := PutReply{}
	err := callRPC(r.client, "KVHandler.Put", &PutArgs{Key: key, Value: value}, &reply)

	if err != nil || !reply.Success {
		return fmt.Errorf("call failed: %v", err)
	}
	return nil
}

// ── Get ───────────────────────────────────────────────────
// Call "KVHandler.Get" with GetArgs{Key: key}
// Return reply.Value, reply.Found, error
//
// TODO: implement
func (r *RPCClient) Get(key string) (string, bool, error) {
	reply := GetReply{}
	err := callRPC(r.client, "KVHandler.Get", &GetArgs{Key: key}, &reply)

	return reply.Value, reply.Found, err
}

// ── Delete ────────────────────────────────────────────────
// Call "KVHandler.Delete" with DeleteArgs{Key: key}
// Return reply.Deleted, error
//
// TODO: implement
func (r *RPCClient) Delete(key string) (bool, error) {
	reply := DeleteReply{}

	err := callRPC(r.client, "KVHandler.Delete", &DeleteArgs{Key: key}, &reply)

	return reply.Deleted, err
}

// ── List ──────────────────────────────────────────────────
// Call "KVHandler.List" with ListArgs{}
// Return reply.Keys, error
//
// TODO: implement
func (r *RPCClient) List() ([]string, error) {
	reply := ListReply{}
	err := callRPC(r.client, "KVHandler.List", &ListArgs{}, &reply)
	return reply.Keys, err
}

// ── Close ─────────────────────────────────────────────────
// Close the connection: r.client.Close()
//
// TODO: implement
func (r *RPCClient) Close() {
	r.client.Close()
}

// NewRPCClient creates a new RPCClient for the given address
func NewRPCClient(addr string) *RPCClient {
	return &RPCClient{addr: addr}
}

func callRPC(client *rpc.Client, method string, args, reply interface{}) error {
	defer client.Close()
	return client.Call(method, args, reply)
}
