package msgserver

import "time"

type Message struct {
	ID      uint64 `json:"ID"`
	Sender  uint32 `json:"Sender"`
	Good    uint8  `json:"Good,omitempty"`
	Content string `json:"Content,omitempty"`
}

type MessageServer interface {
	Handle(receiveTime time.Time, msgRaw []byte)
}
