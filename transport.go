package hhru

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
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
		if opts.RequestHook != nil {
			if err := opts.RequestHook(ctx, req); err != nil {
				return err
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

func chainRoundTripper(base http.RoundTripper, opts Options) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	t := base
	if opts.ResponseHook != nil {
		t = &responseHookTransport{next: t, hook: opts.ResponseHook}
	}
	if opts.MaxRequestsPerSecond > 0 {
		interval := time.Duration(float64(time.Second) / opts.MaxRequestsPerSecond)
		if interval < time.Millisecond {
			interval = time.Millisecond
		}
		t = &pacingTransport{next: t, interval: interval}
	}
	if opts.MaxRetries > 0 {
		t = &retryTransport{next: t, max: opts.MaxRetries, minBackoff: opts.RetryBackoffMin, maxBackoff: opts.RetryBackoffMax}
	}
	return t
}

type responseHookTransport struct {
	next http.RoundTripper
	hook func(ctx context.Context, resp *http.Response)
}

func (h *responseHookTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := h.next.RoundTrip(req)
	if err != nil || resp == nil {
		return resp, err
	}
	if h.hook != nil {
		h.hook(req.Context(), resp)
	}
	return resp, err
}

type pacingTransport struct {
	next     http.RoundTripper
	interval time.Duration
	mu       sync.Mutex
	last     time.Time
}

func (p *pacingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if p.interval <= 0 {
		return p.next.RoundTrip(req)
	}
	p.mu.Lock()
	if !p.last.IsZero() {
		elapsed := time.Since(p.last)
		if elapsed < p.interval {
			wait := p.interval - elapsed
			p.mu.Unlock()
			t := time.NewTimer(wait)
			defer t.Stop()
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-t.C:
			}
			p.mu.Lock()
		}
	}
	p.last = time.Now()
	p.mu.Unlock()
	return p.next.RoundTrip(req)
}

type retryTransport struct {
	next       http.RoundTripper
	max        int
	minBackoff time.Duration
	maxBackoff time.Duration
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
		t.sleepBackoff(attempt, resp.Header.Get("Retry-After"))
	}
	return last, lastErr
}

func (t *retryTransport) sleepBackoff(attempt int, retryAfter string) {
	var d time.Duration
	if retryAfter != "" {
		if sec, err := strconv.Atoi(retryAfter); err == nil && sec > 0 {
			d = time.Duration(sec) * time.Second
		} else if dur, err := time.ParseDuration(retryAfter); err == nil {
			d = dur
		}
	}
	if d == 0 {
		d = time.Duration(1<<min(attempt, 6)) * time.Second
	}
	if t.minBackoff > 0 && d < t.minBackoff {
		d = t.minBackoff
	}
	if t.maxBackoff > 0 && d > t.maxBackoff {
		d = t.maxBackoff
	}
	if d > 0 {
		time.Sleep(d)
	}
}
