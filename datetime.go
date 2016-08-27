package ari

import (
	"encoding/json"
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
		//	Logger.Error("Failed to unmarshal asterisk timestamp", err)
		//	fmt.Printf("data: %s\nstringDate: %s\n", data, stringDate)
		return err
	}

	t, err := time.Parse(DateFormat, stringDate)
	if err != nil {
		//	Logger.Error("Failed to parse asterisk date format", stringDate)
		return err
	}
	*dt = (DateTime)(t)
	return nil
}

func (dt DateTime) String() string {
	t := (time.Time)(dt)
	return t.String()
}
