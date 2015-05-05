package money

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/FoxComm/money/currency"
)

var parseRegex = regexp.MustCompile(`[+-]?[0-9]*[.]?[0-9]*`)

// Money represents an amount of a specific currency as an immutable value
type Money struct {
	amount   int
	currency currency.Currency
}

// Make is the Federal Reserve
func Make(amount int, c currency.Currency) Money {
	return Money{amount, c}
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

	amountStr := strings.Replace(
		strings.Replace(parsed[0], string(c.Decimal), "", 1),
		string(c.Delimiter),
		"",
		0,
	)

	if amount, err := strconv.Atoi(amountStr); err != nil {
		return money, err
	} else {
		return Make(amount, c), nil
	}
}

// Amount is the monetary value in its major unit
func (m Money) Amount() int {
	return m.amount
}

// Amount is the monetary value in its minor unit
func (m Money) AmountMinor() int {
	return m.amount
}

// String represents the amount in a currency context. e.g., for US: "USD $10.00"
func (m Money) String() string {
	c := m.currency
	return fmt.Sprintf("%s %.2f", c.Code, float64(m.Amount()))
}

// Abs returns |amount|
func (m Money) Abs() uint {
	if m.amount < 0 {
		return uint(-m.Amount())
	}
	return uint(m.Amount())
}

// IsEqual is true if other Money is the same amount and currency
func (m Money) IsEqual(other Money) bool {
	return m.currency == other.currency && m.Amount() == other.Amount()
}

// Currency returns the set Currency
func (m Money) Currency() currency.Currency {
	return m.currency
}

// WithCurrency transforms this Money to a different Currency
func (m Money) WithCurrency(c currency.Currency) Money {
	return Make(m.amount, c)
}
