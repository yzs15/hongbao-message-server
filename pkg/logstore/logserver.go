package logstore

import (
	"fmt"
	"net/http"
)

type LogServer struct {
	Addr string

	LogStore *LogStore
}

func (s *LogServer) Run() {
	mux := http.NewServeMux()

	// TODO develop the handle to get all logs
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("not implemented")
	})

	fmt.Printf("log server listen at: %s\n", s.Addr)
	if err := http.ListenAndServe(s.Addr, mux); err != nil {
		panic(err)
	}
}
