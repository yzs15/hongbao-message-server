package fakething

import "ict.ac.cn/hbmsgserver/pkg/mathutils"

const (
	Cycle   Mode = "cycle"   // 周期发送
	Uniform Mode = "uniform" // 均匀分布
	Normal  Mode = "normal"  // 正态分布
)

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
