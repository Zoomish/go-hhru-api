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
		defaultUA = "go-hhru-api/examples/public_position_suggest (mailto:noreply@example.com)"
	}
	ua := flag.String("hh-user-agent", defaultUA, "HH-User-Agent (or set HH_USER_AGENT)")
	text := flag.String("text", "Go developer", "search text (min 2 characters)")
	flag.Parse()

	c, err := hhru.New(hhru.Options{HHUserAgent: *ua})
	if err != nil {
		log.Fatal(err)
	}
	host := public.GetPositionsSuggestionsParamsHostHhRu
	resp, err := c.Public.GetPositionsSuggestionsWithResponse(context.Background(), &public.GetPositionsSuggestionsParams{
		HHUserAgent: *ua,
		Host:        &host,
		Text:        *text,
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.HTTPResponse.StatusCode != 200 || resp.JSON200 == nil {
		log.Fatalf("status %d", resp.HTTPResponse.StatusCode)
	}
	for _, it := range resp.JSON200.Items {
		fmt.Println(it.Text)
	}
}
