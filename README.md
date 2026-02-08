# proof-first-finance-calc

A tiny **proof-first** finance calculator service (Go-first).

![ci](https://github.com/nicholaskarlson/proof-first-finance-calc/actions/workflows/ci.yml/badge.svg)
![license](https://img.shields.io/badge/license-MIT-blue.svg)

> **Book:** *The Deterministic Finance Toolkit*
> This repo is **Project 4 of 4**. The exact code referenced in the manuscript is tagged **[`book-v1`](https://github.com/nicholaskarlson/proof-first-finance-calc/tree/book-v1)**.

## Toolkit navigation

- **[proof-first-recon](https://github.com/nicholaskarlson/proof-first-recon)** — deterministic CSV reconciliation (matched/unmatched + summary JSON)
- **[proof-first-auditpack](https://github.com/nicholaskarlson/proof-first-auditpack)** — deterministic audit packs (manifest.json + sha256 + verify)
- **[proof-first-normalizer](https://github.com/nicholaskarlson/proof-first-normalizer)** — deterministic CSV normalize + validate (schema → normalized.csv/errors.csv/report.json)
- **[proof-first-finance-calc](https://github.com/nicholaskarlson/proof-first-finance-calc)** — proof-first finance calc service (Amortization v1 API + demo)

## What it does

`fincalc` implements **Amortization v1** (fixed-rate, monthly payments).

It exposes the calculator in two ways:

1) **HTTP API**
- `POST /v1/amortize` → JSON response
- `POST /v1/amortize/schedule.csv` → CSV schedule

2) **Local demo**
- `go run ./cmd/fincalc demo --out ./out` writes deterministic outputs derived from fixtures and verifies they match the golden files.

## Quick start

Requirements:
- Go **1.22+**
- GNU Make (optional, but recommended)

```bash
# One-command proof gate
make verify

# Portable proof gate (no Makefile)
go test -count=1 ./...
go run ./cmd/fincalc demo --out ./out
```


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

- `cmd/fincalc/` — CLI entrypoint (`demo`, `serve`, `version`)
- `internal/calc/` — deterministic amortization core + renderers
- `internal/api/` — HTTP handlers
- `fixtures/` — input cases + golden outputs
- `tests/` — golden + API tests

## Determinism contract

This project is intentionally “boring” in the best way: the same inputs must produce the same outputs.

See: **[`docs/CONVENTIONS.md`](docs/CONVENTIONS.md)** (rounding, ordering, LF, atomic writes, stable JSON, etc.).


## Handoff / maintenance

See: **[`docs/HANDOFF.md`](docs/HANDOFF.md)** (acceptance gates, troubleshooting, and “what to change (and what not to)”).


## License

MIT (see `LICENSE`).

