package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/fakething"
)

var wsAddr = flag.String("ws", "172.17.0.1:5544", "http service address")
var zmqAddr = flag.String("zmq", "tcp://172.17.0.1:5543", "czmq service address")
var start = flag.String("start", "", "the start time")

func main() {
	flag.Parse()
	log.SetFlags(0)

	fmt.Println("ThingMS WebSocket: ", *wsAddr)
	fmt.Println("ThingMS CZMQ: ", *zmqAddr)
	fmt.Println("Start At: ", *start)

	thing := &fakething.Thing{
		ID:           rand.Uint32(),
		ExpectedTime: time.Now(),
		MsgWsEnd:     *wsAddr,
		MsgZmqEnd:    *zmqAddr,
	}

	thing.Run()
}
