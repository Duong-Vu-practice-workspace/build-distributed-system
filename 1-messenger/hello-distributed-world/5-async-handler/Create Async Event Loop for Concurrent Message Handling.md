---
title: "Create Async Event Loop for Concurrent Message Handling"
source: "https://builddistributedsystem.com/tracks/messenger/tasks/task-1-5-async-handler"
author:
  - "[[Mohit Mishra]]"
published:
created: 2026-06-14
description: "Real distributed systems handle many messages concurrently. Your current synchronous implementation processes one message at a time, which limits throughput.  R"
tags:
  - "clippings"
---
## Implementation

Network partitions are not rare. AWS us-east-1 has had at least 12 documented network partition events in the past five years. Every major distributed system must handle the case where some nodes can communicate with each other but not others. What you build here determines whether your system degrades gracefully or corrupts data silently.

Used in: AWS, Google Cloud, Azure

Real distributed systems handle many messages concurrently. Your current synchronous implementation processes one message at a time, which limits throughput.

Refactor your node to handle messages concurrently:

1. Read messages in the main thread
2. Dispatch each message to a handler that runs concurrently
3. Ensure thread-safe access to shared state (node\_id, counters, etc.)

This prepares you for more complex workloads where you need to send messages while waiting for responses.

### Sample Test Cases

Handle 3 concurrent echo messagesTimeout: 5000ms

```
{"src":"c0","dest":"n1","body":{"type":"init","msg_id":1,"node_id":"n1","node_ids":["n1"]}}
{"src":"c1","dest":"n1","body":{"type":"echo","msg_id":2,"echo":"test1"}}
{"src":"c2","dest":"n1","body":{"type":"echo","msg_id":3,"echo":"test2"}}
{"src":"c3","dest":"n1","body":{"type":"echo","msg_id":4,"echo":"test3"}}
```

```
{"src": "n1", "dest": "c0", "body": {"type": "init_ok", "in_reply_to": 1, "msg_id": 0}}
{"src": "n1", "dest": "c1", "body": {"type": "echo_ok", "echo": "test1", "in_reply_to": 2, "msg_id": 1}}
{"src": "n1", "dest": "c2", "body": {"type": "echo_ok", "echo": "test2", "in_reply_to": 3, "msg_id": 2}}
{"src": "n1", "dest": "c3", "body": {"type": "echo_ok", "echo": "test3", "in_reply_to": 4, "msg_id": 3}}
```

### Hints

Hint 1▾

Use threading or asyncio for concurrent handling

Hint 2▾

Be careful with shared state

Hint 3▾

Consider using a queue for message processing

Hint 4▾

Go tip: handle init synchronously before spawning goroutines for other messages

Hint 5▾

Buffer output and sort by in\_reply\_to before printing to get deterministic msg\_id ordering

main.go

go

Create Async Event Loop for Concurrent Message Handling - The Messenger | Build Distributed Systems<iframe src="chrome-extension://cnjifjpddelmedmihgijeibhnjfabmlf/side-panel.html?context=iframe"></iframe>