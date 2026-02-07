package calc

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strconv"
)

// RenderResponseJSON emits a stable, indented JSON representation
// with a trailing newline (for checked-in fixtures/goldens).
func RenderResponseJSON(resp AmortizeResponseV1) ([]byte, error) {
	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, err
	}
	b = append(b, '\n')
	return b, nil
}

// RenderScheduleCSV emits a stable CSV schedule (LF line endings).
func RenderScheduleCSV(rows []ScheduleRow) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	// csv.Writer uses \n internally; Go does not auto-convert line endings.
	if err := w.Write([]string{"period", "date", "payment_cents", "principal_cents", "interest_cents", "balance_cents"}); err != nil {
		return nil, err
	}
	for _, r := range rows {
		rec := []string{
			itoa(r.Period),
			r.Date,
			itoa64(r.PaymentCents),
			itoa64(r.PrincipalCents),
			itoa64(r.InterestCents),
			itoa64(r.BalanceCents),
		}
		if err := w.Write(rec); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func itoa(v int) string {
	return strconv.Itoa(v)
}

func itoa64(v int64) string {
	return strconv.FormatInt(v, 10)
}
