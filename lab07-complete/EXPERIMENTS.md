# Lab 07 — Raft Experiments Guide (Windows / docker exec)

> ⚠️ Prerequisite: all 14 tasks complete, **and** the critical bugs fixed
> (self-deadlock in `startElection()`, candidate step-down on same-term
> `AppendEntries`, `leaderID` tracking). Without the deadlock fix in
> particular, no node will ever successfully become leader and none of
> these experiments will behave as described below.

Container/port convention used throughout (adjust if your
`docker-compose.yml` differs):

| Node | Container name | Address |
|------|-----------------|---------|
| 1 | `lab07-node1` | `raft-node1:9400` |
| 2 | `lab07-node2` | `raft-node2:9401` |
| 3 | `lab07-node3` | `raft-node3:9402` |
| 4 | `lab07-node4` | `raft-node4:9403` |
| 5 | `lab07-node5` | `raft-node5:9404` |

Start the cluster:

```powershell
docker-compose up -d
```

---

## Experiment A — First Leader Election

**Commands:**

```powershell
docker exec lab07-node1 /lab07/raft/raft_bin -mode watch
```

Let it run for up to 5 seconds, then `Ctrl+C`.

> **Why `watch` alone won't show you the election happening:** `-mode watch`
> polls each node's `Status` RPC once per **second** (hardcoded in
> `main.go`), but a Raft election resolves in **150–300ms** — faster than
> one poll interval. Unless you attach `watch` at the exact instant the
> containers start, every sample will already show a stable leader. That's
> not a bug, just a limitation of polling for something this fast.
>
> **To get a real, precise election time, read the timestamped container
> logs of whichever node won instead:**
>
> ```powershell
> docker logs -t lab07-node2
> ```
>
> The `-t` flag prefixes each line with a nanosecond-precision timestamp.
> Find the `Starting as Follower`, `Became Candidate`, and `Became LEADER`
> lines and subtract their timestamps.
>
> **Worked example (real run):**
> ```
> 2026-07-27T09:13:57.752531494Z [NODE raft-node2] Starting as Follower (term 0)
> 2026-07-27T09:13:57.752541395Z [RPC] Raft node listening on port 9401
> 2026-07-27T09:13:57.752544595Z [RAFT] Node raft-node2 ready
> 2026-07-27T09:13:57.921173000Z [NODE raft-node2] Became Candidate (term 1) - starting election
> 2026-07-27T09:13:57.924204666Z [NODE raft-node2] Became LEADER (term 1)
> ```
>
> | Phase | Delta |
> |---|---|
> | Follower start → Became Candidate (election timeout firing) | **168.64 ms** |
> | Became Candidate → Became LEADER (vote-gathering) | **3.03 ms** |
> | **Total: boot → leader** | **171.67 ms** |
>
> The 168.64ms matches the expected `150 + rand.Intn(150)` election-timeout
> window exactly, confirming the timer logic is correct. The 3.03ms
> vote-gathering phase is near-instant because all 5 nodes are on the same
> local Docker bridge network with negligible RPC round-trip time.

**Observe and record:**
- Which node became the first leader?
- At what term?
- How long did the election take (rough estimate from watch output)?
- Why did only ONE node become leader if all started simultaneously?

**Sample expected answer:**
- Leader: any one of the 5 nodes — it's essentially random which one (determined by whose randomized election timer fires first). Don't expect a specific node; just confirm exactly one is marked `← LEADER`.
- Term: **1** (occasionally 2, if an early split vote happens and a second round is needed).
- Election time: `watch`'s 1-second polling is too coarse to catch the transition — every sample will simply show the already-elected leader. For a precise figure, use `docker logs -t <leader-container>` and diff the timestamps: in a real run this measured **~172ms total** (168.6ms for the election timeout to fire + 3.0ms to gather votes and become leader), consistent with the expected 150–300ms window.
- Why only one:
  - Each node's election timer is a random duration in `150–300ms`, so they don't all expire at exactly the same instant.
  - Whichever node's timer fires first increments its term and requests votes *before* the others time out.
  - Every other node can cast at most one vote per term (`votedFor` guard in `RequestVote`), so once a majority (3 of 5) has voted for the first candidate, no other candidate in that term can reach a majority — their votes are already spent.
  - If two candidates' timers happen to fire close together (a split vote), neither gets a majority, both time out again, and the randomness on the *next* attempt almost always desyncs them enough to produce a single winner.

