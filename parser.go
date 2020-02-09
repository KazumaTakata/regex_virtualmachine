package regex

import "fmt"

type Repetition int

const (
	ASTERISK Repetition = 0
	QUESTION Repetition = 1
	PLUS     Repetition = 2
	NONE     Repetition = 3
)

var paren_counter = 0
var paren_stack = paren_Stack{}
var group_stack = group_Stack{}

type paren_Stack struct {
	data []int
}

func (s *paren_Stack) pop() int {
	top := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return top
}

func (s *paren_Stack) top() int {
	return s.data[len(s.data)-1]
}

func (s *paren_Stack) push(d int) {
	s.data = append(s.data, d)

}

type group_Stack struct {
	data []string
}

func (s *group_Stack) pop() string {
	top := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return top
}

func (s *group_Stack) top() string {
	return s.data[len(s.data)-1]
}

func (s *group_Stack) push(d string) {
	s.data = append(s.data, d)

}

type regex struct {
	terms []*Term
}

func (re *regex) gen() []Inst {

	var alternated []Inst

	if len(re.terms) == 1 {
		return re.terms[0].gen()
	}

	for i, term := range re.terms {
		if i == 0 {
			alternated = term.gen()
			continue
		}
		code_length := len(alternated)
		new_term := term.gen()
		code_length2 := len(new_term)
		inst := Inst{opcode: Split, jump1: 1, jump2: code_length + 2}
		new_inst := append([]Inst{inst}, alternated...)
		new_inst = append(new_inst, Inst{opcode: Jmp, jump1: code_length2 + 1})
		new_inst = append(new_inst, new_term...)
		alternated = new_inst
	}

	return alternated

}

func appendMatch(insts []Inst) []Inst {
	match := Inst{opcode: Match}
	insts = append(insts, match)
	return insts
}

type Term struct {
	factors []*Factor
}

func (te *Term) gen() []Inst {

	var concat []Inst

	if len(te.factors) == 1 {
		return te.factors[0].gen()
	}

	for i, factor := range te.factors {
		if i == 0 {
			concat = factor.gen()
			continue
		}
		concat = append(concat, factor.gen()...)
	}

	return concat

}

type Factor struct {
	base       *Base
	repitition Repetition
}

func (fa *Factor) gen() []Inst {

	base := fa.base.gen()
	if fa.repitition == QUESTION {

		inst := Inst{opcode: Split, jump1: 1, jump2: len(base) + 1}
		new_inst := append([]Inst{inst}, base...)
		return new_inst

	} else if fa.repitition == ASTERISK {

		inst := Inst{opcode: Split, jump1: 1, jump2: len(base) + 2}
		new_inst := append([]Inst{inst}, base...)
		new_inst = append(new_inst, Inst{opcode: Jmp, jump1: -len(base) - 1})
		return new_inst

	} else if fa.repitition == PLUS {
		inst := Inst{opcode: Split, jump1: -len(base), jump2: 1}
		new_inst := append(base, inst)
		return new_inst
	}

	return base

}

type Base struct {
	char       byte
	if_escaped bool
	regex      *regex
	group_name string
	ch_range   []char_class_range
}

func (ba *Base) gen() []Inst {

	if ba.regex != nil {
		//add save instruction
		save_inst := Inst{opcode: Save, save_id: paren_counter, save_group: ba.group_name}
		paren_stack.push(paren_counter)
		group_stack.push(ba.group_name)
		fmt.Printf("%+v\n", ba.group_name)
		paren_counter = paren_counter + 2

		new_inst := append([]Inst{save_inst}, ba.regex.gen()...)
		paren_id := paren_stack.pop()
		group_name := group_stack.pop()
		save_inst = Inst{opcode: Save, save_id: paren_id + 1, save_group: group_name}
		new_inst = append(new_inst, save_inst)

		return new_inst

	}

	if len(ba.ch_range) > 0 {
		insts := []Inst{Inst{opcode: CharClass, char_class: ba.ch_range}}
		return insts
	}

	if ba.if_escaped {
		if ba.char == 'd' {
			digit_range := char_class_range{begin: '0', end: '9'}
			insts := []Inst{Inst{opcode: CharClass, char_class: []char_class_range{digit_range}}}
			return insts
		} else if ba.char == 'w' {
			under_range := char_class_range{begin: '_', end: '_'}
			small_range := char_class_range{begin: 'a', end: 'z'}
			capital_range := char_class_range{begin: 'A', end: 'Z'}
			digit_range := char_class_range{begin: '0', end: '9'}

			insts := []Inst{Inst{opcode: CharClass, char_class: []char_class_range{under_range, small_range, capital_range, digit_range}}}
			return insts

		}
	}

	inst := []Inst{Inst{opcode: Char, char: ba.char}}
	return inst

}

