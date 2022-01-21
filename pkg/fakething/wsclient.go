package fakething

import (
	"net/url"

	"github.com/gorilla/websocket"
)

func (t *Thing) Connect() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: t.MsgWsEnd[t.SvrIdx], Path: "/"}
	q := u.Query()
	q.Set("mac", t.MacAddr)
	u.RawQuery = q.Encode()

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}
