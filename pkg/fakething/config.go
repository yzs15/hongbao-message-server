package fakething

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"ict.ac.cn/hbmsgserver/pkg/timeutils"
)

type Config struct {
	MacAddr   string
	MsgWsEnd  []string
	MsgZmqEnd []string

	Mode Mode

	LoadNumPer int
	NoisNumPer int
	Period timeutils.Duration

	NumConn   int
	TotalTime timeutils.Duration
	PeakTime  timeutils.Duration
	PeakNum   int
}

func GetConfig(filename string) (Config, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		err = errors.Wrap(err, "read config file failed")
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(raw, &config); err != nil {
		err = errors.Wrap(err, "json decode config failed")
		return Config{}, err
	}

	return config, nil
}
