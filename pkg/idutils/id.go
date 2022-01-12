package idutils

import "fmt"

const idLen = 20

const FullId = (1 << idLen) - 1

func MsgId32(id uint64) uint32 {
	return uint32(id >> (2 * idLen))
}

func CliId32(id uint64) uint32 {
	return uint32((id >> idLen) & ((uint64(1) << idLen) - 1))
}

func SvrId32(id uint64) uint32 {
	return uint32(id & ((uint64(1) << idLen) - 1))
}

func DeviceId(sid uint32, cid uint32) uint64 {
	return uint64(cid)<<idLen | uint64(sid)
}

func MessageID(sid uint32, cid uint32, mid uint32) uint64 {
	return uint64(mid)<<(2*idLen) | uint64(cid)<<idLen | uint64(sid)
}

func String(id uint64) string {
	return fmt.Sprintf("%d.%d.%d", MsgId32(id), CliId32(id), SvrId32(id))
}
