# proof-first-finance-calc

A tiny finance calculator HTTP service built with the **proof-first** pattern from *The Deterministic Finance Toolkit*.

This repo is intentionally small: **one calculator, done right**.

## Quickstart (dev)

```bash
go test -count=1 ./...
go run ./cmd/fincalc demo --out ./out/demo
```

## CLI

- `demo` writes deterministic example outputs to `--out`
- `serve` starts the HTTP server

```bash
go run ./cmd/fincalc --help
```

## Proof gates

- `go test -count=1 ./...` passes
- CI runs the same gate on ubuntu/macos/windows
