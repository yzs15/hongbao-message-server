package czmqutils

import (
	"sync"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/timeutils"

	"github.com/pkg/errors"
	"gopkg.in/zeromq/goczmq.v4"
)

type sockItem struct {
	Sock   *goczmq.Sock
	Type   int
	Active bool
}

func (i *sockItem) Free() {
	i.Active = false
}

var mu sync.Mutex
var sockCache map[string][]*sockItem

func init() {
	sockCache = make(map[string][]*sockItem)
}

func GetSock(endpoint string, typ int) (*sockItem, error) {
	var sock *sockItem = nil
	mu.Lock()
	defer mu.Unlock()

	if _, ok := sockCache[endpoint]; !ok {
		sockCache[endpoint] = make([]*sockItem, 0)
	}

	socks := sockCache[endpoint]
	for _, s := range socks {
		if s.Type != typ {
			continue
		}

		if s.Active == false {
			if s.Active == true {
				continue
			}

			s.Active = true
			sock = s
			break
		}
	}

	if sock == nil {
		s := goczmq.NewSock(typ)
		err := s.Connect(endpoint)
		if err != nil {
			return nil, errors.Wrap(err, "create zmq push Sock failed")
		}

		sock = &sockItem{
			Sock:   s,
			Type:   typ,
			Active: true,
		}
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
