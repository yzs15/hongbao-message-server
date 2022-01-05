package thingms

import (
	"encoding/json"
	"fmt"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	"ict.ac.cn/hbmsgserver/pkg/logstore"
	"ict.ac.cn/hbmsgserver/pkg/wshub"
)

type ThingMS struct {
	LogStore *logstore.LogStore
	WsHub    *wshub.Hub

	ThingMsgHdl ThingMsgHandler
}

func (s *ThingMS) Handle(receiveTime time.Time, msgRaw []byte) {
	if msgRaw[0] == '{' { // receive one JSON message from Wang
		msg := &msgserver.Message{}
		if err := json.Unmarshal(msgRaw, msg); err != nil {
			fmt.Println(err)
			return
		}
		s.handleWang(msg)
		s.LogStore.Add(msg.ID, receiveTime, logstore.SenderMsgSvrReceived)

	} else { // receive one BLOB message from Thing
		task := ParseTask(msgRaw)
		if task.ServiceID == 0 {
			s.LogStore.Add(task.ID, time.Unix(0, int64(task.SendTime)), logstore.ReceiverReceived)

		} else {
			s.handleThing(task)
			s.LogStore.Add(task.ID, time.Unix(0, int64(task.SendTime)), logstore.SenderSended)
			s.LogStore.Add(task.ID, receiveTime, logstore.ReceiverMsgSvrReceived)
		}
	}

	// FIXME 由于二进制消息中可能包含不是 \n 的 \n，因此不能采用 \n 做消息分割
	//msgBytes := bytes.Split(msgRaw, []byte{'\n'})
	//msgBytes := [][]byte{msgRaw}
	//for _, msgByte := range msgBytes {
	//	go func(msgByte []byte) {
	//		if msgByte[0] == '{' { // receive one JSON message from Wang
	//			msg := &msgserver.Message{}
	//			if err := json.Unmarshal(msgByte, msg); err != nil {
	//				fmt.Println(err)
	//				return
	//			}
	//			s.handleWang(msg)
	//			s.LogStore.Add(msg.ID, receiveTime, logstore.SenderMsgSvrReceived)
	//
	//		} else { // receive one BLOB message from Thing
	//			task := ParseTask(msgByte)
	//			if task.ServiceID == 0 {
	//				s.LogStore.Add(task.ID, time.Unix(0, int64(task.SendTime)), logstore.ReceiverReceived)
	//
	//			} else {
	//				s.handleThing(task)
	//				s.LogStore.Add(task.ID, time.Unix(0, int64(task.SendTime)), logstore.SenderSended)
	//				s.LogStore.Add(task.ID, receiveTime, logstore.ReceiverMsgSvrReceived)
	//			}
	//		}
	//	}(msgByte)
	//}
}

func (s *ThingMS) handleWang(msg *msgserver.Message) {
	msgRaw, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	sendTime := time.Now()
	s.WsHub.Broadcast <- msgRaw
	s.LogStore.Add(msg.ID, sendTime, logstore.SenderMsgSvrSended)
}

func (s *ThingMS) handleThing(task *Task) {
	sendTime, err := s.ThingMsgHdl.Handle(task)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.LogStore.Add(task.ID, sendTime, logstore.ReceiverMsgSvrSended)
}
