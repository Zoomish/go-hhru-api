# Contributing

- Regenerate `gen/` after changing [`api/openapi.yml`](api/openapi.yml): `make generate` (requires Go and network for tool modules).
- Run `go test -short ./...` before opening a PR.
- Live API tests live under [`integration/`](integration/) and only compile with `-tags=integration`; they are not part of the default PR workflow in [`.github/workflows/verify-codegen.yml`](.github/workflows/verify-codegen.yml).
