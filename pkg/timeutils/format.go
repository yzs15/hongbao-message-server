package timeutils

import (
	"fmt"
	"time"
)

func Time2string(t time.Time) string {
	ymdhms := t.Format("2006-01-02 15:04:05")
	ms := t.UnixMilli() % 1e3
	us := t.UnixMicro() % 1e3
	ns := t.UnixNano() % 1e3
	return fmt.Sprintf("%s.%d.%d.%d", ymdhms, ms, us, ns)
}
