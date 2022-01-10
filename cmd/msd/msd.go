package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"ict.ac.cn/hbmsgserver/pkg/nameserver"

	"ict.ac.cn/hbmsgserver/pkg/registry"

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

var isNet = flag.Bool("net", false, "in internet environment")
var isSpb = flag.Bool("spb", false, "in superbahn environment")
var spbConfig = flag.String("spbcfg", "./bin/infospb/client.config", "config file for info superbahn")

var kubeEnds StringArrFlag

var nsEnd = flag.String("nsend", "", "name server endpoint")
var zmqOutEnd = flag.String("zmq-out", "", "zmq remote connection endpoint")

func main() {
	flag.Var(&kubeEnds, "kend", "kubernetes' http endpoint")
	flag.Parse()

	if *zmqOutEnd == "" || *nsEnd == "" {
		panic("need '-zmq-out' and '-nsend'")
	}
	me := nameserver.Register(*nsEnd, *zmqOutEnd)
	fmt.Printf("my id is %d\n", me)

	ns := nameserver.NewNameServer(*nsEnd, me)

	logStore := logstore.NewLogStore(me)
	go logStore.Run()

	wsHub := wshub.NewHub()
	go wsHub.Run()

	reg := registry.NewRegistry()
	go reg.Run()

	var taskMsgHdl thingms.TaskMsgHandler
	if *isNet {
		fmt.Println("run in internet")
		svs := buildNetSvs()
		taskMsgHdl = thingms.NewNetThingMsgHandler(kubeEnds, svs, ns)
	} else if *isSpb {
		fmt.Println("run in superbahn")
		taskMsgHdl = thingms.NewSpbThingMsgHandler(*spbConfig)
	} else {
		panic("need to specify environment by '--net' or '--spb'")
	}
	msgServer := &thingms.ThingMS{
		LogStore:   logStore,
		Registry:   reg,
		TaskMsgHdl: taskMsgHdl,
		NameServer: ns,
		Me:         me,
	}

	wsServer := &wsserver.WebSocketServer{
		Addr:      *wsAddr,
		MsgServer: msgServer,
		WsHub:     wsHub,
		Me:        me,
		Registry:  reg,
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
