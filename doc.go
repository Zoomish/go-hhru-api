// Package hhru provides a typed HTTP client for the HeadHunter API (api.hh.ru).
//
// Use [New] with [Options].HHUserAgent set to your application name and contact email
// as required by HH. Optional [Options.TokenSource] adds Bearer authorization to every request.
// Use [ExchangeClientCredentials] and [AccessToken] for the application token flow.
// Use [NewRefreshingTokenSource] (or [NewRefreshingTokenSourceWithOptions] with a [Clock] for tests)
// for access tokens obtained with a refresh_token (user OAuth).
//
// Generated sub-clients (from the published OpenAPI split) live under:
//   - [github.com/Zoomish/go-hhru-api/gen/employer] — employer API;
//   - [github.com/Zoomish/go-hhru-api/gen/applicant] — applicant auth and user flows;
//   - [github.com/Zoomish/go-hhru-api/gen/public] — public dictionaries and search helpers;
//   - [github.com/Zoomish/go-hhru-api/gen/app] — application-scoped endpoints.
//
// Each is exposed on [Client] as Employer, Applicant, Public, and App.
//
// Reliability and observability [Options]: optional [Options.MaxRetries] with [Options.RetryBackoffMin]
// and [Options.RetryBackoffMax], [Options.MaxRequestsPerSecond] pacing, [Options.RequestHook],
// and [Options.ResponseHook]. Parse JSON error bodies with [ParseAPIError].
//
// Runnable programs live under the examples/ directory (see the repository README).
// Live HTTP tests use build tag "integration" in the integration/ package.
//
// Official OpenAPI documentation: https://api.hh.ru/openapi/redoc
package hhru
