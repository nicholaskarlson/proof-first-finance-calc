# HANDOFF — proof-first-finance-calc

This repo is designed to be extended by another developer without guesswork.

## Contract

This tool follows the book’s pattern:

contract → outputs → fixtures/goldens → tests → demo → handoff/docs.

The calculator contract and endpoint(s) are introduced in the README and fixtures once implemented.

## Canonical commands

From repo root:

```bash
go test -count=1 ./...
go run ./cmd/fincalc demo --out ./out/demo
```

Optional build:

```bash
make build
./bin/fincalc demo --out ./out/demo
```

## Repo layout

- `cmd/fincalc/` CLI + HTTP entrypoint
- `internal/` deterministic core logic (no I/O)
- `fixtures/` inputs + goldens
- `tests/` golden tests and invariants
- `docs/CONVENTIONS.md` shared proof-first conventions
