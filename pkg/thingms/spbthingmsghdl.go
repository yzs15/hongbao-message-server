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
	"log"
	"math/rand"
	"os"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/nameserver"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

	"gopkg.in/zeromq/goczmq.v4"
)

func Send(sock *goczmq.Sock, data []byte, flag int) {
	err := sock.SendFrame(data, flag)
	for err != nil {
		//fmt.Println("!!! SafeSendFrame ERROR", err)
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
	return 99999998
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
	Worker_id          uint64
	Task_send_endpoint string
	With_body          bool
	With_body_7        bool
	Task_path          string
	Task_config_file   []string

	Ns *nameserver.NameServer
}

func NewSpbThingMsgHandler(spbConfig string, ns *nameserver.NameServer) TaskMsgHandler {
	rand.Seed(time.Now().UnixNano())
	client_config_file, err := ioutil.ReadFile(spbConfig)
	if err != nil {
		fmt.Print("[C] Open Client Config File Error", err)
	}
	var client_config ClientConfig
	json.Unmarshal(client_config_file, &client_config)

	send_worker, err := goczmq.NewReq(client_config.Worker_send_endpoint)
	defer send_worker.Destroy()
	//send_task, _ := goczmq.NewReq(client_config.Task_send_endpoint)

	worker_id := commit_worker(send_worker, client_config.Worker_path+client_config.Worker_filename, client_config.Worker_priority)

	return &spbThingMsgHandler{
		Worker_id:          worker_id,
		Task_send_endpoint: client_config.Task_send_endpoint,
		With_body:          true,
		With_body_7:        true,
		Task_path:          client_config.Task_path,
		Task_config_file:   client_config.Task_config_list,
		Ns:                 ns,
	}
}

func (h *spbThingMsgHandler) Handle(msg msgserver.Message) (time.Time, error) {
	task := ParseTask(msg.Body())
	worker_id := h.Worker_id
	//send_task := h.Send_task

	var task_config_path string
	var with_body bool
	if task.ServiceID == 1 {
		task_config_path = h.Task_path + h.Task_config_file[0]
		with_body = h.With_body
	} else if task.ServiceID == 7 {
		task_config_path = h.Task_path + h.Task_config_file[1]
		with_body = h.With_body_7
	} else {
		fmt.Println("[C] Can't find Service ID ", task.ServiceID)
	}
	task_config_file, err := ioutil.ReadFile(task_config_path)
	if err != nil {
		fmt.Println("[C] Open Task Config File Error", err)
	}
	var task_config TaskConfig
	json.Unmarshal(task_config_file, &task_config)

	// construct task body
	var task_file_size int64
	var task_buf []byte
	if with_body {
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

	fromSvr, err := h.Ns.GetServer(idutils.SvrId32(msg.Sender()))
	if err != nil {
		fmt.Println(err)
		return time.Time{}, err
	}
	toSvr, err := h.Ns.GetServer(idutils.SvrId32(msg.Receiver()))
	if err != nil {
		fmt.Println(err)
		return time.Time{}, err
	}

	var task_info TaskPackageInfo
	//copy(task_info.From[:], []byte(task_config.From))
	//copy(task_info.To[:], []byte(task_config.To))
	copy(task_info.From[:], fromSvr.ZMQEndpoint)
	copy(task_info.To[:], toSvr.ZMQEndpoint)
	task_info.Task_sub_id = msg.ID()
	task_info.Task_type = FUNCTION
	task_info.Task_body_id = task_config.Task_body_id
	task_info.Worker_id = worker_id

	cid := idutils.CliId32(msg.ID())
	if cid < 4 {
		task_info.Priority = 0
	} else {
		task_info.Priority = task_config.Priority
	}

	task_info.Timestamp = uint64(C.my_GetTime())
	if with_body {
		task_info.Task_body_size = uint64(task_file_size)
	} else {
		task_info.Task_body_size = uint64(0)
	}
	task_info.Task_args_size = uint64(len(task.Args) + 8)

	mid := idutils.MsgId32(msg.ID())
	svrID := idutils.SvrId32(msg.ID())
	if (1 << 19 & mid) == 0 {
		task_info.Priority = 0
	} else {
		task_info.Priority = 1
	}
	if svrID == 1 {
		//beijing
		if mid%2 == 1 {
			task_info.Task_QoS.End_before = uint64(msg.SendTime().UnixNano() + 50)
		} else {
			task_info.Task_QoS.End_before = uint64(msg.SendTime().UnixNano() + 100)
		}
	} else if svrID == 2 {
		//nanjing
		if mid%2 == 0 {
			task_info.Task_QoS.End_before = uint64(msg.SendTime().UnixNano() + 100)
		} else if mid%4 == 1 {
			task_info.Task_QoS.End_before = uint64(msg.SendTime().UnixNano() + 50)
		} else if mid%4 == 3 {
			task_info.Task_QoS.End_before = uint64(msg.SendTime().UnixNano() + 20)
		}
	} else {
		fmt.Println("error svrID", svrID)
		os.Exit(1)
	}
	//task_info.Task_QoS.End_before = uint64(msg.SendTime().UnixNano() + 50)
	if mid < 200 {
		task_info.Task_QoS.End_before = uint64(msg.SendTime().UnixNano() + 100)
	} else if mid < 300 {
		task_info.Task_QoS.End_before = uint64(msg.SendTime().UnixNano() + 50)
	} else {
		task_info.Task_QoS.End_before = uint64(msg.SendTime().UnixNano() + 20)
	}

	task_info.Task_QoS.Estimate_running_time = task_config.Estimate_running_time
	task_info.Task_body_constraint.Cpu = task_config.Cpu
	task_info.Task_body_constraint.Mem = task_config.Mem
	task_info.Task_body_constraint.Gpu = task_config.Gpu
	task_info.Estimate_strat_time = task_config.Estimate_strat_time
	task_info.Estimate_result_size = task_config.Estimate_result_size
	task_info.Estimate_call_time = task_config.Estimate_call_time

	sockItem, err := czmqutils.GetSock(h.Task_send_endpoint, goczmq.Req)
	if err != nil {
		log.Println("czmq get sock failed: ", err)
		return time.Time{}, nil
	}
	defer sockItem.Free()

	task_header := task_info.Info2Header()

	sendTime, _ := czmqutils.Send(sockItem, task_header, goczmq.FlagMore)
	//Send(send_task, task_header, goczmq.FlagMore)
	if task.ServiceID == 1 {
		if with_body {
			czmqutils.Send(sockItem, task_buf, goczmq.FlagMore)
			//Send(send_task, task_buf, goczmq.FlagMore)
			h.With_body = false
		}
	} else if task.ServiceID == 7 {
		if with_body {
			czmqutils.Send(sockItem, task_buf, goczmq.FlagMore)
			h.With_body_7 = false
		}
		//Send(send_task, task_buf, goczmq.FlagMore)
	} else {
		fmt.Println("[C] Can't find Service ID ", task.ServiceID)
	}

	args := make([]byte, 8+len(task.Args))
	binary.LittleEndian.PutUint64(args[:8], msg.Receiver())
	copy(args[8:], task.Args)

	czmqutils.Send(sockItem, args, goczmq.FlagNone)
	//Send(send_task, task.Args, goczmq.FlagNone)

	Recv(sockItem.Sock)
	//fmt.Println("[C] User Task Sent", string(send_ok_msg[0]))
	return sendTime, nil
}
