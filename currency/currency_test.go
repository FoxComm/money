package currency

import (
	"strings"
	"testing"
)

var (
	// XTS is an ISO code reserved for use in testing
	XTS = "XTS"
)

func TestString(t *testing.T) {
	usd := USD
	usd.Code = strings.ToLower(USD.Code)

	if usd.String() != USD.String() {
		t.Errorf("String() => %s, expected %s", usd.String(), USD.String())
	}
}

func TestEquals(t *testing.T) {
	changedUSD := USD
	changedUSD.Code = strings.ToLower(USD.Code)

	var currencies = []struct {
		base   Currency
		other  Currency
		equals bool
	}{
		{USD, USD, true},
		{MXN, MXN, true},
		{changedUSD, USD, false},
		{USD, MXN, false},
	}

	for _, c := range currencies {
		if c.base.Equals(c.other) != c.equals {
			t.Errorf("Expected %v Equals() %v to be %v, got %v", c.base, c.other, c.equals, !c.equals)
		}
	}
}
