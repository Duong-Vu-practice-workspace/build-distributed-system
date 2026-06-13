package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout captures all output written to os.Stdout during the execution of the provided function.
func captureStdout(fn func()) (string, error) {
	// Save the original stdout
	oldStdout := os.Stdout

	// Create a pipe to capture output
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	// Replace os.Stdout with the pipe's writer
	os.Stdout = w

	// Run the function that produces output
	go func() {
		fn()
		w.Close()
	}()

	// Read the captured output
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return "", err
	}

	// Restore the original stdout
	os.Stdout = oldStdout

	return buf.String(), nil
}

func TestHandleInitMessage(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedNodeID string
		expectedIDs    []string
		checkResponse  func(t *testing.T, response map[string]interface{})
	}{
		{
			name:           "basic init message with single node",
			input:          `{"src":"c0","dest":"n1","body":{"type":"init","msg_id":1,"node_id":"n1","node_ids":["n1"]}}`,
			expectedNodeID: "n1",
			expectedIDs:    []string{"n1"},
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["type"] != "init_ok" {
					t.Errorf("expected type 'init_ok', got %v", response["type"])
				}
				if response["in_reply_to"] != float64(1) {
					t.Errorf("expected in_reply_to 1, got %v", response["in_reply_to"])
				}
			},
		},
		{
			name:           "init message with multiple nodes",
			input:          `{"src":"c0","dest":"n1","body":{"type":"init","msg_id":5,"node_id":"n1","node_ids":["n1","n2","n3"]}}`,
			expectedNodeID: "n1",
			expectedIDs:    []string{"n1", "n2", "n3"},
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["type"] != "init_ok" {
					t.Errorf("expected type 'init_ok', got %v", response["type"])
				}
				if response["in_reply_to"] != float64(5) {
					t.Errorf("expected in_reply_to 5, got %v", response["in_reply_to"])
				}
			},
		},
		{
			name:           "init message with msg_id 0",
			input:          `{"src":"c0","dest":"n1","body":{"type":"init","msg_id":0,"node_id":"n1","node_ids":["n1"]}}`,
			expectedNodeID: "n1",
			expectedIDs:    []string{"n1"},
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["type"] != "init_ok" {
					t.Errorf("expected type 'init_ok', got %v", response["type"])
				}
				if response["in_reply_to"] != float64(0) {
					t.Errorf("expected in_reply_to 0, got %v", response["in_reply_to"])
				}
			},
		},
		{
			name:           "init message with large msg_id",
			input:          `{"src":"c0","dest":"n1","body":{"type":"init","msg_id":999999,"node_id":"n1","node_ids":["n1"]}}`,
			expectedNodeID: "n1",
			expectedIDs:    []string{"n1"},
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["type"] != "init_ok" {
					t.Errorf("expected type 'init_ok', got %v", response["type"])
				}
				if response["in_reply_to"] != float64(999999) {
					t.Errorf("expected in_reply_to 999999, got %v", response["in_reply_to"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{}

			// Parse the input message
			var msg Message
			if err := json.Unmarshal([]byte(tt.input), &msg); err != nil {
				t.Fatalf("Failed to unmarshal input: %v", err)
			}

			// Verify parsing
			msgType, _ := msg.Body["type"].(string)
			if msgType != "init" {
				t.Fatalf("Expected type 'init', got %v", msgType)
			}

			// Capture stdout during processing
			output, err := captureStdout(func() {
				// Process the init message
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
					"msg_id":      0,
				}
				node.Reply(msg, responseBody)
			})

			if err != nil {
				t.Fatalf("Failed to capture stdout: %v", err)
			}

			// Verify node metadata was stored
			if node.NodeID != tt.expectedNodeID {
				t.Errorf("Expected NodeID '%s', got '%s'", tt.expectedNodeID, node.NodeID)
			}
			if len(node.NodeIDs) != len(tt.expectedIDs) {
				t.Errorf("Expected %d node IDs, got %d", len(tt.expectedIDs), len(node.NodeIDs))
			}
			for i, expectedID := range tt.expectedIDs {
				if i >= len(node.NodeIDs) || node.NodeIDs[i] != expectedID {
					t.Errorf("Expected NodeIDs[%d] = '%s', got '%s'", i, expectedID, node.NodeIDs[i])
				}
			}

			// Verify the response
			var response Message
			if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v\nOutput: %s", err, output)
			}

			// Verify response structure
			if response.Src != tt.expectedNodeID {
				t.Errorf("Expected response src '%s', got '%s'", tt.expectedNodeID, response.Src)
			}
			if response.Dest != "c0" {
				t.Errorf("Expected response dest 'c0', got '%s'", response.Dest)
			}

			// Run the specific response checks for this test case
			if tt.checkResponse != nil {
				tt.checkResponse(t, response.Body)
			}
		})
	}
}

