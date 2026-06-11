#!/bin/bash
# ============================================================
# Lab 03 — CAP Theorem — Experiment Runner (Part 3, A–E)
#
# Usage (from lab03-complete folder, in Git Bash):
#   ./run_experiments.sh all      # run every experiment A–E
#   ./run_experiments.sh a        # run a single experiment (a|b|c|d|e)
#   ./run_experiments.sh reset    # restart all containers (clean state)
#
# Output is printed AND appended to results.log so you can copy
# answers into the Moodle quiz afterwards.
# ============================================================

# Stop Git Bash from rewriting /lab03/... into a Windows path.
export MSYS_NO_PATHCONV=1

LOG="results.log"

# Auto-detect compose command: "docker-compose" (older) or "docker compose" (newer).
if docker compose version >/dev/null 2>&1; then
  DC="docker compose"
else
  DC="docker-compose"
fi

# Helpers — run the CLI inside a container (internal port is always 7000/8000)
ap() { docker exec "lab03-ap-$1" /lab03/ap/ap_bin -mode cli -port 7000 "${@:2}" 2>&1; }
cp() { docker exec "lab03-cp-$1" /lab03/cp/cp_bin -mode cli -port 8000 "${@:2}" 2>&1; }

say() { echo -e "\n>>> $*" | tee -a "$LOG"; }
run() { echo "\$ $*" | tee -a "$LOG"; eval "$*" | tee -a "$LOG"; }

reset_all() {
  say "RESET — clean state (down then up)"
  $DC -f docker/docker-compose.yml down
  $DC -f docker/docker-compose.yml up -d
  say "Waiting 5s for nodes to start..."
  sleep 5
  run "docker ps --format '{{.Names}}\t{{.Status}}' | grep lab03"
}

exp_a() {
  say "EXPERIMENT A — Baseline (normal operation)"
  say "AP: put 5 keys on node1"
  ap node1 put k1 a1; ap node1 put k2 b2; ap node1 put k3 c3
  ap node1 put k4 d4; ap node1 put k5 e5
  say "Wait 3s for AP sync, then get all 5 from node5"
  sleep 3
  for k in k1 k2 k3 k4 k5; do ap node5 get $k; done
  say "CP: put 5 keys on node1, get all 5 from node5 immediately"
  cp node1 put k1 a1; cp node1 put k2 b2; cp node1 put k3 c3
  cp node1 put k4 d4; cp node1 put k5 e5
  for k in k1 k2 k3 k4 k5; do cp node5 get $k; done
  say "RECORD: AP all 5 ok after 3s? CP all 5 ok immediately?"
}

exp_b() {
  say "EXPERIMENT B — Stale read (core AP vs CP difference)"
  say "AP: put version=1 on node1"
  ap node1 put version 1
  say "AP: IMMEDIATE get version from node5 (expect stale / not found)"
  ap node5 get version
  say "AP: wait 3s, get version from node5 again (expect 1)"
  sleep 3
  ap node5 get version
  say "CP: put version=1 on node1, IMMEDIATE get from node5 (expect 1)"
  cp node1 put version 1
  cp node5 get version
  say "RECORD: AP immediate vs after-3s; CP immediate. Which gave a stale read?"
}

exp_c() {
  say "EXPERIMENT C — Minority partition (2 nodes down)"
  say "AP: stop node4 + node5"
  docker stop lab03-ap-node4 lab03-ap-node5
  say "AP: put status=active on node1 (3 nodes up)"
  ap node1 put status active
  say "AP: get status from node1, node2, node3"
  ap node1 get status; ap node2 get status; ap node3 get status
  say "AP: restart node4 + node5, wait 5s, get from all 5"
  docker start lab03-ap-node4 lab03-ap-node5; sleep 5
  for n in node1 node2 node3 node4 node5; do ap $n get status; done

  say "CP: stop node4 + node5"
  docker stop lab03-cp-node4 lab03-cp-node5
  sleep 2
  say "CP: put status=active on node1 (3 nodes up = quorum)"
  cp node1 put status active
  say "CP: get status from node1, node2, node3"
  cp node1 get status; cp node2 get status; cp node3 get status
  say "CP: restart node4 + node5"
  docker start lab03-cp-node4 lab03-cp-node5; sleep 5
  say "RECORD: did AP accept with 3 up? did CP accept with 3 up (=quorum)? AP converge after restart?"
}

exp_d() {
  say "EXPERIMENT D — Majority partition (3 nodes down)"
  say "AP: stop node3 + node4 + node5 (only 2 up)"
  docker stop lab03-ap-node3 lab03-ap-node4 lab03-ap-node5
  say "AP: put alert=critical on node1"
  ap node1 put alert critical
  say "AP: get alert from node1, node2"
  ap node1 get alert; ap node2 get alert
  say "AP: restart the 3 nodes, wait 5s, get from all 5"
  docker start lab03-ap-node3 lab03-ap-node4 lab03-ap-node5; sleep 5
  for n in node1 node2 node3 node4 node5; do ap $n get alert; done

  say "CP: stop node3 + node4 + node5 (only 2 up = below quorum)"
  docker stop lab03-cp-node3 lab03-cp-node4 lab03-cp-node5
  sleep 2
  say "CP: put alert=critical on node1 (EXPECT FAILURE — record exact error)"
  cp node1 put alert critical
  say "CP: restart the 3 nodes"
  docker start lab03-cp-node3 lab03-cp-node4 lab03-cp-node5; sleep 5
  say "RECORD: AP wrote with only 2 up? CP exact error message? AP recovery consistent?"
}

exp_e() {
  say "EXPERIMENT E — Conflicting writes (AP only)"
  say "Block node3 outgoing traffic (isolate it)"
  docker exec lab03-ap-node3 iptables -A OUTPUT -j DROP
  say "Write score=100 on node1, score=999 on isolated node3"
  ap node1 put score 100
  ap node3 put score 999
  say "Restore node3 traffic, wait 5s, get score from all 5 nodes"
  docker exec lab03-ap-node3 iptables -F
  sleep 5
  for n in node1 node2 node3 node4 node5; do ap $n get score; done
  say "RECORD: which value won (100 or 999) and why? (last-write-wins on timestamp)"
}

echo "=== Lab 03 experiment run ===" | tee "$LOG"
case "${1:-all}" in
  reset) reset_all ;;
  a) exp_a ;;
  b) exp_b ;;
  c) exp_c ;;
  d) exp_d ;;
  e) exp_e ;;
  all)
    reset_all; exp_a; exp_b; exp_c; exp_d; exp_e
    say "ALL DONE — see results.log"
    ;;
  *) echo "Usage: $0 [all|reset|a|b|c|d|e]" ;;
esac
