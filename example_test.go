package hhru_test

import (
	"context"
	"fmt"

	"github.com/Zoomish/go-hhru-api"
)

func ExampleAccessToken() {
	ts := hhru.AccessToken("opaque")
	tok, _ := ts.Token(context.Background())
	fmt.Println(tok)

	// Output: opaque
}

func ExampleNew() {
	_, err := hhru.New(hhru.Options{
		HHUserAgent: "MyBot/1.0 (mailto:you@example.com)",
	})
	fmt.Println(err == nil)

	// Output: true
}

func ExampleNewRefreshingTokenSource() {
	ts, err := hhru.NewRefreshingTokenSource(nil,
		hhru.TokenEndpoint(hhru.DefaultBaseURL),
		"MyBot/1.0 (mailto:you@example.com)",
		"client-id", "client-secret",
		&hhru.TokenResponse{
			AccessToken:  "initial-access",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
		},
	)
	if err != nil {
		fmt.Println("err")
		return
	}
	tok, err := ts.Token(context.Background())
	if err != nil || tok == "" {
		fmt.Println("bad")
		return
	}
	fmt.Println("ok")

	// Output: ok
}

func ExampleTokenEndpoint() {
	u := hhru.TokenEndpoint(hhru.DefaultBaseURL)
	fmt.Println(u)

	// Output: https://api.hh.ru/token
}

func ExamplePagesUntil() {
	ctx := context.Background()
	sum := 0
	_ = hhru.PagesUntil(ctx, 1, func(ctx context.Context, page int) (bool, error) {
		sum += page
		if page >= 3 {
			return false, nil
		}
		return true, nil
	})
	fmt.Println(sum)

	// Output: 6
}
