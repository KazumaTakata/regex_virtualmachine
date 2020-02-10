package regex

import (
	"fmt"
	"strings"
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

func TestLexerArithmetic(t *testing.T) {
	lexer_rules := [][]string{}
	lexer_rules = append(lexer_rules, []string{"NUMBER", "\\d+"})
	lexer_rules = append(lexer_rules, []string{"ADD", "\\+"})
	lexer_rules = append(lexer_rules, []string{"SUB", "\\-"})
	lexer_rules = append(lexer_rules, []string{"MUL", "\\*"})
	lexer_rules = append(lexer_rules, []string{"DIV", "\\/"})

	regex_parts := []string{}

	for _, rule := range lexer_rules {
		regex_parts = append(regex_parts, fmt.Sprintf("(?<%s>%s)", rule[0], rule[1]))
	}

	regex_string := strings.Join(regex_parts, "|")
	fmt.Printf("%s\n", regex_string)
}

func TestDot(t *testing.T) {

	regex_input := "(?<DOT>\\d+\\.\\d*)"
	input := "34.3"

	regex := NewRegexWithParser(regex_input)
	match, ifmatch, named := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}

	fmt.Printf("%+v\n", match)

	fmt.Printf("%+v\n", named["DOT"])

}
