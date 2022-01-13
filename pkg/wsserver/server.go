package wsserver

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/registry"

	"github.com/gorilla/websocket"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"
	"ict.ac.cn/hbmsgserver/pkg/wshub"
)

type WebSocketServer struct {
	Addr      string
	MsgServer msgserver.MessageServer
	WsHub     *wshub.Hub

	Registry *registry.Registry

	Me uint32
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func (s *WebSocketServer) ServeWs(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("mac") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no mac addr in uri"))
		return
	}

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

	go func() {
		cli := &registry.Client{
			Mac:      r.URL.Query().Get("mac"),
			Location: "",
			WsClient: client,
		}
		cid := s.Registry.Register(cli)

		msg := msgserver.NewMessage(uint64(time.Now().UnixNano()), uint64(s.Me), idutils.DeviceId(s.Me, cid),
			msgserver.NameMsg, []byte(cli.Mac))
		client.Send <- msg

		fmt.Printf("register a client with id:%d, Mac:%s\n", cid, cli.Mac)
	}()

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
