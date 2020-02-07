package regex

import (
	"fmt"
	"testing"
)

func TestKeyword(t *testing.T) {

	regex_input := "set"
	input := "set"

	regex := NewRegex(regex_input)
	match, ifmatch := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}
	if match[0] != 0 || match[1] != len(input) {
		t.Errorf("Regex returned position is not correct: expected:[%d,%d], got:[%d,%d]", 0, len(input), match[0], match[1])
	}

}

func TestCharacterClass(t *testing.T) {

	regex_input := "\\d+"
	input := "9455"

	regex := NewRegex(regex_input)
	match, ifmatch := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}
	if match[0] != 0 || match[1] != len(input) {
		t.Errorf("Regex returned position is not correct: expected:[%d,%d], got:[%d,%d]", 0, len(input), match[0], match[1])
	}

}

func TestEscapedCharacter(t *testing.T) {

	regex_input := "bbb\\+a+"
	input := "bbb+aaaaa"

	regex := NewRegex(regex_input)
	match, ifmatch := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}
	if match[0] != 0 || match[1] != len(input) {
		t.Errorf("Regex returned position is not correct: expected:[%d,%d], got:[%d,%d]", 0, len(input), match[0], match[1])
	}

}

func TestUnanchoredMatch(t *testing.T) {

	regex_input := "a+"
	input := "aaaaab"

	regex := NewRegex(regex_input)
	match, ifmatch := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}
	if match[0] != 0 || match[1] != len(input)-1 {
		t.Errorf("Regex returned position is not correct: expected:[%d,%d], got:[%d,%d]", 0, len(input), match[0], match[1])
	}

}

func TestIdentifierMatch(t *testing.T) {

	regex_input := "[a-zA-Z_]\\w*"
	input := "var1"

	regex := NewRegex(regex_input)
	match, ifmatch := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}
	if match[0] != 0 || match[1] != len(input) {
		t.Errorf("Regex returned position is not correct: expected:[%d,%d], got:[%d,%d]", 0, len(input), match[0], match[1])
	}

}

func TestEqualMatch(t *testing.T) {

	regex_input := "[a-zA-Z_]\\w*=\\d"

	input := "var1=1"

	regex := NewRegex(regex_input)
	match, ifmatch := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}
	if match[0] != 0 || match[1] != len(input) {
		t.Errorf("Regex returned position is not correct: expected:[%d,%d], got:[%d,%d]", 0, len(input), match[0], match[1])
	}

}

func TestAlternation(t *testing.T) {

	regex_input := "(set)|(\\d+)|([a-zA-Z_]\\w+)"

	input := "aa1"

	regex := NewRegex(regex_input)
	match, ifmatch := regex.Match(input)

	fmt.Printf("%+v", match)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}
	if match[2] != 0 || match[3] != len(input)-2 {
		t.Errorf("Regex returned position is not correct: expected:[%d,%d], got:[%d,%d]", 0, len(input)-2, match[2], match[3])
	}

}
