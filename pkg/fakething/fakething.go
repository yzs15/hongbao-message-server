package fakething

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"gopkg.in/zeromq/goczmq.v4"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"
	"ict.ac.cn/hbmsgserver/pkg/thingms"
	"ict.ac.cn/hbmsgserver/pkg/timeutils"
)

const Full = ^uint8(0)

const reqWindow = 2 * time.Millisecond
const waitTime = 50 * time.Millisecond

type Mode string

var wangID = idutils.DeviceId(2, 2)

type Thing struct {
	SvrIdx int
	Me     uint64

	Config

	LoadTasks []*thingms.Task
	CongTasks []*thingms.Task

	mid chan uint32
}

func (c *Thing) Run() {
	go func() {
		c.mid = make(chan uint32, 10000)
		var id uint32
		for id = 1; ; id++ {
			c.mid <- id
		}
	}()

	// 生成每个时间窗口要发送的消息数量
	var connDis []int
	var connSum int
	switch c.Mode {
	case Cycle:
		connDis, connSum = c.cycleConnDis()
	case Normal:
		connDis, connSum = c.normalConnDis()
	case Uniform:
		connDis, connSum = c.uniformConnDis()
	}
	if c.Mode != Cycle && float64(connSum) < float64(c.NumConn)*0.9 {
		log.Fatalf("分布设定不合理，请增大峰值，或增长时间，总连接数为%d", connSum)
	}

	c.connAndServe(connDis)
}

func (c *Thing) connAndServe(connDis []int) {
	// 连接 WebScoket
	conn, err := c.Connect()
	if err != nil {
		panic(err)
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
	c.Request(t)

	t = c.CongTasks[0].Clone()
	c.Request(t)
}

func (c *Thing) handleTest(msg msgserver.Message, connDis []int) {
	ran := rand.New(rand.NewSource(time.Now().Unix()))
	wangID = msg.Sender()

	nextTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(connDis))
	for _, connNum := range connDis {
		go func(connNum int) {
			defer wg.Done()
			c.concurrentReq(connNum, ran)
		}(connNum)

		nextTime = nextTime.Add(reqWindow)
		timeutils.SleepUtil(nextTime)
	}
	wg.Wait()
}

func (c *Thing) handleNotice(msg msgserver.Message) {
	task := c.CongTasks[rand.Intn(len(c.CongTasks))].Clone()
	timeutils.SleepUtil(msg.SendTime().Add(waitTime))
	c.Request(task)
}

func (c *Thing) concurrentReq(num int, ran *rand.Rand) {
	numPerMs := num / 10
	numRemain := num % 10

	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < 10; i++ {
		curNum := numPerMs
		if numRemain > 0 {
			curNum++
			numRemain--
		}

		for ri := 0; ri < curNum; ri++ {
			go func() {
				defer wg.Done()
				task := c.LoadTasks[ran.Intn(len(c.LoadTasks))].Clone()
				c.Request(task)
			}()
		}
		time.Sleep(time.Millisecond)
	}
	wg.Wait()
}

func (c *Thing) Request(task *thingms.Task) {
	msg := msgserver.NewMessage(idutils.MessageID(idutils.SvrId32(c.Me), idutils.CliId32(c.Me), <-c.mid),
		c.Me, wangID, msgserver.TaskMsg, task.ToBytes())

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
	log.Printf("[%s] send a message:%s, size: %d\n", timeutils.Time2string(msg.SendTime()), idutils.String(msg.ID()), len(msg))
}
