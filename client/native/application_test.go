package native

import (
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/internal/dockertest"
)

func createClientURL(t *testing.T, baseURL string) (*ari.Client, bool) {
	client, err := New(&Options{
		Application:  "",
		URL:          baseURL,
		WebsocketURL: baseURL + "/events",
		Username:     "admin",
		Password:     "admin",
	})
	if err != nil {
		t.Errorf("Error building connection: %s", err)
		return nil, false
	}

	return client, true
}

func TestApplication(t *testing.T) {

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

	_, err = client.Application.Data("")
	if err == nil {
		t.Errorf("Expected error getting zero-length application")
	}

	_, err = client.Application.Data("test")
	if err == nil {
		t.Errorf("Expected error getting application 'test'")
	}
}
