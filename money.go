package money

import (
	"database/sql/driver"
	"errors"
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

// ErrDifferentCurrency is used for functions which take another money/currency
// whereby the Money's currency != other currency
type ErrDifferentCurrency struct {
	Actual   currency.Currency
	Expected currency.Currency
}

func (e *ErrDifferentCurrency) Error() string {
	return fmt.Sprintf("expected currency %s got %s", e.Expected.Code, e.Actual.Code)
}

// Make is the Federal Reserve
func Make(amount decimal.Decimal, c currency.Currency) Money {
	return Money{amount, c}
}

// MakeFromString is the Federal Reserve
func MakeFromString(amount string, c currency.Currency) (Money, error) {
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return Zero(c), err
	}

	return Money{d, c}, nil
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
	amt := m.Amount().String()
	if !strings.Contains(amt, string(m.currency.Decimal)) {
		amt += fmt.Sprintf("%s00", string(m.currency.Decimal))
	}
	return fmt.Sprintf("%s %s", m.currency.Code, amt)
}

// Equals is true if other Money is the same amount and currency
func (m Money) Equals(other Money) bool {
	return m.currency.Equals(other.currency) && m.Amount().Equals(other.Amount())
}

// Currency returns the set Currency
func (m Money) Currency() currency.Currency {
	return m.currency
}

// WithCurrency transforms this Money to a different Currency
func (m Money) WithCurrency(c currency.Currency) Money {
	return Make(m.amount, c)
}

func (m Money) panicIfDifferentCurrency(c currency.Currency) {
	if !m.currency.Equals(c) {
		panic(fmt.Errorf("expected currency %s, got %s", m.currency.Code, c.Code))
	}
}

// Math, aka, here be dragons

// Cmp comparies monies. errors if currency is different.
func (m Money) Cmp(other Money) (int, error) {
	if !m.currency.Equals(other.currency) {
		return 0, &ErrDifferentCurrency{m.currency, other.currency}
	}
	return m.amount.Cmp(other.amount), nil
}

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

// Add adds monies. errors if currency is different.
func (m Money) Add(other Money) (Money, error) {
	if !m.currency.Equals(other.currency) {
		return Zero(m.currency), &ErrDifferentCurrency{m.currency, other.currency}
	}
	return Make(m.amount.Add(other.amount), m.currency), nil
}

// Sub subtracts monies. errors if currency is different.
func (m Money) Sub(other Money) (Money, error) {
	if !m.currency.Equals(other.currency) {
		return Zero(m.currency), &ErrDifferentCurrency{m.currency, other.currency}
	}
	return Make(m.amount.Sub(other.amount), m.currency), nil
}

// Div divides monies. errors if currency is different.
func (m Money) Div(other Money) (Money, error) {
	if !m.currency.Equals(other.currency) {
		return Zero(m.currency), &ErrDifferentCurrency{m.currency, other.currency}
	}
	return Make(m.amount.Div(other.amount), m.currency), nil
}

// Mul multiplies monies. errors if currency is different.
func (m Money) Mul(other Money) (Money, error) {
	if !m.currency.Equals(other.currency) {
		return Zero(m.currency), &ErrDifferentCurrency{m.currency, other.currency}
	}
	return Make(m.amount.Mul(other.amount), m.currency), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Money) UnmarshalJSON(data []byte) (err error) {
	*m, err = Parse(strings.Trim(string(data), `"`))
	return
}

// MarshalJSON implements the json.Marshaler interface.
func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.String() + `"`), nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (m *Money) Scan(value interface{}) (err error) {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("scan value was not []byte")
	}

	*m, err = Parse(string(asBytes))
	return err
}

// Value implements the driver.Valuer interface for database serialization.
func (m Money) Value() (driver.Value, error) {
	return m.String(), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization.
func (m *Money) UnmarshalText(text []byte) (err error) {
	*m, err = Parse(string(text))
	return
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (m Money) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}
