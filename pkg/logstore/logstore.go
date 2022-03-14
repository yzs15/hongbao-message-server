package logstore

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/nameserver"

	"ict.ac.cn/hbmsgserver/pkg/idutils"
)

type LogStore struct {
	logs []*Log
	add  chan *Log

	ns *nameserver.NameServer

	logPath string
}

func NewLogStore(logPath string, ns *nameserver.NameServer) *LogStore {
	return &LogStore{
		logs:    make([]*Log, 0),
		add:     make(chan *Log, 1024),
		ns:      ns,
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

	for l := range s.add {
		s.logs = append(s.logs, l)

		csvStr := fmt.Sprintf("%s,%s,%s,%d,%s,%s",
			s.id2string(l.Sender), s.id2string(l.Receiver), s.id2string(l.Me),
			l.PtpTimestamp.UnixNano(), s.id2string(l.MID), l.Event)

		logFile.WriteString(csvStr)
		logFile.WriteString("\n")
	}
}

func (s *LogStore) id2string(id uint64) string {
	if idutils.CliId32(id) != 0 || idutils.MsgId32(id) != 0 {
		return idutils.String(id)
	}

	if id >= 100 {
		return strconv.FormatInt(int64(id), 10)
	}

	svr, err := s.ns.GetServer(idutils.SvrId32(id))
	if err != nil {
		fmt.Println(err)
		return idutils.String(id)
	}

	return svr.ZMQEndpoint
}
