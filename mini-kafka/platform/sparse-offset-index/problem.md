## Implementation

A log with millions of messages can't scan from offset 0 every time. Kafka maintains a **sparse offset index** — a sorted array of (offset, file-position) pairs sampled every N messages. A seek binary-searches the index for the nearest entry, then scans forward by at most N messages to reach the exact target.

This achieves O(log n) seeks with modest memory: indexing every 4096 bytes means at most 4096 sequential reads per seek, regardless of log size.

**Commands**

| Command | Output |
| --- | --- |
| `INDEX_INTERVAL <n>` | *(index every n-th message, starting at offset 0)* |
| `APPEND <message>` | `offset:<n>` |
| `SEEK <offset>` | `<message>` at that offset or `ERROR: offset not found` |
| `INDEX_SIZE` | number of entries in the sparse index |

**Example**

```
INDEX_INTERVAL 3
APPEND a   -> offset:0  (indexed, position 0)
APPEND b   -> offset:1
APPEND c   -> offset:2
APPEND d   -> offset:3  (indexed, position 3)
APPEND e   -> offset:4
SEEK 4     -> e  (index has entry at 3; scan from 3 forward to 4)
INDEX_SIZE -> 2  (entries at offsets 0 and 3)
```

### Sample Test Cases

```
INDEX_INTERVAL 3
APPEND a
APPEND b
APPEND c
APPEND d
SEEK 0
SEEK 3
```

```
offset:0
offset:1
offset:2
offset:3
a
d
```

```
INDEX_INTERVAL 4
APPEND zero
APPEND one
APPEND two
APPEND three
APPEND four
SEEK 2
SEEK 4
```

```
offset:0
offset:1
offset:2
offset:3
offset:4
two
four
```

```
INDEX_INTERVAL 4
APPEND a
APPEND b
APPEND c
APPEND d
APPEND e
APPEND f
APPEND g
APPEND h
INDEX_SIZE
```

```
offset:0
offset:1
offset:2
offset:3
offset:4
offset:5
offset:6
offset:7
2
```

```
INDEX_INTERVAL 2
APPEND x
APPEND y
SEEK 99
```

```
offset:0
offset:1
ERROR: offset not found
```

main.py

python

1319

<iframe allow="clipboard-write; web-share" src="chrome-extension://cnjifjpddelmedmihgijeibhnjfabmlf/side-panel.html?context=iframe"></iframe>