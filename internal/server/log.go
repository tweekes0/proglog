package server

import (
	"sync"
)

// Record struct that abstracts a log record
type Record struct {
	Value []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

// Log struct that holds all logging records
type Log struct {
	mu sync.Mutex
	records []Record
}

// Return a blank instance of Log
func NewLog() *Log {
	return &Log{}
}

// Append record to log 
func (l *Log) Append(r Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	r.Offset = uint64(len(l.records))
	l.records = append(l.records, r)
	return r.Offset, nil
}

// Read log record and return error if offset is invalid
func (l *Log) Read(offset uint64) (*Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if offset >= uint64(len(l.records)) {
		return nil, ErrOffsetNotFound
	}

	return &l.records[offset], nil
}