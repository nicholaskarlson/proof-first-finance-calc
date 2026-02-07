package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/nicholaskarlson/proof-first-finance-calc/internal/api"
	"github.com/nicholaskarlson/proof-first-finance-calc/internal/calc"
	"github.com/nicholaskarlson/proof-first-finance-calc/internal/fsutil"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	sub := os.Args[1]
	args := os.Args[2:]

	var err error
	switch sub {
	case "demo":
		err = cmdDemo(args)
	case "serve":
		err = cmdServe(args)
	case "help", "-h", "--help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", sub)
		usage()
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `fincalc - proof-first finance calculators

Usage:
  fincalc demo  --out <dir> [--fixtures fixtures]
  fincalc serve --addr <host:port>

Commands:
  demo   Recompute known cases from fixtures and verify outputs match goldens.
  serve  Run the HTTP API server (v1).
`)
}

func cmdDemo(args []string) error {
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	outDir := fs.String("out", "", "Output directory")
	fixtures := fs.String("fixtures", "fixtures", "Fixtures root")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *outDir == "" {
		return fmt.Errorf("--out is required")
	}

	inRoot := filepath.Join(*fixtures, "input")
	entries, err := os.ReadDir(inRoot)
	if err != nil {
		return fmt.Errorf("read fixtures: %w", err)
	}

	cases := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			cases = append(cases, e.Name())
		}
	}
	sort.Strings(cases)
	if len(cases) == 0 {
		return fmt.Errorf("no fixture cases found under %s", inRoot)
	}

	for _, c := range cases {
		if err := runCase(*fixtures, *outDir, c); err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stdout, "OK: demo outputs match fixtures (%d case(s))\n", len(cases))
	return nil
}

func runCase(fixturesRoot, outRoot, caseName string) error {
	reqPath := filepath.Join(fixturesRoot, "input", caseName, "request.json")
	b, err := os.ReadFile(reqPath)
	if err != nil {
		return fmt.Errorf("%s: read request.json: %w", caseName, err)
	}
	var req calc.AmortizeRequestV1
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return fmt.Errorf("%s: invalid JSON", caseName)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return fmt.Errorf("%s: invalid JSON", caseName)
	}

	expectedDir := filepath.Join(fixturesRoot, "expected", caseName)
	outDir := filepath.Join(outRoot, caseName)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("%s: mkdir out: %w", caseName, err)
	}

	wantErrPath := filepath.Join(expectedDir, "error.txt")
	if wantErr, errRead := os.ReadFile(wantErrPath); errRead == nil {
		// expected-fail case
		_, _, errCalc := calc.AmortizeV1(req)
		if errCalc == nil {
			return fmt.Errorf("%s: expected error, got nil", caseName)
		}
		gotErr := []byte(fmt.Sprintf("error: %s\n", errCalc.Error()))
		if err := fsutil.AtomicWriteFile(filepath.Join(outDir, "error.txt"), gotErr, 0o644); err != nil {
			return fmt.Errorf("%s: write error.txt: %w", caseName, err)
		}
		if !bytes.Equal(gotErr, wantErr) {
			return fmt.Errorf("%s: error.txt mismatch", caseName)
		}
		return nil
	}

	resp, sched, err := calc.AmortizeV1(req)
	if err != nil {
		return fmt.Errorf("%s: %w", caseName, err)
	}

	respJSON, err := calc.RenderResponseJSON(resp)
	if err != nil {
		return fmt.Errorf("%s: render response: %w", caseName, err)
	}
	schedCSV, err := calc.RenderScheduleCSV(sched)
	if err != nil {
		return fmt.Errorf("%s: render schedule: %w", caseName, err)
	}

	// write outputs (for humans)
	if err := fsutil.AtomicWriteFile(filepath.Join(outDir, "response.json"), respJSON, 0o644); err != nil {
		return fmt.Errorf("%s: write response.json: %w", caseName, err)
	}
	if err := fsutil.AtomicWriteFile(filepath.Join(outDir, "schedule.csv"), schedCSV, 0o644); err != nil {
		return fmt.Errorf("%s: write schedule.csv: %w", caseName, err)
	}

	// verify against fixtures
	wantResp, err := os.ReadFile(filepath.Join(expectedDir, "response.json"))
	if err != nil {
		return fmt.Errorf("%s: read expected response.json: %w", caseName, err)
	}
	wantSched, err := os.ReadFile(filepath.Join(expectedDir, "schedule.csv"))
	if err != nil {
		return fmt.Errorf("%s: read expected schedule.csv: %w", caseName, err)
	}
	if !bytes.Equal(respJSON, wantResp) {
		return fmt.Errorf("%s: response.json mismatch", caseName)
	}
	if !bytes.Equal(schedCSV, wantSched) {
		return fmt.Errorf("%s: schedule.csv mismatch", caseName)
	}

	return nil
}

func cmdServe(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	addr := fs.String("addr", "127.0.0.1:8080", "Listen address")
	if err := fs.Parse(args); err != nil {
		return err
	}

	srv := api.NewServer(*addr)
	fmt.Fprintf(os.Stdout, "Listening on http://%s\n", *addr)
	return srv.ListenAndServe()
}
