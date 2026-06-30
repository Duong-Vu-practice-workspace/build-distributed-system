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

func NewPartitionLog() *PartitionLog {
	return &PartitionLog{}
}

func (l *PartitionLog) Append(message string) int {
	offset := len(l.messages)
	l.messages = append(l.messages, message)
	return offset
}

func (l *PartitionLog) Size() int {
	return len(l.messages)
}

// TODO: return message at offset, or error if out of range
func (l *PartitionLog) Read(offset int) (string, error) {
	return "", fmt.Errorf("not implemented")
}

type IndexEntry struct {
	offset   int
	position int
}

type SparseOffsetIndex struct {
	log      *PartitionLog
	index    []IndexEntry
	interval int
}

func NewSparseOffsetIndex(log *PartitionLog, interval int) *SparseOffsetIndex {
	return &SparseOffsetIndex{log: log, interval: interval}
}

func (s *SparseOffsetIndex) Append(message string) int {
	offset := s.log.Append(message)
	// TODO: if offset % interval == 0, add (offset, offset) to index
	return offset
}

func (s *SparseOffsetIndex) Seek(offset int) (string, error) {
	if s.log.Size() == 0 || offset < 0 || offset >= s.log.Size() {
		return "", fmt.Errorf("ERROR: offset not found")
	}
	// TODO:
	// 1. Binary search s.index for largest (indexed_off, pos) where indexed_off <= offset
	// 2. From that position, scan forward via s.log.Read() to find exact offset
	return "", fmt.Errorf("not implemented")
}

func (s *SparseOffsetIndex) IndexSize() int {
	return len(s.index)
}

func main() {
	log := NewPartitionLog()
	sparseIndex := NewSparseOffsetIndex(log, 4)
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
			sparseIndex.interval = n
		case "APPEND":
			fmt.Printf("offset:%d\n", sparseIndex.Append(parts[1]))
		case "SEEK":
			offset, _ := strconv.Atoi(parts[1])
			msg, err := sparseIndex.Seek(offset)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(msg)
			}
		case "INDEX_SIZE":
			fmt.Println(sparseIndex.IndexSize())
		case "SIZE":
			fmt.Println(log.Size())
		}
	}
}
