package logstore

import (
	"os"
	"time"
)

type LogStore struct {
	logs []*Log

	add chan *Log

	logPath string
}

func NewLogStore(logPath string) *LogStore {
	return &LogStore{
		logs:    make([]*Log, 0),
		add:     make(chan *Log, 1024),
		logPath: logPath,
	}
}

func (s *LogStore) Add(mid uint64, sender uint64, receiver uint64, logger uint64, ptpTime time.Time, event EventType) {
	log := &Log{mid, sender, receiver, logger,
		ptpTime, event}
	s.add <- log
}

func (s *LogStore) Run() {
	logFile, err := os.OpenFile(s.logPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	for log := range s.add {
		s.logs = append(s.logs, log)
		csvStr := log.ToCsv()
		logFile.WriteString(csvStr)
		logFile.WriteString("\n")
	}
}
