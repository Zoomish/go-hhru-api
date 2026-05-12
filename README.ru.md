# go-hhru-api

**Русский** | [English](README.en.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/Zoomish/go-hhru-api.svg)](https://pkg.go.dev/github.com/Zoomish/go-hhru-api)
[![Go 1.26](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev/dl/)
[![Test](https://github.com/Zoomish/go-hhru-api/actions/workflows/test.yml/badge.svg)](https://github.com/Zoomish/go-hhru-api/actions/workflows/test.yml)

Клиент на Go для [API HeadHunter](https://api.hh.ru/openapi/redoc): типизированные подклиенты из OpenAPI и фасад `hhru.New` с общими заголовками, опциональным OAuth Bearer, дефолтными query и ретраями на `429` / `503`. Документация API пакета на [pkg.go.dev](https://pkg.go.dev/github.com/Zoomish/go-hhru-api) (бейдж выше).

Совпадает с корневым [`README.md`](README.md); файл для явного суффикса `.ru` — обновляйте вместе с ним.

## Содержание

- [Зачем эта библиотека](#зачем-эта-библиотека)
- [Установка](#установка)
- [Быстрый старт](#быстрый-старт)
- [Примеры `go run`](#примеры-go-run)
- [Токен приложения (client credentials)](#токен-приложения-client-credentials)
- [Пользовательский токен и автообновление](#пользовательский-токен-и-автообновление)
- [Опции: ретраи и лимиты](#опции-ретраи-и-лимиты)
- [Пагинация](#пагинация)
- [Наблюдаемость и ошибки API](#наблюдаемость-и-ошибки-api)
- [Мейнтейнерам: генерация `gen/`](#мейнтейнерам-генерация-gen)
- [Интеграционные тесты](#интеграционные-тесты)
- [Версионирование](#версионирование)

## Зачем эта библиотека

- **Типы и методы под OpenAPI** — меньше ручной сборки URL, query и заголовков; четыре сгенерированных клиента (работодатель, соискатель, публичное, приложение).
- **Единые правила HH** — обязательный `HH-User-Agent`, при необходимости дефолтные `host` / `locale`, Bearer через [`TokenSource`](option.go), опциональные ретраи при перегрузке API, опциональный лимит частоты запросов и хуки запроса/ответа.
- **Долгоживущие сервисы** — [`NewRefreshingTokenSource`](token_refresh.go) обновляет access-токен по `refresh_token` под мьютексом; [`hhru.New`](client.go) остаётся без состояния, передаёте один `TokenSource`.
- **Воспроизводимость** — в репозитории зафиксированы `api/openapi.yml` и `gen/`; CI проверяет совпадение генерации с коммитом.

## Установка

```bash
go get github.com/Zoomish/go-hhru-api
```

## Быстрый старт

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

Заголовок `HH-User-Agent` также подставляется редактором запроса, если его нет; во многих `*Params` поле `HHUserAgent` всё равно ожидается — используйте `c.HHUserAgent()`.

## Примеры `go run`

Из корня репозитория.

Справочники без OAuth:

```bash
go run ./examples/public_countries
go run ./examples/public_locales
go run ./examples/public_areas
go run ./examples/public_industries
go run ./examples/public_languages
go run ./examples/public_position_suggest
go run ./examples/public_position_suggest -text "backend engineer"
```

Опции клиента (`DefaultHost`, `DefaultLocale`, `MaxRetries`):

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

Файл токена должен соответствовать [`TokenResponse`](oauth.go): `access_token`, `refresh_token`, `expires_in` после OAuth пользователя. Флаги `-hh-user-agent` или переменная `HH_USER_AGENT` поддерживаются в примерах.

## Токен приложения (client credentials)

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

Обновление через [`ExchangeRefreshToken`](oauth.go) с `grant_type=refresh_token` (при необходимости `client_id` / `client_secret`).

## Пользовательский токен и автообновление

После OAuth обычно есть `access_token`, `refresh_token`, `expires_in`. Передайте первый ответ в [`NewRefreshingTokenSource`](token_refresh.go) и используйте как `TokenSource` в [`hhru.New`](client.go).

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

Токены только **client_credentials** не обновляются этим путём; перевыпускайте приложение через [`AccessToken`](option.go).

## Опции: ретраи и лимиты

- `Options.MaxRetries` — ретраи для безопасных запросов на `429` / `503` с учётом `Retry-After` и экспоненциальной задержкой; при заданных `RetryBackoffMin` / `RetryBackoffMax` длительность сна ограничивается этим интервалом.
- `Options.MaxRequestsPerSecond` — мягкий лимит частоты исходящих запросов (равномерный интервал между стартами запросов).
- `Options.RequestHook` / `Options.ResponseHook` — хуки до отправки и после ответа (тело ответа не читается; удобно для логов и трассировки).

## Пагинация

[`PagesUntil`](pagination.go) вызывает `page = start, start+1, …`, пока колбэк не вернёт `continueNext == false` или ошибку.

## Наблюдаемость и ошибки API

- Разбор типичного JSON-ошибки HH: [`ParseAPIError`](api_error.go) по телу ответа и [`APIError`](api_error.go) с полями вроде `RequestID` (если есть в JSON).
- Для тестов [`NewRefreshingTokenSourceWithOptions`](token_refresh.go) с полем [`Clock`](clock.go) в [`RefreshingSourceOptions`](token_refresh.go).

## Мейнтейнерам: генерация `gen/`

Нужны Go и сеть для загрузки модулей инструментов:

```bash
make generate
```

См. [CONTRIBUTING.md](CONTRIBUTING.md) про тесты и тег `integration`. Спека: [`api/openapi.yml`](api/openapi.yml). CI: [`.github/workflows/`](.github/workflows/).

## Интеграционные тесты

- `go test -short ./...` — быстрая проверка без живых вызовов HH ([`token_refresh_test.go`](token_refresh_test.go) на `httptest`).
- Живой API: `go test -tags=integration -timeout 5m ./integration` (без `-short`). С `HH_TEST_CLIENT_ID` и `HH_TEST_CLIENT_SECRET` из [dev.hh.ru](https://dev.hh.ru/admin) дополнительно выполняется OAuth-сценарий.
- Вручную: [`.github/workflows/integration-tests.yml`](.github/workflows/integration-tests.yml) (`workflow_dispatch`).

## Версионирование

После `v1` — семантическое версионирование. До этого минорные релизы при изменении `openapi.yml` / `gen/` могут ломать совместимость. См. [CHANGELOG.md](CHANGELOG.md).
