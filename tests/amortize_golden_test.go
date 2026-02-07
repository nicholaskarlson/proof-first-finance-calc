package tests

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/nicholaskarlson/proof-first-finance-calc/internal/calc"
)

func TestAmortizeV1_Goldens(t *testing.T) {
	root := filepath.Join("..", "fixtures")
	inRoot := filepath.Join(root, "input")

	entries, err := os.ReadDir(inRoot)
	if err != nil {
		t.Fatalf("read fixtures input: %v", err)
	}

	caseNames := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			caseNames = append(caseNames, e.Name())
		}
	}
	sort.Strings(caseNames)

	for _, c := range caseNames {
		c := c
		t.Run(c, func(t *testing.T) {
			inPath := filepath.Join(inRoot, c, "request.json")
			inB, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatalf("read input request: %v", err)
			}

			var req calc.AmortizeRequestV1
			if err := json.Unmarshal(inB, &req); err != nil {
				t.Fatalf("unmarshal request: %v", err)
			}

			expDir := filepath.Join(root, "expected", c)
			errPath := filepath.Join(expDir, "error.txt")
			if _, statErr := os.Stat(errPath); statErr == nil {
				// Expected-fail case: compare stable one-line error body.
				_, _, err := calc.AmortizeV1(req)
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				got := []byte("error: " + err.Error() + "\n")
				want, err := os.ReadFile(errPath)
				if err != nil {
					t.Fatalf("read expected error: %v", err)
				}
				if !bytes.Equal(got, want) {
					t.Fatalf("error.txt mismatch\n--- got ---\n%s\n--- want ---\n%s", string(got), string(want))
				}
				return
			}

			wantRespPath := filepath.Join(expDir, "response.json")
			wantCSVPath := filepath.Join(expDir, "schedule.csv")

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
			if !bytes.Equal(gotResp, wantResp) {
				t.Fatalf("response.json mismatch\n--- got ---\n%s\n--- want ---\n%s", string(gotResp), string(wantResp))
			}

			gotCSV, err := calc.RenderScheduleCSV(rows)
			if err != nil {
				t.Fatalf("render schedule csv: %v", err)
			}
			wantCSV, err := os.ReadFile(wantCSVPath)
			if err != nil {
				t.Fatalf("read expected schedule: %v", err)
			}
			if !bytes.Equal(gotCSV, wantCSV) {
				t.Fatalf("schedule.csv mismatch\n--- got ---\n%s\n--- want ---\n%s", string(gotCSV), string(wantCSV))
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
