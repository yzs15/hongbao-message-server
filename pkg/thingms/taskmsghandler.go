package thingms

import (
	"time"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"
)

type TaskMsgHandler interface {
	Handle(msg msgserver.Message) (time.Time, error)
}
