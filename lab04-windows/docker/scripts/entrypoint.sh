#!/bin/bash
echo "============================================"
echo " Lab 04 — RPC and Web Services"
echo " Role: ${ROLE}"
echo "============================================"

if [ "$ROLE" = "server" ]; then
    cd /lab04/server
    ./server_bin &
    echo "[ENTRYPOINT] All 3 servers started"
    echo "[ENTRYPOINT] net/rpc → port 7000"
    echo "[ENTRYPOINT] gRPC    → port 7001"
    echo "[ENTRYPOINT] REST    → port 8080"
elif [ "$ROLE" = "client" ]; then
    echo "[ENTRYPOINT] Client ready"
    echo "[ENTRYPOINT] Use: ./client_bin rpc|grpc|rest|all|bench <cmd>"
    cd /lab04/client
fi

tail -f /dev/null
