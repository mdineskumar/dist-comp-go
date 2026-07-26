FROM golang:1.21-alpine AS builder
WORKDIR /build/raft
COPY student/raft /build/raft
RUN go mod init lab07 && go mod tidy && go build -o raft_bin .

FROM golang:1.21-alpine
RUN apk add --no-cache bash
WORKDIR /lab07/raft
COPY --from=builder /build/raft/raft_bin ./raft_bin
COPY student/raft/*.go ./
COPY docker/scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /lab07/raft/raft_bin

ENV NODE_ID=raft-node1
ENV PORT=9400
ENV PEERS=raft-node2:9401,raft-node3:9402,raft-node4:9403,raft-node5:9404
ENTRYPOINT ["/entrypoint.sh"]
