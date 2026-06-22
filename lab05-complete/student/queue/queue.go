package main

// ============================================================
// Lab 05 — Event Queues and Pub/Sub
// File: queue.go  (Event Queue system)
// Role: Core queue data structure and operations
//
// TASKS IN THIS FILE:
//   Task 1 — NewQueueManager()
//   Task 2 — Enqueue()
//   Task 3 — Dequeue()
//   Task 4 — Ack()
//   Task 5 — Nack()
// ============================================================

import (
	"fmt"
	"sync"
	"time"
)

// Task represents one unit of work in the queue
type Task struct {
	ID        string
	QueueName string
	Payload   string
	Attempts  int
}

// inFlight tracks a task that has been dequeued but not yet acked
type inFlight struct {
	task      Task
	workerID  string
	startTime time.Time
}

// QueueManager manages MULTIPLE named queues.
// Each queue name (e.g. "orders", "emails") has its own
// independent channel of tasks — producers and workers
// specify which queue they want by name.
type QueueManager struct {
	mu        sync.Mutex
	queues    map[string]chan Task    // queueName -> channel of pending tasks
	inFlight  map[string]inFlight     // taskID -> inFlight record (dequeued, awaiting ack)
	queueSize int                     // buffer size per queue channel
}

// ============================================================
// TASK 1 — NewQueueManager
// ============================================================
// Create and return a new QueueManager.
//
// Steps:
//   1. Create a QueueManager with:
//        queues:    make(map[string]chan Task)
//        inFlight:  make(map[string]inFlight)
//        queueSize: 100   (buffer size for each queue channel)
//   2. Return a pointer to it
//
// TODO: implement this function
func NewQueueManager() *QueueManager {
	// YOUR CODE HERE
	return nil
}

// getOrCreateQueue returns the channel for a queue name,
// creating it if it doesn't exist yet. Already implemented —
// use this helper inside your Enqueue/Dequeue implementations.
func (qm *QueueManager) getOrCreateQueue(name string) chan Task {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	if ch, exists := qm.queues[name]; exists {
		return ch
	}
	ch := make(chan Task, qm.queueSize)
	qm.queues[name] = ch
	fmt.Printf("[QUEUE] Created new queue: %q\n", name)
	return ch
}

// ============================================================
// TASK 2 — Enqueue
// ============================================================
// Add a task to the named queue. Called by producers.
//
// Steps:
//   1. Get the queue channel: ch := qm.getOrCreateQueue(queueName)
//   2. Send the task into the channel: ch <- task
//      (this will block if the queue is full — that's OK for now)
//   3. Print: [QUEUE] Enqueued task=ID to queue="..."
//
// TODO: implement this function
func (qm *QueueManager) Enqueue(queueName string, task Task) {
	// YOUR CODE HERE
}

// ============================================================
// TASK 3 — Dequeue
// ============================================================
// Remove and return a task from the named queue. Called by workers.
// This BLOCKS until a task is available (channels do this naturally).
//
// Steps:
//   1. Get the queue channel: ch := qm.getOrCreateQueue(queueName)
//   2. Receive a task: task := <-ch  (blocks if empty)
//   3. Record it as in-flight (awaiting ack) — use a lock:
//        qm.mu.Lock()
//        qm.inFlight[task.ID] = inFlight{task: task, workerID: workerID, startTime: time.Now()}
//        qm.mu.Unlock()
//   4. Print: [QUEUE] Dequeued task=ID by worker=workerID
//   5. Return the task
//
// TODO: implement this function
func (qm *QueueManager) Dequeue(queueName, workerID string) Task {
	// YOUR CODE HERE
	return Task{}
}

// ============================================================
// TASK 4 — Ack
// ============================================================
// Worker confirms a task was processed successfully.
// Remove it from the inFlight map — it is now permanently done.
//
// Steps:
//   1. Lock: qm.mu.Lock() / defer qm.mu.Unlock()
//   2. Check the task exists in qm.inFlight — if not, return false
//   3. Delete it: delete(qm.inFlight, taskID)
//   4. Print: [QUEUE] Acked task=ID
//   5. Return true
//
// ── WHY ACK MATTERS ───────────────────────────────────────
// Ack happens AFTER the worker finishes processing — not before.
// If we acked immediately on Dequeue, a worker crash during
// processing would lose the task forever (no redelivery).
// By acking only after success, we guarantee at-least-once
// delivery — the task is redelivered if the worker never acks.
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func (qm *QueueManager) Ack(taskID string) bool {
	// YOUR CODE HERE
	return false
}

// ============================================================
// TASK 5 — Nack
// ============================================================
// Worker reports failure — task must be redelivered.
//
// Steps:
//   1. Lock: qm.mu.Lock() / defer qm.mu.Unlock()
//   2. Look up the task in qm.inFlight — if not found, return false
//   3. Remove it from inFlight: delete(qm.inFlight, taskID)
//   4. Increment attempt count: record.task.Attempts++
//   5. Put it back in its queue: qm.queues[record.task.QueueName] <- record.task
//      (you'll need to look up the queue channel by name)
//   6. Print: [QUEUE] Nacked task=ID — redelivering (attempt N)
//   7. Return true
//
// TODO: implement this function
func (qm *QueueManager) Nack(taskID string) bool {
	// YOUR CODE HERE
	return false
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// checkStaleDeliveries finds in-flight tasks that have been
// dequeued for too long without an ack — likely the worker crashed.
// Background goroutine calls this periodically.
func (qm *QueueManager) checkStaleDeliveries(timeout time.Duration) {
	qm.mu.Lock()
	var staleIDs []string
	for id, rec := range qm.inFlight {
		if time.Since(rec.startTime) > timeout {
			staleIDs = append(staleIDs, id)
		}
	}
	qm.mu.Unlock()

	for _, id := range staleIDs {
		fmt.Printf("[QUEUE] Task %s appears stuck (worker may have crashed) — redelivering\n", id)
		qm.Nack(id)
	}
}

// startStaleChecker launches a background goroutine that
// periodically redelivers stuck tasks (worker crash detection)
func (qm *QueueManager) startStaleChecker() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			qm.checkStaleDeliveries(5 * time.Second)
		}
	}()
	fmt.Println("[QUEUE] Stale-delivery checker started (5s timeout)")
}

// QueueDepth returns how many pending tasks are in a queue
func (qm *QueueManager) QueueDepth(queueName string) int {
	qm.mu.Lock()
	ch, exists := qm.queues[queueName]
	qm.mu.Unlock()
	if !exists {
		return 0
	}
	return len(ch)
}

// InFlightCount returns how many tasks are currently being processed
func (qm *QueueManager) InFlightCount() int {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	return len(qm.inFlight)
}
