package thingms

import (
	"time"
)

type ThingMsgHandler interface {
	Handle(task *Task) (time.Time, error)
}
