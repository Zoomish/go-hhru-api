package hhru

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const refreshSkew = 90 * time.Second

type RefreshingTokenSource struct {
	HTTPClient   *http.Client
	TokenURL     string
	HHUserAgent  string
	ClientID     string
	ClientSecret string

	mu           sync.Mutex
	accessToken  string
	refreshToken string
	expiry       time.Time
}

func NewRefreshingTokenSource(httpClient *http.Client, tokenURL, hhUserAgent, clientID, clientSecret string, initial *TokenResponse) (TokenSource, error) {
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
	s := &RefreshingTokenSource{
		HTTPClient:   httpClient,
		TokenURL:     tokenURL,
		HHUserAgent:  hhUserAgent,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	s.applyLocked(initial)
	return s, nil
}

func (s *RefreshingTokenSource) applyLocked(out *TokenResponse) {
	s.accessToken = out.AccessToken
	if out.RefreshToken != "" {
		s.refreshToken = out.RefreshToken
	}
	if out.ExpiresIn > 0 {
		s.expiry = time.Now().Add(time.Duration(out.ExpiresIn) * time.Second)
	} else {
		s.expiry = time.Now().Add(1 * time.Hour)
	}
}

func (s *RefreshingTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.accessToken != "" && !time.Now().Add(refreshSkew).After(s.expiry) {
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
