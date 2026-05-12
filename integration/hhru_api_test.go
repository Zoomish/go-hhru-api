//go:build integration

package integration_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/Zoomish/go-hhru-api"
	"github.com/Zoomish/go-hhru-api/gen/public"
)

func TestPublicGetCountries(t *testing.T) {
	if testing.Short() {
		t.Skip("skips live API call when -short is set")
	}
	t.Parallel()
	c, err := hhru.New(hhru.Options{
		HHUserAgent: "go-hhru-api/integration (mailto:noreply@example.com)",
	})
	if err != nil {
		t.Fatal(err)
	}
	host := public.GetCountriesParamsHostHhRu
	resp, err := c.Public.GetCountriesWithResponse(context.Background(), &public.GetCountriesParams{
		HHUserAgent: c.HHUserAgent(),
		Host:        &host,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.HTTPResponse.StatusCode != 200 {
		t.Fatalf("status %d", resp.HTTPResponse.StatusCode)
	}
	if resp.JSON200 == nil || len(*resp.JSON200) == 0 {
		t.Fatal("expected nonempty countries list")
	}
}

func TestExchangeClientCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skips live API call when -short is set")
	}
	id, sec := os.Getenv("HH_TEST_CLIENT_ID"), os.Getenv("HH_TEST_CLIENT_SECRET")
	if id == "" || sec == "" {
		t.Skip("set HH_TEST_CLIENT_ID and HH_TEST_CLIENT_SECRET to run OAuth integration test")
	}
	t.Parallel()
	ua := "go-hhru-api/integration (mailto:noreply@example.com)"
	tok, err := hhru.ExchangeClientCredentials(context.Background(), http.DefaultClient, hhru.TokenEndpoint(hhru.DefaultBaseURL), ua, id, sec)
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken == "" {
		t.Fatal("empty access_token")
	}
	c, err := hhru.New(hhru.Options{
		HHUserAgent: ua,
		TokenSource: hhru.AccessToken(tok.AccessToken),
	})
	if err != nil {
		t.Fatal(err)
	}
	host := public.GetCountriesParamsHostHhRu
	resp, err := c.Public.GetCountriesWithResponse(context.Background(), &public.GetCountriesParams{
		HHUserAgent: ua,
		Host:        &host,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.HTTPResponse.StatusCode != 200 {
		t.Fatalf("status %d", resp.HTTPResponse.StatusCode)
	}
}
