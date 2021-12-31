package thingms

import (
	"fmt"
	"io"
)

type NetService struct {
	Method string
	Query  string
	File   string
}

func (s *NetService) getPath(query string) string {
	if len(s.Query) > 0 {
		return fmt.Sprintf("/?%s=%s", s.Query, query)
	}
	return "/"
}

func (s *NetService) getBody(file string) io.Reader {
	return nil
}
