package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
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

var isThing = flag.Bool("thing", false, "run as one Thing")
var isWang = flag.Bool("wang", false, "run as Wang")

var isNet = flag.Bool("net", false, "in internet environment")
var isSpb = flag.Bool("spb", false, "in superbahn environment")

var wangEnd = flag.String("wend", "tcp://127.0.0.1:5553", "Wang's czmq endpoint")
var thingEnds StringArrFlag
var kubeEnds StringArrFlag

var spbConfig = flag.String("spbcfg", "./bin/infospb/client.config", "config file for info superbahn")

func main() {
	flag.Var(&thingEnds, "tend", "thing's czmq endpoint")
	flag.Var(&kubeEnds, "kend", "kubernetes' http endpoint")
	flag.Parse()

	var me string

	if *isNet {
		fmt.Println("run in internet")
	} else if *isSpb {
		fmt.Println("run in superbahn")
	}
	if *isWang {
		fmt.Printf("this is Wang\n")
		fmt.Printf("thing czmq endpoints: %v\n", thingEnds)
		me = "WangMS"
	} else if *isThing {
		fmt.Printf("this is one Thing\n")
		fmt.Printf("wang czmq endpoint: %v\n", *wangEnd)
		fmt.Printf("kube http endpoints: %v\n", kubeEnds)
		me = "ThingMS"
	}

	logStore := logstore.NewLogStore(me)
	go logStore.Run()

	wsHub := wshub.NewHub()
	go wsHub.Run()

	var msgServer msgserver.MessageServer
	if *isThing {
		var thingMsgHdl thingms.ThingMsgHandler
		if *isNet {
			svs := buildNetSvs()
			thingMsgHdl = thingms.NewNetThingMsgHandler(kubeEnds, *wangEnd, svs)
		} else if *isSpb {
			thingMsgHdl = thingms.NewSpbThingMsgHandler(*spbConfig)
		} else {
			panic("need to specify environment by '--net' or '--spb'")
		}

		msgServer = &thingms.ThingMS{
			LogStore:    logStore,
			WsHub:       wsHub,
			ThingMsgHdl: thingMsgHdl,
		}

	} else if *isWang {
		msgServer = &wangms.WangMS{
			LogStore:    logStore,
			WsHub:       wsHub,
			ThingMsEnds: thingEnds,
		}

	} else {
		panic("need to specify one rule by '-thing' or '-wang'")
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

func buildNetSvs() map[uint8]*thingms.NetService {
	svs := make(map[uint8]*thingms.NetService)

	fibSvs := &thingms.NetService{
		Method: "GET",
		Port:   "32101",
		Path: func(args []byte) string {
			if len(args) > 0 {
				return fmt.Sprintf("/?num=%s", args)
			}
			return "/"
		},
		Body: thingms.NilBody,
	}

	numrecdSvs := &thingms.NetService{
		Method: "POST",
		Port:   "32100",
		Path:   thingms.NilPath,
		Body: func(args []byte) (io.Reader, string) {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)

			fw, err := w.CreateFormFile("img", "num")
			if err != nil {
				fmt.Println(err)
				return nil, ""
			}

			_, err = fw.Write(args)
			if err != nil {
				fmt.Println(err)
				return nil, ""
			}
			w.Close()

			return &b, w.FormDataContentType()
		},
	}

	svs[1] = numrecdSvs
	svs[2] = fibSvs

	return svs
}
