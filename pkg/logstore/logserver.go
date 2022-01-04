package logstore

import (
	"encoding/json"
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
		txt, err := json.Marshal(s.LogStore.logs)
		if err != nil {
			fmt.Println("json encode logs failed: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(txt)
	})

	fmt.Printf("log server listen at: %s\n", s.Addr)
	if err := http.ListenAndServe(s.Addr, mux); err != nil {
		panic(err)
	}
}
