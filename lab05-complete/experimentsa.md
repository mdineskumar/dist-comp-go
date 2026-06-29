Test the Event Queue

```bash
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


## Experiemnts

```bash
docker exec lab05-worker /lab05/queue/queue_bin -mode work -queue orders -workers 3 &

docker exec lab05-producer /lab05/queue/queue_bin -mode produce -queue orders -payload order -count 20

```



# Lab 05 Experiments — Full Guide

## Before Starting — Check Containers Are Running

```bash
docker ps
```
You should see: `lab05-worker`, `lab05-producer`, `lab05-publisher`, `lab05-subscriber1`, `lab05-subscriber2`, `lab05-subscriber3`, `queue-broker`, `pubsub-broker`

---

## Experiment B — Worker Crash Recovery

### Step 1 — Kill any old processes
```bash
docker exec lab05-worker pkill queue_bin
sleep 1
```

### Step 2 — Start 1 worker in foreground (open a NEW terminal tab)
```bash
docker exec -it lab05-worker sh -c '/lab05/queue/queue_bin -mode work -broker queue-broker:9000 -queue orders -workers 1 > /tmp/worker.log 2>&1 & tail -f /tmp/worker.log'
```

### Step 3 — Produce 5 tasks (in your ORIGINAL terminal tab)
```bash
docker exec lab05-producer /lab05/queue/queue_bin -mode produce -broker queue-broker:9000 -queue orders -payload order -count 5
```

### Step 4 — Wait 2-3 seconds then kill the worker
```bash
sleep 3
docker exec lab05-worker pkill queue_bin
```

### Step 5 — Check broker logs to see in-flight tasks
```bash
docker logs lab05-queue-broker --tail 30
```

### Step 6 — Start a NEW worker
```bash
docker exec -d lab05-worker sh -c '/lab05/queue/queue_bin -mode work -broker queue-broker:9000 -queue orders -workers 1 > /tmp/worker2.log 2>&1'
```

### Step 7 — Wait 7 seconds for stale checker, then check logs
```bash
sleep 7
docker exec lab05-worker cat /tmp/worker2.log
docker logs lab05-queue-broker --tail 30
```

### What to record
You should see in broker logs:
```
[QUEUE] Task task-X appears stuck (worker may have crashed) — redelivering
[QUEUE] Nacked task=task-X — redelivering (attempt 2)
```
And in worker2 logs:
```
[WORKER worker-0] Processing task=task-X payload=order-X
[WORKER worker-0] Done task=task-X
```

**Answer the questions:**
- Redelivery takes **5-7 seconds** (stale checker runs every 2s, timeout is 5s)
- If Ack happened BEFORE processing: a worker crash after Ack but during processing would **permanently lose the task** — no redelivery possible since the broker already considers it done

---

## Experiment C — Pub/Sub Fan-Out

### Step 1 — Kill any old subscriber processes
```bash
docker exec lab05-subscriber1 pkill pubsub_bin
docker exec lab05-subscriber2 pkill pubsub_bin
docker exec lab05-subscriber3 pkill pubsub_bin
sleep 1
```

### Step 2 — Start all 3 subscribers on topic 'news'
```bash
docker exec -d lab05-subscriber1 sh -c '/lab05/pubsub/pubsub_bin -mode subscribe -broker pubsub-broker:9001 -topic news -id sub1 -port 9100 -host subscriber1 > /tmp/sub1.log 2>&1'

docker exec -d lab05-subscriber2 sh -c '/lab05/pubsub/pubsub_bin -mode subscribe -broker pubsub-broker:9001 -topic news -id sub2 -port 9100 -host subscriber2 > /tmp/sub2.log 2>&1'

