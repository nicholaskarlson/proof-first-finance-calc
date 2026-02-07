# Handoff

This repo is intentionally small: one calculator (Amortization v1) with a strict contract and fixture-based proof gates.

## Canonical commands

```bash
# Proof gates
go test -count=1 ./...

# Demo (writes outputs and verifies against fixtures)
go run ./cmd/fincalc demo --out ./out/demo

# Serve HTTP API
go run ./cmd/fincalc serve --addr 127.0.0.1:8080
```

## Input contract (Amortize v1)

JSON request body:

- `principal_cents` (int, > 0)
- `annual_rate_bps` (int, >= 0)
- `term_months` (int, > 0)
- `start_date` (YYYY-MM-DD)

## Output contract

### HTTP

- `POST /v1/amortize` returns `application/json` (the amortization summary)
- `POST /v1/amortize/schedule.csv` returns `text/csv` (the payment schedule)

On error, the API responds with status `400` and a stable one-line body:

```
error: <message>
```

### Demo output

`fincalc demo --out <dir>` writes one folder per fixture case:

- `response.json`
- `schedule.csv`

Expected-fail cases write:

- `error.txt`

## Extending safely

- Add a new calculator under `internal/calc/`.
- Add fixtures under `fixtures/input/<case>/request.json`.
- Check in goldens under `fixtures/expected/<case>/`.
- Add tests in `tests/` (goldens first, then API coverage).
