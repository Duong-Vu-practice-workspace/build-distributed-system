// Package maelstrom provides primitives for building Maelstrom-compatible nodes.
//
// Maelstrom communicates with nodes via JSON messages on stdin/stdout. Each
// line is a complete message with the shape:
//
//	{"src": "...", "dest": "...", "body": {"type": "...", ...}}
package maelstrom

import (
	"encoding/json"
)

// Message is the envelope used by Maelstrom.
type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

// Type returns the message type from the body, or an empty string if missing.
func (m Message) Type() string {
	if m.Body == nil {
		return ""
	}
	t, _ := m.Body["type"].(string)
	return t
}

// String returns the JSON encoding of the message.
func (m Message) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

// ValidateMessage checks that msg has the minimal fields Maelstrom requires.
//
// Requests that expect a reply must include a msg_id. Responses include
// in_reply_to instead, so the msg_id check is skipped for those messages.
func ValidateMessage(msg Message) (bool, string) {
	if msg.Src == "" || msg.Dest == "" || msg.Body == nil {
		return false, "message must have src, dest, and body fields"
	}
	if _, ok := msg.Body["type"]; !ok {
		return false, "body must have a type field"
	}
	if _, ok := msg.Body["in_reply_to"]; !ok {
		if _, ok := msg.Body["msg_id"]; !ok {
			return false, "request expecting a response must have a msg_id"
		}
	}
	return true, ""
}
