package msgserver

import "time"

type MessageServer interface {
	Handle(receiveTime time.Time, msgRaw Message)
}
