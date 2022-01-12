package logstore

import (
	"fmt"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/idutils"
)

type Log struct {
	MID          uint64
	Sender       uint64
	Receiver     uint64
	Me           uint64
	PtpTimestamp time.Time
	Event        EventType
}

func (l *Log) ToCsv() string {
	return fmt.Sprintf("%s,%s,%s,%d,%s,%s",
		idutils.String(l.Sender), idutils.String(l.Receiver), idutils.String(l.Me),
		l.PtpTimestamp.UnixNano(), idutils.String(l.MID), l.Event)
}

type EventType string

const (
	Send    EventType = "send"
	Receive EventType = "recv"
)
