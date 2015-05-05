package money

import (
	"errors"
	"testing"

	. "github.com/FoxComm/money/currency"
)

var (
	// XTS is an ISO code reserved for use in testing
	XTS = "XTS"
)

func TestMoneyAbs(t *testing.T) {
	var monies = []struct {
		money    Money
		expected uint
	}{
		{Make(-10, USD), 10},
		{Make(0, USD), 0},
		{Make(10, USD), 10},
	}

	for _, m := range monies {
		if m.money.Abs() != m.expected {
			t.Errorf("Money.Abs() => %d, expected %d", m.money.Abs(), m.expected)
		}
	}
}

func TestMoneyIsEqual(t *testing.T) {
	var monies = []struct {
		money   Money
		other   Money
		isEqual bool
		name    string
	}{
		{Make(-10, USD), Make(-10, USD), true, "same amount, same currency"},
		{Make(-10, USD), Make(-5, USD), false, "different amount, same currency"},
		{Make(-10, USD), Make(-10, MXN), false, "same amount, different currency"},
	}

	for _, m := range monies {
		if isEqual := m.money.IsEqual(m.other); isEqual != m.isEqual {
			t.Errorf("Money.IsEqual() => %v, expected %v. Table: %+v", isEqual, m.isEqual, m)
		}
	}
}

func TestMoneyString(t *testing.T) {
	var monies = []struct {
		money    Money
		expected string
	}{
		{Make(-10, USD), "USD -0.10"},
		{Make(10, USD), "USD 0.10"},
		{Make(1055, USD), "USD 10.55"},
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
		{Make(55, USD).String(), nil},
		{Make(5555, USD).String(), nil},
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

func TestShouldNotUseFloat64(t *testing.T) {
	t.Error("Money.String() must remove float64() coercion")
}
