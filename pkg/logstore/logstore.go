package logstore

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type LogStore struct {
	logs []*Log

	add chan *Log

	me uint32
}

func NewLogStore(me uint32) *LogStore {
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
		Me:        uint64(s.me),
	}
	s.add <- log
}

func (s *LogStore) Run() {
	for log := range s.add {
		s.logs = append(s.logs, log)
		txt, err := json.Marshal(s.logs)
		if err != nil {
			continue
		}
		ioutil.WriteFile("log", txt, 0644)
	}
}
