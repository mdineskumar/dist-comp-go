# Lab 06 — Caching and Replication

## Quick Start
    docker build -t lab06 -f docker/Dockerfile .
    docker-compose -f docker/docker-compose.yml up -d

    # Register replicas with primary
    docker exec lab06-client /lab06/cache/cache_bin -mode register \
        -primary cache-primary:9200 -replica cache-replica1:9200
    docker exec lab06-client /lab06/cache/cache_bin -mode register \
        -primary cache-primary:9200 -replica cache-replica2:9200

    # Write-through set
    docker exec lab06-client /lab06/cache/cache_bin -mode set \
        -cache cache-primary:9200 -key city -value London -strategy write-through

    # Get (from cache or origin)
    docker exec lab06-client /lab06/cache/cache_bin -mode get \
        -cache cache-primary:9200 -key city

    # Run benchmark
    docker exec lab06-client /lab06/cache/cache_bin -mode bench \
        -cache cache-primary:9200

    # Show stats
    docker exec lab06-client /lab06/cache/cache_bin -mode stats \
        -cache cache-primary:9200

## Tasks
    cache.go       Tasks 1-4:  NewCache, Get, Set, evictLRU
    strategy.go    Tasks 5-7:  WriteThrough, WriteBack, flushDirty
    replication.go Tasks 8-10: RegisterReplica, Replicate, ApplyReplica
    rpc.go         Task 11:    Get, Set, ApplyReplica handlers
    stats.go       Tasks 12-13:Hit/Miss/Eviction counters, HitRate, PrintStats
