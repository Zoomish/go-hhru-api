# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- [NewRefreshingTokenSource](token_refresh.go): thread-safe `TokenSource` that refreshes access tokens with `ExchangeRefreshToken` before expiry.
- MIT [`LICENSE`](LICENSE).
- Runnable [`examples/`](examples/) (`public_countries`, `public_locales`, `public_areas`, `public_industries`, `public_languages`, `public_position_suggest`, `custom_options`, `app_token`, `refreshing_token`).
- [`integration/`](integration/) live tests behind `-tags=integration`; optional [`.github/workflows/integration-tests.yml`](.github/workflows/integration-tests.yml) (`workflow_dispatch`).
- Godoc runnable `Example*` tests in [`example_test.go`](example_test.go); [`CONTRIBUTING.md`](CONTRIBUTING.md).
- Repository docs: Russian [`README.md`](README.md) / [`README.ru.md`](README.ru.md) and English [`README.en.md`](README.en.md); TOC; CI badge via [`.github/workflows/ci.yml`](.github/workflows/ci.yml).
- [`SECURITY.md`](SECURITY.md); [Dependabot](.github/dependabot.yml) for Go modules.
- [`ParseAPIError`](api_error.go) / [`APIError`](api_error.go) for typical HH JSON errors with `request_id`.
- [`Options`](option.go): [`RequestHook`](option.go), [`ResponseHook`](option.go), [`MaxRequestsPerSecond`](option.go), [`RetryBackoffMin`](option.go) / [`RetryBackoffMax`](option.go) for transport pacing and retry backoff clamping.
- [`Clock`](clock.go) and [`NewRefreshingTokenSourceWithOptions`](token_refresh.go) with [`RefreshingSourceOptions`](token_refresh.go) for testable token refresh timing.
- [`.github/workflows/release.yml`](.github/workflows/release.yml): on `v*` tag push, runs the [ci-checks composite action](.github/actions/ci-checks/action.yml), then creates a GitHub Release.

### Changed

- Live API tests moved from root into the tagged `integration` package so `go test -short ./...` stays free of HH calls.
- GitHub Actions: merged PR checks into a single [CI](.github/workflows/ci.yml) job that runs the shared [composite action](.github/actions/ci-checks/action.yml); [Release](.github/workflows/release.yml) runs the same action before publishing; removed separate `test.yml` and `verify-codegen.yml`.
- Removed standalone `scripts/patch-applicant` tool (applicant `client.gen.go` patch is applied only from [`scripts/generate`](scripts/generate/main.go)).
- [`New`](client.go) composes optional pacing, response hooks, and retries on the HTTP transport when those options are set.

### For maintainers

- After merging, update [CHANGELOG.md](CHANGELOG.md), tag the merge commit (`git tag v0.2.0` and `git push origin v0.2.0`). [Release](.github/workflows/release.yml) runs the same [ci-checks](.github/actions/ci-checks/action.yml) as CI, then creates a GitHub Release so [proxy.golang.org](https://proxy.golang.org) and [pkg.go.dev](https://pkg.go.dev/github.com/Zoomish/go-hhru-api) can index the version.

## [0.1.0] - 2026-05-12

### Added

- Facade `hhru.New` / `Client` with `Employer`, `Applicant`, `Public`, `App` generated clients.
- `Options`: required `HHUserAgent`, optional `DefaultHost`, `DefaultLocale`, `TokenSource`, `MaxRetries`.
- Global request editor: default query `host` / `locale`, `HH-User-Agent` when absent, `Authorization: Bearer` from `TokenSource`.
- `ExchangeClientCredentials`, `ExchangeRefreshToken`, `TokenEndpoint`, `AccessToken` helper.
- `PagesUntil` pagination helper.
- Optional retry `RoundTripper` for `429` / `503`.
- `make generate` via unified [`scripts/generate`](scripts/generate/main.go).
- GitHub Actions: verify `gen/` drift; optional upstream OpenAPI compare.

[Unreleased]: https://github.com/Zoomish/go-hhru-api/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Zoomish/go-hhru-api/releases/tag/v0.1.0
