package main

import (
	"flag"

	"ict.ac.cn/hbmsgserver/pkg/logstore"

	"ict.ac.cn/hbmsgserver/pkg/czmqserver"

	"ict.ac.cn/hbmsgserver/pkg/wangms"

	"ict.ac.cn/hbmsgserver/pkg/wshub"
	"ict.ac.cn/hbmsgserver/pkg/wsserver"
)

var wsAddr = flag.String("ws", "0.0.0.0:5554", "web socket service address")
var zmqAddr = flag.String("zmq", "tcp://0.0.0.0:5553", "zmq service address")
var logAddr = flag.String("log", "0.0.0.0:5552", "log service address")

func main() {
	flag.Parse()

	logStore := logstore.NewLogStore()
	go logStore.Run()

	wsHub := wshub.NewHub()
	go wsHub.Run()

	msgServer := &wangms.WangMS{
		LogStore:    logStore,
		WsHub:       wsHub,
		MsEndpoints: []string{""},
	}

	wsServer := &wsserver.WebSocketServer{
		Addr:      *wsAddr,
		MsgServer: msgServer,
		WsHub:     wsHub,
	}
	go wsServer.Run()

	zmqServer := &czmqserver.CZMQServer{
		Addr:      *zmqAddr,
		MsgServer: msgServer,
	}
	go zmqServer.Run()

	logServer := &logstore.LogServer{
		Addr:     *logAddr,
		LogStore: logStore,
	}
	logServer.Run()
}
