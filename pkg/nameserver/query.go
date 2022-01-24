package nameserver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"
	"ict.ac.cn/hbmsgserver/pkg/timeutils"
)

func Query(nsEnd string, myID uint64, svrID uint32) *MsgSvr {
	cli := getHttpCli()
	url := fmt.Sprintf("http://%s/query", nsEnd)

	svrIDBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(svrIDBytes, svrID)
	msg := msgserver.NewMessage(uint64(timeutils.GetPtpTime().UnixNano()), myID, 0,
		msgserver.QueryMsg, svrIDBytes)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(msg))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	resp, err := cli.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	resRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	resMsg := msgserver.Message(resRaw)
	if resMsg.Type() != msgserver.ResultMsg {
		fmt.Println("the response from name server is not ResultMsg")
		return nil
	}

	svr := &MsgSvr{}
	if err := json.Unmarshal(resMsg.Body(), svr); err != nil {
		fmt.Println(err)
		return nil
	}
	return svr
}
