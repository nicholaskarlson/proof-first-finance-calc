package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nicholaskarlson/proof-first-finance-calc/internal/calc"
)

// Handler returns an http.Handler serving the v1 API.
func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/v1/amortize", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		req, err := decodeRequest(r)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		resp, _, err := calc.AmortizeV1(req)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		b, err := calc.RenderResponseJSON(resp)
		if err != nil {
			internalError(w)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	})

	mux.HandleFunc("/v1/amortize/schedule.csv", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		req, err := decodeRequest(r)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		_, sched, err := calc.AmortizeV1(req)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		b, err := calc.RenderScheduleCSV(sched)
		if err != nil {
			internalError(w)
			return
		}
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	})

	return mux
}

func decodeRequest(r *http.Request) (calc.AmortizeRequestV1, error) {
	// Tight, stable failures. DisallowUnknownFields gives better signals for users,
	// but we don't want to lock tests to Go's JSON error text. So we keep messages
	// minimal + stable.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req calc.AmortizeRequestV1
	if err := dec.Decode(&req); err != nil {
		return calc.AmortizeRequestV1{}, fmt.Errorf("invalid JSON")
	}
	// Reject trailing tokens.
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return calc.AmortizeRequestV1{}, fmt.Errorf("invalid JSON")
	}
	return req, nil
}

func badRequest(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = fmt.Fprintf(w, "error: %s\n", msg)
}

func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Allow", http.MethodPost)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte("error: method not allowed\n"))
}

func internalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("error: internal server error\n"))
}

// NewServer constructs an http.Server with sensible timeouts.
func NewServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}
