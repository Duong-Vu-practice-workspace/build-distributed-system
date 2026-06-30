package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PartitionLog struct {
	messages []string
}

func (l *PartitionLog) Append(message string) int {
	// TODO
	if l == nil {
		return 0
	}
	l.messages = append(l.messages, message)
	return len(l.messages) - 1
}

func (l *PartitionLog) Read(offset int) (string, error) {
	// TODO: return message at offset, or error if out of range
	if offset > len(l.messages)-1 || offset < 0 {
		return "", fmt.Errorf("ERROR: offset out of range")
	}
	return l.messages[offset], nil
}

func (l *PartitionLog) Tail() int {
	// TODO: return next offset to be assigned
	return len(l.messages)
}

func main() {
	log := &PartitionLog{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		cmd := parts[0]

		switch cmd {
		case "APPEND":
			offset := log.Append(parts[1])
			fmt.Printf("offset:%d\n", offset)
		case "READ":
			offset, _ := strconv.Atoi(parts[1])
			msg, err := log.Read(offset)
			if err != nil {
				fmt.Println("ERROR: offset out of range")
			} else {
				fmt.Println(msg)
			}
		case "TAIL":
			fmt.Println(log.Tail())
		case "SIZE":
			fmt.Println(len(log.messages))
		}
	}
}
