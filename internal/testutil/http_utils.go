package testutil

import (
	"encoding/json"
	"net/http"
	"testing"
)

func DecodeJSON(t *testing.T, resp *http.Response, target interface{}) {
	t.Helper()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(target); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

func AssertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()

	if resp.StatusCode != expected {
		t.Fatalf("expected status %d, got %d", expected, resp.StatusCode)
	}
}
