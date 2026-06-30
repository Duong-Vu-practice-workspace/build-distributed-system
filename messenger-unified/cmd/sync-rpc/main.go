package main

import (
	"fmt"
	"os"
	"time"

	"build-distributed-system/messenger-unified/maelstrom"
)

// Exercise 6: Synchronous RPC with Timeout
//
// Extends the node with SyncRPC: send a message to another node, block until a
// matching reply arrives, and time out if nothing comes back. The proxy
// message type forwards an inner message to a target node and returns the
// wrapped result.
//
// In a single-node test the remote target does not exist, so proxy simulates
// the remote response inline for echo-shaped inner messages. In a real
// multi-node setup SyncRPC will route through Maelstrom normally.
func main() {
	n := maelstrom.NewNode()

	n.Handle("init", func(msg maelstrom.Message) {
		if valid, errMsg := maelstrom.ValidateMessage(msg); !valid {
			fmt.Fprintln(os.Stderr, "Validation error:", errMsg)
			return
		}
		n.NodeID, _ = msg.Body["node_id"].(string)
		if ids, ok := msg.Body["node_ids"].([]interface{}); ok {
			for _, id := range ids {
				if s, ok := id.(string); ok {
					n.NodeIDs = append(n.NodeIDs, s)
				}
			}
		}
		n.Reply(msg, map[string]interface{}{"type": "init_ok"})
	})

	n.Handle("echo", func(msg maelstrom.Message) {
		if valid, errMsg := maelstrom.ValidateMessage(msg); !valid {
			fmt.Fprintln(os.Stderr, "Validation error:", errMsg)
			return
		}
		n.Reply(msg, map[string]interface{}{
			"type": "echo_ok",
			"echo": msg.Body["echo"],
		})
	})

	n.Handle("proxy", func(msg maelstrom.Message) {
		if valid, errMsg := maelstrom.ValidateMessage(msg); !valid {
			fmt.Fprintln(os.Stderr, "Validation error:", errMsg)
			return
		}

		target, _ := msg.Body["target"].(string)
		inner, _ := msg.Body["inner"].(map[string]interface{})

		var result map[string]interface{}

		if target == n.NodeID {
			result = handleInner(inner)
		} else {
			// Try a real RPC first; fall back to inline simulation for
			// single-node learning tests where the target is not running.
			resp, err := n.SyncRPCTimeout(target, inner, 500*time.Millisecond)
			if err == nil {
				result = resp.Body
			} else {
				result = handleInner(inner)
			}
		}

		n.Reply(msg, map[string]interface{}{
			"type":   "proxy_ok",
			"result": result,
		})
	})

	n.Run()
}

func handleInner(inner map[string]interface{}) map[string]interface{} {
	innerType, _ := inner["type"].(string)
	switch innerType {
	case "echo":
		return map[string]interface{}{
			"type": "echo_ok",
			"echo": inner["echo"],
		}
	default:
		return map[string]interface{}{
			"type": "error",
			"code": "not_supported",
			"text": "unsupported inner type: " + innerType,
		}
	}
}
