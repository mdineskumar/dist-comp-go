# Lab 03 — CAP Theorem — Experiment Guide (Part 3)

Step-by-step guide for running Experiments A–E and answering the discussion
questions. The implementation (Tasks 1–15) is already done and all 10 containers
are running, so this covers **only the experiments**.

---

## Manual vs. script?

You **can** do every experiment manually (they are just `docker exec` /
`docker stop` / `docker start` commands), but a helper script is better here:

1. **Git Bash mangles paths.** On Windows, Git Bash rewrites `/lab03/ap/ap_bin`
   into `C:/Program Files/Git/...` and the command fails. The script sets
   `export MSYS_NO_PATHCONV=1` once to stop this. (Run it manually before any
   `docker exec` if you are typing commands by hand.)
2. **Reproducibility + recording.** The script labels every step and writes all
   output to `results.log` so you can fill the Moodle tables accurately.

You still read and record the output yourself — the script just runs the steps
reliably.

---

## Step 0 — Clean state (always first)

From the `lab03-complete` folder in **Git Bash**:

```bash
chmod +x run_experiments.sh        # one time only
./run_experiments.sh reset         # down + up + wait 5s, shows all 10 running
```

---

## Step 1 — Run the experiments

Run them **one at a time** so you can read and record each result:

```bash
./run_experiments.sh a    # Baseline
./run_experiments.sh b    # Stale read   <- the key AP vs CP difference
./run_experiments.sh c    # Minority partition (2 down)
./run_experiments.sh d    # Majority partition (3 down)  <- CP write fails here
./run_experiments.sh e    # Conflicting writes (AP only)
```

Or `./run_experiments.sh all` to chain them with an automatic reset.

---

## Step 2 — Read the output

Everything prints to screen **and** appends to `results.log`. Each experiment
ends with a `RECORD:` line telling you exactly what to note for the Moodle quiz.

---

## Step 3 — Fill the Moodle tables

Map your `results.log` output to the lab's two tables (Results + AP-vs-CP
comparison) and the 4 discussion questions.

---

## What you should expect to see (so you know it works)

| Experiment | AP | CP |
|---|---|---|
| **A** Baseline | all 5 keys readable after 3s | all 5 readable immediately |
| **B** Stale read | node5 immediate = stale / `not found`; after 3s = correct | node5 immediate = correct (consistent) |
| **C** 2 down | write **succeeds** | write **succeeds** (3 = quorum) |
| **D** 3 down | write **succeeds** (only 2 up!) | write **FAILS** — `quorum not reached: got 2 need 3` |
| **E** conflict | node with **higher timestamp wins** (last-write-wins) | N/A |

---

## Notes for the discussion questions

- **Q26 (convergence time):** controlled by `syncInterval: 1 * time.Second` at
  `student/ap/node.go:61`. Converges in ~1–3s; set to 10s → stale reads last up
  to 10s.
- **Q27 (quorum = 3, not 2):** quorum=3 prevents two minority partitions (2+2)
  from both accepting writes → split-brain. Quorum=2 would allow that and break
  consistency.
- **Q28 (`UnixNano()` as version):** clock skew between nodes makes it
  unreliable — a node with a fast clock can overwrite a logically-newer write.
  Vector clocks track causality instead of wall-clock time.
- **Q29 (hospital vs Twitter likes):** hospital records = **CP** (Exp D: refuses
  the write rather than risk inconsistency); Twitter likes = **AP** (Exp B: a
  stale count is fine, availability matters more).

---

## Manual command reference (if not using the script)

Run `export MSYS_NO_PATHCONV=1` first in your Git Bash session.

```bash
# AP CLI (internal port is always 7000)
docker exec lab03-ap-node1 /lab03/ap/ap_bin -mode cli -port 7000 put city London
docker exec lab03-ap-node5 /lab03/ap/ap_bin -mode cli -port 7000 get city

# CP CLI (internal port is always 8000)
docker exec lab03-cp-node1 /lab03/cp/cp_bin -mode cli -port 8000 put city London
docker exec lab03-cp-node5 /lab03/cp/cp_bin -mode cli -port 8000 get city

# Partition control
docker stop  lab03-ap-node4 lab03-ap-node5
docker start lab03-ap-node4 lab03-ap-node5

# Experiment E — isolate a node with iptables (run INSIDE the container)
docker exec lab03-ap-node3 iptables -A OUTPUT -j DROP   # block
docker exec lab03-ap-node3 iptables -F                  # restore
```

---

## Caveat — Experiment E

Experiment E needs `iptables` inside the AP container. If you get
`iptables: command not found`, use a fallback to simulate the partition:

```bash
docker network disconnect docker_lab03 lab03-ap-node3   # isolate
docker network connect    docker_lab03 lab03-ap-node3   # restore
```

(The exact network name comes from `docker network ls | grep lab03`.)
