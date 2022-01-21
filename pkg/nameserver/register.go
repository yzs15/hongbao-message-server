package nameserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/linuxutils"
	"ict.ac.cn/hbmsgserver/pkg/msgserver"
)

func Register(nsEnd, zmqEnd, mac string) uint32 {
	cli := getHttpCli()

	if mac == "" {
		var err error
		mac, err = linuxutils.GetMac()
		if err != nil {
			panic(err)
		}
	}

	msgSvr := &MsgSvr{ZMQEndpoint: zmqEnd, Mac: mac}
	body, err := json.Marshal(msgSvr)
	if err != nil {
		panic(err)
	}

	msg := msgserver.NewMessage(uint64(time.Now().UnixNano()), 0, 0,
		msgserver.RegisterMsg, body)

	url := fmt.Sprintf("http://%s/register", nsEnd)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
	if err != nil {
		panic(err)
	}

	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	resRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	resMsg := msgserver.Message(resRaw)
	if resMsg.Type() != msgserver.NameMsg {
		panic("the response from name server is not NameMsg")
	}

	return uint32(resMsg.Receiver())
}