---

## Experiment B — Log Replication — Consensus in Action

**Commands:**

```powershell
docker exec lab07-node1 /lab07/raft/raft_bin -mode client -node raft-node2:9401 -cmd 'SET a 1'
docker exec lab07-node1 /lab07/raft/raft_bin -mode client -node raft-node2:9401 -cmd 'SET b 2'
docker exec lab07-node1 /lab07/raft/raft_bin -mode client -node raft-node2:9401 -cmd 'SET c 3'
docker exec lab07-node1 /lab07/raft/raft_bin -mode client -node raft-node2:9401 -cmd 'SET d 4'
docker exec lab07-node1 /lab07/raft/raft_bin -mode client -node raft-node2:9401 -cmd 'SET e 5'

# Read from nodes other than wherever the leader turned out to be
docker exec lab07-node3 /lab07/raft/raft_bin -mode client -node raft-node3:9402 -get a
docker exec lab07-node5 /lab07/raft/raft_bin -mode client -node raft-node5:9404 -get e
docker exec lab07-node5 /lab07/raft/raft_bin -mode client -node raft-node5:9404 -get e

docker exec lab07-node2 /lab07/raft/raft_bin -mode watch
```

> Note: `Submit` only succeeds if the node you target is the actual
> leader (the skeleton's client does **not** auto-forward — it just
> reports `❌ Failed: not the leader (leader: <id>)`). Point `-node` at
> whichever address `-mode watch` shows as `← LEADER`.

**Observe and record:**
- Were all 5 values readable from non-leader nodes?
- Did all 5 nodes show the same log length and commit index?
- What does this demonstrate about Raft consensus?

**Sample expected answer:**
- Yes — all 5 key/value pairs (`a=1` … `e=5`) are readable from any node, not just the leader, since `GetValue` reads local state and every node's state machine converges to the same content once entries are committed and applied.
- Yes — `watch` should show `log=5 committed=5` on **all 5 nodes** within roughly one heartbeat interval (50ms) of the last write, since heartbeats carry `LeaderCommit` and followers replicate/apply it.
- This demonstrates the core Raft guarantee: once a majority acknowledges a log entry, it is durably committed and will eventually be applied identically on every node — the "same input log → same state machine output" principle that makes replicated logs equivalent to consensus.

---

## Experiment C — Leader Failure — Automatic Re-election

**Commands:**

```powershell
docker exec lab07-node1 /lab07/raft/raft_bin -mode watch    # note the current leader, e.g. lab07-node2

docker stop lab07-node2     # replace 4 with whichever node is actually leader

docker exec lab07-node1 /lab07/raft/raft_bin -mode watch    # watch remaining 4 nodes re-elect


# once a new leader appears, submit a new command to it:
# new leader is lab07-node3
docker exec lab07-node3 /lab07/raft/raft_bin -mode client -node raft-node3:9402 -cmd 'SET f 6'
```

**Observe and record:**
- How long did it take for a new leader to be elected?
- What term did the new election happen in?
- Were the 5 previously committed entries still readable from the new leader?
- How does Raft guarantee the new leader has all committed entries?

