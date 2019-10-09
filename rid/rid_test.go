package rid

import (
	"strings"
	"testing"
	"time"
)

func TestGeneric(t *testing.T) {
	a := New("")
	t.Logf("Generated ID: (%s)", a)

	b := New("")
	t.Logf("Generated ID: (%s)", b)

	if a == b {
		t.Errorf("consecutive IDs do not differ (%s) (%s)", a, b)
		return
	}
}

func TestKind(t *testing.T) {
	a := New(Channel)
	t.Logf("Generated channel ID: (%s)", a)

	if !strings.HasSuffix(a, "-ch") {
		t.Errorf("failed to apply proper resource suffix (%s); should have be -ch", a)
	}
}

func TestKindClipping(t *testing.T) {
	a := New("hello")
	t.Logf("Generated hello ID: (%s)", a)

	if !strings.HasSuffix(a, "-he") {
		t.Errorf("failed to apply proper resource suffix (%s); should have been -he", a)
	}
}

func TestTimestamp(t *testing.T) {
	a := New(Channel)

	ts, err := Timestamp(a)
	if err != nil {
		t.Error("failed to parse channel resource ID", err)
		return
	}

	t.Log("parsed timestamp", ts.String())

	if time.Since(ts) > time.Second {
		t.Error("timestamp is older than a second")
	}

	if time.Until(ts) > time.Second {
		t.Error("timestamp is in the future")
	}
}

func TestTimestampGeneric(t *testing.T) {
	a := New("")

	ts, err := Timestamp(a)
	if err != nil {
		t.Error("failed to parse channel resource ID", err)
		return
	}

	t.Log("parsed timestamp", ts.String())

	if time.Since(ts) > time.Second {
		t.Error("timestamp is older than a second")
	}

	if time.Until(ts) > time.Second {
		t.Error("timestamp is in the future")
	}
}
