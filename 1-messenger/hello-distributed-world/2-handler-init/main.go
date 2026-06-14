package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

// Node represents a Maelstrom node
type Node struct {
	NodeID    string
	NodeIDs   []string
	NextMsgID int
	mu        sync.Mutex
}

// Message represents a Maelstrom message
type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

// Send sends a message to a destination node
func (n *Node) Send(dest string, body map[string]interface{}) {
	message := Message{
		Src:  n.NodeID,
		Dest: dest,
		Body: body,
	}
	s, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(s))
}

// Reply sends a response to an incoming request
func (n *Node) Reply(request Message, body map[string]interface{}) {
	n.mu.Lock()
	body["msg_id"] = n.NextMsgID
	n.NextMsgID++
	n.mu.Unlock()
	n.Send(request.Src, body)
}

func main() {
	node := &Node{}
	node.NextMsgID = 0
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}

		msgType, _ := msg.Body["type"].(string)
		if msgType == "init" {
			// TODO: Handle init message
			// 1. Store node_id and node_ids
			// 2. Reply with init_ok
			nodeId := msg.Body["node_id"].(string)
			rawNodeIds := msg.Body["node_ids"].([]interface{})
			nodeIds := make([]string, len(rawNodeIds))
			for i, value := range rawNodeIds {
				nodeIds[i] = value.(string)
			}
			node.NodeID = nodeId
			node.NodeIDs = nodeIds
			responseBody := map[string]interface{}{
				"type":        "init_ok",
				"in_reply_to": msg.Body["msg_id"],
			}
			node.Reply(msg, responseBody)
		}
	}
}
