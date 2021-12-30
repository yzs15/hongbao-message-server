package thingms

import "time"

type ThingMsgHandler interface {
	Handle(msg msgserver.Message) time.Time
}
