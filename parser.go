package regex

import "fmt"

type Repetition int

const (
	ASTERISK Repetition = 0
	QUESTION Repetition = 1
	PLUS     Repetition = 2
	NONE     Repetition = 3
)

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
		code_length2 := len(term.gen())
		inst := Inst{opcode: Split, jump1: 1, jump2: code_length + 2}
		new_inst := append([]Inst{inst}, alternated...)
		new_inst = append(new_inst, Inst{opcode: Jmp, jump1: code_length2 + 1})
		new_inst = append(new_inst, term.gen()...)
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
}

func (ba *Base) gen() []Inst {

	if ba.regex != nil {
		return ba.regex.gen()
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
			regex := re.parse_Regex()
			re.eat(')')
			base := &Base{}
			base.regex = &regex
			return base
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
	default:
		{
			ch := re.next()
			base := &Base{}
			base.char = ch
			base.if_escaped = false
			return base
		}

	}

}
