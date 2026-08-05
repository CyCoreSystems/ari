package native

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequestStopsAtCustomHTTPClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		<-request.Context().Done()
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 25 * time.Millisecond}
	client := New(&Options{
		URL:        server.URL,
		HTTPClient: httpClient,
	})

	startedAt := time.Now()
	err := client.get("/channels", nil)

	require.Error(t, err)
	require.Less(t, time.Since(startedAt), 500*time.Millisecond)

	var networkError net.Error
	require.ErrorAs(t, err, &networkError)
	require.True(t, networkError.Timeout())
}

func TestNewUsesProvidedHTTPClient(t *testing.T) {
	httpClient := &http.Client{}

	client := New(&Options{HTTPClient: httpClient})

	require.Same(t, httpClient, client.httpClient)
}

func TestNewUsesPackageDefaultRequestTimeout(t *testing.T) {
	client := New(&Options{})

	require.Equal(t, RequestTimeout, client.httpClient.Timeout)
}
