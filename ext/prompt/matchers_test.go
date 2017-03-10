package prompt

import (
	"testing"
)

func TestMatchHashSuccess(t *testing.T) {
	str, status := MatchHash("Hash#Test")
	if str == "Hash#Test" ||
		status != Complete {
		t.Errorf("Expected status of '%v' and got '%v' and status of '%v'",
			"Complete", str, status)
	}
}

func TestMatchHashFail(t *testing.T) {
	str, status := MatchHash("HashTest")
	if str != "HashTest" ||
		status != Incomplete {
		t.Errorf("Expected status of '%v' and got '%v' and status of '%v'",
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

func TestMatchLenOrTerminatorFuncSuccessTerm(t *testing.T) {
	_, status := MatchLenOrTerminatorFunc(12, "/n")("Test Pass Term /n")
	if status != Complete {
		t.Errorf("Expected status of Complete and received status of '%v'",
			status)
	}
}

func TestMatchLenOrTerminatorFuncSuccessLen(t *testing.T) {
	_, status := MatchLenOrTerminatorFunc(20, "/n")("Length SUccess for check Fail ")
	if status != Complete {
		t.Errorf("Expected status of Incomplete and received status of '%v'",
			status)
	}
}

func TestMatchLenOrTerminatorFuncFail(t *testing.T) {
	_, status := MatchLenOrTerminatorFunc(12, "/n")("Fail len ")
	if status != Incomplete {
		t.Errorf("Expected status of Incomplete and received status of '%v'",
			status)
	}
}
