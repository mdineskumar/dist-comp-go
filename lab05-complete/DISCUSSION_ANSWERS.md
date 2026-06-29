# Lab 05 — Discussion Question Answers

---

## Q1 — Worker Crash and At-Least-Once Delivery

> In Experiment B, a worker crashed before acking a task. Explain how your queue ensured the task wasn't lost. What specifically would happen if Ack happened BEFORE processing instead of after?

### How the task was not lost

- When a worker calls `Dequeue`, the task is removed from the queue channel
- But it is immediately saved into an `inFlight` map with a timestamp — so the broker still knows about it
- A background checker runs every 2 seconds and looks at every task in `inFlight`
- If any task has been sitting there for more than 5 seconds without an Ack, the checker assumes the worker crashed
- It calls `Nack` on that task — which removes it from `inFlight` and puts it back into the queue channel
- The next available worker picks it up and processes it again
- This guarantees the task is never silently lost — it will always be retried

### What if Ack happened BEFORE processing

- The task would be deleted from `inFlight` immediately after Dequeue, before any work is done
- If the worker then crashes during processing, the task is gone from both the channel and the `inFlight` map
- The background checker has nothing to find — it cannot redeliver what it cannot see
- The task is permanently lost with no error or warning
- This is why Ack must happen AFTER processing — it is the confirmation that the work is actually done

---

## Q2 — Late-Joining Subscriber and Replay

> In Experiment D, a late-joining subscriber missed earlier messages. Real systems like Kafka solve this with a persistent log and offsets. Describe at a high level how you would add this capability to your Pub/Sub broker.

### The problem with your current broker

- Your broker only keeps a list of who is currently subscribed
- When `Publish` is called, it delivers to whoever is subscribed right now and forgets the event
- There is no storage of past events — once delivered (or undelivered), the event is gone forever
- A subscriber that joins late has no way to ask for missed events

### How to add replay (high level)

- **Store every published event** — keep an append-only list of events per topic inside the broker, like a diary that never gets erased
- **Give every event a position number** — this is the sequence number (seq) you already have, acting as an offset
- **Each subscriber tracks its own position** — when a subscriber joins, it tells the broker "I want events from position 0" (all history) or "I want only new events from now"
- **On subscribe, send missed events first** — the broker loops through the stored events from the requested position and delivers them before switching to live delivery
- **On new publish** — save the event to the log AND deliver to active subscribers as normal
- **Clean up old events periodically** — delete events older than a set time (e.g. 7 days) so storage does not grow forever

### Why Kafka does this well

- Kafka stores all events on disk so they survive crashes
- Every consumer remembers its own offset (position) independently
- One slow consumer does not affect others — it just falls behind and catches up later
- Any consumer can rewind to offset 0 and replay all history at any time

---

## Q3 — Choosing the Right Pattern

> Based on Experiment E, which pattern fits: (a) process each order exactly once, (b) notify multiple services when an order is placed?

### Event Queue for processing each order exactly once

- In Experiment E, 100 tasks were split across 3 workers — each task was processed by exactly one worker, total = 100
- No task was processed twice, no task was skipped
- This is the right pattern for payment processing, inventory deduction, or any action that must only happen once
- If you used Pub/Sub instead, all 3 worker instances would receive every order and you would charge the customer 3 times

### Pub/Sub for notifying multiple services

- In Experiment E, 100 events were published and each of the 3 subscribers received all 100, total = 300 deliveries
- Every service independently received every event
- This is the right pattern when multiple services (inventory, billing, shipping) all need to react to the same order
- If you used Event Queue instead, only one service would receive each order and the other two would never know about it

### Summary

| Use Case | Pattern | Why |
|---|---|---|
| Charge the customer once | Event Queue | One task, one worker |
| Tell 3 services about the order | Pub/Sub | One event, all subscribers |

---

## Q4 — In-Memory State and CAP Theorem

> What happens to subscriptions if the broker crashes? Is your broker more like an AP or CP system?

### What happens when the broker crashes

- All subscriber information lives only in memory inside the `subscribers` map
- If the broker process crashes and restarts, that map is completely empty — a fresh start
- Every subscriber that was registered is now gone from the broker's perspective
- All subscriber processes must re-subscribe before they can receive any events again
- Any events published between the crash and re-subscription are never delivered to anyone

### Your broker is an AP system

- **A = Available** — the broker responds to every request immediately with no waiting or voting with other nodes
- **P = Partition-tolerant** — it keeps working even if some network connections break
- **Not C = Not Consistent** — state is lost on crash, different clients may see different views after a restart

### Compare to Lab 03 Raft (CP system)

- Raft is a CP system — it sacrifices availability for consistency
- Before accepting a write, Raft waits for a majority of nodes to agree (quorum)
- If the leader crashes, the cluster pauses and holds an election before accepting new requests
- But once a write is accepted, it is guaranteed to survive crashes — it is replicated across nodes
- Your broker does the opposite — it never waits for anyone, so it is fast and always available, but loses everything on a crash

### Real-world comparison

| System | Type | Why |
|---|---|---|
| Your broker | AP | Fast, in-memory only, loses state on crash |
| Redis Pub/Sub | AP | Same — in-memory, no persistence |
| Kafka | Closer to CP | Stores events on disk, replicated, survives crashes |
| Raft (Lab 03) | CP | Majority quorum required, consistent but slower |
