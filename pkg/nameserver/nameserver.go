package nameserver

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type MsgSvr struct {
	Mac         string
	ZMQEndpoint string
}

type NameServer struct {
	id2server map[uint32]*MsgSvr

	Me    uint32
	NsEnd string
}

func NewNameServer(nsEnd string, me uint32) *NameServer {
	return &NameServer{
		id2server: make(map[uint32]*MsgSvr),
		Me:        me,
		NsEnd:     nsEnd,
	}
}

func (r *NameServer) GetServer(id uint32) (*MsgSvr, error) {
	var svr *MsgSvr = nil
	var ok bool
	if svr, ok = r.id2server[id]; !ok {
		svr = Query(r.NsEnd, uint64(r.Me), id)
		if svr != nil {
			r.id2server[id] = svr
		}
	}

	if svr == nil {
		return nil, errors.New(fmt.Sprintf("no this message server: %d", id))
	}
	return svr, nil
}

func (r *NameServer) GetAllServer() ([]*MsgSvr, []uint32) {
	servers := make([]*MsgSvr, 0)
	ids := make([]uint32, 0)
	for k, v := range r.id2server {
		servers = append(servers, v)
		ids = append(ids, k)
	}
	return servers, ids
}

func getHttpCli() *http.Client {
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
	return cli
}
