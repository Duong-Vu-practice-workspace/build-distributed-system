package main

import (
	"fmt"
	"os"

	"build-distributed-system/messenger-unified/maelstrom"
)

// Exercise 4: Envelope Validation
//
// Extends the echo service with structural validation before processing. Any
// message missing required fields is logged to stderr and ignored.
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

	n.Run()
}
