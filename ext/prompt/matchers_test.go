package prompt

import (
	"testing"
)

func TestMatchHashSuccess(t *testing.T) {
	str, status := MatchHash("Hash#Test")
	if str == "Hash#Test" ||
		status != Complete {
		t.Errorf("Expected slice of 2 and status of '%v' and got '%v' and status of '%v'",
			"Complete", str, status)
	}
}

func TestMatchHashFail(t *testing.T) {
	str, status := MatchHash("HashTest")
	if str != "HashTest" ||
		status != Incomplete {
		t.Errorf("Expected slice of 1 and status of '%v' and got '%v' and status of '%v'",
			"Incomplete", str, status)
	}
}

func TestMatchAnySuccess(t *testing.T) {
	_, status := MatchAny("AnyTest")
	if status != Complete {
		t.Errorf("Expected status of Complete and received status of '%v'",
			status)
	}
}

func TestMatchAnyFail(t *testing.T) {
	_, status := MatchAny("")
	if status != Incomplete {
		t.Errorf("Expected status of Incomplete and received status of '%v'",
			status)
	}
}

func TestMatchTerminatorSuccess(t *testing.T) {
	st := MatchTerminatorFunc("/n")
	_, status := st("Test Terminator/nstring")
	if status != Complete {
		t.Errorf("Expected status of Complete and received status of '%v'",
			status)
	}
}

func TestMatchTerminatorFail(t *testing.T) {
	st := MatchTerminatorFunc("/n")
	_, status := st("Test Terminator string")
	if status != Incomplete {
		t.Errorf("Expected status of Incomplete and received status of '%v'",
			status)
	}
}

func TestMatchLenFuncSuccess(t *testing.T) {
	_, status := MatchLenFunc(14)("Test Match Len")
	if status != Complete {
		t.Errorf("Expected status of Complete and received status of '%v'",
			status)
	}
}

func TestMatchLenFuncFail(t *testing.T) {
	_, status := MatchLenFunc(12)("Len Fail")
	if status != Incomplete {
		t.Errorf("Expected status of Incomplete and received status of '%v'",
			status)
	}
}

func TestMatchLenOrTerminatorFuncSuccess(t *testing.T) {
	_, status := MatchLenOrTerminatorFunc(12, "/n")("Test Match Len Terminator /n")
	if status != Complete {
		t.Errorf("Expected status of Complete and received status of '%v'",
			status)
	}
}

func TestMatchLenOrTerminatorFuncFailLen(t *testing.T) {
	_, status := MatchLenOrTerminatorFunc(20, "/n")("Ln Fail /n")
	if status != Incomplete {
		t.Errorf("Expected status of Incomplete and received status of '%v'",
			status)
	}
}

func TestMatchLenOrTerminatorFuncFailTerm(t *testing.T) {
	_, status := MatchLenOrTerminatorFunc(12, "/n")("Terminator Fail len succeeds")
	if status != Incomplete {
		t.Errorf("Expected status of Incomplete and received status of '%v'",
			status)
	}
}
