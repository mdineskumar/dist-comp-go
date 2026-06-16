package main

// ============================================================
// Lab 04 — RPC and Web Services
// File: main.go  (server)
// Role: Start all three servers on different ports
//
// TASKS IN THIS FILE:
//   Task 3 — Start the net/rpc server on port 7000
//   Task 8 — Start the gRPC server on port 7001
// ============================================================

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	pb "lab04server/proto"

	"google.golang.org/grpc"
)

const (
	RPC_PORT  = "7000"
	GRPC_PORT = "7001"
	REST_PORT = "8080"
)

func main() {
	fmt.Println("============================================")
	fmt.Println(" Lab 04 — Key-Value Store Server")
	fmt.Printf(" net/rpc → port %s\n", RPC_PORT)
	fmt.Printf(" gRPC    → port %s\n", GRPC_PORT)
	fmt.Printf(" REST    → port %s\n", REST_PORT)
	fmt.Println("============================================")

	// Shared store — all 3 servers use the same underlying data
	store := NewStore()

	// Start all three servers
	go startRPCServer(store)
	go startGRPCServer(store)
	go startRESTServer(store)

	fmt.Println("\n[MAIN] All servers running. Press Ctrl+C to stop.")
	fmt.Println("[MAIN] Use the client container to test all three.")

	// Block forever
	select {}
}

// ============================================================
// TASK 3 — Start net/rpc Server
// ============================================================
// Steps:
//  1. Create handler: handler := &KVHandler{store: store}
//  2. Register with rpc: rpc.Register(handler)
//  3. Listen on TCP port RPC_PORT:
//     ln, err := net.Listen("tcp", ":"+RPC_PORT)
//  4. Print: [net/rpc] Server listening on port ...
//  5. Accept connections: rpc.Accept(ln)
//     (this blocks — it must run in a goroutine)
//
// TODO: implement this function
func startRPCServer(store *Store) {
	handler := &KVHandler{store: store}

	rpc.Register(handler)
	ln, err := net.Listen("tcp", ":"+RPC_PORT)
	if err != nil {
		log.Fatalf("[net/rpc] failed to listen: %v\n", err)
	}
	fmt.Printf("[net/rpc] Server listening on port: %v\n", RPC_PORT)
	rpc.Accept(ln) // startRPCServer() already launched as goroutine in main
}

// ============================================================
// TASK 8 — Start gRPC Server
// ============================================================
// Steps:
//  1. Listen on TCP port GRPC_PORT:
//     ln, err := net.Listen("tcp", ":"+GRPC_PORT)
//  2. Create gRPC server: s := grpc.NewServer()
//  3. Register your implementation:
//     pb.RegisterKeyValueStoreServer(s, &GRPCServer{store: store})
//  4. Print: [gRPC] Server listening on port ...
//  5. Start serving: s.Serve(ln)
//     (this blocks — it must run in a goroutine)
//
// TODO: implement this function
func startGRPCServer(store *Store) {
	ln, err := net.Listen("tcp", ":"+GRPC_PORT)
	if err != nil {
		log.Fatalf("[gRPC] failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterKeyValueStoreServer(s, &GRPCServer{store: store})

	fmt.Printf("[gRPC] Server listening on port: %v\n", GRPC_PORT)
	s.Serve(ln)
}

// startRESTServer is already implemented — do not change
func startRESTServer(store *Store) {
	rest := &RESTServer{store: store}
	mux := http.NewServeMux()
	rest.SetupRoutes(mux)

	fmt.Printf("[REST] Server listening on port %s\n", REST_PORT)
	if err := http.ListenAndServe(":"+REST_PORT, mux); err != nil {
		log.Fatalf("[REST] Failed to start: %v", err)
	}
}
