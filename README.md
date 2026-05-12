# go-hhru-api

[![Go Reference](https://pkg.go.dev/badge/github.com/Zoomish/go-hhru-api.svg)](https://pkg.go.dev/github.com/Zoomish/go-hhru-api)

Go client for the [HeadHunter API](https://api.hh.ru/openapi/redoc): typed sub-clients generated from OpenAPI, plus a small facade (`hhru.New`) for shared headers, optional OAuth bearer injection, default query parameters, and optional retries on `429` / `503`.

Package documentation and examples: [pkg.go.dev/github.com/Zoomish/go-hhru-api](https://pkg.go.dev/github.com/Zoomish/go-hhru-api).

## Зачем эта библиотека

- **Типы и методы под OpenAPI** — меньше ручной сборки URL, query и заголовков; четыре сгенерированных клиента по смыслу API (работодатель, соискатель, публичное, приложение).
- **Единые правила HH в одном месте** — обязательный `HH-User-Agent`, при желании дефолтные `host` / `locale`, Bearer через [TokenSource](option.go), опциональные ретраи на перегрузку API.
- **Долгоживущие сервисы** — [NewRefreshingTokenSource](token_refresh.go) сам обновляет access-токен по `refresh_token` под мьютексом; `hhru.New` остаётся без своего состояния, вы только передаёте один `TokenSource`.
- **Воспроизводимость** — в репозитории зафиксированы `api/openapi.yml` и `gen/`; CI проверяет, что генерация совпадает с коммитом.

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

`HH-User-Agent` is also applied via a request editor when it is missing on the outgoing request; generated methods still expect `HHUserAgent` in many `*Params` structs—use `c.HHUserAgent()` to reuse the same string.

## Runnable examples (from repo root)

```bash
go run ./examples/public_countries
```

```bash
export HH_CLIENT_ID=… HH_CLIENT_SECRET=…
go run ./examples/app_token
```

```bash
export HH_CLIENT_ID=… HH_CLIENT_SECRET=…
export HH_INITIAL_TOKEN_JSON_PATH=/path/to/token.json
go run ./examples/refreshing_token
```

`token.json` must match [`TokenResponse`](oauth.go) (`access_token`, `refresh_token`, `expires_in`) from the user OAuth redirect flow. Use `-hh-user-agent` or `HH_USER_AGENT` where supported.

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

Refresh flow: `ExchangeRefreshToken` in [`oauth.go`](oauth.go) with `grant_type=refresh_token` (optional `client_id` / `client_secret` if your app requires them).

## User token: automatic refresh (`TokenSource`)

After the browser / OAuth redirect flow you typically have `access_token`, `refresh_token`, and `expires_in`. Pass the first token response into [NewRefreshingTokenSource](token_refresh.go): it keeps the access token in memory, returns it from [TokenSource](option.go) until shortly before expiry, then calls `ExchangeRefreshToken` under a mutex (safe for concurrent API calls). Use it with `hhru.New`:

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

Application-only **client_credentials** tokens from HH do not use this refresh model; keep using [AccessToken](option.go) with a periodically re-issued app token if HH invalidates it.

## Optional: retries

Set `Options.MaxRetries` to a small positive number to retry `GET`-safe style requests on `429` / `503` with backoff / `Retry-After`. Avoid high values for requests with bodies.

## Optional: pagination helper

[PagesUntil](pagination.go) runs `page = start, start+1, …` until the callback returns `continueNext == false` or an error.

## Maintainers: regenerate `gen/`

Requires Go and network for module downloads:

```bash
make generate
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for tests and integration tags. Source spec: [`api/openapi.yml`](api/openapi.yml). CI checks that `gen/` matches the spec (see `.github/workflows/`).

## Integration tests

- `go test -short ./...` — default CI and local quick check; no live HH calls (root [`token_refresh_test.go`](token_refresh_test.go) uses `httptest` only).
- Live API: `go test -tags=integration -timeout 5m ./integration` — hits HH for `TestPublicGetCountries` (omit `-short`). With `HH_TEST_CLIENT_ID` and `HH_TEST_CLIENT_SECRET` from [dev.hh.ru](https://dev.hh.ru/admin), also runs `TestExchangeClientCredentials`.
- Optional GitHub workflow: [`.github/workflows/integration-tests.yml`](.github/workflows/integration-tests.yml) (`workflow_dispatch`; configure repository secrets if you use it).

## Versioning

Semantic versioning after `v1`. Until then, treat minors as potentially breaking when `openapi.yml` or `gen/` changes. See [CHANGELOG.md](CHANGELOG.md).
