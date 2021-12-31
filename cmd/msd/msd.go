package main

import (
	"flag"
	"strings"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"
	"ict.ac.cn/hbmsgserver/pkg/wangms"

	"ict.ac.cn/hbmsgserver/pkg/thingms"

	"ict.ac.cn/hbmsgserver/pkg/logstore"

	"ict.ac.cn/hbmsgserver/pkg/czmqserver"

	"ict.ac.cn/hbmsgserver/pkg/wshub"
	"ict.ac.cn/hbmsgserver/pkg/wsserver"
)

type StringArrFlag []string

func (i *StringArrFlag) String() string {
	return strings.Join(*i, ",")
}

func (i *StringArrFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var wsAddr = flag.String("ws", "0.0.0.0:5554", "web socket service address")
var zmqAddr = flag.String("zmq", "tcp://0.0.0.0:5553", "zmq service address")
var logAddr = flag.String("log", "0.0.0.0:5552", "log service address")

var isThing = flag.Bool("thing", false, "run as thing rule")
var isWang = flag.Bool("wang", false, "run as wang rule")

var wangEnd = flag.String("wend", "tcp://127.0.0.1:5553", "wang endpoint")
var thingEnds StringArrFlag
var kubeEnds StringArrFlag

func main() {
	flag.Var(&thingEnds, "tend", "things endpoint")
	flag.Var(&kubeEnds, "kend", "things endpoint")
	flag.Parse()

	logStore := logstore.NewLogStore()
	go logStore.Run()

	wsHub := wshub.NewHub()
	go wsHub.Run()

	var msgServer msgserver.MessageServer
	if *isThing {
		svs := buildSvs()
		thingMsgHdl := thingms.NewNetThingMsgHandler(kubeEnds, *wangEnd, svs)

		msgServer = &thingms.ThingMS{
			LogStore:    logStore,
			WsHub:       wsHub,
			ThingMsgHdl: thingMsgHdl,
		}

	} else if *isWang {
		msgServer = &wangms.WangMS{
			LogStore:    logStore,
			WsHub:       wsHub,
			MsEndpoints: []string{""},
		}

	} else {
		panic("need to specify one rule by '--thing' or '--wang'")
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
