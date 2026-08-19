# Repository Guidelines

## Project Structure & Module Organization

EasyCache is a Go 1.25 module (`github.com/hugocarreira/easycache`). The public `cache/` package selects an eviction strategy and exposes configuration and metrics. Strategy implementations live in `basic/`, `fifo/`, `lru/`, and `lfu/`; shared behavior is defined in `engine/`. Integration-style tests and benchmarks are in `tests/`. Repository documentation is maintained in `README.md` and `CONTRIBUTING.md`, while GitHub Actions configuration lives in `.github/workflows/tests.yml`.

Keep policy-specific logic inside its strategy package. Changes to the public API belong in `cache/`, with corresponding tests in `tests/`.

## Build, Test, and Development Commands

- `go test ./...` runs every package and the complete test suite.
- `go test ./... -v` runs the complete test suite used by CI.
- `go test -race ./...` checks concurrent cache operations for data races.
- `go test -coverpkg=./... -coverprofile=/tmp/easycache.cover ./... && go tool cover -func=/tmp/easycache.cover` reports aggregate statement coverage across the module.
- `go test -run '^$' -bench=. -benchmem ./tests` runs benchmarks without rerunning tests and includes allocation statistics.
- `go vet ./...` performs standard static analysis.
- `gofmt -w cache/*.go tests/*.go` formats changed Go files; include other edited package paths as needed.

There is no separate build system: use standard Go tooling and run commands from the repository root.

## Coding Style & Naming Conventions

Follow idiomatic Go and let `gofmt` control tabs and layout. Package names are short and lowercase. Exported identifiers use `PascalCase` and require clear GoDoc comments; unexported identifiers use `camelCase`. Keep constructors conventional (`New`) and configuration fields descriptive (`CleanupInterval`, `EvictionPolicy`). Prefer small, policy-focused changes over cross-package duplication.

## Testing Guidelines

Tests use Go's `testing` package with `testify/assert` and `testify/suite`. Name files `*_test.go`, test entry points `TestXxx`, and benchmarks `BenchmarkXxx`. Add coverage for normal behavior, eviction boundaries, expiration, and concurrency when relevant. No numeric coverage threshold is enforced, but every behavior change should include a regression test. Run tests, race detection, and `go vet` before submitting.

## Commit & Pull Request Guidelines

Prefix every commit subject with a parenthesized change type: `(feat)` for features, `(fix)` for bug fixes, `(docs)` for documentation, `(test)` for tests, `(refactor)` for internal restructuring, and `(chore)` for maintenance. Follow the prefix with a short, lowercase, imperative description, for example `(fix) prevent expired cache hits` or `(docs) clarify memory limits`. Keep each commit focused and avoid vague messages. Open pull requests against `main`. Include a concise change summary, motivation or linked issue, and the commands used to verify the work. Document benchmark impact for performance-sensitive changes and update examples when the public API changes.
