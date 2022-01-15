package czmqutils

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/zeromq/goczmq.v4"
	"ict.ac.cn/hbmsgserver/pkg/msgserver"
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

	if _, ok := sockCache[endpoint]; !ok {
		mu.Lock()
		sockCache[endpoint] = make([]*sockItem, 0)
		mu.Unlock()
	}

	socks := sockCache[endpoint]
	for _, s := range socks {
		if s.Type != typ {
			continue
		}

		if s.Active == false {
			mu.Lock()
			if s.Active == true {
				mu.Unlock()
				continue
			}

			s.Active = true
			mu.Unlock()
			sock = s
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
		mu.Lock()
		socks = append(socks, sock)
		mu.Unlock()
	}

	return sock, nil
}

func Send(item *sockItem, data msgserver.Message, flag int) (time.Time, error) {
	sock := item.Sock

	if err := sock.SendFrame(data, flag); err != nil {
		return time.Time{}, errors.Wrap(err, "zmq push Sock send frame failed")
	}

	return data.SendTime(), nil
}
