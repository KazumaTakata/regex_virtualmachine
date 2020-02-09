package regex

import (
	"fmt"
	"testing"
)

func TestRegex(t *testing.T) {

	regex_input := "(a(?<number>a*)b+)"
	input := "aaab"

	regex := NewRegexWithParser(regex_input)
	match, ifmatch, named := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}

	fmt.Printf("%+v\n", match)

	fmt.Printf("%+v\n", named["number"])

}

func TestPLUS(t *testing.T) {

	regex_input := "([0-9]+)"
	input := "3344"

	regex := NewRegexWithParser(regex_input)
	match, ifmatch, _ := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}

	fmt.Printf("%+v\n", match)
}

func TestCharClass(t *testing.T) {

	regex_input := "(\\d+)a"
	input := "3344a"

	regex := NewRegexWithParser(regex_input)
	match, ifmatch, _ := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}

	fmt.Printf("%+v\n", match)
}

func TestIdent(t *testing.T) {

	regex_input := "(?<IDENT>[a-zA-Z_]\\w*)"
	input := "var1"

	regex := NewRegexWithParser(regex_input)
	match, ifmatch, named := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}

	fmt.Printf("%+v\n", match)

	fmt.Printf("%+v\n", named)

}
func TestLexer(t *testing.T) {

	regex_input := "(?<SET>set)|(?<IDENT>[a-zA-Z_]\\w*)|(?<NUMBER>\\d+)|(?<PLUS>\\+)"
	input := "+"

	regex := NewRegexWithParser(regex_input)
	match, ifmatch, named := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}

	fmt.Printf("%+v\n", match)

	fmt.Printf("%+v\n", named)

}
