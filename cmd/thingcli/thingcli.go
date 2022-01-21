package main

import (
	"bytes"
	"embed"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"sync"

	"ict.ac.cn/hbmsgserver/pkg/fakething"
	"ict.ac.cn/hbmsgserver/pkg/linuxutils"
	"ict.ac.cn/hbmsgserver/pkg/thingms"
	"ict.ac.cn/hbmsgserver/pkg/timeutils"
)

const (
	NetEnv = 0
	SpbEnv = 1
)

var macAddr = flag.String("mac", "", "the fake mac addr for this fake thing")

var wsAddrNet = flag.String("ws-net", "", "http service address")
var zmqAddrNet = flag.String("zmq-net", "", "czmq service address")

var wsAddrSpb = flag.String("wsb-spb", "", "http service address")
var zmqAddrSpb = flag.String("zmqb-spb", "", "czmq service address")

var modeRaw = flag.String("mode", "cycle", "distribution mode")
var period timeutils.Duration
var connNum = flag.Int("conn", 300, "total request num")
var totalTime timeutils.Duration
var peakTime timeutils.Duration
var peakNum = flag.Int("peak-n", 300, "the request num at peak time")

var configFile = flag.String("config", "", "the path to config file, if set, omit other parameters")

//go:embed nums/0.png
//go:embed nums/1.png
//go:embed nums/2.png
//go:embed nums/3.png
//go:embed nums/4.png
//go:embed nums/5.png
//go:embed nums/6.png
//go:embed nums/7.png
//go:embed nums/8.png
//go:embed nums/9.png
var f embed.FS

func main() {
	flag.Var(&period, "period", "the request period")
	flag.Var(&totalTime, "duration", "duration")
	flag.Var(&peakTime, "peak-t", "peak time")
	flag.Parse()
	log.SetFlags(0)

	var conf fakething.Config
	if len(*configFile) != 0 {
		var err error
		conf, err = fakething.GetConfig(*configFile)
		if err != nil {
			panic(err)
		}

	} else {
		conf = fakething.Config{
			MacAddr: *macAddr,

			MsgWsEnd:  []string{*wsAddrNet, *wsAddrSpb},
			MsgZmqEnd: []string{*zmqAddrNet, *zmqAddrSpb},

			Mode: fakething.Mode(*modeRaw),

			Period: period,

			NumConn:   *connNum,
			TotalTime: totalTime,
			PeakTime:  peakTime,
			PeakNum:   *peakNum,
		}
	}

	if len(conf.MacAddr) == 0 {
		var err error
		conf.MacAddr, err = linuxutils.GetMac()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("My Mac Addr: ", conf.MacAddr)
	fmt.Println("Msg Server WebSocket: ", conf.MsgWsEnd)
	fmt.Println("Msg Server CZMQ: ", conf.MsgZmqEnd)

	var wg sync.WaitGroup
	if conf.MsgWsEnd[NetEnv] != "" {
		wg.Add(1)
		thingNet := &fakething.Thing{
			SvrIdx:    NetEnv,
			Config:    conf,
			LoadTasks: buildNumrecTasks(),
			CongTasks: buildHongbaoTasks(),
		}
		go func() {
			defer wg.Done()
			thingNet.Run()
		}()
	}

	if conf.MsgWsEnd[SpbEnv] != "" {
		wg.Add(1)
		thingSpb := &fakething.Thing{
			SvrIdx:    SpbEnv,
			Config:    conf,
			LoadTasks: buildNumrecTasks(),
			CongTasks: buildHongbaoTasks(),
		}
		go func() {
			defer wg.Done()
			thingSpb.Run()
		}()
	}

	wg.Wait()
}

func buildFibTask() *thingms.Task {
	return &thingms.Task{
		ServiceID: 2,
		Args:      []byte("3"),
	}
}

func buildNumrecTasks() []*thingms.Task {
	tasks := make([]*thingms.Task, 0)

	for i := 0; i < 10; i++ {
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

func buildHongbaoTasks() []*thingms.Task {
	tasks := make([]*thingms.Task, 0)

	content := "恭喜信息高铁开通了！"
	l := uint32(len(content))
	lenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBytes, l)

	for i := 0; i < 10; i++ {
		pngRaw, err := f.ReadFile(fmt.Sprintf("nums/%d.png", i))
		if err != nil {
			panic(err)
		}

		args := new(bytes.Buffer)
		args.Write(lenBytes)
		args.WriteString(content)
		args.Write(pngRaw)

		task := &thingms.Task{
			ServiceID: 7,
			Args:      args.Bytes(),
		}
		tasks = append(tasks, task)
	}
	return tasks
}
