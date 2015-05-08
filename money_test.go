package money

import (
	"errors"
	"testing"

	. "github.com/FoxComm/money/currency"
	"github.com/shopspring/decimal"
)

var (
	// XTS is an ISO code reserved for use in testing
	XTS = "XTS"
)

func d(value string) decimal.Decimal {
	if v, err := decimal.NewFromString(value); err != nil {
		panic(err)
	} else {
		return v
	}
}

func checkForZeroAndErr(t *testing.T, funcName string, zero Money, err error) {
	if errDiffCurrency, ok := err.(*ErrDifferentCurrency); !ok {
		t.Errorf("%s => (%+v, %s) expected ErrDifferentCurrency, got %s", funcName, zero, err)
	} else if ok && (errDiffCurrency.Actual == errDiffCurrency.Expected) {
		t.Errorf("%s => (%+v, %s) got ErrDifferentCurrency but currencies are same %s", errDiffCurrency)
	}

	if !zero.IsZero() {
		t.Errorf("%s => (%+v, err) expected zero money, got %s", funcName, zero)
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
		expected decimal.Decimal
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
	t.Skip("TestString is broken for now")

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

func TestNegate(t *testing.T) {
	var monies = []struct {
		money    Money
		expected Money
	}{
		{Make(d("-1"), USD), Make(d("1"), USD)},
		{Make(d("-100000"), USD), Make(d("100000"), USD)},
		{Make(d("5555"), USD), Make(d("-5555"), USD)},
		{Make(d("-10.005"), USD), Make(d("10.005"), USD)},
	}

	for _, m := range monies {
		if !m.money.Negate().Equals(m.expected) {
			t.Errorf("Money.Negate() => %s, expected %s", m.money.Amount(), m.expected.Amount())
		}
	}
}

func TestIsPositive(t *testing.T) {
	var monies = []struct {
		money    Money
		expected bool
	}{
		{Make(d("0"), USD), false},
		{Make(d("-1"), USD), false},
		{Make(d("-100000"), USD), false},
		{Make(d("5555"), USD), true},
		{Make(d("1"), USD), true},
		{Make(d("1000000000"), USD), true},
	}

	for _, m := range monies {
		isPos := m.money.IsPositive()
		if isPos != m.expected {
			t.Errorf("Money.IsPositive() => %s, expected %s", isPos, m.expected)
		}
	}
}

func TestIsNegative(t *testing.T) {
	var monies = []struct {
		money    Money
		expected bool
	}{
		{Make(d("0"), USD), false},
		{Make(d("-1"), USD), true},
		{Make(d("-100000"), USD), true},
		{Make(d("5555"), USD), false},
		{Make(d("1"), USD), false},
		{Make(d("1000000000"), USD), false},
	}

	for _, m := range monies {
		isNeg := m.money.IsNegative()
		if isNeg != m.expected {
			t.Errorf("Money.IsNegative() => %s, expected %s", isNeg, m.expected)
		}
	}
}

func TestIsZero(t *testing.T) {
	var monies = []struct {
		money    Money
		expected bool
	}{
		{Make(d("0"), USD), true},
		{Make(d("0.0"), USD), true},
		{Make(d("0.000000"), USD), true},
		{Make(d("-1"), USD), false},
		{Make(d("-100000"), USD), false},
		{Make(d("5555"), USD), false},
		{Make(d("1"), USD), false},
		{Make(d("1000000000"), USD), false},
	}

	for _, m := range monies {
		isZero := m.money.IsZero()
		if isZero != m.expected {
			t.Errorf("Money.IsZero() => %s, expected %s", isZero, m.expected)
		}
	}
}

func TestAdd(t *testing.T) {
	var monies = []struct {
		money    Money
		add      decimal.Decimal
		expected Money
	}{
		{Make(d("0"), USD), d("0"), Make(d("0"), USD)},
		{Make(d("10"), USD), d("0"), Make(d("10"), USD)},
		{Make(d(".5"), USD), d(".2"), Make(d(".7"), USD)},
		{Make(d("0.005"), USD), d("10"), Make(d("10.005"), USD)},
		{Make(d("0.5"), USD), d("-.02"), Make(d("0.48"), USD)},
	}

	for _, m := range monies {
		if actual, err := m.money.Add(Make(m.add, USD)); err != nil {
			t.Errorf("Money.Add() => unexpected error %s", err)
		} else if !actual.Equals(m.expected) {
			t.Errorf("Money.Add() => (%s, nil) expected %s", actual.Amount(), m.expected.Amount())
		}
	}

	zero, err := monies[0].money.Add(Make(d("0"), MXN))
	checkForZeroAndErr(t, "Money.Add()", zero, err)
}

func TestSub(t *testing.T) {
	var monies = []struct {
		money    Money
		sub      decimal.Decimal
		expected Money
	}{
		{Make(d("0"), USD), d("0"), Make(d("0"), USD)},
		{Make(d("10"), USD), d("0"), Make(d("10"), USD)},
		{Make(d(".5"), USD), d(".2"), Make(d(".3"), USD)},
		{Make(d("0.005"), USD), d("10"), Make(d("-9.995"), USD)},
		{Make(d("0.5"), USD), d("-.02"), Make(d("0.52"), USD)},
	}

	for _, m := range monies {
		if actual, err := m.money.Sub(Make(m.sub, USD)); err != nil {
			t.Errorf("Money.Sub() => unexpected error %s", err)
		} else if !actual.Equals(m.expected) {
			t.Errorf("Money.Sub() => (%s, nil) expected %s", actual.Amount(), m.expected.Amount())
		}
	}

	zero, err := monies[0].money.Sub(Make(d("0"), MXN))
	checkForZeroAndErr(t, "Money.Sub()", zero, err)
}

func TestDiv(t *testing.T) {
	var monies = []struct {
		money    Money
		div      decimal.Decimal
		expected Money
	}{
		{Make(d("0"), USD), d("1"), Make(d("0"), USD)},
		{Make(d("10"), USD), d("2"), Make(d("5"), USD)},
		{Make(d(".5"), USD), d(".2"), Make(d("2.5"), USD)},
		{Make(d("0.005"), USD), d("10"), Make(d("0.0005"), USD)},
		{Make(d("0.5"), USD), d("-.02"), Make(d("-25"), USD)},
	}

	for _, m := range monies {
		if actual, err := m.money.Div(Make(m.div, USD)); err != nil {
			t.Errorf("Money.Div() => unexpected error %s", err)
		} else if !actual.Equals(m.expected) {
			t.Errorf("Money.Div() => (%s, nil) expected %s", actual.Amount(), m.expected.Amount())
		}
	}

	zero, err := monies[0].money.Div(Make(d("1"), MXN))
	checkForZeroAndErr(t, "Money.Div()", zero, err)
}

func TestMul(t *testing.T) {
	var monies = []struct {
		money    Money
		mul      decimal.Decimal
		expected Money
	}{
		{Make(d("0"), USD), d("1"), Make(d("0"), USD)},
		{Make(d("10"), USD), d("2"), Make(d("20"), USD)},
		{Make(d(".5"), USD), d(".2"), Make(d(".1"), USD)},
		{Make(d("0.005"), USD), d("10"), Make(d("0.05"), USD)},
		{Make(d("0.5"), USD), d("-.02"), Make(d("-0.01"), USD)},
	}

	for _, m := range monies {
		if actual, err := m.money.Mul(Make(m.mul, USD)); err != nil {
			t.Errorf("Money.Mul() => unexpected error %s", err)
		} else if !actual.Equals(m.expected) {
			t.Errorf("Money.Mul() => (%s, nil) expected %s", actual.Amount(), m.expected.Amount())
		}
	}

	zero, err := monies[0].money.Mul(Make(d("1"), MXN))
	checkForZeroAndErr(t, "Money.Mul()", zero, err)
}

func TestCmp(t *testing.T) {
	var monies = []struct {
		money    Money
		other    Money
		expected int
	}{
		{Make(d("0"), USD), Make(d("0"), USD), 0},
		{Make(d("1"), USD), Make(d("1"), USD), 0},
		{Make(d("5"), USD), Make(d("5.99999"), USD), -1},
		{Make(d("0"), USD), Make(d("1"), USD), -1},
		{Make(d("2"), USD), Make(d("1"), USD), 1},
		{Make(d(".1"), USD), Make(d("0"), USD), 1},
		{Make(d(".1"), USD), Make(d("-.1"), USD), 1},
	}

	for _, m := range monies {
		if actual, err := m.money.Cmp(m.other); err != nil {
			t.Errorf("Money.Cmp() => unexpected error %s", err)
		} else if actual != m.expected {
			t.Errorf("Money.Cmp() => (%d, nil) expected %d %+v", actual, m.expected, m.money.Amount())
		}
	}
}
