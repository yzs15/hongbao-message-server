package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/thingms"

	"ict.ac.cn/hbmsgserver/pkg/fakething"
)

var wsAddr = flag.String("ws", "172.17.0.1:5544", "http service address")
var zmqAddr = flag.String("zmq", "tcp://172.17.0.1:5543", "czmq service address")

var start = flag.String("start", "2022-01-26 00:00:00", "the start time(yyyy-MM-dd HH:mm:ss)")

var modeRaw = flag.String("mode", "cycle", "distribution mode")
var period = flag.Duration("period", time.Second, "the request period")
var connNum = flag.Int("conn", 300, "total request num")
var totalTime = flag.Duration("duration", 3*time.Second, "duration")
var peakTime = flag.Duration("peak-t", time.Second, "peak time")
var peakNum = flag.Int("peak-n", 300, "the request num at peak time")

//go:embed nums/0.png
//go:embed nums/1.png
//go:embed nums/2.png
//go:embed nums/3.png
var f embed.FS

func main() {
	flag.Parse()
	log.SetFlags(0)

	expected, err := time.Parse("2006-01-02 15:04:05", *start)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ThingMS WebSocket: ", *wsAddr)
	fmt.Println("ThingMS CZMQ: ", *zmqAddr)
	fmt.Println("Start At: ", expected)

	thing := &fakething.Thing{
		ID:           rand.Uint32(),
		ExpectedTime: expected,

		MsgWsEnd:  *wsAddr,
		MsgZmqEnd: *zmqAddr,

		Tasks: buildNumrecTasks(),

		Mode: fakething.Mode(*modeRaw),

		Period: *period,

		NumConn:   *connNum,
		TotalTime: *totalTime,
		PeakTime:  *peakTime,
		PeakNum:   *peakNum,
	}

	thing.Run()
}

func buildFibTask() *thingms.Task {
	return &thingms.Task{
		ID:        0,
		Sender:    rand.Uint32(),
		Good:      1,
		ServiceID: 2,
		SendTime:  0,
		Args:      []byte("3"),
	}
}

func buildNumrecTasks() []*thingms.Task {
	pngRaw, err := f.ReadFile("nums/0.png")
	if err != nil {
		panic(err)
	}

	pngRaw1, err := f.ReadFile("nums/1.png")
	if err != nil {
		panic(err)
	}

	pngRaw2, err := f.ReadFile("nums/2.png")
	if err != nil {
		panic(err)
	}

	pngRaw3, err := f.ReadFile("nums/3.png")
	if err != nil {
		panic(err)
	}

	task0 := &thingms.Task{
		ID:        0,
		Sender:    uint32(rand.Int31n(1 << 20)),
		Good:      1,
		ServiceID: 1,
		SendTime:  0,
		Args:      pngRaw,
	}

	task1 := task0.Clone()
	task1.Args = pngRaw1

	task2 := task0.Clone()
	task2.Args = pngRaw2

	task3 := task0.Clone()
	task3.Args = pngRaw3

	return []*thingms.Task{task0, task1, task2, task3}
}
