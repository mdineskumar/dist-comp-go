FROM golang:1.21-alpine AS builder

WORKDIR /build/queue
COPY student/queue /build/queue
RUN go mod init lab05queue && go mod tidy && go build -o queue_bin .

WORKDIR /build/pubsub
COPY student/pubsub /build/pubsub
RUN go mod init lab05pubsub && go mod tidy && go build -o pubsub_bin .

FROM golang:1.21-alpine
RUN apk add --no-cache bash

WORKDIR /lab05/queue
COPY --from=builder /build/queue/queue_bin ./queue_bin
COPY student/queue/*.go ./

WORKDIR /lab05/pubsub
COPY --from=builder /build/pubsub/pubsub_bin ./pubsub_bin
COPY student/pubsub/*.go ./

COPY docker/scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /lab05/queue/queue_bin /lab05/pubsub/pubsub_bin

ENV ROLE=idle
ENTRYPOINT ["/entrypoint.sh"]
