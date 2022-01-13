package main

import (
	"bytes"
	"encoding/binary"
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

var isNet = flag.Bool("net", false, "in internet environment")
var isSpb = flag.Bool("spb", false, "in superbahn environment")
var spbConfig = flag.String("spbcfg", "./bin/infospb/client.config", "config file for info superbahn")

var kubeEnds StringArrFlag

var nsEnd = flag.String("nsend", "", "name server endpoint")
var zmqOutEnd = flag.String("zmq-out", "", "zmq remote connection endpoint")

var logPath = flag.String("log-path", ".", "the directory to save log file")

func main() {
	flag.Var(&kubeEnds, "kend", "kubernetes' http endpoint")
	flag.Parse()

	if *zmqOutEnd == "" || *nsEnd == "" {
		panic("need '-zmq-out' and '-nsend'")
	}
	me := nameserver.Register(*nsEnd, *zmqOutEnd)
	fmt.Printf("my id is %d\n", me)

	ns := nameserver.NewNameServer(*nsEnd, me)

	logStore := logstore.NewLogStore(fmt.Sprintf("%s/%d.log", *logPath, me))
	go logStore.Run()

	wsHub := wshub.NewHub()
	go wsHub.Run()

	reg := registry.NewRegistry(wsHub)
	go reg.Run()

	var taskMsgHdl thingms.TaskMsgHandler
	if *isNet {
		fmt.Println("run in internet")
		svs := buildNetSvs()
		taskMsgHdl = thingms.NewNetThingMsgHandler(me, kubeEnds, svs, reg, ns, logStore)
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
	zmqServer.Run()
}

func buildNetSvs() map[uint8]*thingms.NetService {
	svs := make(map[uint8]*thingms.NetService)

	fibSvs := &thingms.NetService{
		Method: "GET",
		Port:   "32101",
		Path: func(mid uint64, args []byte) string {
			prefix := thingms.NilPath(mid, args)
			if len(args) > 0 {
				return fmt.Sprintf("%s&num=%s", prefix, args)
			}
			return prefix
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

	hongbaoSvs := &thingms.NetService{
		Port:   "32107",
		Method: "POST",
		Path: func(mid uint64, args []byte) string {
			prefix := thingms.NilPath(mid, args)
			l := binary.LittleEndian.Uint32(args[:4])
			if len(args) > 0 {
				return fmt.Sprintf("%s&msg=%s", prefix, args[4:4+l])
			}
			return prefix
		},
		Body: func(args []byte) (io.Reader, string) {
			l := binary.LittleEndian.Uint32(args[:4])

			var b bytes.Buffer
			w := multipart.NewWriter(&b)

			fw, err := w.CreateFormFile("img", "num")
			if err != nil {
				fmt.Println(err)
				return nil, ""
			}

			_, err = fw.Write(args[4+l:])
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
	svs[7] = hongbaoSvs

	return svs
}
