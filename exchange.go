package money

import "github.com/shopspring/decimal"

type RoundingMode func(amount decimal.Decimal) decimal.Decimal

type Exchange func(base, convertTo Money, r RoundingMode)

func RoundUp(amount decimal.Decimal) decimal.Decimal {
	return amount
}

//# Default to bankers rounding
//self.rounding_mode = BigDecimal::ROUND_HALF_EVEN
