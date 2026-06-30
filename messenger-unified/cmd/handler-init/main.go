package main

import (
	"build-distributed-system/messenger-unified/maelstrom"
)

// Exercise 2: Handler Init
//
// Builds a minimal Maelstrom node that responds to the init message with
// init_ok. The node stores its node_id and the list of node_ids from the
// cluster.
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

		n.Reply(msg, map[string]interface{}{
			"type": "init_ok",
		})
	})

	n.Run()
}
