# Handoff

This repo is intentionally small: one calculator (Amortization v1) with a strict contract and fixture-based proof gates.

## Canonical commands

```bash
# Proof gate (one command)
make verify

# Proof gates (portable, no Makefile)
go test -count=1 ./...
go run ./cmd/fincalc demo --out ./out/demo
```

## Input contract (Amortize v1)

JSON request body:

- `principal_cents` (int, > 0)
- `annual_rate_bps` (int, >= 0)
- `term_months` (int, > 0)
- `start_date` (YYYY-MM-DD)

Schedule dates:

- The first schedule row date equals `start_date`.
- Each subsequent row advances by one calendar month using Go's `time.Time.AddDate(0, 1, 0)` semantics.
- No business-day or end-of-month adjustments are applied.

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

## Optional: Python check (stdlib only)

```bash
# Run the Go demo first (writes outputs and verifies goldens)
go run ./cmd/fincalc demo --out ./out/demo

# Then run the optional Python verifier on one case
python3 examples/python/verify_fincalc_case.py --out-root ./out/demo --case case02_interest
```


## Serve the HTTP API

```bash
# Local development server
go run ./cmd/fincalc serve --addr :8080
```
