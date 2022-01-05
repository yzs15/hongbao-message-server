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

	Services map[uint8]*NetService

	httpCli *http.Client
}

func NewNetThingMsgHandler(kubeEnds []string, wangEnd string, svs map[uint8]*NetService) ThingMsgHandler {
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
		Timeout: 3 * time.Second,
	}

	return &netThingMsgHandler{
		KubeEndpoints: kubeEnds,
		WangEndpoint:  wangEnd,
		Services:      svs,
		httpCli:       cli,
	}
}

func (h *netThingMsgHandler) Handle(task *Task) (time.Time, error) {
	result, err := h.httpReq(task.ServiceID, task.Args)
	if err != nil {
		return time.Time{}, err
	}

	resMsg := &msgserver.Message{
		ID:      task.ID,
		Sender:  task.Sender,
		Good:    task.Good,
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

func (h *netThingMsgHandler) httpReq(svcId uint8, args []byte) (string, error) {
	service := h.Services[svcId]

	path := service.Path(args)
	body, contentType := service.Body(args)
	endpoint := h.KubeEndpoints[rand.Intn(len(h.KubeEndpoints))]
	url := fmt.Sprintf("http://%s:%s%s", endpoint, service.Port, path)

	req, err := http.NewRequest(service.Method, url, body)
	if err != nil {
		return "", errors.Wrap(err, "create http request failed")
	}
	req.Header.Set("Content-Type", contentType)

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
