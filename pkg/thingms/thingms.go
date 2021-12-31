package thingms

import (
	"bytes"
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
	msgBytes := bytes.Split(msgRaw, []byte{'\n'})
	for _, msgByte := range msgBytes {
		msg := &msgserver.Message{}
		if err := json.Unmarshal(msgByte, msg); err != nil {
			fmt.Println(err)
			continue
		}

		go func() {
			if msg.Sender == 0 {
				s.handleWang(msg)
				s.LogStore.Add(msg.ID, receiveTime, logstore.SenderMsgSvrReceived)
			} else {
				s.handleThing(msg)
				s.LogStore.Add(msg.ID, receiveTime, logstore.ReceiverMsgSvrReceived)
			}
		}()
	}
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

func (s *ThingMS) handleThing(msg *msgserver.Message) {
	sendTime, err := s.ThingMsgHdl.Handle(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.LogStore.Add(msg.ID, sendTime, logstore.ReceiverMsgSvrSended)
}
