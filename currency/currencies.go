// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.
package currency

// Table holds all compiled currencies in a map ISO-NAME => value
var Table = map[string]Currency{
	"USD": USD,
	"MXN": MXN,
	"CNY": CNY,
	"CAD": CAD,
}

// CNY is the Chinese Renminbi Yuan Currency
var CNY = Currency{
	Code:      "CNY",
	Number:    156,
	Symbol:    '¥',
	Decimal:   '.',
	Delimiter: ',',
	Minor:     100,
}

// CAD is the Canadian Dollar Currency
var CAD = Currency{
	Code:      "CAD",
	Number:    124,
	Symbol:    '$',
	Decimal:   '.',
	Delimiter: ',',
	Minor:     100,
}

// USD is the United States Dollar Currency
var USD = Currency{
	Code:      "USD",
	Number:    840,
	Symbol:    '$',
	Decimal:   '.',
	Delimiter: ',',
	Minor:     100,
}

// MXN is the Mexican Peso Currency
var MXN = Currency{
	Code:      "MXN",
	Number:    484,
	Symbol:    '$',
	Decimal:   '.',
	Delimiter: ',',
	Minor:     100,
}
