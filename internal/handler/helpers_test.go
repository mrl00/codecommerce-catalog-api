package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rr.Code != expected {
		t.Errorf("expected status %d, got %d", expected, rr.Code)
	}
}

func assertContentType(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	t.Helper()
	if ct := rr.Header().Get("Content-Type"); ct != expected {
		t.Errorf("expected Content-Type %q, got %q", expected, ct)
	}
}

func assertJSONError(t *testing.T, rr *httptest.ResponseRecorder, expectedMsg string) {
	t.Helper()
	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp["error"] != expectedMsg {
		t.Errorf("expected error %q, got %q", expectedMsg, resp["error"])
	}
}
