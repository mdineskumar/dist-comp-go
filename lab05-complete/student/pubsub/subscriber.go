package main

// ============================================================
// Lab 05 — Event Queues and Pub/Sub
// File: subscriber.go  (Pub/Sub system)
// Role: Subscriber-side RPC server — receives delivered events
//
// TASK IN THIS FILE:
//   Task 12 — Deliver()
// ============================================================

import "fmt"

// SubscriberRPC is registered on the SUBSCRIBER side, not the broker.
// When the broker calls Publish, it makes an RPC call to EACH
// subscriber's Deliver method — this is how fan-out actually
// reaches each subscriber's process.
type SubscriberRPC struct {
	receivedEvents chan Event
}

type DeliverArgs struct{ Event Event }
type DeliverReply struct{}

// ============================================================
// TASK 12 — Deliver
// ============================================================
// Called BY THE BROKER when an event is published to a topic
// this subscriber is subscribed to.
//
// Steps:
//   1. Send the event into the channel: s.receivedEvents <- args.Event
//   2. Print: [SUBSCRIBER] Received topic="..." key="..." value="..." seq=N
//
// TODO: implement this function
func (s *SubscriberRPC) Deliver(args *DeliverArgs, reply *DeliverReply) error {
	// YOUR CODE HERE
	return nil
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

func NewSubscriberRPC() *SubscriberRPC {
	return &SubscriberRPC{receivedEvents: make(chan Event, 100)}
}

func (s *SubscriberRPC) PrintEvents() {
	count := 0
	for event := range s.receivedEvents {
		count++
		fmt.Printf("[SUBSCRIBER] #%d topic=%-12q key=%-12q value=%-20q seq=%d\n",
			count, event.Topic, event.Key, event.Value, event.Seq)
	}
}
