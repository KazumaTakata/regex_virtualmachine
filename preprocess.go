package main

import "unicode"

func Preprocess(input string) string {

	output := ""
	var prev_rune rune
	for i, cur_rune := range input {
		if i == 0 {
			prev_rune = cur_rune
			output = output + string(cur_rune)
			continue
		}
		if unicode.IsNumber(prev_rune) || unicode.IsLetter(prev_rune) || is_quantifier(prev_rune) {
			if !is_quantifier(cur_rune) && cur_rune != rune('|') && cur_rune != rune(')') {
				output = output + string(',')
			}
		}
		output = output + string(cur_rune)
		prev_rune = cur_rune

	}

	return output
}
