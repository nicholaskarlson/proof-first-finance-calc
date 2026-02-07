package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/nicholaskarlson/proof-first-finance-calc/internal/api"
)

func TestHTTPAPI_V1_Amortize_Fixtures(t *testing.T) {
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

	srv := httptest.NewServer(api.Handler())
	defer srv.Close()

	for _, c := range caseNames {
		c := c
		t.Run(c, func(t *testing.T) {
			reqPath := filepath.Join(inRoot, c, "request.json")
			body, err := os.ReadFile(reqPath)
			if err != nil {
				t.Fatalf("read request: %v", err)
			}

			expDir := filepath.Join(root, "expected", c)
			errPath := filepath.Join(expDir, "error.txt")
			_, errStat := os.Stat(errPath)
			hasErr := errStat == nil

			// JSON endpoint
			check := func(path string, wantFile string, wantStatus int) {
				r, err := http.Post(srv.URL+path, "application/json", bytes.NewReader(body))
				if err != nil {
					t.Fatalf("POST %s: %v", path, err)
				}
				defer r.Body.Close()

				got, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read response: %v", err)
				}

				if r.StatusCode != wantStatus {
					t.Fatalf("status mismatch: got %d want %d\nbody:\n%s", r.StatusCode, wantStatus, string(got))
				}

				want, err := os.ReadFile(filepath.Join(expDir, wantFile))
				if err != nil {
					t.Fatalf("read expected %s: %v", wantFile, err)
				}
				if !bytes.Equal(got, want) {
					t.Fatalf("%s mismatch\n--- got ---\n%s\n--- want ---\n%s", wantFile, string(got), string(want))
				}
			}

			if hasErr {
				check("/v1/amortize", "error.txt", http.StatusBadRequest)
				check("/v1/amortize/schedule.csv", "error.txt", http.StatusBadRequest)
				return
			}

			check("/v1/amortize", "response.json", http.StatusOK)
			check("/v1/amortize/schedule.csv", "schedule.csv", http.StatusOK)
		})
	}
}
