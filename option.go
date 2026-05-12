package hhru

import (
	"context"
	"net/http"
	"time"
)

type TokenSource interface {
	Token(ctx context.Context) (string, error)
}

type Options struct {
	HTTPClient *http.Client
	BaseURL    string

	HHUserAgent string

	DefaultHost   string
	DefaultLocale string

	TokenSource TokenSource

	MaxRetries int

	RetryBackoffMin time.Duration
	RetryBackoffMax time.Duration

	MaxRequestsPerSecond float64

	RequestHook  func(ctx context.Context, req *http.Request) error
	ResponseHook func(ctx context.Context, resp *http.Response)
}

type staticToken string

func (s staticToken) Token(ctx context.Context) (string, error) {
	return string(s), nil
}

func AccessToken(token string) TokenSource {
	return staticToken(token)
}
