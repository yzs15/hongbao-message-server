package czmqutils

import (
	"sync"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/timeutils"

	"gopkg.in/zeromq/goczmq.v4"
)

type sockItem struct {
	Sock   *goczmq.Sock
	Type   int
	Active bool
}

func (i *sockItem) Free() {
	i.Active = false
	<-concurrencyCtl
}

var mu sync.Mutex
var sockCache map[string][]*sockItem

var concurrencyCtl chan struct{}

func init() {
	sockCache = make(map[string][]*sockItem)
	concurrencyCtl = make(chan struct{}, 1000)
}

func createSocketItem(endpoint string, typ int) *sockItem {
	var s *goczmq.Sock
	for {
		s = goczmq.NewSock(typ)
		err := s.Connect(endpoint)
		if err != nil {
			continue
		}
		break
	}

	return &sockItem{
		Sock:   s,
		Type:   typ,
		Active: false,
	}
}

func GetSock(endpoint string, typ int) (*sockItem, error) {
	concurrencyCtl <- struct{}{}

	var sock *sockItem = nil
	mu.Lock()
	defer mu.Unlock()

	if _, ok := sockCache[endpoint]; !ok {
		sock = createSocketItem(endpoint, typ)
		sock.Active = true

		socks := make([]*sockItem, 0)
		socks = append(socks, sock)
		sockCache[endpoint] = socks

		return sock, nil
	}

	socks := sockCache[endpoint]
	for _, s := range socks {
		if s.Type != typ {
			continue
		}

		if s.Active == false {
			s.Active = true
			sock = s
			break
		}
	}

	if sock == nil {
		sock = createSocketItem(endpoint, typ)
		sock.Active = true

		socks = append(socks, sock)
		sockCache[endpoint] = socks
	}

	return sock, nil
}

func Send(item *sockItem, data []byte, flag int) (time.Time, error) {
	sock := item.Sock

	var sendTime time.Time
	for {
		sendTime = timeutils.GetPtpTime()
		if err := sock.SendFrame(data, flag); err != nil {
			continue
			//return time.Time{}, errors.Wrap(err, "zmq push Sock send frame failed")
		}
		break
	}

	return sendTime, nil
}
