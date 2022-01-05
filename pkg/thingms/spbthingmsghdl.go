package thingms

/*
#include <time.h>
#include <stdint.h>
uint64_t my_GetTime() {
    struct timespec current_time;
    clock_gettime(CLOCK_REALTIME, &current_time);
	uint64_t time = current_time.tv_sec*1000000000+current_time.tv_nsec;
	return time;
}
*/
import "C"
import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"gopkg.in/zeromq/goczmq.v4"
)

func Send(sock *goczmq.Sock, data []byte, flag int) {
	err := sock.SendFrame(data, flag)
	for err != nil {
		fmt.Println("!!! SafeSendFrame ERROR", err)
		err = sock.SendFrame(data, flag)
	}
}

func Recv(sock *goczmq.Sock) [][]byte {
	for {
		msg, err := sock.RecvMessage()
		if err == nil {
			return msg
		}
	}
}

type ClientConfig struct {
	Task_send_endpoint      string
	Worker_send_endpoint    string
	Result_receive_endpoint string

	Worker_path     string
	Worker_filename string
	Worker_priority uint32

	Task_path        string
	Task_config_list []string
}

type TaskConfig struct {
	Body_file string
	Args_file string // Not used

	From         string
	To           string
	Task_body_id uint64
	Priority     uint32

	End_before            uint64
	Estimate_running_time uint64

	Cpu uint64
	Mem uint64
	Gpu uint64

	Estimate_strat_time  uint64
	Estimate_result_size uint64
	Estimate_call_time   uint64
}

func get_fake_worker_id() uint64 {
	return uint64(rand.Intn(99999999))
}

func get_fake_task_sub_id() uint64 {
	return uint64(rand.Intn(99999999))
}

func commit_worker(send_worker *goczmq.Sock, worker string, priority uint32) uint64 {
	worker_file, err := os.Open(worker)
	if err != nil {
		fmt.Println("[C] Open Worker File Error")
		return 0
	}
	defer worker_file.Close()
	worker_file_info, _ := worker_file.Stat()
	var worker_file_size int64 = worker_file_info.Size()
	worker_buf := make([]byte, worker_file_size)
	binary.Read(worker_file, binary.LittleEndian, worker_buf)

	var worker_info WorkerPackageInfo
	worker_info.Worker_id = get_fake_worker_id()
	worker_info.Priority = priority
	worker_info.Command_size = 1
	worker_info.Worker_body_size = uint64(worker_file_size)

	worker_header := worker_info.Info2Header()
	Send(send_worker, []byte("worker_package"), goczmq.FlagMore)
	Send(send_worker, worker_header, goczmq.FlagMore)
	Send(send_worker, worker_buf, goczmq.FlagMore)
	Send(send_worker, []byte("\x00"), goczmq.FlagNone)
	fmt.Println("[C] User Worker Sent")

	worker_send_reply := Recv(send_worker)
	fmt.Println("[C] Mechine Replied.", string(worker_send_reply[0]))
	return worker_info.Worker_id
}

type spbThingMsgHandler struct {
	Worker_id        uint64
	Send_task        *goczmq.Sock
	With_body        bool
	Task_config_file string
}

func NewSpbThingMsgHandler(spbConfig string) ThingMsgHandler {
	rand.Seed(time.Now().UnixNano())
	client_config_file, err := ioutil.ReadFile(spbConfig)
	if err != nil {
		fmt.Print("[C] Open Client Config File Error", err)
	}
	var client_config ClientConfig
	json.Unmarshal(client_config_file, &client_config)

	send_worker, err := goczmq.NewReq(client_config.Worker_send_endpoint)
	defer send_worker.Destroy()
	send_task, _ := goczmq.NewReq(client_config.Task_send_endpoint)
	defer send_task.Destroy()

	worker_id := commit_worker(send_worker, client_config.Worker_path+client_config.Worker_filename, client_config.Worker_priority)
	task_config_file := client_config.Task_path + client_config.Task_config_list[0]

	return &spbThingMsgHandler{
		Worker_id:        worker_id,
		Send_task:        send_task,
		With_body:        true,
		Task_config_file: task_config_file,
	}
}

func (h *spbThingMsgHandler) Handle(task *Task) (time.Time, error) {
	worker_id := h.Worker_id
	send_task := h.Send_task
	task_config_file, err := ioutil.ReadFile(h.Task_config_file)
	if err != nil {
		fmt.Println("[C] Open Task Config File Error", err)
	}
	var task_config TaskConfig
	json.Unmarshal(task_config_file, &task_config)

	var task_file_size int64
	var task_buf []byte
	if h.With_body {
		task_file, err := os.Open(task_config.Body_file)
		if err != nil {
			fmt.Println("[C] Open Task File Error")
			return time.Time{}, err
		}
		defer task_file.Close()
		task_file_info, _ := task_file.Stat()
		task_file_size = task_file_info.Size()
		task_buf = make([]byte, task_file_size)
		binary.Read(task_file, binary.LittleEndian, task_buf)
	}

	var task_info TaskPackageInfo
	copy(task_info.From[:], []byte(task_config.From))
	copy(task_info.To[:], []byte(task_config.To))
	task_info.Task_sub_id = get_fake_task_sub_id()
	task_info.Task_type = FUNCTION
	task_info.Task_body_id = task_config.Task_body_id
	task_info.Worker_id = worker_id
	task_info.Priority = task_config.Priority
	task_info.Timestamp = uint64(C.my_GetTime())
	if h.With_body {
		task_info.Task_body_size = uint64(task_file_size)
	} else {
		task_info.Task_body_size = uint64(0)
	}
	task_info.Task_args_size = uint64(len(task.Args))
	task_info.Task_QoS.End_before = task_info.Timestamp + task_config.End_before
	task_info.Task_QoS.Estimate_running_time = task_config.Estimate_running_time
	task_info.Task_body_constraint.Cpu = task_config.Cpu
	task_info.Task_body_constraint.Mem = task_config.Mem
	task_info.Task_body_constraint.Gpu = task_config.Gpu
	task_info.Estimate_strat_time = task_config.Estimate_strat_time
	task_info.Estimate_result_size = task_config.Estimate_result_size
	task_info.Estimate_call_time = task_config.Estimate_call_time

	task_header := task_info.Info2Header()
	Send(send_task, task_header, goczmq.FlagMore)
	if h.With_body {
		Send(send_task, task_buf, goczmq.FlagMore)
		h.With_body = false
	}
	Send(send_task, task.Args, goczmq.FlagNone)

	send_ok_msg := Recv(send_task)
	fmt.Println("[C] User Task Sent", string(send_ok_msg[0]))
	return time.Time{}, nil
}
