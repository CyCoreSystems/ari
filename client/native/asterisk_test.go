package native

import (
	"testing"

	"github.com/CyCoreSystems/ari/internal/dockertest"
)

func TestAsteriskGetInfo(t *testing.T) {

	// start asterisk in docker container
	url, deferFn, err := dockertest.StartAsterisk()
	if err != nil {
		t.Fatalf("cannot start asterisk in container for testing: %s", err)
	}
	defer deferFn()

	client, ok := createClientURL(t, url)
	if !ok {
		return
	}

	_, err = client.Asterisk.Info("some-filter")
	if err != errOnlyUnsupported {
		t.Errorf("Unexpected error: '%s', expected '%s'", err, errOnlyUnsupported)
	}

	_, err = client.Asterisk.Info("")
	if err != nil {
		t.Errorf("Error getting asterisk info: %s", err)
	}
}

func TestAsteriskGetVariable(t *testing.T) {

	// start asterisk in docker container
	url, deferFn, err := dockertest.StartAsterisk()
	if err != nil {
		t.Fatalf("cannot start asterisk in container for testing: %s", err)
	}
	defer deferFn()

	client, ok := createClientURL(t, url)
	if !ok {
		return
	}

	v, err := client.Asterisk.GetVariable("var")
	if err != nil {
		t.Errorf("Error getting asterisk var: %s", err)
	}

	if v != "" {
		t.Errorf("Expected var to be empty, got %s", v)
	}

	v, err = client.Asterisk.GetVariable("")
	if err == nil {
		t.Errorf("Expected error getting variable")
	}
}

func TestAsteriskSetVariable(t *testing.T) {

	// start asterisk in docker container
	url, deferFn, err := dockertest.StartAsterisk()
	if err != nil {
		t.Fatalf("cannot start asterisk in container for testing: %s", err)
	}
	defer deferFn()

	client, ok := createClientURL(t, url)
	if !ok {
		return
	}

	err = client.Asterisk.SetVariable("var", "value")
	if err != nil {
		t.Errorf("Error setting asterisk var: %s", err)
	}

	err = client.Asterisk.SetVariable("", "value")
	if err == nil {
		t.Errorf("Expected error setting variable")
	}

	// clear variable
	client.Asterisk.SetVariable("var", "")
}

func TestAsteriskGetSetVariable(t *testing.T) {

	// start asterisk in docker container
	url, deferFn, err := dockertest.StartAsterisk()
	if err != nil {
		t.Fatalf("cannot start asterisk in container for testing: %s", err)
	}
	defer deferFn()

	client, ok := createClientURL(t, url)
	if !ok {
		return
	}

	v, err := client.Asterisk.GetVariable("var")
	if err != nil {
		t.Errorf("Error getting asterisk var: %s", err)
	}

	if v != "" {
		t.Errorf("Expected var to be empty, got %s", v)
	}

	err = client.Asterisk.SetVariable("var", "value")
	if err != nil {
		t.Errorf("Error setting asterisk var: %s", err)
	}

	v, err = client.Asterisk.GetVariable("var")
	if err != nil {
		t.Errorf("Error getting asterisk var: %s", err)
	}

	if v != "value" {
		t.Errorf("Expected var to be 'value', got %s", v)
	}

	// clear variable
	client.Asterisk.SetVariable("var", "")
}

func TestAsteriskReloadModule(t *testing.T) {
	// start asterisk in docker container
	url, deferFn, err := dockertest.StartAsterisk()
	if err != nil {
		t.Fatalf("cannot start asterisk in container for testing: %s", err)
	}
	defer deferFn()

	client, ok := createClientURL(t, url)
	if !ok {
		return
	}

	err = client.Asterisk.ReloadModule("")
	if err == nil {
		t.Errorf("Expected error reloading zero-length module")
	}

	err = client.Asterisk.ReloadModule("unknown")
	if err == nil {
		t.Errorf("Expected error reloading unknown module")
	}

	err = client.Asterisk.ReloadModule("pbx_config.so")
	if err != nil {
		t.Errorf("Unexpected error reloading http module: %s", err)
	}

}
