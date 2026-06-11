# Lab 03 — CAP Theorem: Answers, Results & Comparison

> MSc Computer Science — Distributed Computing
> AP vs CP key-value stores (Go + Docker)

---

## ⚠️ Two things in the current code that affect the results

Before reading the tables below, note two deviations in the present implementation. The
answers/tables describe what a **correct** implementation produces; the notes flag what *you*
will actually observe until these are fixed.

| Issue | Location | Effect on experiments |
|---|---|---|
| `syncInterval` is `10 * time.Second` | `student/ap/node.go` → `NewNode()` | Task 1 asks for **1 s**. With 10 s, AP convergence takes up to **10 s**, so Experiment B's "get after 3 seconds" will **still be stale**. Change to `1 * time.Second` to match the spec. |
| `Read` RPC handler acquires `RLock()` but never releases it | `student/cp/rpc.go` → `Read()` | Add `defer r.node.mu.RUnlock()`. Without it, after the first CP `get`, every node holds a permanent read-lock; the next `put`'s `Write` RPC blocks on `Lock()` forever → CP hangs. Must be fixed before CP experiments A/B/D work. |

---

## Discussion Questions (Q26–Q29)

### Q26 — AP stale read: convergence time, controlling line, and the effect of 10 s

**How long to converge.** The AP system propagates data through a periodic background sync, so a
value converges across all nodes within **one sync cycle** — i.e. the `syncInterval`. With the
spec value of `1 * time.Second`, a write on node1 reaches node5 within roughly **1 second** (up to
~2 s in the worst case, when the write lands just after a tick has fired). That is why the lab's
"wait 3 seconds, then get" reliably returns the correct value.

**Which line controls it.** The pacing is set in `sync.go`, `startSync()`:

```go
ticker := time.NewTicker(n.syncInterval)   // <-- controls how often peers exchange data
for range ticker.C {
    for _, peer := range n.peers {
        n.syncWith(peer)
    }
}
```

The `time.NewTicker(n.syncInterval)` line is what fires the sync. The interval *value* itself is
set in `node.go` `NewNode()` (`syncInterval: 1 * time.Second` per the spec).

**If you changed it to 10 seconds.** Nodes would exchange data only every 10 s, so the
**staleness window grows to up to 10 seconds**. A `get` issued within those 10 s returns a stale
(or "not found") value. This is the classic AP trade-off: a *larger* interval means **less network
and CPU overhead** (fewer sync messages — O(N) per cycle) but **weaker freshness** (longer
inconsistency window); a *smaller* interval gives faster convergence at the cost of more traffic.

> **Note on your code:** `NewNode()` currently sets `10 * time.Second`. So as written, Experiment B's
> "get after 3 seconds" will **still be stale** — the first sync tick has not fired yet. Set it back
> to `1 * time.Second` to get the intended behaviour.

---

### Q27 — Why quorum = 3, not 2? (`len(peers)/2 + 1`)

**Why 3.** With 5 total nodes, `len(peers)/2 + 1 = 4/2 + 1 = 3` is the **majority**. A majority
quorum guarantees the **overlap property**: any two quorums share at least one node, because
`3 + 3 = 6 > 5`. Concretely, if a write was committed to some set of 3 nodes (W = 3) and a later
read contacts any 3 nodes (R = 3), then `W + R = 6 > N = 5`, so the read set **must** intersect the
write set in at least one node — that node holds the latest value, and the read (which picks the
highest timestamp) returns it. This is what makes the system **Consistent**.

**What breaks if you set quorum = 2.**

1. **Split-brain / divergent writes.** Two *disjoint* groups can each reach "quorum." During a
   partition, `{node1, node2}` and `{node3, node4}` each have 2 reachable nodes, so each side
   accepts a write for the same key — `x = A` on one side, `x = B` on the other. Now two parts of
   the cluster hold **different committed values**. Consistency is gone — exactly what CP is
   supposed to prevent.
2. **Reads can miss the latest write.** With W = 2 and R = 2, `W + R = 4 < 5`, so the overlap
   guarantee no longer holds. A read of 2 nodes can completely miss the 2 nodes that took the most
   recent write, returning a **stale** value.

In short, quorum = 2 violates `W + R > N`, destroying both the no-split-brain and the
read-sees-latest-write guarantees. Majority (3) is the smallest quorum that preserves them.


---

### Q28 — `time.Now().UnixNano()` as a version number: where it goes wrong, and how vector clocks fix it

**The failure: physical clock skew.** Different machines' wall clocks are never perfectly
synchronised; they drift and are corrected (NTP, VM pauses) independently. Last-write-wins by
wall-clock timestamp assumes a higher timestamp = a later write — which is false under skew.

*Concrete scenario.* Suppose node B's clock runs **5 seconds fast**.

1. At real time 10:00:00 a client writes `balance = 100` on **node B** → it is stamped
   `10:00:05` (B is 5 s fast).
