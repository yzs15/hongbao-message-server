package wangms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	"ict.ac.cn/hbmsgserver/pkg/logstore"
	"ict.ac.cn/hbmsgserver/pkg/wshub"
)

type WangMS struct {
	LogStore *logstore.LogStore
	WsHub    *wshub.Hub

	ThingMsEnds []string
}

func (s *WangMS) Handle(receiveTime time.Time, msgRaw []byte) {
	msgBytes := bytes.Split(msgRaw, []byte{'\n'})
	for _, msgByte := range msgBytes {
		msg := &msgserver.Message{}
		if err := json.Unmarshal(msgByte, msg); err != nil {
			fmt.Println(err)
			continue
		}

		go func() {
			if msg.Sender == 0 {
				if msg.Good == 3 {
					s.LogStore.Add(msg.ID, time.Unix(0, int64(msg.SendTime)), logstore.ReceiverReceived)

				} else {
					s.LogStore.Add(msg.ID, time.Unix(0, int64(msg.SendTime)), logstore.SenderSended)
					msg.SendTime = 0
					s.handleWang(msg)
					s.LogStore.Add(msg.ID, receiveTime, logstore.SenderMsgSvrReceived)
				}

			} else {
				s.handleThing(msg)
				s.LogStore.Add(msg.ID, receiveTime, logstore.ReceiverMsgSvrReceived)
			}
		}()
	}
}

func (s *WangMS) handleWang(msg *msgserver.Message) {
	msgRaw, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, endpoint := range s.ThingMsEnds {
		sendTime := time.Now()
		if err := czmqutils.Send(endpoint, msgRaw); err == nil {
			s.LogStore.Add(msg.ID, sendTime, logstore.ReceiverMsgSvrSended)
		}
	}
}

func (s *WangMS) handleThing(msg *msgserver.Message) {
	msgRaw, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	sendTime := time.Now()
	s.WsHub.Broadcast <- msgRaw
	s.LogStore.Add(msg.ID, sendTime, logstore.SenderMsgSvrSended)
}
