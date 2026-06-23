FROM golang:1.21-alpine AS builder

# Install protoc and tools
RUN apk add --no-cache protobuf protobuf-dev git

# Install Go gRPC plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# ── Build server ─────────────────────────────────────────
WORKDIR /build/server
COPY student/server /build/server/
COPY student/server/proto /build/server/proto/

# Generate gRPC code from proto
RUN mkdir -p proto && \
    protoc --go_out=. --go-grpc_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    proto/kvstore.proto

# Initialize module and build
RUN go mod init lab04server && \
    go get google.golang.org/grpc@v1.59.0 && \
    go get google.golang.org/protobuf@v1.31.0 && \
    go mod tidy && \
    go build -o server_bin .

# ── Build client ─────────────────────────────────────────
WORKDIR /build/client
COPY student/client /build/client/
COPY student/server/proto /build/client/proto/

# Generate same proto code for client
RUN mkdir -p proto && \
    protoc --go_out=. --go-grpc_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    proto/kvstore.proto

RUN go mod init lab04client && \
    go get google.golang.org/grpc@v1.59.0 && \
    go get google.golang.org/protobuf@v1.31.0 && \
    go mod tidy && \
    go build -o client_bin .

# ── Final image ───────────────────────────────────────────
FROM golang:1.21-alpine
RUN apk add --no-cache bash curl protobuf

# Install gRPC tools for recompilation inside container
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /lab04

# Server files
COPY --from=builder /build/server/server_bin ./server/server_bin
COPY student/server/ ./server/
COPY --from=builder /build/server/proto/ ./server/proto/

# Client files
COPY --from=builder /build/client/client_bin ./client/client_bin
COPY student/client/ ./client/
COPY --from=builder /build/client/proto/ ./client/proto/

COPY docker/scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /lab04/server/server_bin /lab04/client/client_bin

ENV ROLE=server

ENTRYPOINT ["/entrypoint.sh"]