docker exec -d lab05-subscriber3 sh -c '/lab05/pubsub/pubsub_bin -mode subscribe -broker pubsub-broker:9001 -topic news -id sub3 -port 9100 -host subscriber3 > /tmp/sub3.log 2>&1'
```

### Step 3 — Wait for subscriptions to register
```bash
sleep 2
```

### Step 4 — Publish 5 events (type each manually with straight quotes)
```bash
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic news -key headline -value "Event1"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic news -key headline -value "Event2"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic news -key headline -value "Event3"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic news -key headline -value "Event4"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic news -key headline -value "Event5"
```

### Step 5 — Check all 3 subscriber logs
```bash
sleep 1
docker exec lab05-subscriber1 cat /tmp/sub1.log
docker exec lab05-subscriber2 cat /tmp/sub2.log
docker exec lab05-subscriber3 cat /tmp/sub3.log
```

### What to record
All 3 logs should show **identical output**:
```
[SUBSCRIBER] #1 topic="news" key="headline" value="Event1" seq=1
[SUBSCRIBER] #2 topic="news" key="headline" value="Event2" seq=2
[SUBSCRIBER] #3 topic="news" key="headline" value="Event3" seq=3
[SUBSCRIBER] #4 topic="news" key="headline" value="Event4" seq=4
[SUBSCRIBER] #5 topic="news" key="headline" value="Event5" seq=5
```

**Answer the questions:**
- All 3 subscribers receive all 5 events — that is fan-out
- Sequence numbers match across all 3 (same seq=1 through seq=5)
- vs Experiment A (queue): each task goes to **only one** worker. Here each event goes to **every** subscriber — completely opposite behavior

---

## Experiment D — Late-Joining Subscriber

### Step 1 — Kill any old subscribers
```bash
docker exec lab05-subscriber1 pkill pubsub_bin
sleep 1
```

### Step 2 — Publish 5 events with NO subscribers listening
```bash
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic updates -key update -value "Early1"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic updates -key update -value "Early2"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic updates -key update -value "Early3"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic updates -key update -value "Early4"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic updates -key update -value "Early5"
```

### Step 3 — NOW start subscriber1
```bash
docker exec -d lab05-subscriber1 sh -c '/lab05/pubsub/pubsub_bin -mode subscribe -broker pubsub-broker:9001 -topic updates -id sub1 -port 9100 -host subscriber1 > /tmp/sub1_late.log 2>&1'
sleep 2
```

### Step 4 — Publish 3 MORE events
```bash
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic updates -key update -value "Late1"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic updates -key update -value "Late2"
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic updates -key update -value "Late3"
```

### Step 5 — Check logs
```bash
sleep 1
docker exec lab05-subscriber1 cat /tmp/sub1_late.log
```

### What to record
You will see **only** the 3 late events — Early1 through Early5 are gone:
```
[SUBSCRIBER] #1 topic="updates" key="update" value="Late1" seq=6
[SUBSCRIBER] #2 topic="updates" key="update" value="Late2" seq=7
[SUBSCRIBER] #3 topic="updates" key="update" value="Late3" seq=8
```

**Answer the questions:**
- Early1-5 are **permanently lost** — the broker never stored them, just tried to deliver to zero subscribers and moved on
- Late1-3 are received fine
- **Kafka difference**: Kafka writes every event to a persistent log on disk. A late subscriber just sets its offset to 0 and replays from the beginning — it would receive all 8 events. This is the fundamental difference between a message broker and a log-based system

---

## Experiment E — Compare Both Systems

### Event Queue side — 100 tasks, 3 workers

```bash
# Kill old workers
docker exec lab05-worker pkill queue_bin
sleep 1

# Start 3 workers
docker exec -d lab05-worker sh -c '/lab05/queue/queue_bin -mode work -broker queue-broker:9000 -queue orders2 -workers 3 > /tmp/workers_e.log 2>&1'
sleep 1

# Produce 100 tasks
docker exec lab05-producer /lab05/queue/queue_bin -mode produce -broker queue-broker:9000 -queue orders2 -payload order -count 100

# Wait for processing
sleep 5
docker exec lab05-worker cat /tmp/workers_e.log
```

Count how many each worker processed:
```bash
docker exec lab05-worker grep "Done task" /tmp/workers_e.log | grep "worker-0" | wc -l
docker exec lab05-worker grep "Done task" /tmp/workers_e.log | grep "worker-1" | wc -l
docker exec lab05-worker grep "Done task" /tmp/workers_e.log | grep "worker-2" | wc -l
```

### Pub/Sub side — 3 subscribers, 100 events

```bash
# Kill old subscribers
docker exec lab05-subscriber1 pkill pubsub_bin
docker exec lab05-subscriber2 pkill pubsub_bin
docker exec lab05-subscriber3 pkill pubsub_bin
sleep 1

# Start 3 subscribers on topic orders
docker exec -d lab05-subscriber1 sh -c '/lab05/pubsub/pubsub_bin -mode subscribe -broker pubsub-broker:9001 -topic orders -id sub1 -port 9100 -host subscriber1 > /tmp/sub1_e.log 2>&1'
docker exec -d lab05-subscriber2 sh -c '/lab05/pubsub/pubsub_bin -mode subscribe -broker pubsub-broker:9001 -topic orders -id sub2 -port 9100 -host subscriber2 > /tmp/sub2_e.log 2>&1'
docker exec -d lab05-subscriber3 sh -c '/lab05/pubsub/pubsub_bin -mode subscribe -broker pubsub-broker:9001 -topic orders -id sub3 -port 9100 -host subscriber3 > /tmp/sub3_e.log 2>&1'
sleep 2

# Publish 100 events
docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -broker pubsub-broker:9001 -topic orders -key order -value order-data -count 100
sleep 2

# Count received events per subscriber
docker exec lab05-subscriber1 grep "SUBSCRIBER" /tmp/sub1_e.log | wc -l
docker exec lab05-subscriber2 grep "SUBSCRIBER" /tmp/sub2_e.log | wc -l
docker exec lab05-subscriber3 grep "SUBSCRIBER" /tmp/sub3_e.log | wc -l
```

### What to record

**Event Queue result:**
- 3 workers split ~100 tasks between them, roughly 33 each
- Total across all workers = exactly 100
- Each task processed by exactly **one** worker

**Pub/Sub result:**
- Each subscriber receives all 100 events
- Total across all subscribers = 300 (100 × 3)

**Which is correct for inventory + billing + shipping?**
Pub/Sub is correct. When an order is placed, ALL THREE services need to know about it — inventory needs to reserve stock, billing needs to charge, shipping needs to prepare. Using a queue would mean only one service gets each order, which would break the system. Pub/Sub fan-out ensures every service receives every event independently.