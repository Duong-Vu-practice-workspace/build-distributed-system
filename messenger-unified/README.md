# Messenger Track — Unified Maelstrom Project

This folder unifies all exercises from the **Messenger** track into one Go
project. The duplicated node logic from the original single-file exercises is
extracted into a shared `maelstrom` package, and each exercise becomes a small
standalone command under `cmd/`.

## Project Layout

```
messenger-unified/
├── README.md                      # This file
├── maelstrom/                     # Shared library (one package)
│   ├── message.go                 # Message envelope and validation
│   ├── node.go                    # Node lifecycle, Send/Reply/Run
│   └── rpc.go                     # Synchronous RPC with timeout
└── cmd/                           # One command per exercise
    ├── json-parser/
    │   └── main.go                # Exercise 1: parse Maelstrom JSON
    ├── handler-init/
    │   └── main.go                # Exercise 2: handle init
    ├── echo-service/
    │   └── main.go                # Exercise 3: echo workload
    ├── envelope-validation/
    │   └── main.go                # Exercise 4: validate envelopes
    ├── async-handler/
    │   └── main.go                # Exercise 5: concurrent handler
    └── sync-rpc/
        └── main.go                # Exercise 6: sync RPC + proxy
```

## Build Everything

From the repository root:

```bash
go build ./messenger-unified/...
```

To produce named binaries next to each `main.go`:

```bash
go build -o messenger-unified/cmd/echo-service/echo-service ./messenger-unified/cmd/echo-service
go build -o messenger-unified/cmd/handler-init/handler-init ./messenger-unified/cmd/handler-init
# ... etc
```

Or install all binaries into `$GOPATH/bin` (useful for Maelstrom):

```bash
go install ./messenger-unified/cmd/...
```

## Quick Manual Test

You can drive any binary by piping JSON messages into stdin. Each line must be
a complete Maelstrom message.

```bash
cd messenger-unified/cmd/echo-service
go build -o echo-service .

printf '%s\n%s\n' \
  '{"src":"c0","dest":"n1","body":{"type":"init","msg_id":1,"node_id":"n1","node_ids":["n1"]}}' \
  '{"src":"c1","dest":"n1","body":{"type":"echo","msg_id":2,"echo":"hello"}}' \
  | ./echo-service
```

Expected output:

```json
{"src":"n1","dest":"c0","body":{"in_reply_to":1,"msg_id":0,"type":"init_ok"}}
{"src":"n1","dest":"c1","body":{"echo":"hello","in_reply_to":2,"msg_id":1,"type":"echo_ok"}}
```

## Install Maelstrom

If you do not have Maelstrom yet, download it from the official repo:

```bash
cd ~/bin  # or any folder on your $PATH
curl -L -o maelstrom https://github.com/jepsen-io/maelstrom/releases/latest/download/maelstrom
chmod +x maelstrom
```

Maelstrom requires Java (JRE 11+). Check with:

```bash
java -version
maelstrom test --help
```

## Run the Echo Workload with Maelstrom

The `echo-service` exercise implements the standard Maelstrom `echo` workload.
Build the binary and point Maelstrom at it:

```bash
# Option A: build a binary right next to main.go
go build -o messenger-unified/cmd/echo-service/echo-service ./messenger-unified/cmd/echo-service

./maelstrom test \
  -w echo \
  --bin ./messenger-unified/cmd/echo-service/echo-service \
  --node-count 1 \
  --time-limit 10 \
  --log-stderr
```

If you used `go install`, the binary lives in `$GOPATH/bin` (usually
`~/go/bin`):

```bash
go install ./messenger-unified/cmd/echo-service

./maelstrom test \
  -w echo \
  --bin ~/go/bin/echo-service \
  --node-count 1 \
  --time-limit 10 \
  --log-stderr
```

A successful run ends with something like:

```
Everything looks good! ヽ(‘ー`)ノ
```

## Running Other Exercises

| Exercise | Description | Runnable with Maelstrom workload? |
|----------|-------------|-----------------------------------|
| `json-parser` | Parses stdin lines and prints `PARSED: src\|dest\|type` | No — learning exercise |
| `handler-init` | Replies to `init` with `init_ok` | No — learning exercise |
| `echo-service` | Replies to `echo` with `echo_ok` | **Yes — `echo`** |
| `envelope-validation` | Echo service with message validation | No — learning exercise |
| `async-handler` | Echo service processed concurrently, sorted output | No — learning exercise |
| `sync-rpc` | Echo + proxy with synchronous RPC | No — learning exercise |

All six compile and can be exercised with manual stdin input. Only
`echo-service` maps to a built-in Maelstrom workload.

## Useful Maelstrom Flags

```bash
-w echo            # workload name
--bin <path>       # path to the node binary
--node-count N     # how many node processes to spawn
--time-limit N     # seconds to run the test
--rate N           # requests per second
--concurrency N    # number of clients (try 2n or 4n)
--log-stderr       # include node stderr in Maelstrom logs
```

For example, a slightly more aggressive echo test:

```bash
./maelstrom test \
  -w echo \
  --bin ./messenger-unified/cmd/echo-service/echo-service \
  --node-count 3 \
  --time-limit 20 \
  --rate 10 \
  --concurrency 2n \
  --log-stderr
```

## What's Different from the Original Exercises?

- **Shared library:** `Message`, `Node`, `Send`, `Reply`, `ValidateMessage`,
  and `SyncRPC` live in one place and are reused by every exercise.
- **One binary per exercise:** each `cmd/<exercise>` directory builds a
  runnable program, which is exactly what Maelstrom expects.
- **Clean imports:** exercise files are short and focus on the concept being
  practiced instead of repeating stdin/stdout boilerplate.

## Tips

- Maelstrom expects one JSON message per line on stdin and one JSON message per
  line on stdout. Do not add extra logging to stdout; log to stderr instead.
- The shared `Node.Reply` automatically fills `in_reply_to` from the request's
  `msg_id` and assigns the next `msg_id` before sending.
- `SyncRPC` is safe for concurrent use and cleans up pending requests on
  timeout.
