package ari

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// test data

var dtMarshalTests = []struct {
	Input    dtTest
	Output   string
	HasError bool
}{
	{dtTest{DateTime(time.Date(2005, 0o2, 0o4, 13, 12, 6, 0, time.UTC))}, `{"dt":"2005-02-04T13:12:06.000+0000"}`, false},
}

var dtUnmarshalTests = []struct {
	Input    string
	Output   dtTest
	HasError bool
}{
	{`{"dt":"2005-02-04T13:12:06.000+0000"}`, dtTest{DateTime(time.Date(2005, 0o2, 0o4, 13, 12, 6, 0, time.UTC))}, false},
	{`{"dt":"2x05-02-04T13:12:06.000+0000"}`, dtTest{}, true},
	{`{"dt": 0 }`, dtTest{}, true},
}

var dsUnmarshalTests = []struct {
	Input    string
	Output   dsTest
	HasError bool
}{
	{`{"ds":4}`, dsTest{DurationSec(4 * time.Second)}, false},
	{`{"ds":40}`, dsTest{DurationSec(40 * time.Second)}, false},
	{`{"ds":"4"}`, dsTest{}, true},
	{`{"ds":""}`, dsTest{}, true},
	{`{"ds":"xzsad"}`, dsTest{}, true},
}

var dsMarshalTests = []struct {
	Input    dsTest
	Output   string
	HasError bool
}{
	{dsTest{DurationSec(4 * time.Second)}, `{"ds":4}`, false},
	{dsTest{DurationSec(40 * time.Second)}, `{"ds":40}`, false},
}

// test runners
func TestDateTimeMarshal(t *testing.T) {
	for _, tx := range dtMarshalTests {
		ret := runTestMarshal(tx.Input, tx.Output, tx.HasError)
		if ret != "" {
			t.Errorf(ret)
		}
	}
}

func TestDateTimeUnmarshal(t *testing.T) {
	for _, tx := range dtUnmarshalTests {
		var out dtTest

		ret := runTestUnmarshal(&out, tx.Input, &tx.Output, tx.HasError)
		if ret != "" {
			t.Errorf(ret)
		}
	}
}

func TestDurationSecsMarshal(t *testing.T) {
	for _, tx := range dsMarshalTests {
		ret := runTestMarshal(tx.Input, tx.Output, tx.HasError)
		if ret != "" {
			t.Errorf(ret)
		}
	}
}

func TestDurationSecsUnmarshal(t *testing.T) {
	for _, tx := range dsUnmarshalTests {
		var out dsTest
		if ret := runTestUnmarshal(&out, tx.Input, &tx.Output, tx.HasError); ret != "" {
			t.Errorf(ret)
		}
	}
}

// generalized test functions

func runTestMarshal(input interface{}, output string, hasError bool) (ret string) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(input)
	out := strings.TrimSpace(buf.String())

	failed := false
	failed = failed || (err == nil && hasError)
	failed = failed || (out != output)

	if failed {
		ret = fmt.Sprintf("Marshal(%s) => '%s', 'err != nil => %v'; expected '%s', 'err != nil => %v'.", input, out, err != nil, output, hasError)
	}

	return
}

func runTestUnmarshal(out eq, input string, output eq, hasError bool) (ret string) {
	err := json.NewDecoder(strings.NewReader(input)).Decode(&out)

	failed := false
	failed = failed || (err == nil && hasError)
	failed = failed || (!out.Equal(output))

	if failed {
		ret = fmt.Sprintf("Unmarshal(%s) => '%s', 'err != nil => %v'; expected '%s', 'err != nil => %v'.", input, out, err != nil, output, hasError)
	}

	return
}

// test structures

type dtTest struct {
	Date DateTime `json:"dt"`
}

func (dt *dtTest) Equal(i interface{}) bool {
	o, ok := i.(*dtTest)
	return ok && time.Time(dt.Date).Equal(time.Time(o.Date))
}

type eq interface {
	Equal(i interface{}) bool
}

type dsTest struct {
	Duration DurationSec `json:"ds"`
}

func (ds *dsTest) Equal(i interface{}) bool {
	o, ok := i.(*dsTest)
	return ok && time.Duration(ds.Duration) == time.Duration(o.Duration)
}
