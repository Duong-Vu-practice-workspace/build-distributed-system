
## Implementation

Before processing any workload, Maelstrom sends an init message to each node. This message tells your node its identity and the full cluster membership.

The init message looks like:

```json
{
  "type": "init",
  "msg_id": 1,
  "node_id": "n1",
  "node_ids": ["n1", "n2", "n3"]
}
```

Your task is to handle the init message by storing the node\_id (your identity) and node\_ids (all cluster members). Then respond with an init\_ok message:

```json
{
  "type": "init_ok",
  "in_reply_to": 1
}
```

The in\_reply\_to field must match the msg\_id from the init request.

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

The init message contains node\_id and node\_ids fields

Hint 2▾

Store these values for use in subsequent message handling

Hint 3▾

Reply with init\_ok message type

main.py

python

<iframe src="chrome-extension://cnjifjpddelmedmihgijeibhnjfabmlf/side-panel.html?context=iframe"></iframe>