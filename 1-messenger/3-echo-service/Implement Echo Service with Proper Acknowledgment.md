---
title: "Implement Echo Service with Proper Acknowledgment"
source: "https://builddistributedsystem.com/tracks/messenger/tasks/task-1-3-echo-service"
author:
  - "[[Mohit Mishra]]"
published:
created: 2026-06-13
description: "The echo workload is the simplest Maelstrom workload. Clients send echo messages containing a value, and your node must echo that value back.  Request format:"
tags:
  - "clippings"
---
## Implementation

Retry logic with exponential backoff and jitter is in every production service. When AWS has a partial outage, the difference between a retry storm that worsens the outage and a graceful degradation is jittered backoff. Stripe's payment retries, AWS SDK retries, and gRPC's retry policy all implement what you are building here.

Used in: Stripe, AWS SDK, gRPC

The echo workload is the simplest Maelstrom workload. Clients send echo messages containing a value, and your node must echo that value back.

Request format:

```json
{
  "type": "echo",
  "msg_id": 1,
  "echo": "Please echo 35"
}
```

Expected response:

```json
{
  "type": "echo_ok",
  "msg_id": 1,
  "in_reply_to": 1,
  "echo": "Please echo 35"
}
```

Combine your init handling from the previous task with a new echo handler. Your node should handle both message types.

### Sample Test Cases

```
{"src":"c0","dest":"n1","body":{"type":"init","msg_id":1,"node_id":"n1","node_ids":["n1"]}}
{"src":"c1","dest":"n1","body":{"type":"echo","msg_id":2,"echo":"hello"}}
```

```
{"src": "n1", "dest": "c0", "body": {"type": "init_ok", "in_reply_to": 1, "msg_id": 0}}
{"src": "n1", "dest": "c1", "body": {"type": "echo_ok", "echo": "hello", "in_reply_to": 2, "msg_id": 1}}
```

### Hints

Hint 1▾

Echo messages contain an "echo" field with the value to echo back

Hint 2▾

Response type is "echo\_ok"

Hint 3▾

Include the original echo value in your response

### Resources[Echo Workload](https://fly.io/dist-sys/1/)

[

Fly.io Gossip Glomers Echo challenge walkthrough

](https://fly.io/dist-sys/1/)

main.go

go

<iframe src="chrome-extension://cnjifjpddelmedmihgijeibhnjfabmlf/side-panel.html?context=iframe"></iframe>