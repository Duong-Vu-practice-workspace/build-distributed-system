---
title: "Add Message Envelope Validation"
source: "https://builddistributedsystem.com/tracks/messenger/tasks/task-1-4-envelope-validation"
author:
  - "[[Mohit Mishra]]"
published:
created: 2026-06-13
description: "Production systems must handle malformed input gracefully. Your node should validate that incoming messages have the required structure before processing.  Requ"
tags:
  - "clippings"
---
## Implementation

Broadcast is how database change data capture (CDC) works. When Debezium detects a row change in Postgres, it broadcasts that change event to every downstream consumer. The broadcast semantics you implement here are the same ones that power real-time analytics at Uber and DoorDash.

Used in: Debezium, Uber, DoorDash

Production systems must handle malformed input gracefully. Your node should validate that incoming messages have the required structure before processing.

Required validations:

1. Message must be valid JSON
2. Message must have `src`, `dest`, and `body` fields
3. Body must have a `type` field
4. For requests expecting responses, body should have a `msg_id`

If validation fails, log an error to stderr but do not crash. This defensive programming prevents a single bad message from taking down your node.

### Sample Test Cases

```
{
  "src": "c0",
  "dest": "n1",
  "body": {
    "type": "init",
    "msg_id": 1,
    "node_id": "n1",
    "node_ids": [
      "n1"
    ]
  }
}
```

```
{"src": "n1", "dest": "c0", "body": {"type": "init_ok", "in_reply_to": 1, "msg_id": 0}}
```

### Hints

Hint 1▾

Check that required fields are present

Hint 2▾

Handle malformed JSON gracefully

Hint 3▾

Log validation errors to stderr

main.go

go

535765

Add Message Envelope Validation - The Messenger | Build Distributed Systems<iframe src="chrome-extension://cnjifjpddelmedmihgijeibhnjfabmlf/side-panel.html?context=iframe"></iframe>