2. At real time 10:00:02 — **two seconds later, causally after** — another client writes
   `balance = 120` on **node A** (correct clock) → stamped `10:00:02`.
3. During sync, `merge()` keeps the higher timestamp: `10:00:05` (`= 100`) beats `10:00:02`
   (`= 120`). The system keeps **100** and silently discards the genuinely newer **120**.

The later, correct value is lost purely because of clock skew. (A related problem: two truly
concurrent writes are ordered arbitrarily, and the "loser" is dropped with no way to detect that a
conflict even happened.)

**How vector clocks fix it.** Each node keeps a vector of per-node counters and increments its own
entry on each event, attaching the vector to every write. Comparing two vectors reveals their
**causal relationship** rather than relying on wall-clock time:

- If vector V1 is entry-wise ≤ V2 (and not equal), then V1 *happened-before* V2 → V2 is genuinely
  newer and wins. Clock skew is irrelevant — only causality matters.
- If neither dominates the other, the writes are **concurrent**. Vector clocks **detect** this
  conflict instead of silently overwriting, so the application can resolve it deterministically or
  keep both versions as siblings (as Dynamo/Riak do).

So vector clocks (Week 10) replace "whoever has the bigger clock reading wins" with "whoever is
causally later wins, and genuine conflicts are surfaced, not lost."

---

### Q29 — Hospital patient records vs Twitter like counter: which is CP, which is AP?

**Hospital patient records → CP.** Medical data (allergies, current medication, dosage, blood
type) must **never be stale or divergent** — a clinician acting on out-of-date or conflicting data
can harm a patient. My Experiment B shows the AP store returns a **stale value immediately after a
write** (node5 returned the old/"not found" value before the sync cycle), which is unacceptable
here. Experiment D shows CP's safety mechanism: when a majority is unreachable it **refuses the
write** with `quorum not reached: got 2 need 3` rather than risk inconsistency. For patient
records, "stop and refuse" is far safer than "answer with possibly-wrong data." Consistency >
availability ⇒ **CP**.

**Twitter like counter → AP.** A like count being off by a few for a second is harmless, while
users expect the button to **always work instantly**. Experiment B shows AP answers
**immediately** and then converges within a sync cycle — perfect for a like counter, where low
latency and always-on writes matter more than exactness. Experiment D shows AP keeps **accepting
writes even with 3 nodes down**, so users can always like even during a partition. Availability >
strict consistency ⇒ **AP**.

> **Honest nuance:** this lab's AP store uses last-write-wins, which would *lose* some concurrent
> likes (two simultaneous increments → one overwrites the other). A production like counter uses a
> **CRDT counter** (a conflict-free merge that sums increments) on top of the same AP availability
> model. The CAP *family* choice (AP) is still correct; only the merge function needs upgrading.

---

## Results Recording

### Experiment Results

These are the expected outcomes for a correct implementation (`syncInterval = 1 s`, CP `Read`
lock fixed). Notes flag where your current code differs.

