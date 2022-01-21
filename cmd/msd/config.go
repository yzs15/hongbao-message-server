package main

import (
	"encoding/json"
	"io/ioutil"
)

const (
	SpbEnv = "spb"
	NetEnv = "net"
)

type Config struct {
	Env string

	WsAddr    string
	ZmqAddr   string
	MacAddr   string
	ZmqOutEnd string

	NsEnd string

	SpbConfig string
	KubeEnds  []string

	LogPath string
}

func ParseConfig(filename string) *Config {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	conf := &Config{}
	if err := json.Unmarshal(text, conf); err != nil {
		panic(err)
	}

	return conf
}
