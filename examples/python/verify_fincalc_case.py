#!/usr/bin/env python3
# SPDX-License-Identifier: MIT

from __future__ import annotations

import argparse
import csv
import hashlib
import json
from dataclasses import dataclass
from pathlib import Path


def sha256_file(p: Path) -> str:
    h = hashlib.sha256()
    with p.open("rb") as f:
        for chunk in iter(lambda: f.read(65536), b""):
            h.update(chunk)
    return h.hexdigest()


def read_text_lf(p: Path) -> str:
    s = p.read_text(encoding="utf-8")
    if "\r\n" in s:
        raise AssertionError(f"CRLF found in {p}")
    return s


def load_json(p: Path):
    return json.loads(read_text_lf(p))


@dataclass(frozen=True)
class Row:
    period: int
    date: str
    payment_cents: int
    principal_cents: int
    interest_cents: int
    balance_cents: int


def parse_schedule_csv(p: Path) -> list[Row]:
    rows: list[Row] = []
    with p.open(newline="", encoding="utf-8") as f:
        r = csv.DictReader(f)
        want_fields = [
            "period",
            "date",
            "payment_cents",
            "principal_cents",
            "interest_cents",
            "balance_cents",
        ]
        assert r.fieldnames == want_fields, f"schedule header mismatch: {r.fieldnames}"
        for rec in r:
            rows.append(
                Row(
                    period=int(rec["period"]),
                    date=rec["date"],
                    payment_cents=int(rec["payment_cents"]),
                    principal_cents=int(rec["principal_cents"]),
                    interest_cents=int(rec["interest_cents"]),
                    balance_cents=int(rec["balance_cents"]),
                )
            )
    return rows


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("--out-root", default="out/demo", help="demo output root (contains <case>/...)")
    ap.add_argument("--case", default="case02_interest", help="case folder name")
    ap.add_argument("--compare-goldens", action="store_true", help="also compare output bytes to fixtures/expected")
    args = ap.parse_args()

    repo = Path(__file__).resolve().parents[2]
    case = args.case

    fixtures_exp = repo / "fixtures" / "expected" / case
    fixtures_in = repo / "fixtures" / "input" / case
    assert fixtures_exp.exists(), f"missing fixtures/expected/{case}"
    assert fixtures_in.exists(), f"missing fixtures/input/{case}"

    out_root = repo / args.out_root
    out_dir = out_root / case
    assert out_dir.exists(), f"missing out case dir: {out_dir}"

    # Expected-fail case: demo writes error.txt and compares it to fixtures.
    if (fixtures_exp / "error.txt").exists():
        got = (out_dir / "error.txt")
        assert got.exists(), f"missing error.txt: {got}"
        read_text_lf(got)
        want = (fixtures_exp / "error.txt").read_bytes()
        assert got.read_bytes() == want, "error.txt golden mismatch"
        print("OK: expected-fail case error.txt matches goldens.")
        return

    resp_p = out_dir / "response.json"
    sched_p = out_dir / "schedule.csv"
    for p in (resp_p, sched_p):
        assert p.exists(), f"missing output file: {p}"
        read_text_lf(p)

    req = load_json(fixtures_in / "request.json")
    resp = load_json(resp_p)
    rows = parse_schedule_csv(sched_p)

    # Contract echoes
    assert resp.get("schema_version") == "v1", "schema_version mismatch"
    assert resp.get("calculator") == "amortize", "calculator mismatch"
    for k in ("principal_cents", "annual_rate_bps", "term_months", "start_date"):
        assert resp.get(k) == req.get(k), f"echo field mismatch: {k}"

    term = int(resp["term_months"])
    assert len(rows) == term, "schedule row count != term_months"
    assert [r.period for r in rows] == list(range(1, term + 1)), "period sequence mismatch"

    # Totals and invariants
    principal_sum = sum(r.principal_cents for r in rows)
    interest_sum = sum(r.interest_cents for r in rows)
    payment_sum = sum(r.payment_cents for r in rows)

    assert principal_sum == int(resp["principal_cents"]), "principal sum mismatch"
    assert interest_sum == int(resp["total_interest_cents"]), "interest sum mismatch"
    assert payment_sum == int(resp["total_paid_cents"]), "total paid mismatch"

    assert rows[-1].balance_cents == 0, "final balance must be 0"

    pay = int(resp["payment_cents"])
    last_pay = int(resp["last_payment_cents"])
    for r in rows[:-1]:
        assert r.payment_cents == pay, "non-last payment must equal payment_cents"
    assert rows[-1].payment_cents == last_pay, "last payment must equal last_payment_cents"

    # Optional: byte-for-byte compare to checked-in goldens.
    if args.compare_goldens:
        for name in ("response.json", "schedule.csv"):
            want = (fixtures_exp / name).read_bytes()
            got = (out_dir / name).read_bytes()
            assert got == want, f"golden mismatch: {case}/{name}"

    print("OK: fincalc demo outputs are internally consistent.")


if __name__ == "__main__":
    main()
