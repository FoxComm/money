package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const (
	header     = "// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.\npackage currency\n"
	jsonPath   = "./internal/currencies.json"
	outputPath = "./currency/currencies.go"
)

type Currency struct {
	Code      string `json:"iso_code"`
	Number    string `json:"iso_numeric"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Decimal   string `json:"decimal_mark"`
	Delimiter string `json:"thousands_separator"`
	Minor     int    `json:"subunit_to_unit"`
}

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
}

var currencyTmpl = template.Must(template.New("currency-file").Funcs(funcMap).Parse(`
// {{ .Code }} is the {{ .Name }} Currency
var {{ .Code }} = Currency{
	Code: "{{ .Code | ToUpper }}",
	Number: {{ .Number }},
	Symbol: '{{ .Symbol }}',
	Decimal: '{{ .Decimal }}',
	Delimiter: '{{ .Delimiter }}',
	Minor: {{ .Minor }},
}

`))

var tableTmpl = template.Must(template.New("currencies-table").Parse(`
// Table holds all compiled currencies in a map ISO-NAME => value
var Table = map[string]Currency{
	{{ range $index, $name := . }}"{{$name}}": {{$name}},
	{{end}}
}

`))

func main() {
	file, err := os.Open("./internal/currencies.json")
	if err != nil {
		panic(err)
	}

	var currencies map[string]Currency
	if err = json.NewDecoder(file).Decode(&currencies); err != nil {
		panic(err)
	}

	if len(currencies) == 0 {
		panic("Expected currencies to be > 0")
	}

	buf := new(bytes.Buffer)
	buf.WriteString(header)

	if err = writeTable(buf, currencies); err != nil {
		panic(err)
	}

	if err = writeCurrencies(buf, currencies); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}

func writeCurrencies(buf *bytes.Buffer, currencies map[string]Currency) error {
	for _, currency := range currencies {
		if err := currencyTmpl.Execute(buf, currency); err != nil {
			return err
		}
	}

	return nil
}

func writeTable(buf *bytes.Buffer, currencies map[string]Currency) error {
	keys := make([]string, 0, len(currencies))
	for k, _ := range currencies {
		keys = append(keys, strings.ToUpper(k))
	}

	if err := tableTmpl.Execute(buf, keys); err != nil {
		return err
	}

	return nil
}
