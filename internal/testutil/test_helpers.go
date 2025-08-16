package testutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rickyreddygari/walletsdk/internal/app"
)

// NewTestServer spins up an httptest server with in-memory dependencies configured by app container.
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

// MustDo performs an HTTP request and fails the test on error.
func MustDo(t *testing.T, client *http.Client, req *http.Request) *http.Response {
	t.Helper()

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("http request failed: %v", err)
	}
	return resp
}

// Background returns context.Background() but allows for future expansion.
func Background() context.Context {
	return context.Background()
}
