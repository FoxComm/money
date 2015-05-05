package money

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/FoxComm/money/currency"
	"github.com/shopspring/decimal"
)

var parseRegex = regexp.MustCompile(`[+-]?[0-9]*[.]?[0-9]*`)

// Money represents an amount of a specific currency as an immutable value
type Money struct {
	amount   decimal.Decimal
	currency currency.Currency
}

// Make is the Federal Reserve
func Make(amount decimal.Decimal, c currency.Currency) Money {
	return Money{amount, c}
}

// Zero returns Money with a zero amount
func Zero(c currency.Currency) Money {
	return Make(decimal.New(0, 0), c)
}

// Parse parses a money.String() into Money
func Parse(str string) (money Money, err error) {
	var c currency.Currency
	var ok bool

	if len(str) < 4 {
		err = fmt.Errorf("'%s' cannot be parsed", str)
		return
	}

	parsed := parseRegex.FindStringSubmatch(str[4:])
	if len(parsed) == 0 {
		err = fmt.Errorf("'%s' cannot be parsed", str)
		return
	}

	if c, ok = currency.Table[str[0:3]]; !ok {
		err = fmt.Errorf("could not find currency %s", str[0:3])
		return
	}

	amountStr := strings.Replace(parsed[0], string(c.Delimiter), "", 0)

	if amount, err := decimal.NewFromString(amountStr); err != nil {
		return money, err
	} else {
		return Make(amount, c), nil
	}
}

// Amount is the monetary value in its major unit
func (m Money) Amount() decimal.Decimal {
	return m.amount
}

// Amount is the monetary value in its minor unit
func (m Money) AmountMinor() decimal.Decimal {
	return decimal.Decimal{}
}

// String represents the amount in a currency context. e.g., for US: "USD 10.00"
func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.currency.Code, m.Amount())
}

// Equals is true if other Money is the same amount and currency
func (m Money) Equals(other Money) bool {
	return m.IsCurrency(other.currency) && m.Amount().Equals(other.Amount())
}

// Currency returns the set Currency
func (m Money) Currency() currency.Currency {
	return m.currency
}

// IsCurrency tests if the currency is equivalent to another
func (m Money) IsCurrency(other currency.Currency) bool {
	return m.currency == other
}

// WithCurrency transforms this Money to a different Currency
func (m Money) WithCurrency(c currency.Currency) Money {
	return Make(m.amount, c)
}

func (m Money) panicIfDifferentCurrency(c currency.Currency) {
	if !m.IsCurrency(c) {
		panic(fmt.Errorf("expected currency %s, got %s", m.currency.Code, c.Code))
	}
}

// Math, aka, here be dragons

// Abs returns |amount|
func (m Money) Abs() decimal.Decimal {
	return m.amount.Abs()
}

// Negate negates the sign of the amount
func (m Money) Negate() Money {
	d := decimal.New(-1, 0)
	return Make(m.amount.Mul(d), m.currency)
}

// IsPositive returns true if the amount i > 0
func (m Money) IsPositive() bool {
	zero := decimal.New(0, 0)
	return m.amount.Cmp(zero) == 1
}

// IsPositive returns true if the amount i > 0
func (m Money) IsNegative() bool {
	zero := decimal.New(0, 0)
	return m.amount.Cmp(zero) == -1
}

// IsZero returns true if the amount i == 0
func (m Money) IsZero() bool {
	return m.amount.Equals(decimal.New(0, 0))
}

// Add adds monies. Panics if currency is different.
func (m Money) Add(other Money) Money {
	m.panicIfDifferentCurrency(other.currency)
	return Make(m.amount.Add(other.amount), m.currency)
}

// Sub subtracts monies. Panics if currency is different.
func (m Money) Sub(other Money) Money {
	m.panicIfDifferentCurrency(other.currency)
	return Make(m.amount.Sub(other.amount), m.currency)
}

// Div divides monies. Panics if currency is different.
func (m Money) Div(other Money) Money {
	m.panicIfDifferentCurrency(other.currency)
	return Make(m.amount.Div(other.amount), m.currency)
}
