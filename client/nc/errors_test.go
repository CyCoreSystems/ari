package nc

import (
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
)

var errA = errors.New("Error A")
var errB = errors.New("Error B")

var errC = errors.Wrap(errB, "This is a long parent message")

var errD = errors.Wrap(errC, "Another Parent")

var errF = &codedError{errors.New("Error F"), 200}
var errG = errors.Wrap(errF, "Parent")

var errorToMapTests = []struct {
	Input  error
	Output string
}{
	{errA, `{"message":"Error A"}`},
	{errC, `{"cause":{"message":"Error B"},"message":"This is a long parent message: Error B"}`},
	{errD, `{"cause":{"cause":{"message":"Error B"},"message":"This is a long parent message: Error B"},"message":"Another Parent: This is a long parent message: Error B"}`},
	{errF, `{"code":200,"message":"Error F"}`},
	{errG, `{"cause":{"code":200,"message":"Error F"},"message":"Parent: Error F"}`},
}

func TestErrorToMap(t *testing.T) {

	for _, tx := range errorToMapTests {
		m := ErrorToMap(tx.Input, "")
		body, err := json.Marshal(m)

		failed := err != nil
		failed = failed || string(body) != tx.Output
		if failed {
			t.Errorf("json.Marshal(errorToMap('%v')) => ('%v', '%v'), expected ('%v', '%v')",
				tx.Input,
				string(body), err,
				tx.Output, nil)

		}
	}

}

var mapToErrorTests = []struct {
	Output error
	Input  string
}{
	{errA, `{"message":"Error A"}`},
	{errC, `{"cause":{"message":"Error B"},"message":"This is a long parent message: Error B"}`},
	{errD, `{"cause":{"cause":{"message":"Error B"},"message":"This is a long parent message: Error B"},"message":"Another Parent: This is a long parent message: Error B"}`},
	{errF, `{"code":200,"message":"Error F"}`},
	{errG, `{"cause":{"code":200,"message":"Error F"},"message":"Parent: Error F"}`},
}

func TestMapToError(t *testing.T) {

	for _, tx := range mapToErrorTests {
		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(tx.Input), &m)

		output := MapToError(m)

		failed := err != nil
		failed = failed || output.Error() != tx.Output.Error()
		if failed {
			t.Errorf("json.Unmarshal, mapToError('%v')) => ('%v', '%v'), expected ('%v', '%v')",
				tx.Input,
				output.Error(), err,
				tx.Output, nil)

		}
	}

}
