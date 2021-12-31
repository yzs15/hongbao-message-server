package main

import (
	"flag"

	"ict.ac.cn/hbmsgserver/pkg/thingms"

	"ict.ac.cn/hbmsgserver/pkg/logstore"

	"ict.ac.cn/hbmsgserver/pkg/czmqserver"

	"ict.ac.cn/hbmsgserver/pkg/wshub"
	"ict.ac.cn/hbmsgserver/pkg/wsserver"
)

var wsAddr = flag.String("ws", "0.0.0.0:5554", "web socket service address")
var zmqAddr = flag.String("zmq", "tcp://0.0.0.0:5553", "zmq service address")
var logAddr = flag.String("log", "0.0.0.0:5552", "log service address")

var wangEnd = flag.String("wang", "tcp://127.0.0.1:5553", "wang endpoint")

func main() {
	flag.Parse()

	logStore := logstore.NewLogStore()
	go logStore.Run()

	wsHub := wshub.NewHub()
	go wsHub.Run()

	kubeEnds := []string{"172.16.32.12:32101", "172.16.32.13:32101", "172.16.32.14:32101", "172.16.32.15:32101"}
	svs := buildSvs()
	thingMsgHdl := thingms.NewNetThingMsgHandler(kubeEnds, *wangEnd, svs)

	msgServer := &thingms.ThingMS{
		LogStore:    logStore,
		WsHub:       wsHub,
		ThingMsgHdl: thingMsgHdl,
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

func buildSvs() map[string]*thingms.NetService {
	svs := make(map[string]*thingms.NetService)

	fibSvs := &thingms.NetService{
		Method: "GET",
		Query:  "num",
		File:   "",
	}
	svs["fib"] = fibSvs

	return svs
}
