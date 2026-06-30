package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type IndexedLog struct {
	messages []string
	index    []IndexEntry
	interval int
}

type IndexEntry struct {
	offset   int
	position int
}

func NewIndexedLog() *IndexedLog {
	return &IndexedLog{interval: 4}
}

func (l *IndexedLog) SetInterval(n int) {
	l.interval = n
}

func (l *IndexedLog) Append(message string) int {
	offset := len(l.messages)
	l.messages = append(l.messages, message)
	// TODO: if offset % interval == 0, add (offset, offset) to index
	// (position here equals offset since we store messages in a list by offset)
	return offset
}

func (l *IndexedLog) Seek(offset int) string {
	if len(l.messages) == 0 || offset < 0 || offset >= len(l.messages) {
		return "ERROR: offset not found"
	}
	// TODO:
	// 1. Binary search self.index for largest (indexed_off, pos) where indexed_off <= offset
	// 2. From that position, scan forward to find exact offset
	return ""
}

func main() {
	log := NewIndexedLog()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		cmd := parts[0]

		switch cmd {
		case "INDEX_INTERVAL":
			n, _ := strconv.Atoi(parts[1])
			log.SetInterval(n)
		case "APPEND":
			fmt.Printf("offset:%d\n", log.Append(parts[1]))
		case "SEEK":
			offset, _ := strconv.Atoi(parts[1])
			fmt.Println(log.Seek(offset))
		case "INDEX_SIZE":
			fmt.Println(len(log.index))
		}
	}
}
