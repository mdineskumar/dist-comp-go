#!/bin/bash
echo "============================================"
echo " Lab 03 — CAP Theorem"
echo " Role   : ${ROLE}"
echo " Port   : ${NODE_PORT}"
echo " Peers  : ${PEERS}"
echo "============================================"

if [ "$ROLE" = "ap" ]; then
    cd /lab03/ap
    ./ap_bin -mode server -port "${NODE_PORT}" -peers "${PEERS}" &
    echo "[ENTRYPOINT] AP node started. Use CLI: ./ap_bin -mode cli -port ${NODE_PORT} <cmd>"
elif [ "$ROLE" = "cp" ]; then
    cd /lab03/cp
    ./cp_bin -mode server -port "${NODE_PORT}" -peers "${PEERS}" &
    echo "[ENTRYPOINT] CP node started. Use CLI: ./cp_bin -mode cli -port ${NODE_PORT} <cmd>"
else
    echo "[ENTRYPOINT] ERROR: Unknown ROLE '${ROLE}'. Use 'ap' or 'cp'."
    exit 1
fi

tail -f /dev/null
