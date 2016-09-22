package ari

import (
	"encoding/json"
	"strconv"
	"time"
)

// NOTE: near direct translation from ARI 2.0

// DateFormat is the date format that ARI returns in the JSON bodies
const DateFormat = "2006-01-02T15:04:05.000-0700"

// DateTime is an alias type for attaching a custom
// asterisk unmarshaller and marshaller for JSON
type DateTime time.Time

// MarshalJSON converts the given date object to ARIs date format
func (dt DateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(dt)
	a := []byte("\"" + t.Format(DateFormat) + "\"")
	return a, nil
}

// UnmarshalJSON parses the given date per ARIs date format
func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var stringDate string
	err := json.Unmarshal(data, &stringDate)
	if err != nil {
		return err
	}

	t, err := time.Parse(DateFormat, stringDate)
	if err != nil {
		return err
	}
	*dt = (DateTime)(t)
	return nil
}

func (dt DateTime) String() string {
	t := (time.Time)(dt)
	return t.String()
}

// Duration support functions

// DurationSec is a JSON type for duration in seconds
type DurationSec time.Duration

// MarshalJSON converts the duration into a JSON friendly format
func (ds DurationSec) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(time.Duration(ds) / time.Second))), nil
}

// UnmarshalJSON parses the data into the duration seconds object
func (ds *DurationSec) UnmarshalJSON(data []byte) error {
	s, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}

	*ds = DurationSec(time.Duration(s) * time.Second)
	return nil
}
