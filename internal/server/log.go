package server

import (
	"fmt"
	"sync"
)

// Log stores the Record Entries in a thread safe manner
type Log struct {
	mu      sync.RWMutex
	records []Record
	length  uint64
}

// Record stores bytes and an Offset
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

// ErrOffsetNotFound is the error returned if the offset does not exist
var ErrOffsetNotFound = fmt.Errorf("offset not found")

// NewLog creates a new Log instance
func NewLog() *Log {
	return &Log{}
}

// Append adds a new record to the Log and returns the offset
func (c *Log) Append(record Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record.Offset = c.length
	c.length++
	c.records = append(c.records, record)
	return record.Offset, nil
}

// Read returns the record at `offset` or errors if it does not exist
func (c *Log) Read(offset uint64) (record Record, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if offset >= c.length {
		return Record{}, ErrOffsetNotFound
	}
	return c.records[offset], nil
}
