package regex

import (
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

	regex_input := "\\+a+"
	input := "+aaaaa"

	regex := NewRegex(regex_input)
	match, ifmatch := regex.Match(input)

	if !ifmatch {
		t.Errorf("Regex not matched: regex:%s, input:%s", regex_input, input)
	}
	if match[0] != 0 || match[1] != len(input) {
		t.Errorf("Regex returned position is not correct: expected:[%d,%d], got:[%d,%d]", 0, len(input), match[0], match[1])
	}

}
