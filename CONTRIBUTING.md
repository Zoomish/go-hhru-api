# Contributing

- Основная документация репозитория на русском: [`README.md`](README.md) (и зеркало [`README.ru.md`](README.ru.md)). Английская версия: [`README.en.md`](README.en.md).
- Regenerate `gen/` after changing [`api/openapi.yml`](api/openapi.yml): `make generate` (requires Go and network for tool modules).
- Run `go test -short ./...` before opening a PR.
- Live API tests live under [`integration/`](integration/) and only compile with `-tags=integration`; they are not part of the default [CI](.github/workflows/ci.yml) pipeline.

## Релизы и CI

- **Один конвейер проверок:** [`.github/workflows/ci.yml`](.github/workflows/ci.yml) — job `checks` вызывает общую [composite action](.github/actions/ci-checks/action.yml) (`go test -short`, `go vet`, регенерация `gen/`, сравнение с коммитом). Запускается на `pull_request`, на `push` в ветки и вручную (`workflow_dispatch`).
- **Релиз после принятия изменений:** после мержа в основную ветку поставьте SemVer-тег (`git tag v0.x.y` на нужный коммит и `git push origin v0.x.y`). Workflow [`.github/workflows/release.yml`](.github/workflows/release.yml) сначала выполняет те же проверки, что и CI ([composite action](.github/actions/ci-checks/action.yml)), затем создаёт [GitHub Release](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository) для этого тега.
- Полностью автоматический bump версии при каждом merge без тега в репозиторий не встроен (при необходимости — Release Please / semantic-release).
