# Lab 05 — Experiment Guide

## Before You Start — Build & Launch

Run these once from your machine (not inside a container):

```bash
# Step 1 — build the image
docker build -t lab05 -f docker/Dockerfile .

# Step 2 — start all 8 containers
docker-compose -f docker/docker-compose.yml up -d

# Step 3 — verify brokers are ready
docker logs lab05-queue-broker
# Expected: [RPC] Queue server listening on port 9000

docker logs lab05-pubsub-broker
# Expected: [RPC] Broker listening on port 9001
```

### Container Reference

| Container | Role | System |
|---|---|---|
| `lab05-queue-broker` | Runs the queue broker (port 9000) | Event Queue |
| `lab05-producer` | Run producer commands here | Event Queue |
| `lab05-worker` | Run worker pool here | Event Queue |
| `lab05-pubsub-broker` | Runs the pub/sub broker (port 9001) | Pub/Sub |
| `lab05-publisher` | Run publish commands here | Pub/Sub |
| `lab05-subscriber1/2/3` | Run subscribe commands here | Pub/Sub |

---

## Experiment A — Work Distribution

**Goal:** prove each task goes to exactly ONE worker, distributed across all 3.

```bash
# Terminal 1 — start 3 workers
docker exec -d lab05-worker /lab05/queue/queue_bin \
    -mode work -queue orders -workers 3

# Terminal 2 — produce 20 tasks
docker exec lab05-producer /lab05/queue/queue_bin \
    -mode produce -queue orders -payload order -count 20

# Watch the worker output
docker logs -f lab05-worker
```

### Expected Output
docker build -t lab05 -f docker/Dockerfile .
docker-compose -f docker/docker-compose.yml up -d

docker-compose -f docker/docker-compose.yml down

go build -o server_bin .
pkill server_bin
docker exec -it lab05-worker bash
#fixed above 

# Start workers logging inside the container
docker exec -d lab05-worker sh -c \
    '/lab05/queue/queue_bin -mode work -broker queue-broker:9000 -queue orders -workers 3 > /tmp/worker.log 2>&1'

# Enqueue some tasks (single line, type quotes manually)
docker exec lab05-producer /lab05/queue/queue_bin -mode produce -broker queue-broker:9000 -queue orders -payload "order-123" -count 5

# Check worker logs
docker exec lab05-worker cat /tmp/worker.log

```


```bash
docker exec lab05-subscriber1 pkill pubsub_bin

# Terminal 1, 2, 3 — start subscribers, redirecting output to log files
docker exec -d lab05-subscriber1 sh -c \
    '/lab05/pubsub/pubsub_bin -mode subscribe -topic news -id sub1 -port 9100 -host subscriber1 > /tmp/sub1.log 2>&1'

docker exec lab05-subscriber2 pkill pubsub_bin
docker exec -d lab05-subscriber2 sh -c \
    '/lab05/pubsub/pubsub_bin -mode subscribe -topic news -id sub2 -port 9100 -host subscriber2 > /tmp/sub2.log 2>&1'

docker exec lab05-subscriber3 pkill pubsub_bin
docker exec -d lab05-subscriber3 sh -c \
    '/lab05/pubsub/pubsub_bin -mode subscribe -topic news -id sub3 -port 9100 -host subscriber3 > /tmp/sub3.log 2>&1'

# Wait for subscriptions to register
sleep 2

# Publish
docker exec lab05-publisher /lab05/pubsub/pubsub_bin \
    -mode publish -topic news -key headline -value "Breaking news!"

# Wait briefly for delivery
sleep 1

# Check logs from inside each container
docker exec lab05-subscriber1 cat /tmp/sub1.log
docker exec lab05-subscriber2 cat /tmp/sub2.log
docker exec lab05-subscriber3 cat /tmp/sub3.log


docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -topic news -key headline -value "Breaking news!"
```


## Experiment B — Worker Crash Recovery

### step 1 - increasing the procesing sleep
### step 2 - run experiment


```bash
docker exec lab05-worker /lab05/queue/queue_bin -mode work -queue orders -workers 3 &

docker exec lab05-producer /lab05/queue/queue_bin -mode produce -queue orders -payload order -count 20

```
[WORKER worker-0] Processing task=task-1 payload=order-0
[WORKER worker-2] Processing task=task-2 payload=order-1
[WORKER worker-1] Processing task=task-3 payload=order-2
[WORKER worker-0] Done task=task-1
...
```

### What to Observe and Record

