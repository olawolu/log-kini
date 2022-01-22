package server

import (
	"fmt"
	"sync"
)

var ErrOffsetNotFound = fmt.Errorf("offset not found")

type Log struct {
	mu      sync.Mutex
	records []Record
}

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

func NewLog() *Log {
	return &Log{}
}

func (lg *Log) Append(record Record) (uint64, error) {
	lg.mu.Lock()
	defer lg.mu.Unlock()
	record.Offset = uint64(len(lg.records))
	lg.records = append(lg.records, record)
	return record.Offset, nil
}

func (lg *Log) Read(offset uint64) (Record, error) {
	lg.mu.Lock()
	defer lg.mu.Unlock()
	if offset >= uint64(len(lg.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return lg.records[offset], nil
}
