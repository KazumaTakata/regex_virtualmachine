package regex

import (
	"fmt"
	"testing"
)

func TestRegex(t *testing.T) {

	regex_input := "aa*b+"
	input := "ab"

	regex := NewRegexWithParser(regex_input)
	_, ifmatch := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}

}

func NewRegexWithParser(input_regex string) Regex {
	regex_input := Regex_Input{input: "aa*b+"}
	regex := regex_input.parse_Regex()
	fmt.Printf("%+v\n", regex)

	insts := regex.gen()
	insts = appendMatch(insts)
	for _, inst := range insts {
		fmt.Printf("%+v\n", inst)

	}

	regex_struct := Regex{instructions: insts, group_number: regex_input.paren_count}
	return regex_struct

}
