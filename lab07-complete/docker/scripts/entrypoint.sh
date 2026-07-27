#!/bin/bash
echo "============================================"
echo " Lab 07 — Raft Consensus Node"
echo " ID:    ${NODE_ID}"
echo " Port:  ${PORT}"
echo "============================================"
cd /lab07/raft
# Small random delay so nodes don't all start at exactly the same time
#sleep $(( RANDOM % 3 ))
./raft_bin -mode node -id "${NODE_ID}" -port "${PORT}" -peers "${PEERS}" &
echo "[ENTRYPOINT] Raft node started"
tail -f /dev/null
