package wsserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"
	"ict.ac.cn/hbmsgserver/pkg/wshub"
)

type WebSocketServer struct {
	Addr      string
	MsgServer msgserver.MessageServer
	WsHub     *wshub.Hub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func (s *WebSocketServer) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &wshub.Client{
		Hub: s.WsHub, Conn: conn, Send: make(chan []byte, 256),
		MsgServer: s.MsgServer,
	}
	s.WsHub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

func (s *WebSocketServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.ServeWs)

	fmt.Printf("web socket server listen at: %s\n", s.Addr)
	if err := http.ListenAndServe(s.Addr, mux); err != nil {
		panic(err)
	}
}
