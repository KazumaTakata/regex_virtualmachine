package regex

import (
	"testing"
)

func TestPreprocess(t *testing.T) {

	input := "a+b+"
	expected := "(a+,b+)"

	output := Preprocess(input)

	if output != expected {
		t.Errorf("Preprocess expected:%s, got:%s", expected, output)
	}
}

func TestShorthandCharacterWord(t *testing.T) {

	input := "\\w"
	expected := "((A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z|0|1|2|3|4|5|6|7|8|9|_))"

	output := Preprocess(input)

	if output != expected {
		t.Errorf("Preprocess expected:%s, got:%s", expected, output)
	}

}
func TestShorthandCharacterDigit(t *testing.T) {

	input := "\\da"
	expected := "((0|1|2|3|4|5|6|7|8|9),a)"

	output := Preprocess(input)

	if output != expected {
		t.Errorf("Preprocess expected:%s, got:%s", expected, output)
	}

}
