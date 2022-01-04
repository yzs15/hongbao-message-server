package fakething

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/thingms"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	"github.com/gorilla/websocket"
)

func (t *Thing) waitFirst() time.Time {
	u := url.URL{Scheme: "ws", Host: t.MsgWsEnd, Path: "/"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	_, msgRaw, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return time.Time{}
	}

	receiveTime := time.Now()
	go func() {
		log.Printf("recv: %s", string(msgRaw))
		var msg msgserver.Message
		if err := json.Unmarshal(msgRaw, &msg); err != nil {
			log.Fatal("json decode message")
		}

		sendMsg := &thingms.Task{
			ID:        msg.ID,
			Sender:    t.ID,
			Good:      1,
			ServiceID: 0,
			SendTime:  uint64(receiveTime.UnixNano()),
			Args:      nil,
		}
		if err := czmqutils.Send(t.MsgZmqEnd, sendMsg.ToBytes()); err != nil {
			log.Fatal(err)
		}
	}()

	return receiveTime
}
