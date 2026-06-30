package maelstrom

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Handler processes a single incoming message.
type Handler func(msg Message)

// Node represents a single Maelstrom node. It stores identity, cluster
// metadata, and registered message handlers.
type Node struct {
	NodeID    string
	NodeIDs   []string
	NextMsgID int

	mu    sync.Mutex
	outMu sync.Mutex
	rpcMu sync.Mutex

	handlers map[string]Handler
	rpc      *rpcState
}

// NewNode creates a fresh node with an empty handler table.
func NewNode() *Node {
	return &Node{
		NextMsgID: 0,
		handlers:  make(map[string]Handler),
	}
}

// Handle registers a handler for messages with the given type.
func (n *Node) Handle(msgType string, h Handler) {
	n.handlers[msgType] = h
}

// nextMsgID returns the next message ID and increments the counter.
func (n *Node) nextMsgID() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	id := n.NextMsgID
	n.NextMsgID++
	return id
}

// sendWithID writes a message to stdout using the supplied message ID.
func (n *Node) sendWithID(dest string, body map[string]interface{}, msgID int) {
	body["msg_id"] = msgID
	msg := Message{Src: n.NodeID, Dest: dest, Body: body}
	output, _ := json.Marshal(msg)

	n.outMu.Lock()
	fmt.Println(string(output))
	n.outMu.Unlock()
}

// Send sends a message to dest, assigns the next msg_id automatically, and
// returns the assigned id.
func (n *Node) Send(dest string, body map[string]interface{}) int {
	id := n.nextMsgID()
	n.sendWithID(dest, body, id)
	return id
}

// Reply sends a response to the source of request. It sets in_reply_to from
// the request's msg_id.
func (n *Node) Reply(request Message, body map[string]interface{}) {
	if msgID, ok := request.Body["msg_id"].(float64); ok {
		body["in_reply_to"] = int(msgID)
	}
	n.Send(request.Src, body)
}

// Run starts the stdin processing loop. It reads JSON messages one per line,
// dispatches them to registered handlers, and resolves pending RPC replies.
//
// The loop exits when stdin is closed or an unrecoverable error occurs.
func (n *Node) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
			continue
		}

		// RPC replies are matched by in_reply_to before user handlers run.
		if replyTo, ok := msg.Body["in_reply_to"]; ok {
			if resolved := n.resolveRPC(replyTo, msg); resolved {
				continue
			}
		}

		msgType := msg.Type()
		if h, ok := n.handlers[msgType]; ok {
			h(msg)
		}
	}
}
