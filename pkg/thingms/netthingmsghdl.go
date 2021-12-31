package thingms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"

	"github.com/pkg/errors"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"
)

type netThingMsgHandler struct {
	KubeEndpoints []string
	WangEndpoint  string

	Services map[string]*NetService

	httpCli *http.Client
}

func NewNetThingMsgHandler(kubeEnds []string, wangEnd string, svs map[string]*NetService) ThingMsgHandler {
	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
			MaxIdleConnsPerHost:   -1,
		},
	}

	return &netThingMsgHandler{
		KubeEndpoints: kubeEnds,
		WangEndpoint:  wangEnd,
		Services:      svs,
		httpCli:       cli,
	}
}

func (h *netThingMsgHandler) Handle(msg *msgserver.Message) (time.Time, error) {
	content := msg.Content

	taskArgs := &Task{}
	if err := json.Unmarshal([]byte(content), taskArgs); err != nil {
		return time.Time{}, errors.Wrap(err, "decode task failed")
	}

	result, err := h.httpReq(taskArgs)
	if err != nil {
		return time.Time{}, err
	}

	resMsg := &msgserver.Message{
		ID:      msg.ID,
		Sender:  msg.Sender,
		Good:    msg.Good,
		Content: result,
	}
	resRaw, err := json.Marshal(resMsg)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "encode result failed")
	}

	sendTime := time.Now()
	if err := czmqutils.Send(h.WangEndpoint, resRaw); err != nil {
		return time.Time{}, errors.Wrap(err, "czmq send result to Wang failed")
	}
	return sendTime, nil
}

func (h *netThingMsgHandler) httpReq(taskArgs *Task) (string, error) {
	service := h.Services[taskArgs.Name]

	path := service.getPath(taskArgs.Query)
	body := service.getBody(taskArgs.File)
	endpoint := h.KubeEndpoints[rand.Intn(len(h.KubeEndpoints))]
	url := fmt.Sprintf("http://%s%s", endpoint, path)

	req, err := http.NewRequest(service.Method, url, body)
	if err != nil {
		return "", errors.Wrap(err, "create http request failed")
	}

	resp, err := h.httpCli.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "do http request failed")
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read response body failed")
	}
	return string(result), nil
}