func TestNodeStructFields(t *testing.T) {
	t.Run("verify node struct has required fields", func(t *testing.T) {
		node := &Node{
			NodeID:  "n1",
			NodeIDs: []string{"n1", "n2", "n3"},
		}

		if node.NodeID != "n1" {
			t.Errorf("Expected NodeID 'n1', got '%s'", node.NodeID)
		}
		if len(node.NodeIDs) != 3 {
			t.Errorf("Expected 3 node IDs, got %d", len(node.NodeIDs))
		}
	})
}

func TestReplyMethod(t *testing.T) {
	t.Run("reply sets in_reply_to correctly", func(t *testing.T) {
		node := &Node{NodeID: "n1"}
		request := Message{
			Src:  "c0",
			Dest: "n1",
			Body: map[string]interface{}{"msg_id": 42},
		}

		output, err := captureStdout(func() {
			responseBody := map[string]interface{}{
				"type":        "init_ok",
				"in_reply_to": 42,
				"msg_id":      0,
			}
			node.Reply(request, responseBody)
		})

		if err != nil {
			t.Fatalf("Failed to capture stdout: %v", err)
		}

		var response Message
		if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Body["in_reply_to"] != float64(42) {
			t.Errorf("Expected in_reply_to 42, got %v", response.Body["in_reply_to"])
		}
	})
}

func TestSendMethod(t *testing.T) {
	t.Run("send outputs valid JSON to stdout", func(t *testing.T) {
		node := &Node{NodeID: "n1"}
		output, err := captureStdout(func() {
			node.Send("c0", map[string]interface{}{"type": "test"})
		})

		if err != nil {
			t.Fatalf("Failed to capture stdout: %v", err)
		}

		var msg Message
		if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &msg); err != nil {
			t.Fatalf("Failed to unmarshal output: %v", err)
		}

		if msg.Src != "n1" {
			t.Errorf("Expected src 'n1', got '%s'", msg.Src)
		}
		if msg.Dest != "c0" {
			t.Errorf("Expected dest 'c0', got '%s'", msg.Dest)
		}
	})
}

func TestJSONMessageParsing(t *testing.T) {
	t.Run("parse init message correctly", func(t *testing.T) {
		input := `{"src":"c0","dest":"n1","body":{"type":"init","msg_id":1,"node_id":"n1","node_ids":["n1","n2"]}}`
		var msg Message
		if err := json.Unmarshal([]byte(input), &msg); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if msg.Src != "c0" {
			t.Errorf("Expected src 'c0', got '%s'", msg.Src)
		}
		if msg.Dest != "n1" {
			t.Errorf("Expected dest 'n1', got '%s'", msg.Dest)
		}

		msgType, _ := msg.Body["type"].(string)
		if msgType != "init" {
			t.Errorf("Expected type 'init', got '%s'", msgType)
		}
	})
}

func TestMainFunctionWithMockedStdin(t *testing.T) {
	// This test simulates the main loop with a mocked stdin
	input := "{\"src\":\"c0\",\"dest\":\"n1\",\"body\":{\"type\":\"init\",\"msg_id\":1,\"node_id\":\"n1\",\"node_ids\":[\"n1\"]}}\n"

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Write test input to the pipe
	go func() {
		w.WriteString(input)
		w.Close()
	}()

	// Save original stdin and replace with our pipe
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// We can't easily test main() directly, but we can test the logic inside
	node := &Node{}
	scanner := bufio.NewScanner(os.Stdin)
	var responseOutput string

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			t.Fatalf("Error parsing JSON: %v", err)
		}

		msgType, _ := msg.Body["type"].(string)
		if msgType == "init" {
			nodeId := msg.Body["node_id"].(string)
			rawNodeIds := msg.Body["node_ids"].([]interface{})
			nodeIds := make([]string, len(rawNodeIds))
			for i, value := range rawNodeIds {
				nodeIds[i] = value.(string)
			}
			node.NodeID = nodeId
			node.NodeIDs = nodeIds

			// Capture the response
			responseOutput, _ = captureStdout(func() {
				responseBody := map[string]interface{}{
					"type":        "init_ok",
					"in_reply_to": msg.Body["msg_id"],
					"msg_id":      0,
				}
				node.Reply(msg, responseBody)
			})
		}
	}

	// Verify the response
	if responseOutput == "" {
		t.Fatal("Expected a response output, got nothing")
	}

	var response Message
	if err := json.Unmarshal([]byte(strings.TrimSpace(responseOutput)), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Body["type"] != "init_ok" {
		t.Errorf("Expected type 'init_ok', got %v", response.Body["type"])
	}
	if response.Body["in_reply_to"] != float64(1) {
		t.Errorf("Expected in_reply_to 1, got %v", response.Body["in_reply_to"])
	}
}
