package thingms

import "encoding/binary"

type TaskType uint8

type WorkerPackageInfo struct {
	Worker_id        uint64
	Priority         uint32
	Command_size     uint64
	Worker_body_size uint64
}

func (info *WorkerPackageInfo) Info2Header() []byte {
	header := make([]byte, 28)
	i := 0

	b_Worker_id := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Worker_id, info.Worker_id)
	i += copy(header[i:], b_Worker_id)

	b_Priority := make([]byte, 4)
	binary.LittleEndian.PutUint32(b_Priority, info.Priority)
	i += copy(header[i:], b_Priority)

	b_Command_size := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Command_size, info.Command_size)
	i += copy(header[i:], b_Command_size)

	b_Worker_body_size := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Worker_body_size, info.Worker_body_size)
	i += copy(header[i:], b_Worker_body_size)

	return header
}

func (info *WorkerPackageInfo) Header2Info(header []byte) {
	i := 0

	info.Worker_id = binary.LittleEndian.Uint64(header[i : i+4])
	i += 4

	info.Priority = binary.LittleEndian.Uint32(header[i : i+4])
	i += 4

	info.Command_size = binary.LittleEndian.Uint64(header[i : i+4])
	i += 4

	info.Worker_body_size = binary.LittleEndian.Uint64(header[i : i+4])
	i += 4

	return
}

const (
	ACTOR    TaskType = 0
	FUNCTION TaskType = 1
)

type TaskBodyConstraint struct {
	Cpu uint64
	Mem uint64
	Gpu uint64
}

type TaskQoS struct {
	End_before            uint64
	Estimate_running_time uint64
}

type TaskPackageInfo struct {
	From                 [32]byte
	To                   [32]byte
	Task_sub_id          uint64
	Task_type            TaskType
	Task_body_id         uint64
	Worker_id            uint64
	Priority             uint32
	Timestamp            uint64
	Task_body_size       uint64
	Task_args_size       uint64
	Task_QoS             TaskQoS
	Task_body_constraint TaskBodyConstraint
	Estimate_strat_time  uint64
	Estimate_result_size uint64
	Estimate_call_time   uint64
}

func (info *TaskPackageInfo) Info2Header() []byte {
	header := make([]byte, 181)
	i := 0

	i += copy(header[i:], info.From[:])
	i += copy(header[i:], info.To[:])

	b_Task_sub_id := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Task_sub_id, info.Task_sub_id)
	i += copy(header[i:], b_Task_sub_id)

	b_Task_type := make([]byte, 1)
	b_Task_type[0] = uint8(info.Task_type)
	i += copy(header[i:], b_Task_type)

	b_Task_body_id := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Task_body_id, info.Task_body_id)
	i += copy(header[i:], b_Task_body_id)

	b_Worker_id := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Worker_id, info.Worker_id)
	i += copy(header[i:], b_Worker_id)

	b_Priority := make([]byte, 4)
	binary.LittleEndian.PutUint32(b_Priority, info.Priority)
	i += copy(header[i:], b_Priority)

	b_Timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Timestamp, info.Timestamp)
	i += copy(header[i:], b_Timestamp)

	b_Task_body_size := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Task_body_size, info.Task_body_size)
	i += copy(header[i:], b_Task_body_size)

	b_Task_args_size := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Task_args_size, info.Task_args_size)
	i += copy(header[i:], b_Task_args_size)

	b_End_before := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_End_before, info.Task_QoS.End_before)
	i += copy(header[i:], b_End_before)

	b_Estimate_running_time := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Estimate_running_time, info.Task_QoS.Estimate_running_time)
	i += copy(header[i:], b_Estimate_running_time)

	b_Cpu := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Cpu, info.Task_body_constraint.Cpu)
	i += copy(header[i:], b_Cpu)

	b_Mem := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Mem, info.Task_body_constraint.Mem)
	i += copy(header[i:], b_Mem)

	b_Gpu := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Gpu, info.Task_body_constraint.Gpu)
	i += copy(header[i:], b_Gpu)

	b_Estimate_strat_time := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Estimate_strat_time, info.Estimate_strat_time)
	i += copy(header[i:], b_Estimate_strat_time)

	b_Estimate_result_size := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Estimate_result_size, info.Estimate_result_size)
	i += copy(header[i:], b_Estimate_result_size)

	b_Estimate_call_time := make([]byte, 8)
	binary.LittleEndian.PutUint64(b_Estimate_call_time, info.Estimate_call_time)
	i += copy(header[i:], b_Estimate_call_time)

	return header
}

func (info *TaskPackageInfo) Header2Info(header []byte) {
	i := 0

	i += copy(info.From[:], header[i:])
	i += copy(info.To[:], header[i:])

	info.Task_sub_id = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Task_type = TaskType(uint8(header[i]))
	i += 1

	info.Task_body_id = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Worker_id = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Priority = binary.LittleEndian.Uint32(header[i : i+4])
	i += 4

	info.Timestamp = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Task_body_size = binary.LittleEndian.Uint64(header[i : i+4])
	i += 8

	info.Task_args_size = binary.LittleEndian.Uint64(header[i : i+4])
	i += 8

	info.Task_QoS.End_before = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Task_QoS.Estimate_running_time = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Task_body_constraint.Cpu = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Task_body_constraint.Mem = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Task_body_constraint.Gpu = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Estimate_strat_time = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Estimate_result_size = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8

	info.Estimate_call_time = binary.LittleEndian.Uint64(header[i : i+8])
	i += 8
}
