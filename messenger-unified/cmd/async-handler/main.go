package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"

	"build-distributed-system/messenger-unified/maelstrom"
)

// pendingResponse captures an outgoing response before it is assigned a final
// sequential msg_id and printed.
type pendingResponse struct {
	Dest    string
	Body    map[string]interface{}
	ReplyTo int
}

// Exercise 5: Async Handler
//
// Demonstrates concurrent message handling. Messages are processed in
// goroutines; their responses are collected, sorted by the original
// in_reply_to value, assigned sequential msg_ids, and then printed.
//
// Note: This pattern is useful for learning, but real Maelstrom nodes usually
// reply immediately to keep latency low.
func main() {
	n := maelstrom.NewNode()
	scanner := bufio.NewScanner(os.Stdin)

	responses := make(chan pendingResponse, 100)
	var wg sync.WaitGroup

	for scanner.Scan() {
		var msg maelstrom.Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
			continue
		}

		msgType := msg.Type()

		switch msgType {
		case "init":
			// Init is handled synchronously so the cluster metadata is ready
			// before any concurrent handlers run.
			if valid, errMsg := maelstrom.ValidateMessage(msg); !valid {
				fmt.Fprintln(os.Stderr, "Validation error:", errMsg)
				continue
			}
			n.NodeID, _ = msg.Body["node_id"].(string)
			if ids, ok := msg.Body["node_ids"].([]interface{}); ok {
				for _, id := range ids {
					if s, ok := id.(string); ok {
						n.NodeIDs = append(n.NodeIDs, s)
					}
				}
			}
			responses <- pendingResponse{
				Dest:    msg.Src,
				Body:    map[string]interface{}{"type": "init_ok"},
				ReplyTo: replyTo(msg),
			}

		default:
			wg.Add(1)
			go func(m maelstrom.Message) {
				defer wg.Done()

				if valid, errMsg := maelstrom.ValidateMessage(m); !valid {
					fmt.Fprintln(os.Stderr, "Validation error:", errMsg)
					return
				}

				switch m.Type() {
				case "echo":
					responses <- pendingResponse{
						Dest: m.Src,
						Body: map[string]interface{}{
							"type": "echo_ok",
							"echo": m.Body["echo"],
						},
						ReplyTo: replyTo(m),
					}
				}
			}(msg)
		}
	}

	wg.Wait()
	close(responses)

	var pending []pendingResponse
	for r := range responses {
		pending = append(pending, r)
	}

	sort.Slice(pending, func(i, j int) bool {
		return pending[i].ReplyTo < pending[j].ReplyTo
	})

	msgID := 0
	for _, r := range pending {
		r.Body["msg_id"] = msgID
		msgID++
		n.Send(r.Dest, r.Body)
	}
}

func replyTo(msg maelstrom.Message) int {
	if msgID, ok := msg.Body["msg_id"].(float64); ok {
		return int(msgID)
	}
	return 0
}
