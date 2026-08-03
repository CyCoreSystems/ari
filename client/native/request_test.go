package native

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequestStopsAtConfiguredTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		<-request.Context().Done()
	}))
	defer server.Close()

	client := New(&Options{
		URL:            server.URL,
		RequestTimeout: 25 * time.Millisecond,
	})

	startedAt := time.Now()
	err := client.get("/channels", nil)

	require.Error(t, err)
	require.Less(t, time.Since(startedAt), 500*time.Millisecond)

	var networkError net.Error
	require.ErrorAs(t, err, &networkError)
	require.True(t, networkError.Timeout())
}

func TestRequestUsesPackageDefaultTimeout(t *testing.T) {
	client := New(&Options{})

	require.Equal(t, RequestTimeout, client.httpClient.Timeout)
}
