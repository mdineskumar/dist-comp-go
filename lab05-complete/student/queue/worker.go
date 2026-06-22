package main

// ============================================================
// Lab 05 — Event Queues and Pub/Sub
// File: worker.go  (Event Queue system)
// Role: Worker pool — multiple workers pulling concurrently
//
// TASK IN THIS FILE:
//   Task 7 — RunWorkerPool()
// ============================================================

import (
	"fmt"
	"math/rand"
	"time"
)

// ============================================================
// TASK 7 — RunWorkerPool
// ============================================================
// Launch numWorkers goroutines that each continuously:
//   1. Dequeue a task from queueName (via RPC to the broker)
//   2. "Process" it (simulate work with a short sleep)
//   3. Ack it (via RPC)
//
// This must guarantee NO TWO WORKERS process the same task —
// because Dequeue removes the task from the channel, only one
// worker can ever receive a given task. Your job is to make
// sure each worker runs this loop independently and concurrently.
//
// Steps:
//   1. For i := 0 to numWorkers-1:
//        launch a goroutine (go func(workerID string) { ... }(...))
//        each goroutine loops forever:
//          a. Call Dequeue RPC with a unique workerID (e.g. "worker-1")
//          b. Print: [WORKER workerID] Processing task=ID payload="..."
//          c. Simulate work: time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
//          d. Call Ack RPC with the task ID
//          e. Print: [WORKER workerID] Done task=ID
//   2. This function should return immediately after launching
//      the goroutines (it does not block)
//
// HINT: each goroutine needs its OWN connection or use callRPC()
//       which dials fresh each time (simpler, slightly less efficient)
//
// TODO: implement this function
func RunWorkerPool(brokerAddr, queueName string, numWorkers int) {
	// YOUR CODE HERE
	_ = rand.Intn   // remove when implementing
	_ = time.Sleep  // remove when implementing
	fmt.Println("[WORKERPOOL] Not implemented yet — complete Task 7")
}
