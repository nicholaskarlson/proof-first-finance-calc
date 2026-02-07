# Optional Python verification (stdlib only)

These scripts are **optional** and use only the Python standard library.
They provide a second, independent check that demo outputs are internally consistent.

Run the Go demo first:

```bash
go run ./cmd/fincalc demo --out ./out/demo
```

Then verify one case:

```bash
python3 examples/python/verify_fincalc_case.py --out-root ./out/demo --case case02_interest
```
