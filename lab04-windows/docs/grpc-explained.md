# gRPC & Protocol Buffers — Explained Simply

> A reference for Lab 04. Built around the lab's own files
> (`proto/kvstore.proto`, `grpc_server.go`, `grpc_client.go`, `main.go`).

## The one-sentence idea

> **gRPC lets you call a function on another computer as if it were a
> local function — and Protocol Buffers is the format used to pack the
> arguments and results into bytes for the network trip.**

Two separate things working together:

- **Protocol Buffers (protobuf)** = the *data format* (how to turn a struct
  into bytes and back).
- **gRPC** = the *transport + calling mechanism* (how the bytes travel and
  how the call is dispatched).

---

## Part 1: The `.proto` file is the single source of truth

Everything starts from **one file** you write: `kvstore.proto`. This is the
*contract* — it describes the service and its messages in a language-neutral way.

```
┌─────────────────────────────────────────────┐
│  kvstore.proto   (the CONTRACT)              │
│                                              │
│  service KeyValueStore {                     │
│    rpc Put(PutRequest) returns (PutResponse);│  ← the functions
│    rpc Get(GetRequest) returns (GetResponse);│
│    ...                                        │
│  }                                            │
│                                              │
│  message PutRequest {                        │  ← the data shapes
│    string key   = 1;                          │
│    string value = 2;                          │
│  }                                            │
│  message PutResponse { bool success = 1; }   │
└─────────────────────────────────────────────┘
```

Two kinds of things in it (see the real file, `proto/kvstore.proto`):

- **`service`** = the list of remote functions (lines 9-14).
- **`message`** = the structs/data shapes those functions use (lines 16-23).

### What are the numbers (`= 1`, `= 2`)?

Those are **field tags**, not values. Protobuf doesn't send field *names* over
the wire (unlike JSON, which sends `"key":"city"`). It sends the **tag number**
instead. So `key` travels as field `1`, `value` as field `2`. This is why
protobuf is smaller and faster than JSON — and why you must never renumber
existing fields.

```
JSON sends:      {"key":"city","value":"London"}   ← names + values (verbose)
Protobuf sends:  [1→"city"][2→"London"]            ← tags + values (compact, binary)
```

---

## Part 2: `protoc` generates code FROM the proto

You don't hand-write the networking. A compiler called **`protoc`** reads
`kvstore.proto` and **generates Go code** for you. This step happens inside the
Dockerfile (lines 16-20).

```
                          ┌──────────────┐
   kvstore.proto  ───────►│    protoc    │───────► generated Go files
   (you write this)       │  (compiler)  │
                          └──────────────┘
                                 │
            ┌────────────────────┴─────────────────────┐
            ▼                                           ▼
  kvstore.pb.go                            kvstore_grpc.pb.go
  ─────────────                            ──────────────────
  • struct PutRequest  { Key, Value }      • interface KeyValueStoreServer
  • struct PutResponse { Success }           (what YOU implement — Task 6)
  • struct GetRequest  ...                  • KeyValueStoreClient stub
  • (marshal/unmarshal to bytes)             (what the CLIENT calls — Task 7)
                                           • RegisterKeyValueStoreServer()
                                             (used in main — Task 8)
```

Every gRPC name used in the lab came from here:

- `pb.PutRequest`, `pb.GetResponse` → from `kvstore.pb.go`
- `pb.KeyValueStoreServer` (the interface implemented in Task 6) → from `kvstore_grpc.pb.go`
- `pb.NewKeyValueStoreClient` (Task 7) and `pb.RegisterKeyValueStoreServer` (Task 8) → also generated

**Key insight:** protoc capitalizes names for Go. The proto says
`string key = 1;` → the generated Go field is `req.Key`. That's why you wrote
`req.Key`, not `req.key`.

---

## Part 3: How a single call travels (the full picture)

What happens when the client runs
`client.Put(ctx, &pb.PutRequest{Key:"city", Value:"London"})`:

