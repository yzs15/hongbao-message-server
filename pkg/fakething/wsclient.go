package fakething

import (
	"fmt"
	"net/url"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	"github.com/gorilla/websocket"
)

func (t *Thing) waitID() uint64 {
	u := url.URL{Scheme: "ws", Host: t.MsgWsEnd, Path: "/"}
	q := u.Query()
	q.Set("mac", t.MacAddr)
	u.RawQuery = q.Encode()

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// Wait For ID
	_, msgRaw, err := c.ReadMessage()
	if err != nil {
		panic(err)
	}
	msg := msgserver.Message(msgRaw)

	if msg.Type() != msgserver.NameMsg {
		panic("don't receive the NameMsg")
	}

	// Wait For Start
	_, msgRaw, err = c.ReadMessage()
	if err != nil {
		panic(err)
	}
	msg = msgRaw

	if msg.Type() != msgserver.TextMsg {
		panic("don't receive the NameMsg")
	}

	fmt.Printf("%s\n", msg.Body())

	//receiveTime := time.Now()
	//go func() {
	//	log.Printf("[%s] recv: %s", timeutils.Time2string(receiveTime), string(msgRaw))
	//
	//	myID := msg.Receiver()
	//	svrID := idutils.SvrId32(myID)
	//	sendMsg := msgserver.NewMessage(msg.ID(), myID, uint64(svrID),
	//		msgserver.LogMsg, nil)
	//	if _, err := czmqutils.Send(t.MsgZmqEnd, sendMsg); err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	return msg.Receiver()
}
