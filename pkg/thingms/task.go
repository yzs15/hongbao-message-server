package thingms

import (
	"bytes"
	"encoding/binary"
)

type Task struct {
	ID        uint64
	Sender    uint32
	Good      uint8
	ServiceID uint8
	SendTime  uint64
	Args      []byte
}

func (t *Task) ToBytes() []byte {
	buf := new(bytes.Buffer)

	idBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idBytes, t.ID)

	senderBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(senderBytes, t.Sender)

	sendTimeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(sendTimeBytes, t.SendTime)

	buf.Write(idBytes)
	buf.Write(senderBytes)
	buf.WriteByte(t.Good)
	buf.WriteByte(t.ServiceID)
	buf.Write(sendTimeBytes)
	buf.Write(t.Args)
	return buf.Bytes()
}

func ParseTask(raw []byte) *Task {
	i := 0
	task := &Task{}

	task.ID = binary.LittleEndian.Uint64(raw[i : i+8])
	i += 8

	task.Sender = binary.LittleEndian.Uint32(raw[i : i+4])
	i += 4

	task.Good = raw[i]
	i += 1

	task.ServiceID = raw[i]
	i += 1

	task.SendTime = binary.LittleEndian.Uint64(raw[i : i+8])
	i += 8

	task.Args = raw[i:]

	return task
}
