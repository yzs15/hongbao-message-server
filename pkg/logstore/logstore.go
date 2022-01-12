package logstore

import (
	"fmt"
	"os"
	"time"
)

type LogStore struct {
	logs []*Log

	add chan *Log

	me uint32

	logPath string
}

func NewLogStore(me uint32, logPath string) *LogStore {
	return &LogStore{
		logs:    make([]*Log, 0),
		add:     make(chan *Log, 1024),
		me:      me,
		logPath: logPath,
	}
}

func (s *LogStore) Add(mid uint64, sender uint64, receiver uint64, logger uint64, ptpTime time.Time, event EventType) {
	log := &Log{mid, sender, receiver, logger,
		ptpTime, event}
	s.add <- log
}

func (s *LogStore) Run() {
	filename := fmt.Sprintf("%s/%d.log", s.logPath, s.me)
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	for log := range s.add {
		s.logs = append(s.logs, log)
		csvStr := log.ToCsv()
		logFile.WriteString(csvStr)
		logFile.WriteString("\n")
	}
}
