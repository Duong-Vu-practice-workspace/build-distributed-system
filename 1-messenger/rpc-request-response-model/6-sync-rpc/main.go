package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Node struct {
	NodeID          string
	NodeIDs         []string
	NextMsgID       int
	mu              sync.Mutex
	outMu           sync.Mutex
	pendingRequests map[int]chan map[string]interface{}
	timeout         time.Duration
}

type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

func (n *Node) Send(dest string, body map[string]interface{}) int {
	n.mu.Lock()
	body["msg_id"] = n.NextMsgID
	n.NextMsgID++
	n.mu.Unlock()

	msg := Message{Src: n.NodeID, Dest: dest, Body: body}
	output, _ := json.Marshal(msg)

	n.outMu.Lock()
	fmt.Println(string(output))
	n.outMu.Unlock()
	return body["msg_id"].(int)
}

func (n *Node) Reply(request Message, body map[string]interface{}) {
	if msgID, ok := request.Body["msg_id"].(float64); ok {
		body["in_reply_to"] = int(msgID)
	}
	n.Send(request.Src, body)

}

// ValidateMessage checks if a message has required structure
func ValidateMessage(msg Message) (bool, string) {
	// TODO: Validate message structure
	// Return true if valid, false with error message otherwise
	if len(msg.Src) == 0 || len(msg.Dest) == 0 || msg.Body == nil {
		errString := "Message must have src, dest, body fields"
		fmt.Errorf(errString)
		return false, errString
	}
	if _, ok := msg.Body["type"]; !ok {
		errString := "Body must have a type field"
		fmt.Errorf(errString)
		return false, errString
	}
	if _, ok := msg.Body["in_reply_to"]; !ok {
		_, isOk := msg.Body["msg_id"]
		if !isOk {
			errString := "requests expecting responses, body should have a msg_id"
			fmt.Errorf(errString)
			return false, errString
		}
	}

	return true, ""
}
func (n *Node) SyncRpc(request Message, body map[string]interface{}) (map[string]interface{}, error) {
	msg_id := n.Send(request.Dest, body)

}
func main() {
	node := &Node{
		pendingRequests: make(map[int]chan map[string]interface{}),
		timeout:         1 * time.Second,
	}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Invalid JSON:", err)
			continue
		}

		// TODO: Validate message before processing
		if valid, errMsg := ValidateMessage(msg); !valid {
			fmt.Fprintln(os.Stderr, "Validation error:", errMsg)
			continue
		}

		msgType, _ := msg.Body["type"].(string)
		switch msgType {
		case "init":
			node.NodeID, _ = msg.Body["node_id"].(string)
			if ids, ok := msg.Body["node_ids"].([]interface{}); ok {
				for _, id := range ids {
					node.NodeIDs = append(node.NodeIDs, id.(string))
				}
			}
			node.Reply(msg, map[string]interface{}{"type": "init_ok"})
		case "echo":
			node.Reply(msg, map[string]interface{}{
				"type": "echo_ok",
				"echo": msg.Body["echo"],
			})
		}
	}
}
