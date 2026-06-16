package main

// ============================================================
// Lab 04 — RPC and Web Services
// File: types.go  (net/rpc)
// Role: Request and reply types for the net/rpc server
//
// TASK IN THIS FILE:
//   Task 1 — Define all request and reply structs
// ============================================================

// ============================================================
// TASK 1 — Define net/rpc Request and Reply Types
// ============================================================
// net/rpc requires every method argument and return value to
// be a struct. Define the following 8 structs:
//
// PutArgs   { Key string, Value string }
// PutReply  { Success bool }
//
// GetArgs   { Key string }
// GetReply  { Value string, Found bool }
//
// DeleteArgs  { Key string }
// DeleteReply { Deleted bool }
//
// ListArgs  {}  (empty — no arguments needed)
// ListReply { Keys []string }
//
// IMPORTANT: All field names must start with a CAPITAL letter.
// net/rpc uses encoding/gob for serialisation — unexported
// (lowercase) fields are silently ignored and cause subtle bugs.
//
// TODO: define all 8 structs below

// YOUR CODE HERE
type PutArgs struct {
	Key   string
	Value string
}

type PutReply struct {
	Success bool
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Value string
	Found bool
}

type DeleteArgs struct {
	Key string
}

type DeleteReply struct {
	Deleted bool
}

type ListArgs struct{}

type ListReply struct {
	Keys []string
}
