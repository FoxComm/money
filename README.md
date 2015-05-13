# Money

[![Build status](https://badge.buildkite.com/4adcfafd46e900e1c20a92112ff00b84cee1bad3b4b55a3672.svg)](https://buildkite.com/foxcommerce/money)

A Go library handling money, currencies, and exchange conversion.

## Goals

1. Correctness
2. Performant
3. Extensible

## Usage

### Installation

```bash
go get -u github.com/FoxComm/money
```

### Examples

```go
m := money.Make(5000, currencies.USD)
m.String() => "USD $50.00"
```

### Internal

The `internal/` dir has some internal tooling with a corresponding
[README](internal/README.md).

## Inspiration

Ideas and inspiration drawn from:

* [Joda Money](http://www.joda.org/joda-money/)
* [RubyMoney](https://github.com/RubyMoney/money). A special thank you
  for the `currencies-iso.json` since it served as the basis for our
  json currencies and the idea of templating.

## Resources

[ISO 4217](http://en.wikipedia.org/wiki/ISO_4217)
[Currency conversions](http://www.yacoset.com/how-to-handle-currency-conversions)
[Exchange rate](http://en.wikipedia.org/wiki/Exchange_rate)
[Circulating currencies](http://en.wikipedia.org/wiki/List_of_circulating_currencies)
[Decimal pkg post](http://engineering.shopspring.com/2015/03/03/decimal/)

## License

Released under the [MIT License](License).
