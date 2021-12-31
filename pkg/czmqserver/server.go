package czmqserver

import (
	"bytes"
	"fmt"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	//"github.com/sirupsen/logrus"
	"gopkg.in/zeromq/goczmq.v4"
)

type CZMQServer struct {
	Addr      string
	MsgServer msgserver.MessageServer
}

func (s *CZMQServer) Run() {
	pullSock := goczmq.NewSock(goczmq.Pull)
	defer pullSock.Destroy()

	_, err := pullSock.Bind(s.Addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("zmq server listen at: %s\n", s.Addr)

	for {
		msgRaws, err := pullSock.RecvMessage()
		if err != nil {
			panic(err)
		}
		msgRaw := bytes.Join(msgRaws, []byte{})

		receiveTime := time.Now()
		go s.MsgServer.Handle(receiveTime, msgRaw)
		fmt.Printf("recevie a zmq message: %s\n", msgRaw)
	}
}
