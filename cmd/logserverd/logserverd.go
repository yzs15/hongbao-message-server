package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type StringArrFlag []string

func (i *StringArrFlag) String() string {
	return strings.Join(*i, ",")
}

func (i *StringArrFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var Addr = flag.String("addr", "0.0.0.0:5552", "log service address")
var logFilenames StringArrFlag

func main() {
	flag.Var(&logFilenames, "f", "the path to log file, if directory, read all files")
	flag.Parse()

	if len(logFilenames) == 0 {
		fmt.Println("need at least one log file, specify by '-f'")
		return
	}

	allExist := true
	for _, filename := range logFilenames {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			allExist = false
			fmt.Printf("file not exist: %s\n", filename)
		}
	}
	if !allExist {
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, logFilename := range logFilenames {
			if stat, _ := os.Stat(logFilename); stat.IsDir() {
				files, err := ioutil.ReadDir(logFilename)
				if err != nil {
					fmt.Printf("read log dir failed: %s", logFilename)
					continue
				}

				for _, file := range files {
					if file.IsDir() {
						continue
					}
					completePath := fmt.Sprintf("%s/%s", logFilename, file.Name())

					log, err := ioutil.ReadFile(completePath)
					if err != nil {
						fmt.Printf("read log file failed: %s\n", completePath)
						continue
					}

					if len(log) > 0 {
						w.Write(log)
						w.Write([]byte("\n"))
					}
				}

			} else {
				log, err := ioutil.ReadFile(logFilename)
				if err != nil {
					fmt.Printf("read log file failed: %s\n", logFilename)
					continue
				}

				if len(log) > 0 {
					w.Write(log)
					w.Write([]byte("\n"))
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("log server listen at: %s\n", *Addr)
	if err := http.ListenAndServe(*Addr, mux); err != nil {
		panic(err)
	}
}
