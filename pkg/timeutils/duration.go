package timeutils

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Duration struct {
	time.Duration
}

func (i *Duration) String() string {
	return i.Duration.String()
}

func (i *Duration) Set(value string) error {
	if value == "" {
		i.Duration = 0
		return nil
	}

	var err error
	i.Duration, err = time.ParseDuration(value)
	if err != nil {
		return err
	}
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil

	case string:
		if value == "" {
			d.Duration = 0
			return nil
		}

		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil

	default:
		return errors.New("invalid duration")
	}
}
