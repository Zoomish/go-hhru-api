package hhru

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestRefreshingTokenSource_CachesAccessToken(t *testing.T) {
	var refreshCalls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "from-server",
			RefreshToken: "new-refresh",
			ExpiresIn:    7200,
		})
	}))
	t.Cleanup(srv.Close)

	initial := &TokenResponse{
		AccessToken:  "cached",
		RefreshToken: "refresh-seed",
		ExpiresIn:    7200,
	}
	ts, err := NewRefreshingTokenSource(http.DefaultClient, srv.URL, "ua", "id", "secret", initial)
	if err != nil {
		t.Fatal(err)
	}
	tok1, err := ts.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if tok1 != "cached" {
		t.Fatalf("got %q want cached", tok1)
	}
	tok2, err := ts.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if tok2 != "cached" {
		t.Fatalf("got %q want cached", tok2)
	}
	if refreshCalls.Load() != 0 {
		t.Fatalf("unexpected refresh calls: %d", refreshCalls.Load())
	}
}

func TestRefreshingTokenSource_RefreshesWhenEmptyAccess(t *testing.T) {
	var refreshCalls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "fresh",
			RefreshToken: "rotated",
			ExpiresIn:    3600,
		})
	}))
	t.Cleanup(srv.Close)

	ts, err := NewRefreshingTokenSource(http.DefaultClient, srv.URL, "ua", "id", "secret", &TokenResponse{
		RefreshToken: "only-refresh",
	})
	if err != nil {
		t.Fatal(err)
	}
	tok, err := ts.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if tok != "fresh" {
		t.Fatalf("got %q want fresh", tok)
	}
	if refreshCalls.Load() != 1 {
		t.Fatalf("refresh calls: %d", refreshCalls.Load())
	}
}
