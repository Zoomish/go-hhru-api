package main

import (
	"context"
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
		defaultUA = "go-hhru-api/examples/public_locales (mailto:noreply@example.com)"
	}
	ua := flag.String("hh-user-agent", defaultUA, "HH-User-Agent (or set HH_USER_AGENT)")
	flag.Parse()

	c, err := hhru.New(hhru.Options{HHUserAgent: *ua})
	if err != nil {
		log.Fatal(err)
	}
	host := public.GetLocalesParamsHostHhRu
	resp, err := c.Public.GetLocalesWithResponse(context.Background(), &public.GetLocalesParams{
		HHUserAgent: *ua,
		Host:        &host,
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.HTTPResponse.StatusCode != 200 || resp.JSON200 == nil {
		log.Fatalf("status %d", resp.HTTPResponse.StatusCode)
	}
	for _, loc := range *resp.JSON200 {
		fmt.Println(loc.Id, loc.Name)
	}
}
