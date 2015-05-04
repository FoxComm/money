package currency

// Currency represents fiat money
type Currency struct {
	// Code is the ISO 4217 alpha-3 name for the currency
	Code string

	// Number is the ISO 3166-1 numeric code
	Number int

	// Symbol is the shorthand used for a currency's name
	Symbol rune

	// Decimal is a rune which separates the decimals
	Decimal rune

	// Delimiter is a rune which delimits integer thousands
	Delimiter rune

	// Minor is the 'exponent' of a currency unit. Assume base 10.
	Minor int
}

// Table holds all compiled currencies in a map ISO-NAME => value
var Table = map[string]Currency{
	"USD": USD,
	"MXN": MXN,
}
