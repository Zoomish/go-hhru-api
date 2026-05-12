package hhru

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const refreshSkew = 90 * time.Second

type RefreshingSourceOptions struct {
	Clock Clock
}

type RefreshingTokenSource struct {
	HTTPClient   *http.Client
	TokenURL     string
	HHUserAgent  string
	ClientID     string
	ClientSecret string
	clock        Clock

	mu           sync.Mutex
	accessToken  string
	refreshToken string
	expiry       time.Time
}

func NewRefreshingTokenSource(httpClient *http.Client, tokenURL, hhUserAgent, clientID, clientSecret string, initial *TokenResponse) (TokenSource, error) {
	return NewRefreshingTokenSourceWithOptions(httpClient, tokenURL, hhUserAgent, clientID, clientSecret, initial, RefreshingSourceOptions{})
}

func NewRefreshingTokenSourceWithOptions(httpClient *http.Client, tokenURL, hhUserAgent, clientID, clientSecret string, initial *TokenResponse, opts RefreshingSourceOptions) (TokenSource, error) {
	if initial == nil {
		return nil, fmt.Errorf("hhru: initial token response is nil")
	}
	if initial.RefreshToken == "" {
		return nil, fmt.Errorf("hhru: initial token must include refresh_token (user OAuth flow); application-only client_credentials tokens do not refresh this way")
	}
	if tokenURL == "" {
		return nil, fmt.Errorf("hhru: tokenURL is empty")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	clk := opts.Clock
	if clk == nil {
		clk = wallClock{}
	}
	s := &RefreshingTokenSource{
		HTTPClient:   httpClient,
		TokenURL:     tokenURL,
		HHUserAgent:  hhUserAgent,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		clock:        clk,
	}
	s.applyLocked(initial)
	return s, nil
}

func (s *RefreshingTokenSource) applyLocked(out *TokenResponse) {
	s.accessToken = out.AccessToken
	if out.RefreshToken != "" {
		s.refreshToken = out.RefreshToken
	}
	now := s.clock.Now()
	if out.ExpiresIn > 0 {
		s.expiry = now.Add(time.Duration(out.ExpiresIn) * time.Second)
	} else {
		s.expiry = now.Add(1 * time.Hour)
	}
}

func (s *RefreshingTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.clock.Now()
	if s.accessToken != "" && !now.Add(refreshSkew).After(s.expiry) {
		return s.accessToken, nil
	}
	out, err := ExchangeRefreshToken(ctx, s.HTTPClient, s.TokenURL, s.HHUserAgent, s.refreshToken, s.ClientID, s.ClientSecret)
	if err != nil {
		return "", err
	}
	s.applyLocked(out)
	if s.accessToken == "" {
		return "", fmt.Errorf("hhru: refresh returned empty access_token")
	}
	return s.accessToken, nil
}
