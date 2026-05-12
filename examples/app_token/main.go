package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Zoomish/go-hhru-api"
	"github.com/Zoomish/go-hhru-api/gen/public"
)

func main() {
	defaultUA := os.Getenv("HH_USER_AGENT")
	if defaultUA == "" {
		defaultUA = "go-hhru-api/examples/app_token (mailto:noreply@example.com)"
	}
	ua := flag.String("hh-user-agent", defaultUA, "HH-User-Agent (or set HH_USER_AGENT)")
	flag.Parse()

	clientID := os.Getenv("HH_CLIENT_ID")
	clientSecret := os.Getenv("HH_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("set HH_CLIENT_ID and HH_CLIENT_SECRET (application credentials from dev.hh.ru)")
	}

	tok, err := hhru.ExchangeClientCredentials(context.Background(), nil, hhru.TokenEndpoint(hhru.DefaultBaseURL), *ua, clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}
	c, err := hhru.New(hhru.Options{
		HHUserAgent: *ua,
		TokenSource: hhru.AccessToken(tok.AccessToken),
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
