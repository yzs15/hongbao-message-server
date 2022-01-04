package fakething

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/thingms"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"
)

type Thing struct {
	ID           uint32
	ExpectedTime time.Time

	MsgWsEnd  string
	MsgZmqEnd string
}

func (t *Thing) Run() {
	receiveTime := t.waitFirst()

	good := receiveTime.Before(t.ExpectedTime)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	msg := t.makeMsg(good)
	for {
		select {
		case <-ticker.C:
			sendTime := time.Now()
			msg.SendTime = uint64(sendTime.UnixNano())

			if err := czmqutils.Send(t.MsgZmqEnd, msg.ToBytes()); err != nil {
				log.Println("czmq send failed: ", err)
			}

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			select {
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (t *Thing) makeMsg(good bool) *thingms.Task {
	msg := &thingms.Task{
		ID:        uint64(t.ID)<<32 | uint64(rand.Uint32()),
		Sender:    t.ID,
		ServiceID: 1,
		Args:      nil,
	}
	if good {
		msg.Good = 1
	} else {
		msg.Good = 0
	}
	return msg
}
