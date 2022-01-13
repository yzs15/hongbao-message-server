package thingms

import (
	"fmt"
	"io"
)

type NetService struct {
	Port   string
	Method string
	Path   func(mid uint64, args []byte) string
	Body   func(args []byte) (io.Reader, string)
}

func NilBody(args []byte) (io.Reader, string) {
	return nil, "application/json"
}

func NilPath(mid uint64, args []byte) string {
	return fmt.Sprintf("/?mid=%d", mid)
}
