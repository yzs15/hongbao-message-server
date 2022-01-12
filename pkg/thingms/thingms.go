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

	if idutils.CliId32(msg.Receiver()) == idutils.FullId {
		s.handleBroadcast(msg, receiveTime)
		return
	}

	svrID := idutils.SvrId32(msg.Receiver())
	if svrID != s.Me && msg.Type() != msgserver.TaskMsg {
		s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), msg.Sender(), msg.SendTime(), logstore.Send)
		s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), uint64(s.Me), receiveTime, logstore.Receive)

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

		s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), uint64(s.Me), sendTime, logstore.Send)
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

func (s *ThingMS) handleBroadcast(msg msgserver.Message, receiveTime time.Time) {
	s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), uint64(s.Me), receiveTime, logstore.Receive)

	// 客户端发布广播消息到的首个Message Server
	if idutils.SvrId32(msg.Receiver()) == idutils.FullId {
		s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), msg.Sender(), msg.SendTime(), logstore.Send)

		servers, ids := s.NameServer.GetAllServer()
		for idx := range servers {
			newMsg := msg.Clone()
			go func(idx int) {
				if ids[idx] == s.Me {
					return
				}

				newMsg.SetReceiver(idutils.DeviceId(ids[idx], idutils.FullId))
				sendTime, err := czmqutils.Send(servers[idx].ZMQEndpoint, newMsg)
				if err != nil {
					fmt.Println(err)
					return
				}
				s.LogStore.Add(newMsg.ID(), newMsg.Sender(), newMsg.Receiver(), uint64(s.Me), sendTime, logstore.Send)
			}(idx)
		}
	}

	sendTime := s.Registry.Broadcast(msg)
	s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), uint64(s.Me), sendTime, logstore.Send)
}

func (s *ThingMS) handleText(msg msgserver.Message, receiveTime time.Time) {
	s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), uint64(s.Me), receiveTime, logstore.Receive)

	sendTime, err := s.Registry.Send(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), uint64(s.Me), sendTime, logstore.Send)
}

func (s *ThingMS) handleTask(msg msgserver.Message, receiveTime time.Time) {
	s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), msg.Sender(), msg.SendTime(), logstore.Send)
	s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), uint64(s.Me), receiveTime, logstore.Receive)

	sendTime, err := s.TaskMsgHdl.Handle(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), uint64(s.Me), sendTime, logstore.Send)
}

func (s *ThingMS) handleLog(msg msgserver.Message, receiveTime time.Time) {
	s.LogStore.Add(msg.ID(), msg.Sender(), msg.Receiver(), msg.Sender(), msg.SendTime(), logstore.Receive)
}
