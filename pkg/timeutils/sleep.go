package timeutils

import "time"

const sleepThreshold = 50 * time.Microsecond

func SleepUtil(t time.Time, env int) {
	curTime := GetSysTime(env)

	sleepTime := t.Sub(curTime.Add(sleepThreshold))
	time.Sleep(sleepTime)

	for {
		curTime = GetSysTime(env)
		if !curTime.Before(t) {
			break
		}
	}
}
