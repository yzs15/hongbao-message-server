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

var macAddr = flag.String("mac", "02:42:ac:12:00:00", "the fake mac addr for this fake thing")

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
		MacAddr:      *macAddr,
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
		ServiceID: 2,
		Args:      []byte("3"),
	}
}

func buildNumrecTasks() []*thingms.Task {
	tasks := make([]*thingms.Task, 0)

	for i := 0; i < 4; i++ {
		pngRaw, err := f.ReadFile(fmt.Sprintf("nums/%d.png", i))
		if err != nil {
			panic(err)
		}

		task := &thingms.Task{
			ServiceID: 1,
			Args:      pngRaw,
		}
		tasks = append(tasks, task)
	}
	return tasks
}