- Every task-1 through task-20 appears exactly once — no duplicates
- Tasks spread roughly evenly (~6-7 per worker) — not exactly equal since it is concurrent, not round-robin
- Total "Done" lines = 20, no tasks lost

---

## Experiment B — Worker Crash Recovery

**Goal:** prove a crashed worker's in-flight task gets redelivered automatically.

> **Why the fix is needed:** the default worker sleep is `rand.Intn(200)` milliseconds (0–200ms).
> By the time you type `pkill`, all tasks are already done and acked — nothing is in-flight to recover.
> The fix increases the sleep to 8 seconds so you have time to kill the worker mid-processing.
> The stale checker timeout is 5 seconds, so any sleep longer than 5s is enough.

### Step 1 — Increase the processing sleep

Enter the worker container and edit `worker.go`:

```bash
docker exec -it lab05-worker bash
cd /lab05/queue
nano worker.go
```

Find the sleep line and change it from:

```go
// original — too fast to kill in time
time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
```

To:

```go
// experiment B — slow enough to kill mid-processing
time.Sleep(time.Duration(8000) * time.Millisecond)
```

Save (`Ctrl+O`, `Enter`, `Ctrl+X`) and rebuild:

```bash
go build -o queue_bin .
exit
```

### Step 2 — Run the experiment (two terminals)

**Terminal 1 — start worker and watch logs:**

```bash
docker exec -d lab05-worker /lab05/queue/queue_bin \
    -mode work -queue orders -workers 1

docker logs -f lab05-worker
```

**Terminal 2 — produce tasks:**

```bash
docker exec lab05-producer /lab05/queue/queue_bin \
    -mode produce -queue orders -payload crash-test -count 3
```

The moment you see this in Terminal 1:

```
[WORKER worker-0] Processing task=task-1 payload=crash-test
```

Immediately run in Terminal 2:

```bash
docker exec lab05-worker pkill queue_bin
```

### Step 3 — Start a new worker and wait for redelivery

```bash
docker exec -d lab05-worker /lab05/queue/queue_bin \
    -mode work -queue orders -workers 1

# wait ~7 seconds for stale checker (5s timeout + 2s tick interval)
docker logs lab05-worker
```

### Expected Output

```
[QUEUE] Task task-1 appears stuck (worker may have crashed) — redelivering
[WORKER worker-0] Processing task=task-1 payload=crash-test
[WORKER worker-0] Done task=task-1
```

### What to Observe and Record

- The task that was in-flight when the worker died reappears after ~5-7 seconds
- Redelivery time: between 5s (stale timeout) and 7s (timeout + next checker tick)
- The new worker picks it up from where the crashed worker left off

### Step 4 — Restore the original sleep after the experiment

```bash
docker exec -it lab05-worker bash
cd /lab05/queue
nano worker.go
# change back to:
# time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
go build -o queue_bin .
exit
```

---

## Experiment C — Pub/Sub Fan-Out

**Goal:** prove ALL 3 subscribers receive ALL 5 events with matching sequence numbers.

```bash
# Step 1 — start all 3 subscribers
docker exec -d lab05-subscriber1 /lab05/pubsub/pubsub_bin \
    -mode subscribe -topic news -id sub1 -port 9100 -host subscriber1

docker exec -d lab05-subscriber2 /lab05/pubsub/pubsub_bin \
    -mode subscribe -topic news -id sub2 -port 9100 -host subscriber2

docker exec -d lab05-subscriber3 /lab05/pubsub/pubsub_bin \
    -mode subscribe -topic news -id sub3 -port 9100 -host subscriber3

# Step 2 — wait for subscriptions to register
sleep 2

# Step 3 — publish 5 events
docker exec lab05-publisher /lab05/pubsub/pubsub_bin \
    -mode publish -topic news -key headline -value "story" -count 5

# Step 4 — check ALL THREE subscriber logs
docker logs lab05-subscriber1
docker logs lab05-subscriber2
docker logs lab05-subscriber3
```

### Expected Output (identical on all 3)

```
[SUBSCRIBER] Received topic="news" key="headline-0" value="story-0" seq=1
[SUBSCRIBER] Received topic="news" key="headline-1" value="story-1" seq=2
[SUBSCRIBER] Received topic="news" key="headline-2" value="story-2" seq=3
[SUBSCRIBER] Received topic="news" key="headline-3" value="story-3" seq=4
[SUBSCRIBER] Received topic="news" key="headline-4" value="story-4" seq=5
```

### What to Observe and Record

