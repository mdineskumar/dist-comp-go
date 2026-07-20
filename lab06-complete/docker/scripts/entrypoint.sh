#!/bin/bash
echo "============================================"
echo " Lab 06 — Caching and Replication"
echo " Role: ${ROLE}"
echo "============================================"

cd /lab06/cache

case "$ROLE" in
  origin)
    ./cache_bin -mode origin -port "${PORT:-9300}" &
    ;;
  cache)
    sleep 2  # wait for origin
    ./cache_bin -mode cache -port "${PORT:-9200}" -origin "${ORIGIN}" -capacity "${CAPACITY}" &
    ;;
  idle)
    echo "[ENTRYPOINT] Client container ready"
    ;;
esac

tail -f /dev/null
