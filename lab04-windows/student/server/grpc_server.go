package main

// ============================================================
// Lab 04 — RPC and Web Services
// File: grpc_server.go  (gRPC)
// Role: gRPC server — implements the generated service interface
//
// TASK IN THIS FILE:
//   Task 6 — Implement GRPCServer methods
// ============================================================

// ── HOW gRPC WORKS ────────────────────────────────────────
//
// 1. You wrote kvstore.proto — the service contract
// 2. protoc read the .proto file and generated two Go files:
//      proto/kvstore.pb.go       ← message structs (PutRequest etc.)
//      proto/kvstore_grpc.pb.go  ← service interface + client stub
//
// 3. You implement the server interface defined in kvstore_grpc.pb.go
//    The interface looks like:
//
//    type KeyValueStoreServer interface {
//        Put(context.Context, *pb.PutRequest) (*pb.PutResponse, error)
//        Get(context.Context, *pb.GetRequest) (*pb.GetResponse, error)
//        Delete(context.Context, *pb.DeleteRequest) (*pb.DeleteResponse, error)
//        List(context.Context, *pb.ListRequest) (*pb.ListResponse, error)
//    }
//
// 4. gRPC handles all the networking, serialisation, and HTTP/2 transport
//    You just implement the business logic — same as net/rpc but with
//    a defined contract and language-independent protocol
// ──────────────────────────────────────────────────────────

import (
	"context"
	"fmt"

	pb "lab04server/proto"
)

// GRPCServer implements the generated KeyValueStoreServer interface.
// The embedded UnimplementedKeyValueStoreServer provides default
// implementations that return "not implemented" errors — you override
// them one by one as you complete the tasks.
type GRPCServer struct {
	pb.UnimplementedKeyValueStoreServer
	store *Store
}

// ============================================================
// TASK 6 — Implement gRPC Server Methods
// ============================================================
// Each method receives a context and a protobuf request message.
// Call the corresponding store method and return a protobuf reply.
//
// Notice how similar this is to Task 2 (net/rpc handlers) —
// the only differences are:
//   - Methods take a context.Context (for cancellation/deadlines)
//   - Args/replies are protobuf message structs (not your own structs)
//   - Methods return the reply directly (not via pointer parameter)
//
// ── Put ───────────────────────────────────────────────────
// Call s.store.Put(req.Key, req.Value)
// Return &pb.PutResponse{Success: true}, nil
// Print: [gRPC] Put key="..." value="..."
//
// TODO: implement
func (s *GRPCServer) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	s.store.Put(req.Key, req.Value)
	fmt.Printf("[gRPC] Put key=%v value=%v\n", req.Key, req.Value)
	return &pb.PutResponse{Success: true}, nil
}

// ── Get ───────────────────────────────────────────────────
// Call s.store.Get(req.Key)
// Return &pb.GetResponse{Value: val, Found: ok}, nil
// Print: [gRPC] Get key="..." → found=true/false
//
// TODO: implement
func (s *GRPCServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	val, ok := s.store.Get(req.Key)
	fmt.Printf("[gRPC] Get key=%v -> found=%v\n", req.Key, ok)
	return &pb.GetResponse{Value: val, Found: ok}, nil
}

// ── Delete ────────────────────────────────────────────────
// Call s.store.Delete(req.Key)
// Return &pb.DeleteResponse{Deleted: deleted}, nil
// Print: [gRPC] Delete key="..."
//
// TODO: implement
func (s *GRPCServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	deleted := s.store.Delete(req.Key)
	fmt.Printf("[gRPC] Delete key=%v\n", req.Key)
	return &pb.DeleteResponse{Deleted: deleted}, nil
}

// ── List ──────────────────────────────────────────────────
// Call s.store.List()
// Return &pb.ListResponse{Keys: keys}, nil
// Print: [gRPC] List → N keys
//
// TODO: implement
func (s *GRPCServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	keys := s.store.List()
	fmt.Printf("[gRPC] List -> %v keys\n", len(keys))
	return &pb.ListResponse{Keys: keys}, nil
}
