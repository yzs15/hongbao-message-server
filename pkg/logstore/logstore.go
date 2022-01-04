package logstore

import (
	"time"
)

type LogStore struct {
	logs []*Log

	add chan *Log

	me string
}

func NewLogStore(me string) *LogStore {
	return &LogStore{
		logs: make([]*Log, 0),
		add:  make(chan *Log, 1024),
		me:   me,
	}
}

func (s *LogStore) Add(mid uint64, timestamp time.Time, event EventType) {
	log := &Log{
		MID:       mid,
		Timestamp: timestamp,
		Event:     event,
		Me:        s.me,
	}
	s.add <- log
}

func (s *LogStore) Run() {
	for log := range s.add {
		s.logs = append(s.logs, log)
	}
}
