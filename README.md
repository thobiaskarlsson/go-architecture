# Go Hexagonal Telemetry Playground

This repository exists to practice **hexagonal architecture** and experiment with design/testing patterns in a small, safe codebase.

## Why This Project Exists

- Learn by doing, not by building production infrastructure.
- Keep the domain tiny (vehicle telemetry) so architecture decisions stay visible.
- Prefer unit tests and clear boundaries over framework-heavy setup.

## Architecture At A Glance

- **Domain**: `internal/domain`
- **Application services**: `internal/application`
- **Ports**: `internal/ports`
- **Adapters**: `internal/adapters/inbound|outbound|shared`
- **Composition root**: `cmd/telemetry-playground/main.go`

Dependency direction is intentional:
- `domain` has no adapter dependencies.
- `ports` define interfaces used by application services.
- inbound/outbound adapters depend on ports.
- `cmd/.../main.go` wires concrete implementations.

## Key Decisions And Whys

### 1) Domain vs Ports

- `domain` contains domain structs and domain-level logic only.
- `ports` contains interfaces (input/output contracts).

**Why**: keeps the business core framework-agnostic and easy to unit test.

### 2) This App Is A Proxy, Not A Read Model

The service forwards telemetry northbound instead of persisting/querying it.
- Outbound port is `TelemetryForwarderPort` (`internal/ports/outbound.go`).

**Why**: the project focus shifted from storage to middleware/proxy behavior.

### 3) Adapter Naming (`http`, `nats`) Even In Both Directions

Both inbound and outbound adapters use transport-oriented names:
- `internal/adapters/inbound/http`
- `internal/adapters/outbound/http`
- `internal/adapters/inbound/nats`
- `internal/adapters/outbound/nats`

**Why**:
- In Go, package identity is the full import path, not only the package name.
- Direction is already explicit in folder path (`inbound` vs `outbound`).
- In wiring code we alias imports for readability (`httpin`, `httpout`, `natsin`, `natsout`).

### 4) Shared Transport Packages

Shared transport constants/helpers live under:
- `internal/adapters/shared/http`
- `internal/adapters/shared/nats`

Examples:
- HTTP paths/headers/constants
- NATS subject constants and subject helpers

**Why**: avoids duplication while preventing inbound/outbound adapters from importing each other directly.

### 5) Service Dependency Pattern (`serviceDependencies`)

`internal/application/telemetry/service.go` uses:
- `newServiceWithDependencies(forwarder, deps)`
- `serviceDependencies` struct for collaborators

Current collaborator:
- `normalizer`

Defaults are applied inside constructor if collaborators are omitted.

**Why**:
- Tests can override exactly one collaborator without noisy setup.
- Constructor remains explicit and close to production wiring.

### 6) Fail Fast On Invalid Wiring

`newServiceWithDependencies` panics if `forwarder` is `nil`.

**Why**: wiring mistakes are programmer errors; fail immediately instead of panicking later deep in request handling.

### 7) Test Strategy And Naming

Telemetry service has two test styles:
- `service_test.go` (package `telemetry`): white-box/internal tests for collaborators and constructor behavior.
- `service_blackbox_test.go` (package `telemetry_test`): public API behavior from outside the package boundary.

**Why**:
- `service_test.go` can directly exercise internal constructors/collaborators.
- `service_blackbox_test.go` protects public behavior and avoids over-coupling to internals.
- Name `service_blackbox_test.go` was chosen intentionally (clearer than `integration` for this scope).

### 8) `cmd/<app>/main.go` Layout

Entrypoint is in `cmd/telemetry-playground/main.go` (not `cmd/main.go`).

**Why**: idiomatic Go layout that scales to multiple binaries.

## How To Extend This Playground

- Add new collaborator to `serviceDependencies` when needed.
- Keep defaults in constructor for low-friction tests.
- Add outbound adapters implementing `TelemetryForwarderPort`.
- Reuse `internal/adapters/shared/*` only for true transport-shared concerns.
- Keep domain free from adapter imports.

## Commands

Run tests:

```bash
env -u GOROOT go test ./...
```

Run playground:

```bash
go run ./cmd/telemetry-playground
```
