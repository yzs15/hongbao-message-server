package thingms

import (
	"fmt"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/nameserver"

	"ict.ac.cn/hbmsgserver/pkg/registry"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"
	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	"ict.ac.cn/hbmsgserver/pkg/logstore"
)

type ThingMS struct {
	Me         uint32
	Registry   *registry.Registry
	NameServer *nameserver.NameServer

	TaskMsgHdl TaskMsgHandler

	LogStore *logstore.LogStore
}

func (s *ThingMS) Handle(receiveTime time.Time, msg msgserver.Message) {
	fmt.Printf("\t%+v\n", msg)
	svrID := idutils.SvrId32(msg.Receiver())
	if svrID != s.Me && msg.Type() != msgserver.TaskMsg {
		s.LogStore.Add(msg.ID(), msg.SendTime(), logstore.SenderSended)
		s.LogStore.Add(msg.ID(), receiveTime, logstore.SenderMsgSvrReceived)

		svr, err := s.NameServer.GetServer(svrID)
		if err != nil {
			fmt.Println(err)
			return
		}
		var sendTime time.Time
		if sendTime, err = czmqutils.Send(svr.ZMQEndpoint, msg); err != nil {
			fmt.Println(err)
			return
		}

		s.LogStore.Add(msg.ID(), sendTime, logstore.SenderMsgSvrSended)
		return
	}

	switch msg.Type() {
	case msgserver.TextMsg:
		s.handleText(msg, receiveTime)

	case msgserver.TaskMsg:
		s.handleTask(msg, receiveTime)

	case msgserver.LogMsg:
		s.handleLog(msg, receiveTime)
	}
}

func (s *ThingMS) handleText(msg msgserver.Message, receiveTime time.Time) {
	cliID := idutils.CliId32(msg.Receiver())
	cli, err := s.Registry.GetClient(cliID)
	if err != nil {
		fmt.Println(err)
		return
	}
	cli.WsClient.Send <- msg
	s.LogStore.Add(msg.ID(), receiveTime, logstore.ReceiverMsgSvrReceived)
}

func (s *ThingMS) handleTask(msg msgserver.Message, receiveTime time.Time) {
	s.LogStore.Add(msg.ID(), msg.SendTime(), logstore.SenderSended)
	s.LogStore.Add(msg.ID(), receiveTime, logstore.SenderMsgSvrReceived)

	sendTime, err := s.TaskMsgHdl.Handle(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.LogStore.Add(msg.ID(), sendTime, logstore.SenderMsgSvrSended)
}

func (s *ThingMS) handleLog(msg msgserver.Message, receiveTime time.Time) {
	s.LogStore.Add(msg.ID(), msg.SendTime(), logstore.ReceiverReceived)
}
