//TODO: This needs refactoring
package stocks

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

type Stock struct {
	Symbol   string
	FullName string
	Exchange string   `json:",omitempty"`
	ETF      bool     `json:",omitempty"`
	Aliases  []string `json:",omitempty"`
}

func simpleBoolConvert(s string) bool {
	return s == "Y"
}

var (
	exchangeNameMap = map[string]string{
		"A": "NYSE MKT",
		"N": "NYSE",
		"P": "NYSE ARCA",
		"Z": "BATS",
		"V": "IEXG",
	}
)

// lastLineReached returns weather or not the error message represents the last line in the NASDAQ file
func lastLineReached(msg string) bool {

	parts := strings.Split(msg, ": ")
	if parts[1] == "wrong number of fields" {
		locationParts := strings.Split(parts[0], " ")
		lineNumber, err := strconv.Atoi(locationParts[len(locationParts)-1])
		if err != nil {
			return false
		}
		return lineNumber > 1000
	}

	return false
}

func newReader(f io.Reader) *csv.Reader {

	r := csv.NewReader(f)
	r.Comma = '|'

	return r
}

func parseOtherListedRow(record []string) Stock {
	return Stock{
		Symbol:   record[0],
		FullName: record[1],
		Exchange: exchangeNameMap[record[2]],
		ETF:      simpleBoolConvert(record[4]),
	}
}

func parseNasdaqListedRow(record []string) Stock {
	return Stock{
		Symbol:   record[0],
		FullName: record[1],
		Exchange: "NASDAQ",
		ETF:      simpleBoolConvert(record[6]),
	}
}

func parseFileWithRowFunc(f io.Reader, rowFunc func([]string) Stock) ([]Stock, error) {

	csvReader := newReader(f)

	stocks := []Stock{}
	for {
		v, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil && !lastLineReached(err.Error()) {
			return nil, err
		}

		stocks = append(stocks, rowFunc(v))
	}

	return stocks[1:], nil
}

func FromNasdaqOtherListed(f io.Reader) ([]Stock, error) {

	return parseFileWithRowFunc(f, parseOtherListedRow)
}

func FromNasdaqListed(f io.Reader) ([]Stock, error) {

	return parseFileWithRowFunc(f, parseNasdaqListedRow)
}
