package calc

// AmortizeRequestV1 is the input contract for the v1 amortization calculator.
//
// Money is expressed in integer cents (no floats).
// Rate is expressed in basis points (bps), where 100 bps = 1.00%.
// StartDate is ISO-8601 (YYYY-MM-DD) and is used only for schedule dates.
//
// This contract is intentionally small and strict.
// If a field is invalid, the calculator returns a stable, user-facing error.
type AmortizeRequestV1 struct {
	PrincipalCents int64  `json:"principal_cents"`
	AnnualRateBps  int64  `json:"annual_rate_bps"`
	TermMonths     int    `json:"term_months"`
	StartDate      string `json:"start_date"`
}

// AmortizeResponseV1 is the versioned JSON response for the v1 amortization calculator.
//
// Notes:
// - payment_cents is the scheduled payment (most periods)
// - last_payment_cents may differ by 1 cent due to final payoff rounding
// - totals are deterministic and derived from the computed schedule
//
// JSON is emitted from a struct (not a map) so key ordering is stable.
type AmortizeResponseV1 struct {
	SchemaVersion      string `json:"schema_version"`
	Calculator         string `json:"calculator"`
	PrincipalCents     int64  `json:"principal_cents"`
	AnnualRateBps      int64  `json:"annual_rate_bps"`
	TermMonths         int    `json:"term_months"`
	StartDate          string `json:"start_date"`
	PaymentCents       int64  `json:"payment_cents"`
	LastPaymentCents   int64  `json:"last_payment_cents"`
	TotalInterestCents int64  `json:"total_interest_cents"`
	TotalPaidCents     int64  `json:"total_paid_cents"`
}

// ScheduleRow is one amortization schedule row.
//
// Date is ISO-8601 (YYYY-MM-DD). Money is integer cents.
type ScheduleRow struct {
	Period         int
	Date           string
	PaymentCents   int64
	PrincipalCents int64
	InterestCents  int64
	BalanceCents   int64
}
