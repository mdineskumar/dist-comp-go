# Lab 03 — CAP Theorem

MSc in CS: Distributed Computing — Lab 03 (Take-Home Assignment)

## Overview
Students implement and compare two distributed key-value stores:
- AP system (Available + Partition-tolerant) — always accepts reads/writes
- CP system (Consistent + Partition-tolerant) — requires quorum for reads/writes

## Prerequisites
- Docker Desktop installed (Windows/Mac/Linux)
- Lab 02 completed (Chord DHT)

## Quick Start
    # Step 1 — build the image
    docker build -t lab03 -f docker/Dockerfile .

    # Step 2 — start all 10 nodes
    docker-compose -f docker/docker-compose.yml up -d

    # Step 3 — enter AP node 1
    docker exec -it lab03-ap-node1 bash
    cd /lab03/ap
    ./ap_bin -mode cli -port 7000 ping

    # Step 4 — enter CP node 1
    docker exec -it lab03-cp-node1 bash
    cd /lab03/cp
    ./cp_bin -mode cli -port 8000 ping

## Ports
  AP nodes: localhost:7000-7004
  CP nodes: localhost:8000-8004

## Tasks
  AP system: Tasks 1-8  (node.go, store.go, sync.go, rpc.go)
  CP system: Tasks 9-15 (node.go, store.go, quorum.go, rpc.go)

## Recompile after editing
    # AP system
    cd /lab03/ap
    go mod init lab03ap && go build -o ap_bin .

    # CP system
    cd /lab03/cp
    go mod init lab03cp && go build -o cp_bin .
