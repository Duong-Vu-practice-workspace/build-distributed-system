package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

type Node struct {
	NodeID    string
	NodeIDs   []string
	NextMsgID int
	mu        sync.Mutex
	//outMu     sync.Mutex
}

type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}
type pendingResponse struct {
	Dest    string
	Body    map[string]interface{}
	ReplyTo int //value of in_reply_to
}

func (n *Node) Send(dest string, body map[string]interface{}, responses chan<- pendingResponse, replyTo int) {
	n.mu.Lock()
	body["msg_id"] = n.NextMsgID
	n.NextMsgID++
	n.mu.Unlock()

	responses <- pendingResponse{
		Dest:    dest,
		Body:    body,
		ReplyTo: replyTo,
	}
}

func (n *Node) Reply(request Message, body map[string]interface{}, responses chan<- pendingResponse) {
	replyTo := 0
	if msgID, ok := request.Body["msg_id"].(float64); ok {
		body["in_reply_to"] = int(msgID)
		replyTo = int(msgID)
	}
	n.Send(request.Src, body, responses, replyTo)
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

func (n *Node) HandleInit(msg Message, responses chan<- pendingResponse) {
	n.mu.Lock()
	n.NodeID, _ = msg.Body["node_id"].(string)
	if ids, ok := msg.Body["node_ids"].([]interface{}); ok {
		for _, id := range ids {
			n.NodeIDs = append(n.NodeIDs, id.(string))
		}
	}
	n.mu.Unlock()
	n.Reply(msg, map[string]interface{}{"type": "init_ok"}, responses)
}
func (n *Node) HandleMessage(msg Message, responses chan<- pendingResponse) {
	if valid, errMsg := ValidateMessage(msg); !valid {
		fmt.Fprintln(os.Stderr, "Validation error:", errMsg)
		return
	}

	msgType, _ := msg.Body["type"].(string)
	switch msgType {
	case "echo":
		n.Reply(msg, map[string]interface{}{
			"type": "echo_ok",
			"echo": msg.Body["echo"],
		}, responses)
	}
}

func main() {
	node := &Node{}
	scanner := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup
	responses := make(chan pendingResponse, 100) // buffered channel

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}
		msgType, _ := msg.Body["type"].(string)
		switch msgType {
		case "init":
			// Handle init synchronously (Hint 4)
			node.HandleInit(msg, responses)
		default:
			wg.Add(1)
			go func(m Message) {
				defer wg.Done()
				node.HandleMessage(m, responses)
			}(msg)
		}
	}
	wg.Wait()
	close(responses)

	// Collect all responses
	var pending []pendingResponse
	for r := range responses {
		pending = append(pending, r)
	}

	// Sort by in_reply_to (Hint 5)
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].ReplyTo < pending[j].ReplyTo
	})

	// Assign sequential msg_ids and print
	msgID := 0
	for _, r := range pending {
		r.Body["msg_id"] = msgID
		msgID++
		m := Message{Src: node.NodeID, Dest: r.Dest, Body: r.Body}
		output, _ := json.Marshal(m)
		fmt.Println(string(output))
	}
}
