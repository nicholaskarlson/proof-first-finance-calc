package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholaskarlson/proof-first-finance-calc/internal/calc"
)

type goldenCase struct {
	name string
}

func TestAmortizeV1_Goldens(t *testing.T) {
	cases := []goldenCase{
		{name: "case01_zero_rate"},
		{name: "case02_interest"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			inPath := filepath.Join("..", "fixtures", "input", tc.name, "request.json")
			wantRespPath := filepath.Join("..", "fixtures", "expected", tc.name, "response.json")
			wantCSVPath := filepath.Join("..", "fixtures", "expected", tc.name, "schedule.csv")

			inB, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatalf("read input request: %v", err)
			}
			var req calc.AmortizeRequestV1
			if err := json.Unmarshal(inB, &req); err != nil {
				t.Fatalf("unmarshal request: %v", err)
			}

			resp, rows, err := calc.AmortizeV1(req)
			if err != nil {
				t.Fatalf("AmortizeV1: %v", err)
			}

			// Invariants (proof-first): totals must tie out.
			assertScheduleInvariants(t, req, resp, rows)

			gotResp, err := calc.RenderResponseJSON(resp)
			if err != nil {
				t.Fatalf("render response json: %v", err)
			}
			wantResp, err := os.ReadFile(wantRespPath)
			if err != nil {
				t.Fatalf("read expected response: %v", err)
			}
			if string(gotResp) != string(wantResp) {
				t.Fatalf("response.json mismatch\n--- got ---\n%s\n--- want ---\n%s", gotResp, wantResp)
			}

			gotCSV, err := calc.RenderScheduleCSV(rows)
			if err != nil {
				t.Fatalf("render schedule csv: %v", err)
			}
			wantCSV, err := os.ReadFile(wantCSVPath)
			if err != nil {
				t.Fatalf("read expected schedule: %v", err)
			}
			if string(gotCSV) != string(wantCSV) {
				t.Fatalf("schedule.csv mismatch\n--- got ---\n%s\n--- want ---\n%s", gotCSV, wantCSV)
			}
		})
	}
}

func assertScheduleInvariants(t *testing.T, req calc.AmortizeRequestV1, resp calc.AmortizeResponseV1, rows []calc.ScheduleRow) {
	t.Helper()
	if len(rows) != req.TermMonths {
		t.Fatalf("expected %d rows, got %d", req.TermMonths, len(rows))
	}
	if rows[len(rows)-1].BalanceCents != 0 {
		t.Fatalf("final balance must be 0, got %d", rows[len(rows)-1].BalanceCents)
	}
	var sumPrincipal, sumInterest, sumPaid int64
	prevBal := req.PrincipalCents
	for _, r := range rows {
		sumPrincipal += r.PrincipalCents
		sumInterest += r.InterestCents
		sumPaid += r.PaymentCents
		if r.BalanceCents > prevBal {
			t.Fatalf("balance must be non-increasing, saw %d -> %d", prevBal, r.BalanceCents)
		}
		prevBal = r.BalanceCents
	}
	if sumPrincipal != req.PrincipalCents {
		t.Fatalf("principal tie-out failed: sum principal %d != principal %d", sumPrincipal, req.PrincipalCents)
	}
	if sumInterest != resp.TotalInterestCents {
		t.Fatalf("interest tie-out failed: sum interest %d != resp.total_interest_cents %d", sumInterest, resp.TotalInterestCents)
	}
	if sumPaid != resp.TotalPaidCents {
		t.Fatalf("paid tie-out failed: sum paid %d != resp.total_paid_cents %d", sumPaid, resp.TotalPaidCents)
	}
	if rows[len(rows)-1].PaymentCents != resp.LastPaymentCents {
		t.Fatalf("last payment mismatch: schedule %d != resp.last_payment_cents %d", rows[len(rows)-1].PaymentCents, resp.LastPaymentCents)
	}
}
