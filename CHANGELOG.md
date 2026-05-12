# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- [NewRefreshingTokenSource](token_refresh.go): thread-safe `TokenSource` that refreshes access tokens with `ExchangeRefreshToken` before expiry.

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
