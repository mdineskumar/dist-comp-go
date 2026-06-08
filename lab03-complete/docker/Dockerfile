FROM golang:1.21-alpine AS builder

# Build AP binary
WORKDIR /build/ap
COPY student/ap /build/ap
RUN go mod init lab03ap && go mod tidy && go build -o ap_bin .

# Build CP binary
WORKDIR /build/cp
COPY student/cp /build/cp
RUN go mod init lab03cp && go mod tidy && go build -o cp_bin .

FROM golang:1.21-alpine
RUN apk add --no-cache bash iptables iproute2

# AP system
WORKDIR /lab03/ap
COPY --from=builder /build/ap/ap_bin ./ap_bin
COPY student/ap/node.go   ./node.go
COPY student/ap/store.go  ./store.go
COPY student/ap/sync.go   ./sync.go
COPY student/ap/rpc.go    ./rpc.go
COPY student/ap/main.go   ./main.go

# CP system
WORKDIR /lab03/cp
COPY --from=builder /build/cp/cp_bin ./cp_bin
COPY student/cp/node.go    ./node.go
COPY student/cp/store.go   ./store.go
COPY student/cp/quorum.go  ./quorum.go
COPY student/cp/rpc.go     ./rpc.go
COPY student/cp/main.go    ./main.go

COPY docker/scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /lab03/ap/ap_bin /lab03/cp/cp_bin

ENV ROLE=ap
ENV NODE_PORT=7000
ENV PEERS=""

ENTRYPOINT ["/entrypoint.sh"]
