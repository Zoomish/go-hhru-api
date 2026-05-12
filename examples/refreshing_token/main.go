package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/Zoomish/go-hhru-api"
	"github.com/Zoomish/go-hhru-api/gen/public"
)

func main() {
	defaultUA := os.Getenv("HH_USER_AGENT")
	if defaultUA == "" {
		defaultUA = "go-hhru-api/examples/refreshing_token (mailto:noreply@example.com)"
	}
	ua := flag.String("hh-user-agent", defaultUA, "HH-User-Agent (or set HH_USER_AGENT)")
	tokenPath := flag.String("token-json", os.Getenv("HH_INITIAL_TOKEN_JSON_PATH"), "path to JSON with access_token, refresh_token, expires_in (or set HH_INITIAL_TOKEN_JSON_PATH)")
	flag.Parse()

	clientID := os.Getenv("HH_CLIENT_ID")
	clientSecret := os.Getenv("HH_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("set HH_CLIENT_ID and HH_CLIENT_SECRET")
	}
	if *tokenPath == "" {
		log.Fatal("set -token-json or HH_INITIAL_TOKEN_JSON_PATH to a file produced after user OAuth (access_token, refresh_token, expires_in)")
	}
	raw, err := os.ReadFile(*tokenPath)
	if err != nil {
		log.Fatal(err)
	}
	var initial hhru.TokenResponse
	if err := json.Unmarshal(raw, &initial); err != nil {
		log.Fatal(err)
	}
	ts, err := hhru.NewRefreshingTokenSource(nil, hhru.TokenEndpoint(hhru.DefaultBaseURL), *ua, clientID, clientSecret, &initial)
	if err != nil {
		log.Fatal(err)
	}
	c, err := hhru.New(hhru.Options{
		HHUserAgent: *ua,
		TokenSource: ts,
	})
	if err != nil {
		log.Fatal(err)
	}
	host := public.GetCountriesParamsHostHhRu
	resp, err := c.Public.GetCountriesWithResponse(context.Background(), &public.GetCountriesParams{
		HHUserAgent: *ua,
		Host:        &host,
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.HTTPResponse.StatusCode != 200 {
		log.Fatalf("status %d", resp.HTTPResponse.StatusCode)
	}
	log.Printf("ok, countries=%d", len(*resp.JSON200))
}
