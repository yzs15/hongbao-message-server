package msgserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/timeutils"
)

type MessageType uint8

const (
	RegisterMsg MessageType = iota
	NameMsg
	QueryMsg
	ResultMsg
	TextMsg
	TaskMsg
	LogMsg
)

type Message []byte

func NewMessage(id, sender, receiver uint64, typ MessageType, body []byte) Message {
	buf := new(bytes.Buffer)

	idBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idBytes, id)

	senderBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(senderBytes, sender)

	receiverBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(receiverBytes, receiver)

	sendTimeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(sendTimeBytes, uint64(time.Now().UnixNano()))

	buf.Write(idBytes)
	buf.Write(senderBytes)
	buf.Write(receiverBytes)
	buf.WriteByte(uint8(typ))
	buf.Write(body)
	buf.Write(sendTimeBytes)
	return buf.Bytes()
}

func (m Message) ID() uint64 {
	return binary.LittleEndian.Uint64(m[:8])
}

func (m Message) Sender() uint64 {
	return binary.LittleEndian.Uint64(m[8:16])
}

func (m Message) Receiver() uint64 {
	return binary.LittleEndian.Uint64(m[16:24])
}

func (m Message) Type() MessageType {
	return MessageType(m[24])
}

func (m Message) Body() []byte {
	l := len(m)
	return m[25 : l-8]
}

func (m Message) SendTime() time.Time {
	l := len(m)
	return time.Unix(0, int64(binary.LittleEndian.Uint64(m[l-8:])))
}

func (m Message) SetReceiver(id uint64) {
	binary.LittleEndian.PutUint64(m[16:24], id)
}

func (m Message) SetSendTime() {
	l := len(m)
	t := time.Now()
	binary.LittleEndian.PutUint64(m[l-8:], uint64(t.UnixNano()))
}

func (m Message) Clone() Message {
	l := len(m)
	newMsg := make([]byte, l)
	copy(newMsg, m)
	return newMsg
}

func (m Message) String() string {
	var content string
	if len(m.Body()) > 0 {
		l := 5
		if len(m.Body()) < l {
			l = len(m.Body())
		}
		content = string(m.Body()[:l])
	}
	return fmt.Sprintf("Message{%v, %v, %v, %v, %v, %s}",
		idutils.String(m.ID()), idutils.String(m.Sender()), idutils.String(m.Receiver()),
		m.Type(), timeutils.Time2string(m.SendTime()), content)
}
