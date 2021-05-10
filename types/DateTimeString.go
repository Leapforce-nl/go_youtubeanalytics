package youtube

import (
	"encoding/json"
	"fmt"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

const (
	layout string = "2006-01-02T15:04:05Z"
)

type DateTimeString time.Time

func (d *DateTimeString) UnmarshalJSON(b []byte) error {
	var returnError = func() error {
		errortools.CaptureError(fmt.Sprintf("Cannot parse '%s' to DateTimeString", string(b)))
		return nil
	}

	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return returnError()
	}

	if s == "" || s == "0000-00-00 00:00:00" {
		d = nil
		return nil
	}

	_t, err := time.Parse(layout, s)
	if err != nil {
		return returnError()
	}

	*d = DateTimeString(_t)
	return nil
}

func (d *DateTimeString) MarshalJSON() ([]byte, error) {
	if d == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(time.Time(*d).Format(layout))
}

func (d *DateTimeString) ValuePtr() *time.Time {
	if d == nil {
		return nil
	}

	_d := time.Time(*d)
	return &_d
}

func (d DateTimeString) Value() time.Time {
	return time.Time(d)
}

func (d *DateTimeString) String() string {
	if d == nil {
		return ""
	}

	return (time.Time(*d)).Format(layout)
}
