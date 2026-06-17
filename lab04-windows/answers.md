# Lab 04 — Experiments & Discussion Answers

> Based on actual benchmark run: net/rpc 57ms (286µs/op), gRPC 64ms (318µs/op),
> REST 158ms (792µs/op).

## Benchmark Results (100 put+get pairs)

| Protocol | Total time | Avg/op | Ops/sec |
| -------- | ---------- | ------ | ------- |
| net/rpc  | 57ms       | 286µs  | 3499    |
| gRPC     | 64ms       | 318µs  | 3145    |
| REST     | 158ms      | 792µs  | 1262    |

---

# Experiments

## Experiment A — Cross-Protocol Data Sharing
- **Read via gRPC a key written via net/rpc?** Yes — verified `rpc put country UK` then `grpc get country` returns `UK`.
- **Read via REST a key written via gRPC?** Yes — all 9 read/write combinations work.
- **What it tells you:** all 3 servers share **one `Store` instance** (created once in `main`, passed to each). The protocol is just the *door*; the data behind every door is the same map. Protocol = *how* you talk, not *what* is stored.

## Experiment B — Benchmark
- **Fastest:** net/rpc (57ms) ≈ gRPC (64ms) — effectively tied, ~32µs/op apart (measurement noise).
- **Slowest:** REST (158ms), ~2.5–3× slower.
- **Why REST is slower:** text JSON encoding (vs compact binary), per-request HTTP/1.1 overhead (headers, status line), and a new connection per call (RPC/gRPC reuse one persistent connection).
- **Does ranking change 100 → 500?** REST stays clearly slowest; net/rpc and gRPC stay near-tied (their order may flip — it's noise).

## Experiment C — REST from Browser and curl
- **Browser/curl call REST?** Yes — plain HTTP. curl returned `{"found":true,"value":"London"}`, `{"success":true}`, `{"deleted":true}`, etc.
- **Why not net/rpc or gRPC?** Browsers/curl speak HTTP. net/rpc uses **gob over raw TCP** (not HTTP); gRPC uses **HTTP/2 + protobuf** (browsers need a special library/gRPC-web proxy). Neither is reachable from a plain browser.
- **What it means for web apps:** a browser frontend needs a **REST (or gRPC-web) backend** — you can't point a browser straight at a net/rpc service.

## Experiment D — Large Values (~10,000 chars)
- **Fastest for large values:** binary protocols (net/rpc, gRPC) — no text escaping, less bloat.
- **Order change?** REST's disadvantage **grows** with size — JSON must escape/encode the big string as text, while gob/protobuf ship raw bytes. Binary's lead widens.

---

# Discussion Questions

## Q1 — Which implementation was easiest / hardest to write?
- **Easiest: gRPC.** The `.proto` generated the structs, server interface, and client stub — you only filled method bodies (`s.store.Put(req.Key, req.Value); return &pb.PutResponse{Success: true}`). Type-safe, minimal boilerplate.
- **Hardest: net/rpc.** Had to **hand-write all 8 request/reply structs twice** (server `types.go` *and* client), match field names exactly for gob, and call methods by **string name** (`"KVHandler.Put"`) with no compile-time check.
- **REST: middle.** No code generation, but manual HTTP routing, JSON encode/decode, and status codes (201/404).

## Q2 — Why can a browser call REST but not the others? Which for a web backend?
- **Why:** REST = standard HTTP/1.1 + JSON, which browsers speak natively. net/rpc = gob over raw TCP (not HTTP at all); gRPC = HTTP/2 + protobuf (needs a special client or gRPC-web proxy).
- **Choice for a web backend: REST** — universally accessible (browser, mobile, curl, any language), easy to debug, no client library needed. (Use gRPC for internal service-to-service calls behind the API.)

## Q3 — Why do lowercase fields silently become empty? (gob)
- net/rpc serializes structs with **gob**, which uses **reflection** and can only access **exported (capitalized) fields**.
- A lowercase field is **unexported** → gob can't see it → it's never put on the wire → arrives as the zero value (`""`, `false`) on the other side. **No error**, just silent data loss.
- **Lesson for distributed systems:** anything crossing a process/network boundary must be **explicitly serializable**; silent serialization gaps cause bugs that compile and run but corrupt data. Always verify what's actually on the wire.

## Q4 — Rebuilding Lab 02 Chord DHT for Python/Java clients — which protocol?
- **Choose gRPC.**
- **Why:** language-independent — the `.proto` contract generates Python and Java clients automatically, so non-Go clients can call your Go nodes (net/rpc is **Go-only** and couldn't talk to them at all).
- **Plus:** near-net/rpc speed (binary protobuf + HTTP/2), a type-safe contract enforced at compile time, and built-in features (deadlines/timeouts via `context`, streaming) suited to production. REST would also be cross-language but slower and weaker-typed for high-frequency internal DHT calls.

---

# Protocol Comparison Summary

| Property | net/rpc | gRPC | REST |
| -------- | ------- | ---- | ---- |
| Language support | Go only | Any language | Any language / browser |
| Encoding | gob (binary) | protobuf (binary) | JSON (text) |
| Transport | raw TCP | HTTP/2 | HTTP/1.1 |
| Connection | persistent | persistent | per-request |
| Fastest for high-throughput | yes (tied) | yes (tied) | no |
| Works from a browser | no | no | yes |
| Needs code generation | no | yes (protoc) | no |
| Best use case | internal Go services | cross-language microservices | public APIs / web backends |
