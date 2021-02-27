package stocks

import (
	"bytes"
	"testing"
)

func TestFromNasdaqListed(t *testing.T) {

	records := []string{
		"SYMBOL",
		"FULLNAME",
		"",
		"",
		"",
		"",
		"Y",
	}

	s := parseNasdaqListedRow(records)

	if s.Symbol != "SYMBOL" {
		t.Errorf("Symbol not parsed correctly: %s != %s", s.Symbol, "SYMBOL")
	}

	if s.Exchange != "NASDAQ" {
		t.Errorf("Symbol not parsed correctly: %s != %s", s.Symbol, "NASDAQ")
	}

}

func TestFromOtherListed(t *testing.T) {

	records := []string{
		"SYMBOL",
		"FULLNAME",
		"V",
		"",
		"N",
	}

	s := parseOtherListedRow(records)

	if s.Symbol != "SYMBOL" {
		t.Errorf("Symbol not parsed correctly: %s != %s", s.Symbol, "SYMBOL")
	}

	if s.FullName != "FULLNAME" {
		t.Errorf("Full Name not parsed correctly: %s != %s", s.FullName, "FULLNAME")
	}

	if s.Exchange != "IEXG" {
		t.Errorf("Exchange not parsed correctly: %s != %s", s.Exchange, "IEXG")
	}

}

func TestFileReader(t *testing.T) {

	b := []byte{}
	f := bytes.NewReader(b)
	r := newReader(f)

	if r.Comma != '|' {
		t.Errorf("Incorrect delim found: %c != %s", r.Comma, "|")
	}

}

func TestLastLineError(t *testing.T) {

	msg := "incorrect test error message: nothing"

	if lastLineReached(msg) {
		t.Errorf("Detectect last line condition when message was incorrect: msg = %s", msg)
	}

	msg = "line 1: wrong number of fields"
	if lastLineReached(msg) {
		t.Errorf("Detectect last line condition when message was incorrect: msg = %s", msg)
	}

	msg = "line 99999: wrong number of fields"
	if !lastLineReached(msg) {
		t.Errorf("Failed to detectect last line condition when message was correct: msg = %s", msg)
	}

}

func TestHigherLevelNasdaqListed(t *testing.T) {

	rows := `Symbol|Security Name|Market Category|Test Issue|Financial Status|Round Lot Size|ETF|NextShares
AACG|ATA Creativity Global - American Depositary Shares, each representing two common shares|G|N|N|100|N|N
`

	r := bytes.NewReader([]byte(rows))

	s, err := FromNasdaqListed(r)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(s) == 0 {
		t.Fatal("empty stock slice returned")
	}

	if s[0].Symbol != "AACG" {
		t.Errorf("Symbol not parsed correctly: %s != %s", s[0].Symbol, "AACG")
	}
}

func TestHigherLevelOtherListed(t *testing.T) {

	rows := `ACT Symbol|Security Name|Exchange|CQS Symbol|ETF|Round Lot Size|Test Issue|NASDAQ Symbol
A|Agilent Technologies, Inc. Common Stock|N|A|N|100|N|A
`

	r := bytes.NewReader([]byte(rows))

	s, err := FromNasdaqOtherListed(r)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(s) == 0 {
		t.Fatal("empty stock slice returned")
	}

	if s[0].Symbol != "A" {
		t.Errorf("Symbol not parsed correctly: %s != %s", s[0].Symbol, "A")
	}
}
