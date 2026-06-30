package main

import (
	"build-distributed-system/messenger-unified/maelstrom"
)

// Exercise 3: Echo Service
//
// Implements the Maelstrom echo workload. After init, the node replies to
// each echo message with echo_ok, preserving the original echo value.
//
// Run with Maelstrom:
//
//	maelstrom test -w echo --bin ./messenger-unified/cmd/echo-service/echo-service --node-count 1 --time-limit 10 --log-stderr
func main() {
	n := maelstrom.NewNode()

	n.Handle("init", func(msg maelstrom.Message) {
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
		n.Reply(msg, map[string]interface{}{
			"type": "echo_ok",
			"echo": msg.Body["echo"],
		})
	})

	n.Run()
}