- All 3 logs show identical events — fan-out is working
- Sequence numbers 1-5 match across all subscribers
- Compare to Exp A: 20 tasks were split (~7 each). Here, 5 events appear in ALL 3 subscribers (5 each, 15 total deliveries) — same message, completely different behaviour

---

## Experiment D — Late-Joining Subscriber

**Goal:** prove a subscriber that joins late misses earlier events (no replay in your system).

```bash
# Step 1 — publish 5 events with NO subscribers listening
docker exec lab05-publisher /lab05/pubsub/pubsub_bin \
    -mode publish -topic updates -key event -value "early" -count 5

# Step 2 — NOW start subscriber1 AFTER those 5 events
docker exec -d lab05-subscriber1 /lab05/pubsub/pubsub_bin \
    -mode subscribe -topic updates -id sub1 -port 9100 -host subscriber1

sleep 2

# Step 3 — publish 3 MORE events
docker exec lab05-publisher /lab05/pubsub/pubsub_bin \
    -mode publish -topic updates -key event -value "late" -count 3

# Step 4 — check what subscriber1 received
docker logs lab05-subscriber1
```

### Expected Output

```
[SUBSCRIBER] Received topic="updates" key="event-0" value="late-0" seq=6
[SUBSCRIBER] Received topic="updates" key="event-1" value="late-1" seq=7
[SUBSCRIBER] Received topic="updates" key="event-2" value="late-2" seq=8
```

The first 5 events (seq 1-5) are completely absent — they were published before the subscriber registered so the broker had nobody to deliver to and discarded them.

### What to Observe and Record

- Late subscriber received the 3 events after joining: YES
- Late subscriber received the 5 events before joining: NO — they are gone forever
- Seq numbers start at 6, not 1 — the counter kept incrementing even with no subscribers

---

## Experiment E — Same Scenario, Both Systems

**Goal:** side-by-side comparison of 100 tasks/events across 3 workers vs 3 subscribers.

```bash
# ── Event Queue side ──────────────────────────────────────────

# Start 3 workers
docker exec -d lab05-worker /lab05/queue/queue_bin \
    -mode work -queue orders -workers 3

# Produce 100 tasks
docker exec lab05-producer /lab05/queue/queue_bin \
    -mode produce -queue orders -payload order -count 100

# Wait for processing to finish, then count per worker
docker logs lab05-worker | grep "Done task" | grep "worker-0" | wc -l
docker logs lab05-worker | grep "Done task" | grep "worker-1" | wc -l
docker logs lab05-worker | grep "Done task" | grep "worker-2" | wc -l

# ── Pub/Sub side ──────────────────────────────────────────────

docker exec -d lab05-subscriber1 /lab05/pubsub/pubsub_bin \
    -mode subscribe -topic orders -id sub1 -port 9100 -host subscriber1

docker exec -d lab05-subscriber2 /lab05/pubsub/pubsub_bin \
    -mode subscribe -topic orders -id sub2 -port 9100 -host subscriber2

docker exec -d lab05-subscriber3 /lab05/pubsub/pubsub_bin \
    -mode subscribe -topic orders -id sub3 -port 9100 -host subscriber3

sleep 2

docker exec lab05-publisher /lab05/pubsub/pubsub_bin \
    -mode publish -topic orders -key order -value placed -count 100

docker logs lab05-subscriber1 | grep "Received" | wc -l
docker logs lab05-subscriber2 | grep "Received" | wc -l
docker logs lab05-subscriber3 | grep "Received" | wc -l
```

### Expected Results

| | Worker/Sub 1 | Worker/Sub 2 | Worker/Sub 3 | Total |
|---|---|---|---|---|
| Event Queue | ~33 | ~33 | ~34 | **100** |
| Pub/Sub | 100 | 100 | 100 | **300** |

Event Queue: 100 tasks split, each processed by exactly one worker.
Pub/Sub: every subscriber receives every event independently.

---

## Discussion Question Answers

### Q1 — How the queue ensured the crashed worker's task was not lost

When `Dequeue` is called, the task is removed from the channel but immediately added to the `inFlight` map with a timestamp. It is gone from the queue channel but still tracked in `inFlight`.

The `startStaleChecker` goroutine runs every 2 seconds. It locks `inFlight`, finds any task where `time.Since(startTime) > 5s`, then calls `Nack` on each — which removes it from `inFlight` and puts the task back into the channel. The next worker to call `Dequeue` receives it again.

If Ack happened **before** processing: the task would be deleted from `inFlight` immediately on Dequeue. If the worker then crashed mid-processing, the task exists in neither the channel nor `inFlight`. The stale checker has nothing to find. The task is permanently and silently lost. This breaks the at-least-once delivery guarantee.

