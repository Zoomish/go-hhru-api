package hhru

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func newGlobalRequestEditor(opts Options) func(context.Context, *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		if err := mergeDefaultQuery(req, opts.DefaultHost, opts.DefaultLocale); err != nil {
			return err
		}
		if opts.HHUserAgent != "" && req.Header.Get("HH-User-Agent") == "" {
			req.Header.Set("HH-User-Agent", opts.HHUserAgent)
		}
		if opts.TokenSource != nil {
			tok, err := opts.TokenSource.Token(ctx)
			if err != nil {
				return err
			}
			if tok != "" {
				req.Header.Set("Authorization", "Bearer "+tok)
			}
		}
		return nil
	}
}

func mergeDefaultQuery(req *http.Request, defaultHost, defaultLocale string) error {
	if defaultHost == "" && defaultLocale == "" {
		return nil
	}
	u := req.URL
	if u == nil {
		return nil
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return err
	}
	if defaultHost != "" && q.Get("host") == "" {
		q.Set("host", defaultHost)
	}
	if defaultLocale != "" && q.Get("locale") == "" {
		q.Set("locale", defaultLocale)
	}
	u.RawQuery = q.Encode()
	return nil
}

type retryTransport struct {
	next http.RoundTripper
	max  int
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	next := t.next
	if next == nil {
		next = http.DefaultTransport
	}
	var last *http.Response
	var lastErr error
	for attempt := 0; attempt <= t.max; attempt++ {
		resp, err := next.RoundTrip(req.Clone(req.Context()))
		lastErr = err
		last = resp
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode != http.StatusServiceUnavailable {
			return resp, nil
		}
		if attempt == t.max {
			return resp, nil
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if sec, err := strconv.Atoi(ra); err == nil && sec > 0 {
				time.Sleep(time.Duration(sec) * time.Second)
			} else if d, err := time.ParseDuration(ra); err == nil {
				time.Sleep(d)
			} else {
				time.Sleep(time.Duration(1<<min(attempt, 6)) * time.Second)
			}
		} else {
			time.Sleep(time.Duration(1<<min(attempt, 6)) * time.Second)
		}
	}
	return last, lastErr
}

func wrapRetryTransport(base http.RoundTripper, maxRetries int) http.RoundTripper {
	if maxRetries <= 0 {
		return base
	}
	if base == nil {
		base = http.DefaultTransport
	}
	return &retryTransport{next: base, max: maxRetries}
}
