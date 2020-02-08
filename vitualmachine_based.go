package regex

import (
	"fmt"
	"github.com/KazumaTakata/shunting-yard"
)

type Inst struct {
	opcode     Opcode
	char       byte
	jump1      int
	jump2      int
	save_id    int
	save_group string
	char_class []char_class_range
}

type char_class_range struct {
	begin byte
	end   byte
}

type Opcode int

const (
	Char      Opcode = 0
	Jmp       Opcode = 1
	Split     Opcode = 2
	Save      Opcode = 3
	Match     Opcode = 4
	CharClass Opcode = 5
)

type Regex struct {
	instructions []Inst
	group_number int
	actual_group []int
}

func (re *Regex) Match(input string) ([]int, bool, map[string]*group_cap) {

	saved := make([]int, (re.group_number)*2)

	saved_group := map[string]*group_cap{}

	matched := Execute(re.instructions, input, 0, 0, saved, saved_group)

	return saved, matched, saved_group

}

func NewRegex(input_regex string) Regex {

	operators := []shunting.Operator{}
	operators = append(operators, shunting.Operator{Value: '|', Precedence: 0, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: ',', Precedence: 1, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: '+', Precedence: 2, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: '*', Precedence: 2, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: '?', Precedence: 2, IsLeftAssociative: true})

	i2p := shunting.NewIn2Post(operators, true)

	fmt.Printf("%s\n", input_regex)

	preprocessed := Preprocess(input_regex)

	fmt.Printf("%s\n", preprocessed)

	postfix := i2p.Parse(preprocessed)
	fmt.Printf("%s\n", postfix)

	postfix = []byte(postfix)
	insts, paren_count := compileToBytecode(postfix)
	regex := Regex{instructions: insts, group_number: paren_count}

	return regex

}

type group_cap struct {
	begin int
	end   int
}

func Execute(instructions []Inst, input string, pc, sp int, saved []int, named_saved map[string]*group_cap) bool {

	for {
		switch instructions[pc].opcode {
		case Char:
			{

				if sp > len(input)-1 {
					return false
				}
				if instructions[pc].char != input[sp] {
					return false
				}

				pc++
				sp++
				continue

			}
		case CharClass:
			{

				if sp > len(input)-1 {
					return false
				}

				matched := false
				for _, class_range := range instructions[pc].char_class {
					if class_range.begin <= input[sp] && input[sp] <= class_range.end {
						matched = true
					}
				}

				if !matched {
					return false
				}

				pc++
				sp++
				continue

			}
		case Match:
			{
				return true
			}
		case Jmp:
			{
				pc = pc + instructions[pc].jump1
				continue
			}
		case Split:
			{
				if Execute(instructions, input, pc+instructions[pc].jump1, sp, saved, named_saved) {
					return true
				}
				pc = pc + instructions[pc].jump2
				continue
			}
		case Save:
			{
				old := saved[instructions[pc].save_id]
				saved[instructions[pc].save_id] = sp

				old_n := named_saved[instructions[pc].save_group]

				if instructions[pc].save_group != "" {
					if _, ok := named_saved[instructions[pc].save_group]; !ok {
						named_saved[instructions[pc].save_group] = &group_cap{}
					}
					if instructions[pc].save_id%2 == 0 {
						named_saved[instructions[pc].save_group].begin = sp
					} else {
						named_saved[instructions[pc].save_group].end = sp
					}
				}

				if Execute(instructions, input, pc+1, sp, saved, named_saved) {
					return true
				}

				saved[instructions[pc].save_id] = old
				named_saved[instructions[pc].save_group] = old_n

				return false

			}
		}

	}
}

type InstStack struct {
	stack [][]Inst
}

func (s *InstStack) push(inst []Inst) {
	s.stack = append(s.stack, inst)
}

func (s *InstStack) pop() []Inst {
	top := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return top
}

func (s *InstStack) empty() bool {
	if len(s.stack) > 0 {
		return false
	}
	return true
}

func compileToBytecode(postfix []byte) ([]Inst, int) {

	inst_stack := InstStack{}
	paren_counter := 0
	paren_stack := []int{}
	group_number := 0

	for i := 0; i < len(postfix); i++ {
		regex_ch := postfix[i]
		switch regex_ch {

		case ',':
			{
				prev_inst := inst_stack.pop()
				prev_inst2 := inst_stack.pop()
				new_inst := append(prev_inst2, prev_inst...)
				inst_stack.push(new_inst)

			}

		case '|':
			{
				prev_inst2 := inst_stack.pop()
				prev_inst := inst_stack.pop()
				code_length := len(prev_inst)
				code_length2 := len(prev_inst2)
				inst := Inst{opcode: Split, jump1: 1, jump2: code_length + 2}
				new_inst := append([]Inst{inst}, prev_inst...)
				new_inst = append(new_inst, Inst{opcode: Jmp, jump1: code_length2 + 1})
				new_inst = append(new_inst, prev_inst2...)
				inst_stack.push(new_inst)

			}

		case '?':
			{
				prev_inst := inst_stack.pop()
				code_length := len(prev_inst)

				inst := Inst{opcode: Split, jump1: 1, jump2: code_length + 1}
				new_inst := append([]Inst{inst}, prev_inst...)
				inst_stack.push(new_inst)

			}

		case '*':
			{
				prev_inst := inst_stack.pop()
				code_length := len(prev_inst)

				inst := Inst{opcode: Split, jump1: 1, jump2: code_length + 2}
				new_inst := append([]Inst{inst}, prev_inst...)
				new_inst = append(new_inst, Inst{opcode: Jmp, jump1: -code_length - 1})
				inst_stack.push(new_inst)

			}

		case '+':
			{
				prev_inst := inst_stack.pop()
				inst := Inst{opcode: Split, jump1: -len(prev_inst), jump2: 1}

				new_inst := append(prev_inst, inst)
				inst_stack.push(new_inst)
			}
		case '(':
			{
				group_number++
				inst_stack.push([]Inst{Inst{opcode: Save, save_id: paren_counter}})
				paren_stack = append(paren_stack, paren_counter)
				paren_counter += 2
			}
		case '[':
			{
				i++
				range_group := []char_class_range{}
				for postfix[i] != ']' {
					if postfix[i+1] != '-' {
						char_range := char_class_range{begin: postfix[i], end: postfix[i]}
						range_group = append(range_group, char_range)
						i++
					} else {
						start := postfix[i]
						end := postfix[i+2]
						char_range := char_class_range{begin: start, end: end}
						range_group = append(range_group, char_range)
						i += 3
					}
				}
				inst := Inst{opcode: CharClass, char_class: range_group}
				inst_stack.push([]Inst{inst})

			}
		case ')':
			{
				paren_index := paren_stack[len(paren_stack)-1]
				paren_stack = paren_stack[:len(paren_stack)-1]

				prev_inst := inst_stack.pop()
				prev_paren := inst_stack.pop()
				new_inst := append(prev_inst, Inst{opcode: Save, save_id: paren_index + 1})
				new_inst = append(prev_paren, new_inst...)
				inst_stack.push(new_inst)

			}
		default:
			{
				if regex_ch == '\\' {
					i++
					regex_ch = postfix[i]
				}
				inst := Inst{opcode: Char, char: regex_ch}
				inst_stack.push([]Inst{inst})
			}
		}

	}

	instructios := inst_stack.pop()
	match := Inst{opcode: Match}
	instructios = append(instructios, match)

	return instructios, group_number
}
