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

	ok = NewKey("channel", "id1").Match(NewKey("channel", "id2"))
	if ok {
		t.Errorf("Differing IDs should not match")
	}
}

func TestKeysFilter(t *testing.T) {
	keys := Keys{
		NewKey(ApplicationKey, "app1"),
		NewKey(ChannelKey, "ch1"),
		NewKey(BridgeKey, "br1"),
	}

	newKeys := keys.Filter(KindKey(ApplicationKey))

	if len(newKeys) != 1 {
		t.Errorf("Expected filters keys by app to be of length 1, got %d", len(newKeys))
	} else {
		if newKeys[0].Kind != ApplicationKey && newKeys[0].ID != "app1" {
			t.Errorf("Unexpected first index %v", newKeys[0])
		}
	}

	newKeys = keys.Without(KindKey(ApplicationKey))

	if len(newKeys) != 2 {
		t.Errorf("Expected without keys by app to be of length 2, got %d", len(newKeys))
	} else {
		if newKeys[0].Kind != ChannelKey && newKeys[0].ID != "ch1" {
			t.Errorf("Unexpected first index %v", newKeys[0])
		}
		if newKeys[1].Kind != BridgeKey && newKeys[1].ID != "br1" {
			t.Errorf("Unexpected second index %v", newKeys[1])
		}
	}

	newKeys = keys.Filter(KindKey(ChannelKey), KindKey(BridgeKey))

	if len(newKeys) != 2 {
		t.Errorf("Expected without keys by app to be of length 2, got %d", len(newKeys))
	} else {
		if newKeys[0].Kind != ChannelKey && newKeys[0].ID != "ch1" {
			t.Errorf("Unexpected first index %v", newKeys[0])
		}
		if newKeys[1].Kind != BridgeKey && newKeys[1].ID != "br1" {
			t.Errorf("Unexpected second index %v", newKeys[1])
		}
	}

	newKeys = keys.Filter(KindKey(ChannelKey), KindKey(BridgeKey)).Without(MatchFunc(func(k *Key) bool {
		return k.ID == "br1"
	}))

	if len(newKeys) != 1 {
		t.Errorf("Expected without keys by app to be of length 2, got %d", len(newKeys))
	} else {
		if newKeys[0].Kind != ChannelKey && newKeys[0].ID != "ch1" {
			t.Errorf("Unexpected first index %v", newKeys[0])
		}
	}
}
