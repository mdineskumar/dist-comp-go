FROM golang:1.21-alpine AS builder

WORKDIR /build/cache
COPY student/cache /build/cache
RUN go mod init lab06 && go mod tidy && go build -o cache_bin .

FROM golang:1.21-alpine
RUN apk add --no-cache bash

WORKDIR /lab06/cache
COPY --from=builder /build/cache/cache_bin ./cache_bin
COPY student/cache/*.go ./

COPY docker/scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /lab06/cache/cache_bin

ENV ROLE=cache
ENV PORT=9200
ENV ORIGIN=origin:9300
ENV CAPACITY=10

ENTRYPOINT ["/entrypoint.sh"]
