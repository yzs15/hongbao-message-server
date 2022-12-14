package fakething

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"gopkg.in/zeromq/goczmq.v4"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"
	"ict.ac.cn/hbmsgserver/pkg/thingms"
	"ict.ac.cn/hbmsgserver/pkg/timeutils"
)

const Full = ^uint8(0)

const reqWindow = 1 * time.Millisecond
const waitTime = 100 * time.Millisecond

// const maxTaskPerSec = 120
// const linearRatio = 4 //每秒增加的任务数目

type Mode string

var wangID = idutils.DeviceId(2, 2)

type Thing struct {
	SvrIdx int
	Me     uint64

	Config

	LoadTasks []*thingms.Task
	NoisTasks []*thingms.Task
	CongTasks []*thingms.Task

	mid      chan uint32
	midNoise chan uint32
}

func (c *Thing) Run() {
	go func() {
		c.mid = make(chan uint32, 10000)
		var id uint32
		for id = 1; ; id++ {
			c.mid <- id
		}
	}()

	go func() {
		c.midNoise = make(chan uint32, 10000)
		var id uint32
		for id = 1; ; id++ {
			c.midNoise <- id
		}
	}()

	fmt.Println("wait time:", waitTime.String())

	// 生成每个时间窗口要发送的消息数量
	var connDis []int
	var connSum int
	if disFunc, ok := disFuncMap[c.Mode]; ok {
		connDis, connSum = disFunc(c.Config)
	} else {
		log.Fatalf("distribution not found: %s", c.Mode)
	}
	if c.Mode != Cycle && float64(connSum) < float64(c.NumConn)*0.9 {
		log.Fatalf("分布设定不合理，请增大峰值，或增长时间，总连接数为%d", connSum)
	}

	c.connAndServe(connDis)
}

func (c *Thing) connAndServe(connDis []int) {
	// 连接 WebSocket
	var conn *websocket.Conn
	for {
		var err error
		conn, err = c.Connect()
		if err == nil {
			break
		}
	}

	for {
		_, msgRaw, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("recv msg failed: ", err.Error())
			break
		}
		msg := msgserver.Message(msgRaw)
		fmt.Println(msg.String())

		go func() {
			// ID Message
			if msg.Type() == msgserver.NameMsg {
				c.handleName(msg)
				return
			}

			if msg.Type() != msgserver.TextMsg {
				fmt.Println("recv wrong type message: ", msg.String())
				return
			}

			// Dredge Notice
			if msg.Body()[0] == Full-1 {
				c.handleNotice(msg)

			} else if msg.Body()[0] == Full-3 {
				c.handleTest(msg, connDis)

			} else {
				fmt.Println("receive wrong message")
			}
		}()
	}
}

func (c *Thing) handleName(msg msgserver.Message) {
	c.Me = msg.Receiver()
	fmt.Println("My ID: ", idutils.String(c.Me))

	t := c.LoadTasks[0].Clone()
	c.Request(t, c.Me, false)

	// t = c.CongTasks[0].Clone()
	// c.Request(t, c.Me, false)

	mid := <-c.midNoise
	task := c.NoisTasks[0].Clone()
	go c.RequestV2(task, c.Me, true, mid)
}

func (c *Thing) handleTest(msg msgserver.Message, connDis []int) {
	ran := rand.New(rand.NewSource(time.Now().Unix()))
	wangID = msg.Sender()

	nextTime := timeutils.GetSysTime(c.SvrIdx)
	var wg sync.WaitGroup
	wg.Add(len(connDis))
	for _, connNum := range connDis {
		go func(connNum int) {
			defer wg.Done()
			c.concurrentReq(connNum, ran)
		}(connNum)

		nextTime = nextTime.Add(reqWindow)
		timeutils.SleepUtil(nextTime, c.SvrIdx)
	}
	wg.Wait()
}

