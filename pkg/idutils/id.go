package idutils

const idLen = 30

func CliId32(id uint64) uint32 {
	return uint32(id >> idLen)
}

func SvrId32(id uint64) uint32 {
	return uint32(id & ((uint64(1) << idLen) - 1))
}

func CompleteId(sid uint32, cid uint32) uint64 {
	return uint64(cid)<<idLen | uint64(sid)
}
