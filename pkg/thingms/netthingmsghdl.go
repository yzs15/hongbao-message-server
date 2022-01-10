package thingms

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/nameserver"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"

	"github.com/pkg/errors"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"
)

type netThingMsgHandler struct {
	KubeEndpoints []string
	Services      map[uint8]*NetService

	httpCli *http.Client

	NameServer *nameserver.NameServer
}

func NewNetThingMsgHandler(kubeEnds []string, svs map[uint8]*NetService, ns *nameserver.NameServer) TaskMsgHandler {
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
		Services:      svs,
		httpCli:       cli,

		NameServer: ns,
	}
}

func (h *netThingMsgHandler) Handle(msg msgserver.Message) (time.Time, error) {
	task := ParseTask(msg.Body())

	result, err := h.httpReq(task.ServiceID, task.Args)
	if err != nil {
		return time.Time{}, err
	}

	resMsg := msgserver.NewMessage(msg.ID(), msg.Sender(), msg.Receiver(),
		msgserver.TextMsg, []byte(result))

	svrID := idutils.SvrId32(msg.Receiver())
	svr, err := h.NameServer.GetServer(svrID)
	if err != nil {
		return time.Time{}, err
	}
	var sendTime time.Time
	if sendTime, err = czmqutils.Send(svr.ZMQEndpoint, resMsg); err != nil {
		return time.Time{}, err
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
