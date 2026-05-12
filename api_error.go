package hhru

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	StatusCode int
	RequestID  string
	RawBody    []byte
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.RequestID != "" {
		return fmt.Sprintf("hhru: API status %d request_id=%s", e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("hhru: API status %d", e.StatusCode)
}

func ParseAPIError(statusCode int, body []byte) *APIError {
	out := &APIError{StatusCode: statusCode}
	if len(body) == 0 {
		return out
	}
	out.RawBody = append([]byte(nil), body...)
	var meta struct {
		RequestID string `json:"request_id"`
	}
	if err := json.Unmarshal(body, &meta); err == nil {
		out.RequestID = meta.RequestID
	}
	return out
}
