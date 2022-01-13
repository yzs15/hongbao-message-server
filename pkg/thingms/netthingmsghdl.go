package thingms

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/registry"

	"ict.ac.cn/hbmsgserver/pkg/logstore"

	"ict.ac.cn/hbmsgserver/pkg/nameserver"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/czmqutils"

	"github.com/pkg/errors"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"
)

type netThingMsgHandler struct {
	Me uint32

	KubeEndpoints []string
	Services      map[uint8]*NetService

	httpCli *http.Client

	Registry   *registry.Registry
	NameServer *nameserver.NameServer

	logStore *logstore.LogStore
}

func NewNetThingMsgHandler(
	me uint32, kubeEnds []string, svs map[uint8]*NetService,
	registry *registry.Registry, ns *nameserver.NameServer,
	logStore *logstore.LogStore,
) TaskMsgHandler {
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
		Me:            me,
		KubeEndpoints: kubeEnds,
		Services:      svs,
		httpCli:       cli,

		Registry:   registry,
		NameServer: ns,

		logStore: logStore,
	}
}

func (h *netThingMsgHandler) Handle(msg msgserver.Message) (time.Time, error) {
	task := ParseTask(msg.Body())

	result, err, reqTime, respTime := h.httpReq(msg.ID(), task.ServiceID, task.Args)
	if err != nil {
		return time.Time{}, err
	}
	h.logStore.Add(msg.ID(), uint64(h.Me), idutils.DeviceId(h.Me, 1<<19), uint64(h.Me), reqTime, logstore.Send)
	h.logStore.Add(msg.ID(), uint64(h.Me), idutils.DeviceId(h.Me, 1<<19), uint64(h.Me), respTime, logstore.Receive)

	resMsg := msgserver.NewMessage(msg.ID(), msg.Sender(), msg.Receiver(),
		msgserver.TextMsg, result)

	var sendTime time.Time
	svrID := idutils.SvrId32(msg.Receiver())
	if svrID == h.Me { // 接受者就在自己链接的客户端内
		sendTime, err = h.Registry.Send(resMsg)
		if err != nil {
			return time.Time{}, err
		}

	} else { // 接受者在另一个 Message Server 内
		svr, err := h.NameServer.GetServer(svrID)
		if err != nil {
			return time.Time{}, err
		}
		if sendTime, err = czmqutils.Send(svr.ZMQEndpoint, resMsg); err != nil {
			return time.Time{}, err
		}
	}

	return sendTime, nil
}

func (h *netThingMsgHandler) httpReq(mid uint64, svcId uint8, args []byte) (result []byte, err error, reqTime, respTime time.Time) {
	service := h.Services[svcId]

	path := service.Path(mid, args)
	body, contentType := service.Body(args)
	endpoint := h.KubeEndpoints[rand.Intn(len(h.KubeEndpoints))]
	url := fmt.Sprintf("http://%s:%s%s", endpoint, service.Port, path)

	req, err := http.NewRequest(service.Method, url, body)
	if err != nil {
		err = errors.Wrap(err, "create http request failed")
		return
	}
	req.Header.Set("Content-Type", contentType)

	reqTime = time.Now()
	resp, err := h.httpCli.Do(req)
	respTime = time.Now()
	if err != nil {
		err = errors.Wrap(err, "do http request failed")
		return
	}
	defer resp.Body.Close()

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "read response body failed")
		return
	}
	return
}
