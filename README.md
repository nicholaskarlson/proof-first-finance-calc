# proof-first-finance-calc

A tiny **proof-first** finance calculator service.

This repo is Project 4 in *The Deterministic Finance Toolkit*.

## What it does

`fincalc` implements **Amortization v1** (fixed-rate, monthly payments).

It exposes the calculator in two ways:

1) **HTTP API**
- `POST /v1/amortize` → JSON response
- `POST /v1/amortize/schedule.csv` → CSV schedule

2) **Local demo**
- `go run ./cmd/fincalc demo --out ./out/demo` writes deterministic outputs derived from fixtures and verifies they match the golden files.

## Canonical commands

```bash
# Proof gate (one command)
make verify

# Proof gates (portable, no Makefile)
go test -count=1 ./...
go run ./cmd/fincalc demo --out ./out/demo
```

## Demo

```bash
go run ./cmd/fincalc demo --out ./out/demo
```

On success it writes one folder per fixture case under `./out/demo/`.

## Run the HTTP API

```bash
go run ./cmd/fincalc serve --addr 127.0.0.1:8080
```

Example request (JSON):

```bash
 curl -sS -X POST http://127.0.0.1:8080/v1/amortize \
   -H 'Content-Type: application/json' \
   --data-binary @fixtures/input/case02_interest/request.json
```

Example request (CSV):

```bash
 curl -sS -X POST http://127.0.0.1:8080/v1/amortize/schedule.csv \
   -H 'Content-Type: application/json' \
   --data-binary @fixtures/input/case02_interest/request.json
```

## Repo layout

- `cmd/fincalc/` — CLI entrypoint (`demo`, `serve`)
- `internal/calc/` — deterministic amortization core + renderers
- `internal/api/` — HTTP handlers
- `fixtures/` — input cases + golden outputs
- `tests/` — golden + API tests
