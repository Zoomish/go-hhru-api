// Package hhru provides a typed HTTP client for the HeadHunter API (api.hh.ru).
//
// Use [New] with [Options].HHUserAgent set to your application name and contact email
// as required by HH. Optional [Options.TokenSource] adds Bearer authorization to every request.
// Use [ExchangeClientCredentials] and [AccessToken] for the application token flow.
// Use [NewRefreshingTokenSource] for access tokens obtained with a refresh_token (user OAuth).
//
// Low-level generated clients live under gen/ (employer, applicant, public, app).
//
// Official OpenAPI documentation: https://api.hh.ru/openapi/redoc
package hhru
