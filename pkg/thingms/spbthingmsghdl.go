package thingms

import (
	"time"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"
)

type spbThingMsgHandler struct {
}

func (h *spbThingMsgHandler) Handle(msg *msgserver.Message) (time.Time, error) {
	return time.Time{}, nil
}
