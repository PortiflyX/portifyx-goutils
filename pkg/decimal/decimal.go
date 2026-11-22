package decimal

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// Zero returns a Decimal representing 0.
func Zero() decimal.Decimal {
	return decimal.Zero
}

// FromString parses a string to decimal.Decimal, rejecting special keywords like NaN or Inf.
func FromString(s string) (decimal.Decimal, error) {
	clean := strings.TrimSpace(s)
	lower := strings.ToLower(clean)
	if lower == "nan" || strings.Contains(lower, "inf") {
		return decimal.Zero, fmt.Errorf("invalid decimal string: %s", s)
	}
	d, err := decimal.NewFromString(clean)
	if err != nil {
		return decimal.Zero, fmt.Errorf("parse decimal: %w", err)
	}
	return d, nil
}

// MustFromString parses decimal and panics on invalid input (useful for constants and test fixtures).
func MustFromString(s string) decimal.Decimal {
	d, err := FromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

// RoundFiat rounds to 2 decimal places (e.g. for THB/USD balances).
func RoundFiat(d decimal.Decimal) decimal.Decimal {
	return d.Round(2)
}

// RoundCrypto rounds to 8 decimal places (standard satoshi scale).
func RoundCrypto(d decimal.Decimal) decimal.Decimal {
	return d.Round(8)
}