**Sample expected answer:**
- Roughly **150–300ms** after the leader stops sending heartbeats — that's exactly one election-timeout window on whichever follower notices first.
- New term = old leader's term **+ 1** (elections always start by incrementing `currentTerm`).
- Yes — `a` through `e` are still readable from the new leader (and from all 4 surviving nodes).
- Because `RequestVote` only grants a vote if the candidate's log is at least as up-to-date as the voter's own (`LastLogTerm`/`LastLogIndex` check in Task 6). Since a committed entry exists on a majority of nodes, any candidate that manages to win a majority of votes must have gotten at least one vote from a node holding that committed entry — and that voter would have refused to vote for a candidate with a *less* up-to-date log. So the winner is guaranteed to already have every committed entry.

---

## Experiment D — Minority Failure — Cluster Still Works

**Commands:**

```powershell
# with 5 running and a leader elected, stop 2 NON-leader nodes, e.g. node2 and node3
docker stop lab07-node2 lab07-node3

docker exec lab07-node5 /lab07/raft/raft_bin -mode client -node raft-node5:9404 -cmd 'SET g 7'
docker exec lab07-node5 /lab07/raft/raft_bin -mode client -node raft-node5:9404 -cmd 'SET h 8'
docker exec lab07-node5 /lab07/raft/raft_bin -mode client -node raft-node5:9404 -cmd 'SET i 9'

# read from the two remaining followers (not the leader, not the two you stopped)
docker exec lab07-node4 /lab07/raft/raft_bin -mode client -node raft-node4:9403 -get g
```

