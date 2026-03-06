# Repository Guidelines

## Project Structure & Module Organization
This repository is currently scaffold-only (IDE metadata under `.idea/`). As code is added, keep a standard Go layout to keep responsibilities clear:
- `cmd/<app>/` for executable entrypoints.
- `internal/` for private application logic.
- `pkg/` for reusable public packages.
- `test/` for integration/e2e tests.
- `configs/` for versioned config templates (no secrets).
- `docs/` for architecture notes and ADRs.

Example: `cmd/api/main.go`, `internal/service/user_service.go`, `test/api_integration_test.go`.

## Build, Test, and Development Commands
Use Go toolchain commands from repo root:
- `go mod tidy` syncs module dependencies.
- `go test ./...` runs all unit tests.
- `go test -race ./...` checks for race conditions.
- `go vet ./...` catches suspicious constructs.
- `golangci-lint run` runs static analysis (if configured).
- `go run ./cmd/<app>` runs a local service.

## Coding Style & Naming Conventions
Follow idiomatic Go:
- Format with `gofmt` (or `go fmt ./...`) before committing.
- Keep package names short, lowercase, and noun-based.
- Exported identifiers use `PascalCase`; internal identifiers use `camelCase`.
- File names use lowercase and descriptive suffixes (e.g., `user_repo.go`, `handler_test.go`).
- Prefer small interfaces near usage sites, not broad shared interfaces.

## Testing Guidelines
- Use the standard `testing` package with table-driven tests where practical.
- Name tests `Test<Behavior>` and benchmarks `Benchmark<Behavior>`.
- Keep fast unit tests close to source as `*_test.go`; place cross-module tests in `test/`.
- For changed packages, target meaningful coverage (aim ~80% where feasible) and include edge cases.

## Commit & Pull Request Guidelines
No established git history is present yet, so use Conventional Commits:
- `feat: add user repository`
- `fix: handle nil config in bootstrap`
- `test: cover retry backoff behavior`

PRs should include:
- Clear summary and scope.
- Linked issue/ticket.
- Test evidence (command + result).
- Notes on config, migrations, or breaking changes.

## Security & Configuration Tips
- Never commit secrets, private keys, or `.env` files.
- Commit only templates like `.env.example`.
- Validate all external input at boundaries (HTTP handlers, message consumers).
