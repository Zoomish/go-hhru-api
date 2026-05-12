package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Zoomish/go-hhru-api"
	"github.com/Zoomish/go-hhru-api/gen/public"
)

func main() {
	defaultUA := os.Getenv("HH_USER_AGENT")
	if defaultUA == "" {
		defaultUA = "go-hhru-api/examples/public_countries (mailto:noreply@example.com)"
	}
	ua := flag.String("hh-user-agent", defaultUA, "HH-User-Agent (or set HH_USER_AGENT)")
	flag.Parse()

	c, err := hhru.New(hhru.Options{HHUserAgent: *ua})
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
	if resp.JSON200 == nil {
		log.Fatal("empty body")
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resp.JSON200); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stderr, "ok")
}
