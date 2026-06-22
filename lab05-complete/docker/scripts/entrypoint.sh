#!/bin/bash
echo "============================================"
echo " Lab 05 — Event Queues and Pub/Sub"
echo " Role: ${ROLE}"
echo "============================================"

case "$ROLE" in
  queue-broker)
    cd /lab05/queue
    ./queue_bin -mode broker -port 9000 &
    ;;
  pubsub-broker)
    cd /lab05/pubsub
    ./pubsub_bin -mode broker -port 9001 &
    ;;
  idle)
    echo "[ENTRYPOINT] Idle container — use docker exec to run commands"
    ;;
  *)
    echo "[ENTRYPOINT] Unknown ROLE '${ROLE}' — staying idle"
    ;;
esac

tail -f /dev/null
