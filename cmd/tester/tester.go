package main

import (
	"encoding/json"
	"fmt"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"
	"ict.ac.cn/hbmsgserver/pkg/msgserver"
	"ict.ac.cn/hbmsgserver/pkg/thingms"
)

func main() {
	thingEnd := "tcp://127.0.0.1:5553"

	task := &thingms.Task{
		Name:  "fib",
		Query: "3",
		File:  "",
	}
	taskRaw, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}

	msg := &msgserver.Message{
		ID:      1,
		Sender:  1,
		Good:    true,
		Content: string(taskRaw),
	}
	msgRaw, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(msgRaw))

	if err := czmqutils.Send(thingEnd, msgRaw); err != nil {
		panic(err)
	}

	for {

	}
}
