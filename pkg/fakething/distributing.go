package fakething

import (
	"fmt"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/mathutils"
	"ict.ac.cn/hbmsgserver/pkg/timeutils"
)

const (
	Cycle   Mode = "cycle"   // 周期发送
	Uniform Mode = "uniform" // 均匀分布
	Normal  Mode = "normal"  // 正态分布
	Linear  Mode = "linear"  //线性增加

)

var disFuncMap map[Mode]func(Config) ([]int, int)

func init() {
	disFuncMap = make(map[Mode]func(Config) ([]int, int))
	disFuncMap[Cycle] = cycleConnDis
	disFuncMap[Normal] = normalConnDis
	disFuncMap[Uniform] = uniformConnDis
	disFuncMap[Linear] = linearConDis
}

func linearConDis(c Config) ([]int, int) {
	up_time := c.MaxTaskPerSec / c.LinearRatio
	mid_time := 10
	down_time := up_time
	seconds := up_time + mid_time + down_time
	c.TotalTime = timeutils.Duration{time.Duration(seconds) * time.Second}
	reqWindowNumPerSec := time.Second.Milliseconds() / reqWindow.Milliseconds()
	totalTimeSlice := c.TotalTime.Milliseconds() / reqWindow.Milliseconds()
	fmt.Println("total time ", seconds)
	connSum := 0
	connDis := make([]int, totalTimeSlice)
	var i int64
	for secondIndex := 0; secondIndex < up_time; secondIndex++ {
		taskPerSecond := (secondIndex + 1) * c.LinearRatio
		taskLeft := taskPerSecond
		connSum += taskPerSecond
		period := reqWindowNumPerSec / int64(taskPerSecond)
		for i = int64(secondIndex) * reqWindowNumPerSec; i < int64(secondIndex+1)*reqWindowNumPerSec; i = i + period {
			connDis[i] = 1
			taskLeft--
			if taskLeft <= 0 {
				break
			}
		}
		if taskLeft > 0 {
			connDis[int64(secondIndex+1)*reqWindowNumPerSec-1] = taskLeft
		}
		fmt.Println(secondIndex+1, "taskPerSecond: ", taskPerSecond, "period", period)

	}

	for secondIndex := up_time; secondIndex < up_time+mid_time; secondIndex++ {
		taskPerSecond := c.MaxTaskPerSec
		taskLeft := taskPerSecond
		connSum += taskPerSecond
		period := reqWindowNumPerSec / int64(taskPerSecond)
		for i = int64(secondIndex) * reqWindowNumPerSec; i < int64(secondIndex+1)*reqWindowNumPerSec; i = i + period {
			connDis[i] = 1
			taskLeft--
			if taskLeft <= 0 {
				break
			}
		}
		if taskLeft > 0 {
			connDis[int64(secondIndex+1)*reqWindowNumPerSec-1] = taskLeft
		}
		fmt.Println(secondIndex+1, "taskPerSecond: ", taskPerSecond, "period", period)
	}

	for secondIndex := up_time + mid_time; secondIndex < up_time+mid_time+down_time; secondIndex++ {
		taskPerSecond := (seconds - secondIndex) * c.LinearRatio
		taskLeft := taskPerSecond
		connSum += taskPerSecond
		period := reqWindowNumPerSec / int64(taskPerSecond)
		for i = int64(secondIndex) * reqWindowNumPerSec; i < int64(secondIndex+1)*reqWindowNumPerSec; i = i + period {
			connDis[i] = 1
			taskLeft--
			if taskLeft <= 0 {
				break
			}
		}
		if taskLeft > 0 {
			connDis[int64(secondIndex+1)*reqWindowNumPerSec-1] = taskLeft
		}
		fmt.Println(secondIndex+1, "taskPerSecond: ", taskPerSecond, "period", period)
	}

	fmt.Println("taskTotal: ", connSum)
	num_will_send := 0
	for i = 0; i < totalTimeSlice; i++ {
		num_will_send += connDis[i]
	}
	fmt.Println("num_will_send: ", num_will_send)
	return connDis, connSum
}

func cycleConnDis(c Config) ([]int, int) {
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

func normalConnDis(c Config) ([]int, int) {
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

func uniformConnDis(c Config) ([]int, int) {
	totalTimeSlice := int(c.TotalTime.Milliseconds() / reqWindow.Milliseconds())

	connNumSlice := c.NumConn / totalTimeSlice
	connSum := connNumSlice * totalTimeSlice

	connDis := make([]int, totalTimeSlice)
	for i := range connDis {
		connDis[i] = connNumSlice
	}
	return connDis, connSum
}
