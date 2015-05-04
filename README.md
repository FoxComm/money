# Money

A Go library handling money, currencies, and exchange conversion.

## Goals

1. Correctness
2. Performant
3. Extensible

## Usage

#### Installation

```bash
go get -u github.com/FoxComm/money
```

### Examples

```go
m := money.Make(5000, currencies.USD)
m.String() => "USD $50.00"
```
