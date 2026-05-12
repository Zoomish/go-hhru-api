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
		defaultUA = "go-hhru-api/examples/custom_options (mailto:noreply@example.com)"
	}
	ua := flag.String("hh-user-agent", defaultUA, "HH-User-Agent (or set HH_USER_AGENT)")
	flag.Parse()

	c, err := hhru.New(hhru.Options{
		HHUserAgent:   *ua,
		DefaultHost:   "hh.ru",
		DefaultLocale: "RU",
		MaxRetries:    2,
	})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := c.Public.GetCountriesWithResponse(context.Background(), &public.GetCountriesParams{
		HHUserAgent: *ua,
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.HTTPResponse.StatusCode != 200 || resp.JSON200 == nil {
		log.Fatalf("status %d", resp.HTTPResponse.StatusCode)
	}
	log.Printf("ok, countries=%d (default host/locale merged by client)", len(*resp.JSON200))
}
