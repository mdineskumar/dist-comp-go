package main

// ============================================================
// Lab 05 — Event Queues and Pub/Sub
// File: rpc.go  (Pub/Sub system)
// Role: Broker-side RPC server and handlers
//
// TASK IN THIS FILE:
//   Task 13 — Implement broker RPC handlers
// ============================================================

import (
	"fmt"
	"net"
	"net/rpc"
)

type SubscribeArgs  struct{ Topic, SubID, SubAddr string }
type SubscribeReply struct{}

type UnsubscribeArgs  struct{ Topic, SubID string }
type UnsubscribeReply struct{}

type PublishArgs  struct{ Topic, Key, Value string }
type PublishReply struct{ DeliveredTo int }

// BrokerRPC is the RPC handler — registered on the BROKER side
type BrokerRPC struct{ broker *Broker }

// ============================================================
// TASK 13 — Broker RPC Handlers
// ============================================================
//
// ── Subscribe ─────────────────────────────────────────────
// Call r.broker.Subscribe(args.Topic, Subscriber{ID: args.SubID, Addr: args.SubAddr})
//
// TODO: implement
func (r *BrokerRPC) Subscribe(args *SubscribeArgs, reply *SubscribeReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Unsubscribe ───────────────────────────────────────────
// Call r.broker.Unsubscribe(args.Topic, args.SubID)
//
// TODO: implement
func (r *BrokerRPC) Unsubscribe(args *UnsubscribeArgs, reply *UnsubscribeReply) error {
	// YOUR CODE HERE
	return nil
}

// ── Publish ───────────────────────────────────────────────
// Call r.broker.Publish(args.Topic, args.Key, args.Value)
// Set reply.DeliveredTo to the returned subscriber count
//
// TODO: implement
func (r *BrokerRPC) Publish(args *PublishArgs, reply *PublishReply) error {
	// YOUR CODE HERE
	return nil
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

func startBrokerRPCServer(b *Broker, port string) error {
	handler := &BrokerRPC{broker: b}
	server := rpc.NewServer()
	if err := server.Register(handler); err != nil {
		return fmt.Errorf("register failed: %v", err)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen failed: %v", err)
	}
	fmt.Printf("[RPC] Broker listening on port %s\n", port)
	go server.Accept(ln)
	return nil
}

func startSubscriberRPCServer(s *SubscriberRPC, port string) error {
	server := rpc.NewServer()
	if err := server.Register(s); err != nil {
		return fmt.Errorf("register failed: %v", err)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen failed: %v", err)
	}
	fmt.Printf("[RPC] Subscriber listening on port %s\n", port)
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
