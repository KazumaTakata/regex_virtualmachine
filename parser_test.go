package regex

import (
	"fmt"
	"testing"
)

func TestRegex(t *testing.T) {
	regex_input := Regex_Input{input: "aaa"}
	regex := regex_input.parse_Regex()
	fmt.Printf("%+v", regex)
}
