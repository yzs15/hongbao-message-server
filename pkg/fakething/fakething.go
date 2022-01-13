package fakething

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"
	"ict.ac.cn/hbmsgserver/pkg/mathutils"
	"ict.ac.cn/hbmsgserver/pkg/thingms"
	"ict.ac.cn/hbmsgserver/pkg/timeutils"
)

const reqWindow = 10 * time.Millisecond

type Mode string

var wangID = idutils.DeviceId(1, 2)

const (
	Cycle   Mode = "cycle"   // 周期发送
	Uniform Mode = "uniform" // 均匀分布
	Normal  Mode = "normal"  // 正态分布
)

type Thing struct {
	ID      uint32
	MacAddr string
	Me      uint64

	ExpectedTime time.Time

	MsgWsEnd  string
	MsgZmqEnd string

	Tasks []*thingms.Task

	Mode Mode

	Period time.Duration

	NumConn   int
	TotalTime time.Duration
	PeakTime  time.Duration
	PeakNum   int

	mid chan uint32

	wsConn *websocket.Conn
}

func (c *Thing) Run() {
	go func() {
		c.mid = make(chan uint32, 10000)
		var id uint32
		for id = 1; ; id++ {
			c.mid <- id
		}
	}()

	c.Me = c.waitID()
	c.waitNextMessage(msgserver.TextMsg)

	//jpgRaw, err := ioutil.ReadFile("test.jpg")
	//if err != nil {
	//	panic(err)
	//}
	//
	//buf := new(bytes.Buffer)
	//buf.WriteByte(^uint8(0))
	//buf.Write(jpgRaw)
	//msg := msgserver.NewMessage(1, c.Me, wangID, msgserver.TextMsg, buf.Bytes())
	//msg.SetSendTime()
	//fmt.Println("send a jpg message")
	//if _, err := czmqutils.Send(c.MsgZmqEnd, msg); err != nil {
	//	log.Println("czmq send failed: ", err)
	//}

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

	nextTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(connDis))
	for _, connNum := range connDis {
		go func(connNum int) {
			c.concurrentReq(connNum)
			wg.Done()
		}(connNum)

		nextTime = nextTime.Add(reqWindow)
		timeutils.SleepUtil(nextTime)
	}
	wg.Wait()

	msg := c.waitNextMessage(msgserver.TextMsg)
	timeutils.SleepUtil(msg.SendTime().Add(50 * time.Millisecond))
	c.concurrentReq(1)
}

func (c *Thing) concurrentReq(num int) {
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
				task := c.Tasks[rand.Intn(len(c.Tasks))].Clone()
				msg := msgserver.NewMessage(idutils.MessageID(idutils.SvrId32(c.Me), idutils.CliId32(c.Me), <-c.mid),
					c.Me, wangID, msgserver.TaskMsg, task.ToBytes())

				msg.SetSendTime()
				if _, err := czmqutils.Send(c.MsgZmqEnd, msg); err != nil {
					log.Println("czmq send failed: ", err)
				}
				log.Printf("[%s] send a message, size: %d\n",
					timeutils.Time2string(msg.SendTime()),
					len(msg))
				wg.Done()
			}()
		}
		time.Sleep(time.Millisecond)
	}
	wg.Wait()
}

func (c *Thing) cycleConnDis() ([]int, int) {
	totalTimeSlice := c.TotalTime.Milliseconds() / reqWindow.Milliseconds()
	periodSlice := c.Period.Milliseconds() / reqWindow.Milliseconds()

	connSum := int(totalTimeSlice / periodSlice)
	connDis := make([]int, totalTimeSlice)
	var i int64
	for i = 0; i < totalTimeSlice; i += periodSlice {
		connDis[i] = 1
	}
	return connDis, connSum
}

func (c *Thing) normalConnDis() ([]int, int) {
	totalTimeSlice := c.TotalTime.Milliseconds() / reqWindow.Milliseconds()
	peakTimePos := c.PeakTime.Milliseconds() / reqWindow.Milliseconds()
	peakPro := float64(c.PeakNum) / float64(c.NumConn)

	var_ := mathutils.CalVariance(peakTimePos, peakPro)

	connSum := 0
	connDis := make([]int, totalTimeSlice)
	for i := range connDis {
		connDis[i] = int((mathutils.NormalFunc(peakTimePos, var_, int64(i)) * float64(c.NumConn)) + 0.5)
		connSum += connDis[i]
	}
	return connDis, connSum
}

func (c *Thing) uniformConnDis() ([]int, int) {
	totalTimeSlice := int(c.TotalTime.Milliseconds() / reqWindow.Milliseconds())

	connNumSlice := c.NumConn / totalTimeSlice
	connSum := connNumSlice * totalTimeSlice

	connDis := make([]int, totalTimeSlice)
	for i := range connDis {
		connDis[i] = connNumSlice
	}
	return connDis, connSum
}