type Regex_Input struct {
	input       string
	paren_count int
}

func (re *Regex_Input) peek() byte {
	return re.input[0]
}

func (re *Regex_Input) empty() bool {
	return len(re.input) == 0
}

func (re *Regex_Input) peek2() byte {
	return re.input[1]
}

func (re *Regex_Input) next() byte {
	ch := re.peek()
	re.eat(ch)
	return ch
}

func (re *Regex_Input) eat(ch byte) {
	if re.peek() == ch {
		re.input = re.input[1:]
	} else {
		fmt.Errorf("eat is not match:got %c, expected %c", ch, re.peek())
	}
}
func (re *Regex_Input) parse_Regex() regex {

	terms := []*Term{}

	term := re.parse_Term()
	terms = append(terms, term)

	for !re.empty() && re.peek() == '|' {
		re.eat('|')
		term = re.parse_Term()
		terms = append(terms, term)
	}
	regex := regex{terms: terms}

	return regex
}

func (re *Regex_Input) parse_Term() *Term {

	factors := []*Factor{}

	for len(re.input) > 0 && re.peek() != ')' && re.peek() != '|' {
		factor := re.parse_Factor()
		factors = append(factors, factor)
	}
	term := &Term{factors: factors}

	return term

}

func (re *Regex_Input) parse_Factor() *Factor {

	factor := &Factor{}
	base := re.parse_Base()

	factor.base = base

	if len(re.input) > 0 {
		if re.peek() == '*' {
			re.eat('*')
			factor.repitition = ASTERISK
		} else if re.peek() == '+' {
			re.eat('+')
			factor.repitition = PLUS
		} else if re.peek() == '?' {
			re.eat('?')
			factor.repitition = QUESTION
		} else {
			factor.repitition = NONE
		}
	} else {
		factor.repitition = NONE

	}
	return factor
}

func (re *Regex_Input) parse_Base() *Base {

	switch re.peek() {
	case '(':
		{
			re.paren_count++
			re.eat('(')

			if re.peek() == '?' {
				re.eat('?')
				if re.peek() == '<' {
					re.eat('<')
					group_name := ""
					for re.peek() != '>' {
						group_name += string(re.next())
					}
					re.eat('>')

					regex := re.parse_Regex()
					re.eat(')')
					base := &Base{}
					base.group_name = group_name
					base.regex = &regex
					return base

				}

			} else {

				regex := re.parse_Regex()
				re.eat(')')
				base := &Base{}
				base.regex = &regex
				return base
			}
		}
	case '\\':
		{
			re.eat('\\')
			escaped := re.next()
			base := &Base{}
			base.char = escaped
			base.if_escaped = true
			return base

		}
	case '[':
		{
			re.eat('[')

			base := &Base{}
			ch_range_list := []char_class_range{}

			for re.peek() != ']' {
				if re.peek2() == '-' {
					begin := re.next()
					re.eat('-')
					end := re.next()
					ch_range := char_class_range{begin: begin, end: end}
					ch_range_list = append(ch_range_list, ch_range)
				} else {
					ch := re.next()
					ch_range := char_class_range{begin: ch, end: ch}
					ch_range_list = append(ch_range_list, ch_range)
				}

			}
			re.eat(']')

			base.ch_range = ch_range_list
			return base
		}
	default:
		{
			ch := re.next()
			base := &Base{}
			base.char = ch
			base.if_escaped = false
			return base
		}

	}
	return nil
}

func NewRegexWithParser(input_regex string) Regex {

	// init global variable
	paren_counter = 0
	paren_stack = paren_Stack{}
	group_stack = group_Stack{}

	regex_input := Regex_Input{input: input_regex}
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
