package hhru

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func TokenEndpoint(base string) string {
	base = strings.TrimRight(base, "/")
	return base + "/token"
}

func ExchangeClientCredentials(ctx context.Context, c *http.Client, tokenURL, hhUserAgent, clientID, clientSecret string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	return postToken(ctx, c, tokenURL, hhUserAgent, form)
}

func ExchangeRefreshToken(ctx context.Context, c *http.Client, tokenURL, hhUserAgent, refreshToken, clientID, clientSecret string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	if clientID != "" {
		form.Set("client_id", clientID)
	}
	if clientSecret != "" {
		form.Set("client_secret", clientSecret)
	}
	return postToken(ctx, c, tokenURL, hhUserAgent, form)
}

func postToken(ctx context.Context, c *http.Client, tokenURL, hhUserAgent string, form url.Values) (*TokenResponse, error) {
	if c == nil {
		c = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if hhUserAgent != "" {
		req.Header.Set("HH-User-Agent", hhUserAgent)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("token endpoint: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var out TokenResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if out.AccessToken == "" {
		return nil, fmt.Errorf("token response: empty access_token")
	}
	return &out, nil
}
