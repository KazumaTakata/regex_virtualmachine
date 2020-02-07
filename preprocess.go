package regex

import (
	"strings"
	"unicode"
)

func is_quantifier(ch rune) bool {
	if ch == '+' || ch == '*' || ch == '?' {
		return true
	}
	return false
}
func Expand_character_classes(input string) string {

	output := ""

	for len(input) > 0 {
		if input[0] == '[' {
			input = input[1:]
			alternate_element := []string{}
			for input[0] != ']' {
				if input[1] != '-' {
					alternate_element = append(alternate_element, string(input[0]))
				} else {
					start := input[0]
					end := input[2]
					for start != end+1 {
						alternate_element = append(alternate_element, string(start))
						start = start + 1
					}
					input = input[2:]
				}
				input = input[1:]
			}
			input = input[1:]
			character_class := "(" + strings.Join(alternate_element, "|") + ")"
			output = output + character_class

		} else {
			output = output + string(input[0])
			input = input[1:]
		}

	}

	return output

}

func expand_shorthand_character(input string) string {
	output := ""

	for len(input) > 0 {
		if input[0] == '\\' {
			input = input[1:]
			if input[0] == 'd' {
				output = output + "[0-9]"
			} else if input[0] == 'w' {
				output = output + "[A-Za-z0-9_]"
			} else {
				output = output + "\\" + string(input[0])

			}

			input = input[1:]
		} else {
			output = output + string(input[0])
			input = input[1:]
		}
	}
	return output

}

func add_concat_regex(input string) string {

	output := ""
	var prev_rune rune
	for i, cur_rune := range input {
		if i == 0 {
			prev_rune = cur_rune
			output = output + string(cur_rune)
			continue
		}
		if unicode.IsNumber(prev_rune) || unicode.IsLetter(prev_rune) || is_quantifier(prev_rune) || prev_rune == ')' || prev_rune == '=' {
			if !is_quantifier(cur_rune) && cur_rune != rune('|') && cur_rune != rune(')') {
				output = output + string(',')
			}
		}
		output = output + string(cur_rune)
		prev_rune = cur_rune

	}

	return output
}

func Preprocess(input string) string {
	output := expand_shorthand_character(input)
	output = Expand_character_classes(output)
	output = add_concat_regex(output)
	output = "(" + output + ")"

	return output
}
