package thingms

import (
	"io"
)

type NetService struct {
	Method string
	Path   func(args []byte) string
	Body   func(args []byte) io.Reader
}

func NilBody(args []byte) io.Reader {
	return nil
}

func NilPath(args []byte) string {
	return "/"
}
