# Lab 05 — Event Queues and Pub/Sub

MSc in CS Distributed Computing — Lab 05 (Take-Home Assignment)

## Overview
Two messaging systems built from scratch:
- Event Queue: point-to-point work distribution (like RabbitMQ/SQS)
- Pub/Sub: fan-out broadcast (like Kafka/Redis Pub/Sub)

## Quick Start
    docker build -t lab05 -f docker/Dockerfile .
    docker-compose -f docker/docker-compose.yml up -d

    # Event Queue
    docker exec lab05-producer /lab05/queue/queue_bin -mode produce -queue orders -payload order1 -count 10
    docker exec lab05-worker /lab05/queue/queue_bin -mode work -queue orders -workers 3

    # Pub/Sub
    docker exec -d lab05-subscriber1 /lab05/pubsub/pubsub_bin -mode subscribe -topic news -id sub1 -port 9100 -host subscriber1
    docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -topic news -key headline -value "Breaking news"

## Tasks
    queue/queue.go    Tasks 1-5: NewQueueManager, Enqueue, Dequeue, Ack, Nack
    queue/rpc.go      Task 6:    RPC handlers
    queue/worker.go   Task 7:    Worker pool
    pubsub/broker.go    Tasks 8-11,14: NewBroker, Subscribe, Unsubscribe, Publish
    pubsub/subscriber.go Task 12: Deliver
    pubsub/rpc.go       Task 13:  Broker RPC handlers
