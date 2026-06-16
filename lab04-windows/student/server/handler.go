package main

// ============================================================
// Lab 04 — RPC and Web Services
// File: handler.go  (net/rpc)
// Role: net/rpc handler — implements the 4 RPC methods
//
// TASK IN THIS FILE:
//   Task 2 — Implement KVHandler methods
// ============================================================

import "fmt"

// KVHandler is the net/rpc handler object.
// net/rpc calls methods on this struct when a remote client
// sends an RPC request. Each method must have this exact signature:
//
//	func (h *KVHandler) MethodName(args *ArgsType, reply *ReplyType) error
//
// The method must be exported (capital letter) and return error.
type KVHandler struct {
	store *Store
}

// ============================================================
// TASK 2 — Implement the 4 RPC Handler Methods
// ============================================================
// Each method receives args from the client, performs the
// operation on h.store, fills in the reply, and returns nil.
//
// ── Put ───────────────────────────────────────────────────
// Call h.store.Put(args.Key, args.Value)
// Set reply.Success = true
// Print: [RPC] Put key="..." value="..."
//
// TODO: implement
func (h *KVHandler) Put(args *PutArgs, reply *PutReply) error {
	h.store.Put(args.Key, args.Value)
	reply.Success = true
	fmt.Printf("[RPC] Put key=%s value=%s\n", args.Key, args.Value)
	return nil
}

// ── Get ───────────────────────────────────────────────────
// Call h.store.Get(args.Key)
// Set reply.Value and reply.Found
// Print: [RPC] Get key="..." → found=true/false
//
// TODO: implement
func (h *KVHandler) Get(args *GetArgs, reply *GetReply) error {
	value, found := h.store.Get(args.Key)

	reply.Value = value
	reply.Found = found

	fmt.Printf("[RPC] Get key=%s -> found= %t\n", args.Key, found)

	return nil
}

// ── Delete ────────────────────────────────────────────────
// Call h.store.Delete(args.Key)
// Set reply.Deleted to the return value
// Print: [RPC] Delete key="..."
//
// TODO: implement
func (h *KVHandler) Delete(args *DeleteArgs, reply *DeleteReply) error {
	deleted := h.store.Delete(args.Key)
	reply.Deleted = deleted
	fmt.Printf("[RPC] Delete key=%s\n", args.Key)
	return nil
}

// ── List ──────────────────────────────────────────────────
// Call h.store.List()
// Set reply.Keys to the returned slice
// Print: [RPC] List → N keys
//
// TODO: implement
func (h *KVHandler) List(args *ListArgs, reply *ListReply) error {
	keys := h.store.List()
	reply.Keys = keys
	fmt.Printf("[RPC] List -> %d keys\n", len(keys))
	return nil
}
