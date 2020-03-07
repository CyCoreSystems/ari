package ari

import "testing"

func TestNullSubscription(t *testing.T) {
	sub := NewNullSubscription()

	select {
	case <-sub.Events():
		t.Error("received event from NullSubscription")
	default:
	}

	sub.Cancel()

	select {
	case <-sub.Events():
	default:
		t.Error("NullSubscription failed to close")
	}

	// Make sure subsequent Cancel doesn't break
	sub.Cancel()

	select {
	case <-sub.Events():
	default:
		t.Error("NullSubscription failed to close")
	}
}
