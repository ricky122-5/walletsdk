package testutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rickyreddygari/walletsdk/internal/app"
)

func NewTestServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	container, err := app.NewContainer()
	if err != nil {
		t.Fatalf("bootstrap container: %v", err)
	}

	server := httptest.NewServer(container.HTTPServer)

	cleanup := func() {
		server.Close()
	}

	return server, cleanup
}

func MustDo(t *testing.T, client *http.Client, req *http.Request) *http.Response {
	t.Helper()

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("http request failed: %v", err)
	}
	return resp
}

func Background() context.Context {
	return context.Background()
}
