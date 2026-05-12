# go-hhru-api

Go client for the [HeadHunter API](https://api.hh.ru/openapi/redoc): typed sub-clients generated from OpenAPI, plus a small facade (`hhru.New`) for shared headers, optional OAuth bearer injection, default query parameters, and optional retries on `429` / `503`.

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

Source spec: [`api/openapi.yml`](api/openapi.yml). CI checks that `gen/` matches the spec (see `.github/workflows/`).

## Integration tests

- `go test -short ./...` — skips network.
- Full: `go test ./...` — hits HH for `TestPublicGetCountries`.
- OAuth: set `HH_TEST_CLIENT_ID` and `HH_TEST_CLIENT_SECRET` (from [dev.hh.ru](https://dev.hh.ru/admin)) for `TestExchangeClientCredentials`.

## Versioning

Semantic versioning after `v1`. Until then, treat minors as potentially breaking when `openapi.yml` or `gen/` changes. See [CHANGELOG.md](CHANGELOG.md).
