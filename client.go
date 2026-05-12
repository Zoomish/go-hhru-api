package hhru

import (
	"fmt"
	"net/http"

	"github.com/Zoomish/go-hhru-api/gen/app"
	"github.com/Zoomish/go-hhru-api/gen/applicant"
	"github.com/Zoomish/go-hhru-api/gen/employer"
	"github.com/Zoomish/go-hhru-api/gen/public"
)

const DefaultBaseURL = "https://api.hh.ru"

type Client struct {
	Employer  *employer.ClientWithResponses
	Applicant *applicant.ClientWithResponses
	Public    *public.ClientWithResponses
	App       *app.ClientWithResponses

	hhUserAgent string
}

func (c *Client) HHUserAgent() string {
	return c.hhUserAgent
}

func New(opts Options) (*Client, error) {
	if opts.HHUserAgent == "" {
		return nil, fmt.Errorf("hhru: Options.HHUserAgent is required (see https://api.hh.ru/openapi/redoc client identification)")
	}
	base := opts.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	hc := opts.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	if opts.MaxRetries > 0 || opts.MaxRequestsPerSecond > 0 || opts.ResponseHook != nil {
		base := hc.Transport
		if base == nil {
			base = http.DefaultTransport
		}
		hc = &http.Client{
			Transport:     chainRoundTripper(base, opts),
			Timeout:       hc.Timeout,
			CheckRedirect: hc.CheckRedirect,
			Jar:           hc.Jar,
		}
	}
	editor := newGlobalRequestEditor(opts)
	e, err := employer.NewClientWithResponses(base,
		employer.WithHTTPClient(hc),
		employer.WithRequestEditorFn(employer.RequestEditorFn(editor)),
	)
	if err != nil {
		return nil, fmt.Errorf("employer client: %w", err)
	}
	a, err := applicant.NewClientWithResponses(base,
		applicant.WithHTTPClient(hc),
		applicant.WithRequestEditorFn(applicant.RequestEditorFn(editor)),
	)
	if err != nil {
		return nil, fmt.Errorf("applicant client: %w", err)
	}
	p, err := public.NewClientWithResponses(base,
		public.WithHTTPClient(hc),
		public.WithRequestEditorFn(public.RequestEditorFn(editor)),
	)
	if err != nil {
		return nil, fmt.Errorf("public client: %w", err)
	}
	ap, err := app.NewClientWithResponses(base,
		app.WithHTTPClient(hc),
		app.WithRequestEditorFn(app.RequestEditorFn(editor)),
	)
	if err != nil {
		return nil, fmt.Errorf("app client: %w", err)
	}
	return &Client{
		Employer:    e,
		Applicant:   a,
		Public:      p,
		App:         ap,
		hhUserAgent: opts.HHUserAgent,
	}, nil
}
