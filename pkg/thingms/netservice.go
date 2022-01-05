package thingms

import (
	"io"
)

type NetService struct {
	Port   string
	Method string
	Path   func(args []byte) string
	Body   func(args []byte) (io.Reader, string)
}

func NilBody(args []byte) (io.Reader, string) {
	return nil, "application/json"
}

func NilPath(args []byte) string {
	return "/"
}
