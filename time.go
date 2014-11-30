package ari

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"
)

// Asterisk ARI does not supply a format that the built-in
// golang json date parser can use, so we have to write our
// own unmarshal routine for it
type AsteriskDate time.Time

const AsteriskDateFormat = "2006-01-02T15:04:05.000-0700"

func (d AsteriskDate) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	a := []byte("\"" + t.Format(AsteriskDateFormat) + "\"")
	return a, nil
}
func (d *AsteriskDate) UnmarshalJSON(data []byte) error {
	var stringDate string
	err := json.Unmarshal(data, &stringDate)
	if err != nil {
		glog.Errorln("Failed to unmarshal asterisk timestamp", err)
		fmt.Printf("data: %s\nstringDate: %s\n", data, stringDate)
		return err
	}

	t, err := time.Parse(AsteriskDateFormat, stringDate)
	if err != nil {
		glog.Errorln("Failed to parse asterisk date format", stringDate)
		return err
	}
	*d = (AsteriskDate)(t)
	return nil
}

func (d AsteriskDate) String() string {
	t := (time.Time)(d)
	return t.String()
}
