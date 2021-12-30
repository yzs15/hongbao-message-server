package logstore

import "time"

type Log struct {
	MID       uint64
	Timestamp time.Time
	Event     EventType
}

type EventType string

const (
	SenderSended           EventType = "SenderSended"
	SenderMsgSvrReceived   EventType = "SenderMsgSvrReceived"
	SenderMsgSvrSended     EventType = "SenderMsgSvrSended"
	ReceiverMsgSvrReceived EventType = "ReceiverMsgSvrReceived"
	ReceiverMsgSvrSended   EventType = "ReceiverMsgSvrSended"
	ReceiverReceived       EventType = "ReceiverReceived"
)
