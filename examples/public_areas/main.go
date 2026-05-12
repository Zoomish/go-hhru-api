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
		defaultUA = "go-hhru-api/examples/public_areas (mailto:noreply@example.com)"
	}
	ua := flag.String("hh-user-agent", defaultUA, "HH-User-Agent (or set HH_USER_AGENT)")
	flag.Parse()

	c, err := hhru.New(hhru.Options{HHUserAgent: *ua})
	if err != nil {
		log.Fatal(err)
	}
	host := public.GetAreasParamsHostHhRu
	resp, err := c.Public.GetAreasWithResponse(context.Background(), &public.GetAreasParams{
		HHUserAgent: *ua,
		Host:        &host,
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.HTTPResponse.StatusCode != 200 || resp.JSON200 == nil {
		log.Fatalf("status %d", resp.HTTPResponse.StatusCode)
	}
	for _, a := range *resp.JSON200 {
		fmt.Println(a.Id, a.Name)
	}
}
