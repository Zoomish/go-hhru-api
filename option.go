package hhru

import (
	"context"
	"net/http"
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
}

type staticToken string

func (s staticToken) Token(ctx context.Context) (string, error) {
	return string(s), nil
}

func AccessToken(token string) TokenSource {
	return staticToken(token)
}
