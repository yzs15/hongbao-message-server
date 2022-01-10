package main

import (
	"fmt"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"
	"ict.ac.cn/hbmsgserver/pkg/msgserver"
)

func main() {
	msg := msgserver.NewMessage(uint64(time.Now().UnixNano()), idutils.CompleteId(2, 1), idutils.CompleteId(1, 1), msgserver.TextMsg, []byte("Hello"))

	fmt.Println(msg)

	end := "tcp://58.213.121.2:10025"
	if _, err := czmqutils.Send(end, msg); err != nil {
		fmt.Println(err)
	}

	for {

	}
}
