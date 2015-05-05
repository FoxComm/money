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

// Abs returns |amount|
func (m Money) Abs() decimal.Decimal {
	return m.amount.Abs()
}

// Equals is true if other Money is the same amount and currency
func (m Money) Equals(other Money) bool {
	return m.currency == other.currency && m.Amount().Equals(other.Amount())
}

// Currency returns the set Currency
func (m Money) Currency() currency.Currency {
	return m.currency
}

// WithCurrency transforms this Money to a different Currency
func (m Money) WithCurrency(c currency.Currency) Money {
	return Make(m.amount, c)
}
