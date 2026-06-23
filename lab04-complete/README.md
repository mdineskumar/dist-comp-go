# Lab 04 — RPC and Web Services

MSc Distributed Computing — Lab 04 (Take-Home Assignment)

## Overview
Implement the same key-value store using three different communication protocols:
- net/rpc  (port 7000) — Go standard library RPC
- gRPC     (port 7001) — Google's language-independent RPC
- REST     (port 8080) — HTTP API

## Quick Start
    docker build -t lab04 -f docker/Dockerfile .
    docker-compose -f docker/docker-compose.yml up -d

    # Enter server container
    docker exec -it lab04-server bash
    cd /lab04/server

    # Enter client container  
    docker exec -it lab04-client bash
    cd /lab04/client
    ./client_bin rpc put city London
    ./client_bin grpc get city
    ./client_bin rest list
    ./client_bin bench 100

## Tasks
    server/types.go      Task 1:  net/rpc request/reply structs
    server/handler.go    Task 2:  net/rpc handler methods
    server/main.go       Task 3:  start net/rpc server
    client/rpc_client.go Task 4:  net/rpc client
    server/grpc_server.go Task 6: gRPC server methods
    client/grpc_client.go Task 7: gRPC client
    server/main.go       Task 8:  start gRPC server
    server/rest_server.go Tasks 9-12: REST handlers
    client/rest_client.go Task 13: REST client

## Test with curl (REST only)
    curl -X PUT http://localhost:8080/keys/city \
         -H "Content-Type: application/json" \
         -d '{"value":"London"}'
    curl http://localhost:8080/keys/city
    curl http://localhost:8080/keys
    curl -X DELETE http://localhost:8080/keys/city
