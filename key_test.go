package ari

import "testing"

func TestKeyMatch(t *testing.T) {

	// two empty keys should match
	ok := NewKey("", "").Match(NewKey("", ""))
	if !ok {
		t.Errorf("Two empty keys should match")
	}

	ok = AppKey("app").Match(NewKey("", ""))
	if !ok {
		t.Errorf("App key should match any subkey")
	}

	ok = AppKey("app").Match(AppKey("app2"))
	if ok {
		t.Errorf("Two separate app keys should not match")
	}

	ok = NodeKey("app", "node").Match(NewKey("", ""))
	if !ok {
		t.Errorf("Node key should match any subkey")
	}

	ok = NewKey("application", "id1").Match(NewKey("application", "id1"))
	if !ok {
		t.Errorf("Application/id1 should match")
	}

	ok = NewKey("application", "").Match(NewKey("application", "id1"))
	if !ok {
		t.Errorf("Application/* should match application/id")
	}
}