---

### Q2 — How to add replay capability (Kafka-style)

Your broker currently holds only a subscriber list with no memory of past events. To add replay:

1. **Add a persistent log per topic** — alongside `subscribers map[string][]Subscriber`, keep `eventLog map[string][]Event`. Every `Publish` appends to the log before fan-out.
2. **Subscribers track their own offset** — when a subscriber registers, it sends a `startOffset` (0 = from the beginning, -1 = new events only).
3. **On subscribe, replay missed events** — the broker sends all events from `startOffset` to `len(log)-1` to the new subscriber before switching to live delivery.
4. **On new publish** — append to log, increment seq, deliver to all active subscribers.
5. **Log compaction** — periodically delete events older than a retention window so the log does not grow forever.

This is exactly what Kafka's partition log does. The offset is owned by the consumer, not the broker — so any consumer can replay from any point without affecting other consumers.

---

### Q3 — Which pattern for which order-processing use case

**Event Queue** for processing each order exactly once (inventory deduction, payment charge). Experimental evidence from Exp E: 100 tasks split across workers, total = 100. No order processed twice. This prevents double-charging a customer.

**Pub/Sub** for notifying multiple services when an order is placed (inventory + billing + shipping all need to react). Experimental evidence from Exp E: 100 events, each subscriber received all 100, total = 300 deliveries. All 3 services independently received every event.

Using Event Queue for multi-service notification would be wrong — only one service would receive each order. Using Pub/Sub for payment processing would be wrong — all 3 instances would charge the customer for the same order.

---

### Q4 — In-memory subscriptions, broker crash, and CAP theorem

If the broker crashes and restarts, `subscribers` is a fresh empty map — all subscriptions are lost. Every subscriber process must re-subscribe before it will receive any events again. Any events published between the crash and re-subscription are never delivered.

This makes your broker an **AP system** (Available + Partition-tolerant, sacrifices Consistency):

- **Available**: the broker responds immediately to every request with no coordination overhead — no quorum, no consensus, no disk write
- **Not Consistent/Durable**: state is lost on crash; a subscriber registered before the crash is not registered after

Compare to **Lab 03 Raft (CP)**: Raft sacrifices availability (leader election pause, majority quorum required) to guarantee consistency — all nodes agree on state and state survives crashes through replication.

Real systems handle this at both layers: Kafka stores topic metadata in ZooKeeper/KRaft (replicated, durable — CP for metadata), while Redis Pub/Sub is purely in-memory like your broker (AP, no persistence). Your broker is closest to Redis Pub/Sub.

---

## Results Recording Tables

### Event Queue Results

| Experiment | Observation |
|---|---|
| A — Tasks per worker (out of 20) | W0: ___ W1: ___ W2: ___ |
| A — Any duplicates or lost tasks? | |
| B — Redelivery time (seconds) | |
| B — Which worker got the redelivered task? | |

### Pub/Sub Results

| Experiment | Observation |
|---|---|
| C — All 3 subscribers got all 5 events? | |
| C — Sequence numbers matched? | |
| D — Late subscriber received earlier events? | |
| D — Late subscriber received later events? | |

### Pattern Comparison

| Property | Event Queue | Pub/Sub |
|---|---|---|
| Each message delivered to | ONE worker | ALL subscribers |
| Best for | Distributing work | Broadcasting events |
| New consumer sees old messages? | No | No (without replay log) |

---

## Troubleshooting

| Problem | Fix |
|---|---|
| Build fails | Run `docker build` from inside `lab05-complete` folder |
| `NewQueueManager returned nil` | Task 1 not implemented — return `&QueueManager{...}` |
| `NewBroker returned nil` | Task 8 not implemented — return `&Broker{...}` |
| Producer hangs forever | Queue channel full (size 100) — workers not running yet |
| Worker never receives tasks | Check `Dequeue` (Task 3) actually receives from channel |
| Task processed twice | Check you are not calling Ack twice or Nack incorrectly |
| Subscriber receives nothing | Check `Subscribe` (Task 9) actually appends to the map |
| Publish says `delivered to 0 subscribers` | Subscriber not registered yet — `sleep 2` after subscribing before publishing |
| Deliver never called | Check subscriber's own RPC server port matches what was registered with the broker |
| Experiment B — worker finishes before you can kill it | Increase sleep to 8s in `worker.go` as described in Experiment B Step 1 |
| Container won't start | `docker logs <container>` to see the error |
