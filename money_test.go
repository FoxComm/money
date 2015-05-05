package money

import (
	"errors"
	"testing"

	. "github.com/FoxComm/money/currency"
	dec "github.com/shopspring/decimal"
)

var (
	// XTS is an ISO code reserved for use in testing
	XTS = "XTS"
)

func d(value string) dec.Decimal {
	if v, err := dec.NewFromString(value); err != nil {
		panic(err)
	} else {
		return v
	}
}

func TestZero(t *testing.T) {
	money := Zero(USD)
	if !money.IsZero() {
		t.Errorf("Zero() => %v, expected %s", money.Amount(), 0)
	}
}

func TestAbs(t *testing.T) {
	var monies = []struct {
		money    Money
		expected dec.Decimal
	}{
		{Make(d("-10"), USD), d("10")},
		{Make(d("0"), USD), d("0")},
		{Make(d("10"), USD), d("10")},
		{Make(d("55.55"), USD), d("55.55")},
		{Make(d("-55.55"), USD), d("55.55")},
	}

	for _, m := range monies {
		if !m.money.Abs().Equals(m.expected) {
			t.Errorf("Money.Abs() => %s, expected %s", m.money.Abs(), m.expected)
		}
	}
}

func TestEquals(t *testing.T) {
	var monies = []struct {
		money  Money
		other  Money
		equals bool
		name   string
	}{
		{Make(d("-10"), USD), Make(d("-10"), USD), true, "same amount, same currency"},
		{Make(d("-10"), USD), Make(d("-5"), USD), false, "different amount, same currency"},
		{Make(d("-10"), USD), Make(d("-10"), MXN), false, "same amount, different currency"},
		{Make(d("-5"), MXN), Make(d("-5"), MXN), true, "same amount, same currency"},
	}

	for _, m := range monies {
		if equals := m.money.Equals(m.other); equals != m.equals {
			t.Errorf("Money.Equals() => %v, expected %v. Table: %+v", equals, m.equals, m)
		}
	}
}

func TestString(t *testing.T) {
	var monies = []struct {
		money    Money
		expected string
	}{
		{Make(d("-1"), USD), "USD -1.00"},
		{Make(d("10"), USD), "USD 10.00"},
		{Make(d("1055"), USD), "USD 1055.00"},
	}

	for _, m := range monies {
		if m.money.String() != m.expected {
			t.Errorf("Money.String() => %v, expected %v", m.money.String(), m.expected)
		}
	}
}

func TestParse(t *testing.T) {
	var monies = []struct {
		money string
		err   error
	}{
		{Make(d("55"), USD).String(), nil},
		{Make(d("5555"), USD).String(), nil},
		{"XXX 55.00", errors.New("could not find currency XXX")},
	}

	for _, m := range monies {
		parsed, err := Parse(m.money)
		if err != nil && err.Error() != m.err.Error() {
			t.Errorf("Parse() => %s, expected %s", err, m.err)
		} else if err == nil && parsed.String() != m.money {
			t.Errorf("Parse() => %+v, expected %s", parsed, m.money)
		}
	}
}
