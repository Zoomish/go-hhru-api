package hhru

import "testing"

func TestParseAPIError_requestID(t *testing.T) {
	body := []byte(`{"request_id":"abc-123","description":"bad"}`)
	e := ParseAPIError(400, body)
	if e.StatusCode != 400 {
		t.Fatalf("status %d", e.StatusCode)
	}
	if e.RequestID != "abc-123" {
		t.Fatalf("request_id %q", e.RequestID)
	}
	if len(e.RawBody) == 0 {
		t.Fatal("raw body")
	}
	if e.Error() == "" {
		t.Fatal("error string")
	}
}