```
   CLIENT process                                    SERVER process
   (Task 7)                                          (Task 6)
┌───────────────────┐                            ┌────────────────────┐
│ client.Put(ctx,   │                            │ func (s *GRPCServer)│
│   &PutRequest{...})│                            │   Put(ctx, req) {   │
│        │          │                            │     s.store.Put(...)│
│        ▼          │                            │     return resp     │
│  generated STUB   │                            │   }                 │
│  (marshals struct │                            │        ▲            │
│   to protobuf     │                            │        │            │
│   bytes)          │                            │  generated code     │
│        │          │                            │  (unmarshals bytes  │
│        ▼          │                            │   back to a struct) │
│  [01 63 69 74 79  │   ===== HTTP/2 =====►      │        ▲            │
│   12 06 4C ...]   │   binary bytes over TCP    │        │            │
└────────┼──────────┘   (port 7001)             └────────┼────────────┘
         └──────────────────────────────────────────────┘
              the SAME bytes travel across the network
```

Step by step:

1. You call the **stub** method `client.Put(...)` — looks like a normal function call.
2. The generated stub **marshals** your `PutRequest` struct into compact protobuf **bytes**.
3. Bytes travel over **HTTP/2** (gRPC's transport) to the server on port 7001.
4. The server's generated code **unmarshals** the bytes back into a `PutRequest` struct.
5. gRPC calls **your** `Put` method (Task 6) with that struct — you run business logic (`s.store.Put`).
6. You return a `PutResponse`; the same marshal → send → unmarshal happens in reverse.

**You only wrote steps 1 and 5.** Everything else (marshal, transport,
unmarshal, dispatch) is generated/handled by gRPC. That's the whole value
proposition.

---

## Part 4: Why this beats net/rpc (Tasks 1-4)

You implemented *both* protocols, so this contrast is the real lesson:

|                | net/rpc (Tasks 1-4)                                            | gRPC (Tasks 6-8)                                  |
| -------------- | -------------------------------------------------------------- | ------------------------------------------------- |
| Define types   | **Twice** by hand — `server/types.go` AND client (duplicated all 8 structs) | **Once** in `.proto`, generated for both          |
| Call a method  | `client.Call("KVHandler.Put", ...)` — a **string**, typo = runtime crash | `client.Put(...)` — real method, typo = **compile error** |
| Wire format    | Go-only `gob`                                                  | protobuf — **any language** (Python, Java, Go interoperate) |
| Transport      | plain TCP                                                      | HTTP/2 (multiplexing, streaming)                  |
| Schema         | lives in your head / comments                                  | lives in `.proto`, enforced                       |

The big wins: **type safety** (the compiler checks your calls) and **language
independence** (a Python client could call your Go server because they share the
`.proto`, not Go-specific `gob`).

---

## The mental model to keep

```
   .proto  ──protoc──►  generated code  ──►  you fill in 2 small parts:
  (contract)            (the plumbing)        • server method bodies (Task 6)
                                              • client just calls the stub (Task 7)
```

Think of the `.proto` as a **Java interface** that both sides agree on — except
protoc *generates the implementation of all the boring network/serialization
parts* on both client and server, so neither side hand-writes them.

---

## Glossary

| Term            | Meaning                                                                 |
| --------------- | ----------------------------------------------------------------------- |
| `.proto`        | The contract file: defines services (functions) and messages (structs). |
| protobuf        | Binary serialization format. Sends tag numbers, not field names.        |
| `protoc`        | The compiler that turns `.proto` into generated Go code.                 |
| stub            | Client-side generated object whose methods look local but call remotely. |
| marshal         | Turn a struct into bytes.                                                |
| unmarshal       | Turn bytes back into a struct.                                           |
| HTTP/2          | The transport gRPC runs over (supports streaming + multiplexing).       |
| `context.Context` | Carries deadlines/cancellation for each call (the 5s timeout in Task 7). |