func (c *Thing) handleNotice(msg msgserver.Message) {
	var wg sync.WaitGroup
	wg.Add(5)

	sid := idutils.SvrId32(c.Me)
	cliPrefix := idutils.CliId32(c.Me) * 100

	timeutils.SleepUtil(msg.SendTime().Add(waitTime), c.SvrIdx)

	var i uint32
	for i = 0; i < 5; i++ {
		go func(idx uint32) {
			defer wg.Done()

			var task *thingms.Task
			if c.SvrIdx == timeutils.SpbEnv {
				task = c.CongTasks[rand.Intn(len(c.CongTasks))].Clone()
			} else {
				task = c.LoadTasks[rand.Intn(len(c.LoadTasks))].Clone()
			}

			c.Request(task, idutils.DeviceId(sid, cliPrefix+idx), false)
		}(i)
		time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()
}

// 一个时间窗口里的请求平均在每一毫秒发出
func (c *Thing) concurrentReq(num int, ran *rand.Rand) {
	msNum := int(reqWindow.Milliseconds() / time.Millisecond.Milliseconds())
	numPerMs := num / msNum
	numRemain := num % msNum

	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < msNum; i++ {
		curNum := numPerMs
		if numRemain > 0 {
			curNum++
			numRemain--
		}

		for ri := 0; ri < curNum; ri++ {
			go func() {
				defer wg.Done()

				for li := 0; li < c.LoadNumPer; li++ {
					mid := <-c.mid
					task := c.LoadTasks[ran.Intn(len(c.LoadTasks))].Clone()
					go c.RequestV2(task, c.Me, false, mid)
				}
				for ni := 0; ni < c.NoisNumPer; ni++ {
					mid := <-c.midNoise
					task := c.NoisTasks[ran.Intn(len(c.NoisTasks))].Clone()
					go c.RequestV2(task, c.Me, true, mid)
				}
			}()
		}
		time.Sleep(time.Millisecond)
	}
	wg.Wait()
}

func (c *Thing) Request(task *thingms.Task, sender uint64, isNoise bool) {
	mid := <-c.mid
	if isNoise {
		mid = (1 << 19) | mid
	}

	msg := msgserver.NewMessage(idutils.MessageID(idutils.SvrId32(sender), idutils.CliId32(sender), mid),
		sender, wangID, msgserver.TaskMsg, task.ToBytes())

	sockItem, err := czmqutils.GetSock(c.MsgZmqEnd[c.SvrIdx], goczmq.Push)
	if err != nil {
		log.Println("czmq get sock failed: ", err)
		return
	}
	defer sockItem.Free()

	msg.SetSendTime()
	if _, err := czmqutils.Send(sockItem, msg, goczmq.FlagNone); err != nil {
		log.Println("czmq send failed: ", err)
	}
	// log.Printf("[%s] send a message:%s, size: %d\n", timeutils.Time2string(msg.SendTime()), idutils.String(msg.ID()), len(msg))
}

func (c *Thing) RequestV2(task *thingms.Task, sender uint64, isNoise bool, mid uint32) {
	if isNoise {
		mid = (1 << 19) | mid
	}

	msg := msgserver.NewMessage(idutils.MessageID(idutils.SvrId32(sender), idutils.CliId32(sender), mid),
		sender, wangID, msgserver.TaskMsg, task.ToBytes())

	sockItem, err := czmqutils.GetSock(c.MsgZmqEnd[c.SvrIdx], goczmq.Push)
	if err != nil {
		log.Println("czmq get sock failed: ", err)
		return
	}
	defer sockItem.Free()

	msg.SetSendTime()
	if _, err := czmqutils.Send(sockItem, msg, goczmq.FlagNone); err != nil {
		log.Println("czmq send failed: ", err)
	}
	// log.Printf("[%s] send a message:%s, size: %d\n", timeutils.Time2string(msg.SendTime()), idutils.String(msg.ID()), len(msg))
}
