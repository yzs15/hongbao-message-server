package czmqutils

import (
	"github.com/pkg/errors"
	"gopkg.in/zeromq/goczmq.v4"
)

var socks map[string]*goczmq.Sock

func init() {
	socks = make(map[string]*goczmq.Sock)
}

func getSock(endpoint string) (*goczmq.Sock, error) {
	var sock *goczmq.Sock
	var ok bool
	var err error

	if sock, ok = socks[endpoint]; !ok {
		sock, err = goczmq.NewPush(endpoint)
		if err != nil {
			return nil, errors.Wrap(err, "create zmq push sock failed")
		}
		socks[endpoint] = sock
	}

	return sock, nil
}

func Send(endpoint string, data []byte) error {
	sock, err := getSock(endpoint)
	if err != nil {
		return err
	}

	if err := sock.SendFrame(data, goczmq.FlagNone); err != nil {
		return errors.Wrap(err, "zmq push sock send frame failed")
	}
	return nil
}
