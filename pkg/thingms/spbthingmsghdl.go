package thingms

import (
	"time"
)

type spbThingMsgHandler struct {
}

func NewSpbThingMsgHandler() ThingMsgHandler {
	return &spbThingMsgHandler{}
}

func (h *spbThingMsgHandler) Handle(task *Task) (time.Time, error) {
	return time.Time{}, nil
}
