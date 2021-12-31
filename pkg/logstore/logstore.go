package logstore

import (
	"sync"
	"time"
)

type LogStore struct {
	logs []*Log

	mu sync.Mutex

	add chan *Log
}

func NewLogStore() *LogStore {
	return &LogStore{
		logs: make([]*Log, 0),
		add:  make(chan *Log, 1024),
	}
}

func (s *LogStore) Add(mid uint64, timestamp time.Time, event EventType) {
	log := &Log{
		MID:       mid,
		Timestamp: timestamp,
		Event:     event,
	}
	s.add <- log
}

func (s *LogStore) Run() {
	for log := range s.add {
		s.logs = append(s.logs, log)
	}
}
