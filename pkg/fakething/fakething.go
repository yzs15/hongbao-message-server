package fakething

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"
	"ict.ac.cn/hbmsgserver/pkg/mathutils"
	"ict.ac.cn/hbmsgserver/pkg/thingms"
	"ict.ac.cn/hbmsgserver/pkg/timeutils"
)

const reqWindow = 10 * time.Millisecond

type Mode string

const (
	Cycle   Mode = "cycle"   // 周期发送
	Uniform Mode = "uniform" // 均匀分布
	Normal  Mode = "normal"  // 正态分布
)

type Thing struct {
	ID           uint32
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
}

func (c *Thing) Run() {
	//receiveTime := c.waitFirst()
	//good := receiveTime.Before(c.ExpectedTime)
	//if good {
	//	c.Task.Good = 1
	//	timeutils.SleepUtil(c.ExpectedTime)
	//} else {
	//	c.Task.Good = 0
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
				task.ID = thingms.GenerateTID(task.Sender)
				task.SendTime = uint64(time.Now().UnixNano())

				if err := czmqutils.Send(c.MsgZmqEnd, task.ToBytes()); err != nil {
					log.Println("czmq send failed: ", err)
				}
				fmt.Printf("send a message: %+v\n", task)
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
