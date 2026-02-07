package calc

import "math/big"

func roundRatHalfUpToInt64(r *big.Rat) int64 {
	// r must be non-negative.
	num := new(big.Int).Set(r.Num())
	den := new(big.Int).Set(r.Denom())
	if den.Sign() == 0 {
		return 0
	}
	// (num + den/2) / den
	half := new(big.Int).Rsh(den, 1)
	num.Add(num, half)
	q := new(big.Int).Quo(num, den)
	return q.Int64()
}

func powRat(x *big.Rat, n int) *big.Rat {
	// exponentiation by squaring
	res := big.NewRat(1, 1)
	base := new(big.Rat).Set(x)
	exp := n
	for exp > 0 {
		if exp&1 == 1 {
			res.Mul(res, base)
		}
		exp >>= 1
		if exp > 0 {
			base.Mul(base, base)
		}
	}
	return res
}

func roundDivHalfUp(numer, denom int64) int64 {
	// For this repo's use, denom must be > 0 and numer must be >= 0.
	if denom <= 0 {
		return 0
	}
	return (numer + denom/2) / denom
}
