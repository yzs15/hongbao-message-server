package czmqserver

import (
	"bytes"
	"fmt"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/timeutils"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	//"github.com/sirupsen/logrus"
	"gopkg.in/zeromq/goczmq.v4"
)

type CZMQServer struct {
	Addr      string
	MsgServer msgserver.MessageServer
}

func (s *CZMQServer) Run() {
	pullSock, err := goczmq.NewPull(s.Addr)
	if err != nil {
		panic(err)
	}
	defer pullSock.Destroy()
	fmt.Printf("zmq server listen at: %s\n", s.Addr)

	pullSock.SetTcpKeepalive(1)
	pullSock.SetTcpKeepaliveIdle(120)

	cnt := 0
	for {
		cnt++
		fmt.Println("Receiving ", cnt)
		// FIXME 随机出现 'recv frame error' 错误，目前采用忽略的方式
		msgRaws, err := pullSock.RecvMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}
		msgRaw := msgserver.Message(bytes.Join(msgRaws, []byte{}))

		receiveTime := time.Now()
		go s.MsgServer.Handle(receiveTime, msgRaw)
		fmt.Printf("[%s] recevie a zmq message:%s, size: %d\n", timeutils.Time2string(receiveTime), idutils.String(msgRaw.ID()), len(msgRaw))
	}
}
