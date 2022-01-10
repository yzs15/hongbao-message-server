package thingms

import (
	"bytes"
	"math/rand"
)

type Task struct {
	ServiceID uint8
	Args      []byte
}

func (t *Task) Clone() *Task {
	task := &Task{ServiceID: t.ServiceID}
	task.Args = make([]byte, len(t.Args))
	copy(task.Args, t.Args)
	return task
}

func (t *Task) ToBytes() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(t.ServiceID)
	buf.Write(t.Args)
	return buf.Bytes()
}

func ParseTask(raw []byte) *Task {
	i := 0
	task := &Task{}

	task.ServiceID = raw[i]
	i += 1

	task.Args = raw[i:]

	return task
}

func GenerateTID(sender uint32) uint64 {
	return uint64(sender)<<32 | uint64(rand.Uint32())
}
