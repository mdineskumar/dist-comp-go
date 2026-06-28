Test the Event Queue

```bash
docker exec lab05-worker /lab05/queue/queue_bin -mode work -queue orders -workers 3 &
 
# Terminal 2 — produce 10 tasks
docker exec lab05-producer /lab05/queue/queue_bin -mode produce -queue orders -payload order -count 10
 
# Watch the worker logs
docker logs -f lab05-worker
# Expected: each task processed by exactly one worker
# [WORKER worker-1] Processing task=task-1 payload="order-0"
# [WORKER worker-2] Processing task=task-2 payload="order-1"
# ... tasks distributed across all 3 workers


#fixed above 

# Start workers logging inside the container
docker exec -d lab05-worker sh -c \
    '/lab05/queue/queue_bin -mode work -broker queue-broker:9000 -queue orders -workers 3 > /tmp/worker.log 2>&1'

# Enqueue some tasks (single line, type quotes manually)
docker exec lab05-producer /lab05/queue/queue_bin -mode produce -broker queue-broker:9000 -queue orders -payload "order-123" -count 5

# Check worker logs
docker exec lab05-worker cat /tmp/worker.log

```


```bash
docker exec lab05-subscriber1 pkill pubsub_bin

# Terminal 1, 2, 3 — start subscribers, redirecting output to log files
docker exec -d lab05-subscriber1 sh -c \
    '/lab05/pubsub/pubsub_bin -mode subscribe -topic news -id sub1 -port 9100 -host subscriber1 > /tmp/sub1.log 2>&1'

docker exec lab05-subscriber2 pkill pubsub_bin
docker exec -d lab05-subscriber2 sh -c \
    '/lab05/pubsub/pubsub_bin -mode subscribe -topic news -id sub2 -port 9100 -host subscriber2 > /tmp/sub2.log 2>&1'

docker exec lab05-subscriber3 pkill pubsub_bin
docker exec -d lab05-subscriber3 sh -c \
    '/lab05/pubsub/pubsub_bin -mode subscribe -topic news -id sub3 -port 9100 -host subscriber3 > /tmp/sub3.log 2>&1'

# Wait for subscriptions to register
sleep 2

# Publish
docker exec lab05-publisher /lab05/pubsub/pubsub_bin \
    -mode publish -topic news -key headline -value "Breaking news!"

# Wait briefly for delivery
sleep 1

# Check logs from inside each container
docker exec lab05-subscriber1 cat /tmp/sub1.log
docker exec lab05-subscriber2 cat /tmp/sub2.log
docker exec lab05-subscriber3 cat /tmp/sub3.log


docker exec lab05-publisher /lab05/pubsub/pubsub_bin -mode publish -topic news -key headline -value "Breaking news!"
```


## Experiemnts

```bash
docker exec lab05-worker /lab05/queue/queue_bin -mode work -queue orders -workers 3 &

docker exec lab05-producer /lab05/queue/queue_bin -mode produce -queue orders -payload order -count 20

```