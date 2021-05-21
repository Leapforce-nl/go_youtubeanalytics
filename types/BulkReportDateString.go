package youtube

import (
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	errortools "github.com/leapforce-libraries/go_errortools"
)

const (
	layout      string = "20060102"
	layoutParse string = "2006-01-02"
)

type BulkReportDateString civil.Date

func (d *BulkReportDateString) UnmarshalJSON(b []byte) error {

	var returnError = func() error {
		errortools.CaptureError(fmt.Sprintf("Cannot parse '%s' to BulkReportDateString", string(b)))
		return nil
	}

	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println(err)
		return returnError()
	}

	if len(s) != len(layout) {
		return returnError()
	}

	if s == "" || s == "00000000" {
		d = nil
		return nil
	}

	_t, err := time.Parse(layoutParse, s[:4]+"-"+s[4:6]+"-"+s[6:])
	if err != nil {
		return err
	}

	*d = BulkReportDateString(civil.DateOf(_t))
	return nil
}

func (d *BulkReportDateString) ValuePtr() *civil.Date {
	if d == nil {
		return nil
	}

	_d := civil.Date(*d)
	return &_d
}

func (d BulkReportDateString) Value() civil.Date {
	return civil.Date(d)
}
