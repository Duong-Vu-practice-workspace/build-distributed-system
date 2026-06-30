package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"build-distributed-system/messenger-unified/maelstrom"
)

// Exercise 1: JSON Parser
//
// Reads Maelstrom messages from stdin, parses each line as JSON, and prints
// the parsed components. This exercise focuses on understanding the message
// envelope before building a full node.
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var msg maelstrom.Message
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

		fmt.Printf("PARSED: %s|%s|%s\n", msg.Src, msg.Dest, bodyType)
	}
}
