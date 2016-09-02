package ari

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

type dtTest struct {
	Date DateTime `json:"dt"`
}

func (dt *dtTest) Equal(o *dtTest) bool {
	return time.Time(dt.Date).Equal(time.Time(o.Date))
}

var dtMarshalTests = []struct {
	Input    dtTest
	Output   string
	HasError bool
}{
	{dtTest{DateTime(time.Date(2005, 02, 04, 13, 12, 6, 0, time.UTC))}, `{"dt":"2005-02-04T13:12:06.000+0000"}`, false},
}

func TestDateTimeMarshal(t *testing.T) {
	for _, tx := range dtMarshalTests {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(tx.Input)
		out := strings.TrimSpace(buf.String())

		failed := false
		failed = failed || (err == nil && tx.HasError)
		failed = failed || (out != tx.Output)

		if failed {
			t.Errorf("Marshal(%s) => '%s', 'err != nil => %v'; expected '%s', 'err != nil => %v'.", tx.Input, out, err != nil, tx.Output, tx.HasError)
		}
	}
}

var dtUnmarshalTests = []struct {
	Input    string
	Output   dtTest
	HasError bool
}{
	{`{"dt":"2005-02-04T13:12:06.000+0000"}`, dtTest{DateTime(time.Date(2005, 02, 04, 13, 12, 6, 0, time.UTC))}, false},
	{`{"dt":"2x05-02-04T13:12:06.000+0000"}`, dtTest{}, true},
}

func TestDateTimeUnmarshal(t *testing.T) {
	for _, tx := range dtUnmarshalTests {
		var out dtTest
		err := json.NewDecoder(strings.NewReader(tx.Input)).Decode(&out)

		failed := false
		failed = failed || (err == nil && tx.HasError)
		failed = failed || (!out.Equal(&tx.Output))

		if failed {
			t.Errorf("Unmarshal(%s) => '%s', 'err != nil => %v'; expected '%s', 'err != nil => %v'.", tx.Input, out, err != nil, tx.Output, tx.HasError)
		}
	}
}
