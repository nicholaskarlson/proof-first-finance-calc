package calc

import (
	"errors"
	"fmt"
	"math/big"
	"time"
)

const (
	calcNameV1   = "amortize"
	schemaV1     = "v1"
	bpsDenom     = int64(10000)
	monthsPerYr  = int64(12)
	monthlyDenom = bpsDenom * monthsPerYr
)

// AmortizeV1 computes a deterministic amortization schedule using:
// - integer cents for all money
// - basis points for annual nominal rate
// - monthly rate r = annual_rate / 12
// - interest rounded half-up to cents each period
// - payment rounded half-up to cents
// - last payment adjusted to bring balance to exactly zero
func AmortizeV1(req AmortizeRequestV1) (AmortizeResponseV1, []ScheduleRow, error) {
	if err := validateReq(req); err != nil {
		return AmortizeResponseV1{}, nil, err
	}
	start, _ := time.Parse("2006-01-02", req.StartDate)
	start = start.UTC()

	pmt := scheduledPaymentCents(req.PrincipalCents, req.AnnualRateBps, req.TermMonths)
	bal := req.PrincipalCents

	rows := make([]ScheduleRow, 0, req.TermMonths)
	var totalInt, totalPaid int64

	for i := 1; i <= req.TermMonths; i++ {
		interest := interestCents(bal, req.AnnualRateBps)
		principal := pmt - interest
		payThis := pmt

		if principal > bal {
			principal = bal
			payThis = interest + principal
		}
		bal -= principal
		totalInt += interest
		totalPaid += payThis

		dt := start.AddDate(0, i-1, 0)
		rows = append(rows, ScheduleRow{
			Period:         i,
			Date:           dt.Format("2006-01-02"),
			PaymentCents:   payThis,
			PrincipalCents: principal,
			InterestCents:  interest,
			BalanceCents:   bal,
		})
	}

	resp := AmortizeResponseV1{
		SchemaVersion:      schemaV1,
		Calculator:         calcNameV1,
		PrincipalCents:     req.PrincipalCents,
		AnnualRateBps:      req.AnnualRateBps,
		TermMonths:         req.TermMonths,
		StartDate:          req.StartDate,
		PaymentCents:       pmt,
		LastPaymentCents:   rows[len(rows)-1].PaymentCents,
		TotalInterestCents: totalInt,
		TotalPaidCents:     totalPaid,
	}
	return resp, rows, nil
}

func validateReq(req AmortizeRequestV1) error {
	if req.PrincipalCents <= 0 {
		return errors.New("principal_cents must be > 0")
	}
	if req.TermMonths <= 0 {
		return errors.New("term_months must be > 0")
	}
	if req.AnnualRateBps < 0 {
		return errors.New("annual_rate_bps must be >= 0")
	}
	if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
		return fmt.Errorf("start_date must be YYYY-MM-DD: %w", err)
	}
	return nil
}

func interestCents(balanceCents, annualRateBps int64) int64 {
	if annualRateBps == 0 || balanceCents == 0 {
		return 0
	}
	// interest = round_half_up(balance * annual_bps / (10000*12))
	return roundDivHalfUp(balanceCents*annualRateBps, monthlyDenom)
}

func scheduledPaymentCents(principalCents, annualRateBps int64, termMonths int) int64 {
	if annualRateBps == 0 {
		// round_half_up(P / n)
		return roundDivHalfUp(principalCents, int64(termMonths))
	}

	// r = annualRateBps / (10000*12)
	r := new(big.Rat).SetFrac(big.NewInt(annualRateBps), big.NewInt(monthlyDenom))
	one := big.NewRat(1, 1)
	onePlus := new(big.Rat).Add(one, r)
	pow := powRat(onePlus, termMonths)

	// payment = P * r * pow / (pow - 1)
	num := new(big.Rat).Mul(new(big.Rat).SetInt64(principalCents), r)
	num.Mul(num, pow)
	den := new(big.Rat).Sub(pow, one)
	pmt := new(big.Rat).Quo(num, den)
	return roundRatHalfUpToInt64(pmt)
}
