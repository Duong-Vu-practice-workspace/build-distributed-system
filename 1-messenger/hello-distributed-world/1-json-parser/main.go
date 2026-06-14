package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Message represents a Maelstrom message
type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

func main() {
	// TODO: Read JSON messages from stdin
	// Each line is a complete JSON message
	// Parse and print: PARSED: src|dest|body_type
	// Log details to stderr for debugging

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
			continue
		}
		bodyType := "unknown"
		if msg.Body != nil {
			if v, ok := msg.Body["type"]; ok && v != nil {
				bodyType = fmt.Sprint(v)
			}
		}
		// Your code here
		fmt.Printf("PARSED: %s|%s|%s\n", msg.Src, msg.Dest, bodyType)

	}
}
