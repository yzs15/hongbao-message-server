package msgserver

import "time"

type Message struct {
	ID      uint64
	Sender  uint32
	Good    bool
	Content string
}

type MessageServer interface {
	Handle(receiveTime time.Time, msgRaw []byte)
}