(Adjust node numbers so you're targeting the actual leader and actual surviving followers.)

**Observe and record:**
- Did writes succeed with only 3 nodes (majority of 5)?
- Were values readable from the remaining followers?
- How many nodes can fail before writes stop working?

**Sample expected answer:**
- Yes — 3 out of 5 is still a majority, so `commitEntries()` reaches `count > len(peers)/2` and commits normally.
- Yes — `g`, `h`, `i` are readable from any surviving node.
- With 5 nodes, up to **2** can fail and writes still work (3 remaining = majority). A 3rd failure (down to 2 of 5) drops below majority and writes stall — that's Experiment E.

---

## Experiment E — Majority Failure — Cluster Stalls

**Commands:**

```powershell
# stop a 3rd node (now 3 down, only 2 remain)
docker stop lab07-node1

docker exec lab07-node5 /lab07/raft/raft_bin -mode client -node raft-node5:9404 -cmd 'SET x 99'
# (this will hang for ~500ms then report a timeout — see explanation below)

docker start lab07-node3
timeout 2
docker exec lab07-node5 /lab07/raft/raft_bin -mode client -node raft-node5:9404 -cmd 'SET x 99'
```

**Observe and record:**
- What happened when you tried to write with only 2 nodes?
- After restarting one node (3 total), did writes work again?
- How does this connect to the CAP theorem from Lab 03?

**Sample expected answer:**
- The write fails after roughly **500ms** with `❌ Failed: commit timeout` (or `not the leader` if the remaining 2 nodes couldn't even elect one). This is `AppendEntry`'s polling loop in `log.go` giving up after `50 × 10ms = 500ms` because `commitIndex` never advances — 2 out of 5 is not a majority, so `commitEntries()` can never satisfy `count > len(peers)/2`. Note the entry *is* appended to the leader's own local log (uncommitted) even though the client sees a failure.
- Yes — once a 3rd node rejoins (3 of 5 reachable again = majority), the next submitted command commits successfully and `watch` shows `committed` advancing again.
- This is exactly the CAP tradeoff from Lab 03: with a majority partition, Raft chooses **Consistency + Partition tolerance over Availability** (it's a CP system) — it refuses to accept writes it can't safely commit, rather than risk two nodes independently accepting conflicting writes that could never be reconciled.

---

## Discussion Questions — Sample Answers

**24. Why does only one node become leader when all 5 start elections "simultaneously"?**

- Election timeouts are randomized per node (`150 + rand.Intn(150)` ms), so in practice they never fire at exactly the same instant.
- Whichever node's timer expires first becomes a candidate first, increments its term, and requests votes before the others have even timed out.
- Each node grants at most one vote per term (guarded by `votedFor`), so once that first candidate collects a majority (3 of 5, including itself), no other candidate can reach a majority in that same term — the votes are already committed.
- In the rare case that two candidates' timers fire close enough together to cause a genuine split vote (each getting some but not a majority), the election fails for that term, both revert to waiting on a new random timeout, and the wider randomness on the retry almost always resolves the tie on the next attempt.

**25. Which condition in `RequestVote` guarantees the new leader has all committed entries?**

- The log up-to-date check:
  ```go
  args.LastLogTerm > n.lastLogTerm() ||
  (args.LastLogTerm == n.lastLogTerm() && args.LastLogIndex >= n.lastLogIndex())
  ```
- A voter refuses to vote for any candidate whose log is *less* up-to-date than its own.
- Since a committed entry is, by definition, present on a majority of nodes, any candidate that manages to win an election must receive a vote from at least one node in that "has the committed entry" majority.
- That node would only vote for a candidate whose log is at least as current as its own — therefore the winning candidate's log must already contain every previously committed entry.

**26. CAP tradeoff observed in Experiment E**

- Raft sacrifices **Availability** to preserve **Consistency** (a CP system per the Lab 03 framing).
- When it can't reach a majority, it deliberately refuses new writes rather than risk two disconnected partitions each committing different, irreconcilable histories.
- You'd accept this tradeoff whenever correctness matters more than uptime — e.g., a leader-election/lock service, financial ledger, or config store — where serving a stale or conflicting write would be far worse than a temporary outage.

**27. Data-loss scenario from the in-memory log**

- Because `n.log` lives only in process memory (nothing is fsynced to disk), a crash-and-restart wipes it clean.
- Client submits `SET balance 100`. Leader appends it, replicates to a majority, commits it, and returns success to the client — the client now believes this write is durable.
- Shortly after, the **leader's container restarts** (crash, OOM kill, `docker restart`, host reboot — anything that re-runs `NewNode()`).
- `NewNode()` reinitializes `log: []LogEntry{}`, `currentTerm: 0`, `commitIndex: 0` from scratch — the entry is gone from that node with no way to recover it, because nothing was ever written to stable storage.
- If a **majority** of nodes restart around the same time (e.g., the whole docker-compose stack is brought down and back up), *every* node's log is empty simultaneously — the previously committed `balance=100` is permanently lost cluster-wide, even though the client was already told the write succeeded.
- Real Raft avoids this by fsyncing each log entry (and `currentTerm`/`votedFor`) to disk *before* replying to `AppendEntries`/casting a vote, and reloading that persisted state on startup — so a crash-restarted node rejoins with everything it had committed before crashing.

---

## Results Recording — what to expect when you fill in Moodle

### Experiments A and B

| Measurement | Expected |
|---|---|
| A — Which node became first leader? | any one node (random) — record whichever it actually was |
| A — At what term? | 1 (rarely 2) |
| A — Approximate election time | ~150–300ms |
| B — All 5 values readable from non-leader nodes? | Yes |
| B — All nodes showed same log=5 committed=5? | Yes |

### Experiments C, D, E — Fault Tolerance

| Measurement | Expected |
|---|---|
| C — Time for new leader after crash | ~150–300ms |
| C — Term of new election | previous term + 1 |
| C — Old committed entries still readable? | Yes |
| D — Writes succeeded with 3 nodes? | Yes |
| E — Write result with only 2 nodes | Fails — `commit timeout` after ~500ms |
| E — Writes resumed after restarting 1 node? | Yes |

---

## Troubleshooting

If experiment A never converges (term climbs into double/triple digits with no `← LEADER` ever appearing in `watch`), you still have the `startElection()`/`becomeLeader()` self-deadlock bug — fix that first; nothing else in this guide will work until a leader can actually be elected.