| Experiment | AP Result | CP Result |
|---|---|---|
| **A** — All gets succeeded? | ✅ Yes — all 5 gets return the correct values after the sync cycle (≈ 1 s; *up to 10 s with current code*) | ✅ Yes — all 5 gets succeed **immediately** with correct values (quorum read always returns latest) |
| **B** — Immediate read from node5 | ❌ **Stale** — `not found` / old value (sync hasn't propagated yet) | ✅ `1` — latest value returned immediately (always consistent) |
| **B** — Read after 3 seconds | ✅ `1` — converged after one sync cycle *(with `syncInterval = 1 s`; **still stale** at 3 s with current 10 s setting)* | ✅ `1` — still correct (every read is consistent) |
| **C** — Write with 2 nodes down? (3 up) | ✅ Yes — AP always accepts the write | ✅ Yes — `self + 2 peers = 3 = quorum`, write commits |
| **D** — Write with 3 nodes down? (2 up) | ✅ Yes — AP still accepts the write (stays available) | ❌ No — write rejected (only 2 reachable < quorum 3) |
| **D** — Error message (if any) | N/A (write succeeded) | `quorum not reached: got 2 need 3` |
| **D** — Recovery time (seconds) | ≈ **1–3 s** (one sync cycle) to reconverge after restart *(up to ~10 s with current 10 s interval)* | N/A — never became inconsistent (write was refused, no divergence to heal) |
| **E** — Winning value after conflict | **`999`** — the write from ap-node3 (step 23) has the **later timestamp** than `100` (step 22); after `iptables -F`, last-write-wins propagates `999` to all 5 nodes | N/A (CP cannot conflict) |

**Experiment E reasoning.** `score=100` is written first (step 22) and `score=999` second (step 23)
while node3 is isolated. `999` therefore carries the **higher `time.Now().UnixNano()` timestamp**.
When node3's traffic is restored (`iptables -F`) and sync runs, `merge()` keeps the entry with the
higher timestamp, so **`999` wins on every node**. (This also illustrates Q28: the winner is decided
purely by wall-clock timestamp.)

---

### AP vs CP Comparison

| Property | AP | CP |
|---|---|---|
| **Always accepts writes** | ✅ Yes — never refuses, even during a partition | ❌ No — only if a majority (quorum) acknowledges |
| **Read always returns latest value** | ❌ No — may return a stale value until sync converges | ✅ Yes — quorum read always returns the most recent committed value |
| **Handles conflicting writes** | ✅ Resolves them via last-write-wins (highest timestamp) — but can silently *lose* concurrent writes | ➖ By design **cannot** produce conflicts (majority quorum + overlap prevents divergence) |
| **Works when majority of nodes down** | ✅ Yes — remains available, keeps serving reads/writes | ❌ No — writes (and consistent reads) fail without a quorum |
| **Best use case** | High-availability, latency-sensitive, staleness-tolerant systems — **DNS, shopping carts, social-media likes** | Strong-consistency-critical systems — **bank balances, inventory, medical records** |

---

## Summary in one line

**Partitions are unavoidable (P is mandatory), so the real choice is what to sacrifice during one:**
AP keeps **answering** and accepts temporary staleness; CP keeps data **correct** and refuses to
answer when it cannot guarantee consistency.

---
Q27
  The setup

  You have 5 nodes: node1, node2, node3, node4, node5.

  Quorum = the minimum number of nodes that must agree before a write or read "counts."

  Formula: len(peers)/2 + 1 = 4/2 + 1 = 3. So you need 3 out of 5 nodes to agree.

  ---
  Why does it have to be 3? The "shared node" trick

  Here's the key idea in one sentence:

  ▎ If every write touches 3 nodes, and every read touches 3 nodes, then a read and a write can never completely miss each other — they always share at least one node.

  Why? Because 3 + 3 = 6, and there are only 5 nodes. You can't fit two groups of 3 into 5 nodes without overlapping. There's always at least one node in both groups.

  Example with values

  Write balance = 100. It goes to 3 nodes:

  WRITE touches:  [node1] [node2] [node3]   ← these now have balance=100
                   node4   node5            ← these are still old

  Read balance later. It asks 3 nodes:

  READ asks:       node3  [node4] [node5]

  Notice node3 is in both groups. So the read sees node3, which has the new value 100. The read picks the value with the newest timestamp → returns 100. ✅ Correct.

  No matter which 3 nodes the read picks, at least one of them will always be a node that got the write. That's the guarantee. That's why the system stays consistent.

  ---
  What goes wrong if you drop quorum to 2

  Now imagine you set quorum = 2. Two separate things break.

  Problem 1: Two groups can both "win" at the same time (split-brain)

  Suppose the network splits the cluster into two halves that can't talk to each other:

  GROUP LEFT:   node1  node2          GROUP RIGHT:   node3  node4  node5

  - A user writes balance = 100 on the left. Left has 2 nodes → "that's quorum (2), accepted!" ✅
  - Another user writes balance = 999 on the right. Right has 2 nodes → "that's quorum (2), accepted!" ✅

  Now you have two different committed values for the same key at the same time:

  node1, node2  →  balance = 100
  node3, node4  →  balance = 999

  Which one is correct? Nobody knows. The cluster disagrees with itself. This is "split-brain," and it's exactly the disaster CP is supposed to prevent.

  With quorum = 3 this is impossible: the left group only has 2 nodes, so it can't reach 3 → it refuses the write (quorum not reached: got 2 need 3). Only a group with a real
  majority can write, and there can only ever be one majority group. No split-brain.

  Problem 2: A read can completely miss the latest write

  With quorum = 2, 2 + 2 = 4, which is less than 5. So a write-group of 2 and a read-group of 2 might not overlap at all:

  WRITE touches:  [node1] [node2]              ← balance=100 written here
  READ asks:                    [node4] [node5] ← never got the write!

  The read asks node4 and node5 — neither has the new value — so it returns the old/stale balance. The reader sees wrong data even though no node is down.

  With quorum = 3, the read group of 3 is forced to overlap the write group of 3 (because 3+3 > 5), so this can't happen.

  ---
  The one rule to remember

  WRITE quorum + READ quorum  >  total nodes
          3      +     3       >        5      ✅  safe
          2      +     2       >        5      ❌  4 is not > 5, broken

  As long as write quorum + read quorum is greater than the total number of nodes, every read is guaranteed to overlap every write — so you always read the latest value, and
  you can never have two conflicting "winners." 3 is the smallest number that satisfies this for 5 nodes. That's why it's 3, not 2.

  Want me to replace the Q27 section in your LAB03_ANSWERS.md with this simpler version (or add it alongside the formal one)?