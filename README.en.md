# go-hhru-api

[Russian](README.md) | **English**

[![Go Reference](https://pkg.go.dev/badge/github.com/Zoomish/go-hhru-api.svg)](https://pkg.go.dev/github.com/Zoomish/go-hhru-api)
[![Go 1.26](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev/dl/)
[![Test](https://github.com/Zoomish/go-hhru-api/actions/workflows/test.yml/badge.svg)](https://github.com/Zoomish/go-hhru-api/actions/workflows/test.yml)

Go client for the [HeadHunter API](https://api.hh.ru/openapi/redoc): typed OpenAPI sub-clients and an `hhru.New` facade for shared headers, optional OAuth bearer, default query parameters, and retries on `429` / `503`. Package docs on [pkg.go.dev](https://pkg.go.dev/github.com/Zoomish/go-hhru-api) (badge above).

Russian documentation: [README.md](README.md) and [README.ru.md](README.ru.md).

## Contents

- [Why this library](#why-this-library)
- [Install](#install)
- [Quick start](#quick-start)
- [Runnable examples](#runnable-examples)
- [Application token (client credentials)](#application-token-client-credentials)
- [User token and auto-refresh](#user-token-and-auto-refresh)
- [Options: retries and limits](#options-retries-and-limits)
- [Pagination](#pagination)
- [Observability and API errors](#observability-and-api-errors)
- [Maintainers: regenerate `gen/`](#maintainers-regenerate-gen)
- [Integration tests](#integration-tests)
- [Versioning](#versioning)

## Why this library

- **OpenAPI types and methods** — less manual URL/query/header wiring; four generated clients (employer, applicant, public, app).
- **HH rules in one place** — required `HH-User-Agent`, optional default `host` / `locale`, Bearer via [`TokenSource`](option.go), optional overload retries, optional request rate limit, request/response hooks.
- **Long-running services** — [`NewRefreshingTokenSource`](token_refresh.go) refreshes access tokens under a mutex; [`hhru.New`](client.go) stays stateless; you pass one `TokenSource`.
- **Reproducibility** — `api/openapi.yml` and `gen/` are committed; CI checks codegen drift.

## Install

```bash
go get github.com/Zoomish/go-hhru-api
```

## Quick start

```go
import (
    "context"
    "github.com/Zoomish/go-hhru-api"
    "github.com/Zoomish/go-hhru-api/gen/public"
)

c, err := hhru.New(hhru.Options{
    HHUserAgent: "MyService/1.0 (mailto:you@example.com)",
    DefaultHost: "hh.ru",
})
if err != nil {
    panic(err)
}
host := public.GetCountriesParamsHostHhRu
resp, err := c.Public.GetCountriesWithResponse(context.Background(), &public.GetCountriesParams{
    HHUserAgent: c.HHUserAgent(),
    Host:        &host,
})
```

`HH-User-Agent` is also applied by the request editor when missing; many `*Params` still expect `HHUserAgent` — use `c.HHUserAgent()`.

## Runnable examples

From the repository root.

Public dictionaries (no OAuth):

```bash
go run ./examples/public_countries
go run ./examples/public_locales
go run ./examples/public_areas
go run ./examples/public_industries
go run ./examples/public_languages
go run ./examples/public_position_suggest
go run ./examples/public_position_suggest -text "backend engineer"
```

Client options (`DefaultHost`, `DefaultLocale`, `MaxRetries`):

```bash
go run ./examples/custom_options
```

OAuth:

```bash
export HH_CLIENT_ID=… HH_CLIENT_SECRET=…
go run ./examples/app_token
```

```bash
export HH_CLIENT_ID=… HH_CLIENT_SECRET=…
export HH_INITIAL_TOKEN_JSON_PATH=/path/to/token.json
go run ./examples/refreshing_token
```

The token file must match [`TokenResponse`](oauth.go) (`access_token`, `refresh_token`, `expires_in`) after the user OAuth flow. Examples support `-hh-user-agent` or `HH_USER_AGENT`.

## Application token (client credentials)

```go
tok, err := hhru.ExchangeClientCredentials(ctx, http.DefaultClient,
    hhru.TokenEndpoint(hhru.DefaultBaseURL),
    "MyService/1.0 (mailto:you@example.com)",
    clientID, clientSecret,
)
c, err := hhru.New(hhru.Options{
    HHUserAgent: "MyService/1.0 (mailto:you@example.com)",
    TokenSource: hhru.AccessToken(tok.AccessToken),
})
```

Refresh via [`ExchangeRefreshToken`](oauth.go) with `grant_type=refresh_token` (optional `client_id` / `client_secret`).

## User token and auto-refresh

After OAuth you typically have `access_token`, `refresh_token`, and `expires_in`. Pass the first token response into [`NewRefreshingTokenSource`](token_refresh.go) and use it as `TokenSource` with [`hhru.New`](client.go).

```go
initial := &hhru.TokenResponse{
    AccessToken:  accessFromOAuth,
    RefreshToken: refreshFromOAuth,
    ExpiresIn:    expiresInSeconds,
}
ts, err := hhru.NewRefreshingTokenSource(http.DefaultClient,
    hhru.TokenEndpoint(hhru.DefaultBaseURL),
    "MyService/1.0 (mailto:you@example.com)",
    clientID, clientSecret,
    initial,
)
if err != nil {
    panic(err)
}
c, err := hhru.New(hhru.Options{
    HHUserAgent: "MyService/1.0 (mailto:you@example.com)",
    TokenSource: ts,
})
```

Application-only **client_credentials** tokens do not use this refresh path; re-issue with [`AccessToken`](option.go).

## Options: retries and limits

- `Options.MaxRetries` — retries for safe requests on `429` / `503` honoring `Retry-After` and exponential backoff; if `RetryBackoffMin` / `RetryBackoffMax` are set, sleep duration is clamped to that range.
- `Options.MaxRequestsPerSecond` — soft limit on outgoing request starts (even spacing).
- `Options.RequestHook` / `Options.ResponseHook` — hooks before send and after response headers (body not read; useful for logging and tracing).

## Pagination

[`PagesUntil`](pagination.go) runs `page = start, start+1, …` until the callback returns `continueNext == false` or an error.

## Observability and API errors

- Parse typical HH JSON errors with [`ParseAPIError`](api_error.go) and [`APIError`](api_error.go) (e.g. `RequestID` when present).
- For tests, use [`NewRefreshingTokenSourceWithOptions`](token_refresh.go) with [`Clock`](clock.go) in [`RefreshingSourceOptions`](token_refresh.go).

## Maintainers: regenerate `gen/`

Requires Go and network for tool modules:

```bash
make generate
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for tests and the `integration` build tag. Spec: [`api/openapi.yml`](api/openapi.yml). CI: [`.github/workflows/`](.github/workflows/).

## Integration tests

- `go test -short ./...` — fast path without live HH calls ([`token_refresh_test.go`](token_refresh_test.go) uses `httptest`).
- Live API: `go test -tags=integration -timeout 5m ./integration` (omit `-short`). With `HH_TEST_CLIENT_ID` and `HH_TEST_CLIENT_SECRET` from [dev.hh.ru](https://dev.hh.ru/admin), the OAuth scenario runs as well.
- Manual workflow: [`.github/workflows/integration-tests.yml`](.github/workflows/integration-tests.yml) (`workflow_dispatch`).

## Versioning

Semantic versioning after `v1`. Until then, minors may break when `openapi.yml` or `gen/` changes. See [CHANGELOG.md](CHANGELOG.md).
