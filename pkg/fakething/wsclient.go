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
	t.wsConn = c

	// Wait For ID
	_, msgRaw, err := c.ReadMessage()
	if err != nil {
		panic(err)
	}
	msg := msgserver.Message(msgRaw)

	if msg.Type() != msgserver.NameMsg {
		panic("don't receive the NameMsg")
	}

	return msg.Receiver()
}

func (t *Thing) waitNextMessage(typ msgserver.MessageType) msgserver.Message {
	_, msgRaw, err := t.wsConn.ReadMessage()
	if err != nil {
		panic(err)
	}
	msg := msgserver.Message(msgRaw)

	if msg.Type() != typ {
		panic(fmt.Sprintf("don't receive the %d", typ))
	}

	fmt.Printf("%s\n", msg.Body())
	return msg
}
