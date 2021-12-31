package msgserver

import "time"

type Message struct {
	ID      uint64 `json:"ID,omitempty"`
	Sender  uint32 `json:"Sender,omitempty"`
	Good    bool   `json:"Good,omitempty"`
	Content string `json:"Content,omitempty"`
}

type MessageServer interface {
	Handle(receiveTime time.Time, msgRaw []byte)
}
