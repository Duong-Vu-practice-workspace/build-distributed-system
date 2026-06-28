I'm in plan mode - let me provide a comprehensive implementation plan.

## Detailed Implementation Plan

### Current State Analysis
- Existing code has basic `Node` struct with `NodeID`, `NodeIDs`, `NextMsgID`, mutexes
- Has `Send()` and `Reply()` methods
- Handles `init` and `echo` message types
- Missing: `sync_rpc`, `proxy` handling, timeout mechanism

### Implementation Steps

**1. Update Node struct** (add fields):
- `pendingRequests map[int]chan map[string]interface{}` - stores response channels keyed by msg_id
- `timeout time.Duration` - configurable timeout (default 1s)

**2. Initialize in main()**:
- `pendingRequests: make(map[int]chan map[string]interface{})`
- `timeout: 1 * time.Second`

**3. Add `sync_rpc` method**:
```
sync_rpc(dest, body) -> (map[string]interface{}, error)
  1. Generate msg_id via Send()
  2. Create response channel
  3. Store channel in pendingRequests[msg_id]
  4. Send message to dest
  5. Wait: select {
       case resp := <-ch: return resp, nil
       case <-time.After(timeout): delete pendingRequests[msg_id]; return nil, timeout error
     }
```

**4. Add `resolvePending` method**:
- Called when response arrives with `in_reply_to`
- Look up channel by `in_reply_to`, send response, delete entry

**5. Add `proxy` message handler**:
```
case "proxy":
  target := msg.Body["target"].(string)
  inner := msg.Body["inner"].(map[string]interface{})
  result, err := node.sync_rpc(target, inner)
  if err != nil { /* handle timeout */ }
  node.Reply(msg, map[string]interface{}{
    "type": "proxy_ok",
    "result": result,
  })
```

**6. Update main loop** to resolve pending requests before processing:
- Extract `in_reply_to` from incoming messages
- Call `resolvePending(msg)` to unblock sync_rpc callers

**7. Update ValidateMessage**: Remove the incorrect validation that requires `in_reply_to` for all messages

Would you like me to proceed with implementation?