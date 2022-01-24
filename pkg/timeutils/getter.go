package timeutils

/*
#include <stdio.h>
#include <stdint.h>
#include <time.h>
#include <errno.h>
#include <fcntl.h>
#include <linux/ptp_clock.h>

#define CLOCK_INVALID -1

#define CLOCKFD 3
#define FD_TO_CLOCKID(fd)   ((clockid_t) ((((unsigned int) ~fd) << 3) | CLOCKFD))
#define CLOCKID_TO_FD(clk)  ((unsigned int) ~((clk) >> 3))

uint64_t utils_get_current_sys_ns_timestamp()
{
    struct timespec current_time;
    clock_gettime(CLOCK_REALTIME, &current_time);
    time_t timestamp = current_time.tv_sec*1000000000+current_time.tv_nsec;
    return timestamp;
}

clockid_t utils_get_clockid(const char *dev_ptp)
{
    printf("dev_ptp %s\n", dev_ptp);
    int fd;
    clockid_t clkid;
    fd = open(dev_ptp, O_RDWR);
    if (fd < 0)
    {
        printf("ptp open error.\n");
        return CLOCK_INVALID;
    }
    clkid = FD_TO_CLOCKID(fd);
    if ((clkid & CLOCKFD) != 3)
    {
        printf("ptp clock id error.\n");
        return CLOCK_INVALID;
    }
    return clkid;
}

uint64_t utils_get_current_ptp_ns_timestamp(clockid_t clkid)
{
    struct timespec current_time;
    clock_gettime(clkid, &current_time);
    time_t timestamp = current_time.tv_sec*1000000000+current_time.tv_nsec;
    return timestamp;
}
*/
import "C"
import (
	"bytes"
	"fmt"
	"time"
	"unsafe"
)

const (
	NetEnv = 0
	SpbEnv = 1
)

var clkId int = -1

const ptpLoss = uint64(37) * 1000000000

type MyString struct {
	Str *C.char
	Len int
}

func SetClkID(dev string) {
	var buf = new(bytes.Buffer)
	buf.WriteString(dev)
	buf.WriteByte(0)
	devName := buf.String()
	devC := (*MyString)(unsafe.Pointer(&devName))
	clkId = int(C.utils_get_clockid(devC.Str))
	fmt.Println("Clock ID: ", clkId)
}

func GetSysTime(env int) time.Time {
	if NetEnv == env {
		return time.Now()
	} else {
		stamp := uint64(C.utils_get_current_ptp_ns_timestamp(C.int(clkId))) - ptpLoss
		return time.Unix(0, int64(stamp))
	}
}

func GetPtpTime() time.Time {
	stamp := uint64(C.utils_get_current_ptp_ns_timestamp(C.int(clkId))) - ptpLoss
	return time.Unix(0, int64(stamp))
}
