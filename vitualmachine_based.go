package main

import (
	"fmt"
	"github.com/KazumaTakata/shunting-yard"
)

type Inst struct {
	opcode  Opcode
	char    byte
	jump1   int
	jump2   int
	save_id int
}

type Opcode int

const (
	Char  Opcode = 0
	Jmp   Opcode = 1
	Split Opcode = 2
	Save  Opcode = 3
	Match Opcode = 4
)

type Regex struct {
	instructions []Inst
	group_number int
}

func (re *Regex) Match(input string) ([]int, bool) {

	saved := make([]int, (re.group_number)*2)

	matched := Execute(re.instructions, input, 0, 0, saved)

	return saved, matched

}

func NewRegex(input_regex string) Regex {

	operators := []shunting.Operator{}
	operators = append(operators, shunting.Operator{Value: '|', Precedence: 0, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: ',', Precedence: 1, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: '+', Precedence: 2, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: '*', Precedence: 2, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: '?', Precedence: 2, IsLeftAssociative: true})

	i2p := shunting.NewIn2Post(operators, true)

	preprocessed := Preprocess(input_regex)
	fmt.Printf("%s\n", preprocessed)

	postfix := i2p.Parse(preprocessed)
	fmt.Printf("%s\n", postfix)

	postfix = []byte(postfix)
	insts, paren_count := compileToBytecode(postfix)

	regex := Regex{instructions: insts, group_number: paren_count}

	return regex

}

func main() {

	regex := NewRegex("a+|[0-9]+")
	match, ifmatch := regex.Match("0034")
	if ifmatch {
		fmt.Printf("%v", match)
	} else {
		fmt.Printf("not match")
	}
}

func Execute(instructions []Inst, input string, pc, sp int, saved []int) bool {

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
				if Execute(instructions, input, pc+instructions[pc].jump1, sp, saved) {
					return true
				}
				pc = pc + instructions[pc].jump2
				continue
			}
		case Save:
			{
				old := saved[instructions[pc].save_id]
				saved[instructions[pc].save_id] = sp

				if Execute(instructions, input, pc+1, sp, saved) {
					return true
				}

				saved[instructions[pc].save_id] = old
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
	group_number := 0

	for _, regex_ch := range postfix {
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
				if !inst_stack.empty() {
					prev_inst := inst_stack.pop()
					new_inst := append(prev_inst, Inst{opcode: Save, save_id: paren_counter})
					inst_stack.push(new_inst)
				} else {
					inst_stack.push([]Inst{Inst{opcode: Save, save_id: paren_counter}})
				}

				paren_counter += 2
			}
		case ')':
			{
				paren_counter -= 2
				prev_inst := inst_stack.pop()
				new_inst := append(prev_inst, Inst{opcode: Save, save_id: paren_counter + 1})
				inst_stack.push(new_inst)

			}

		default:
			{
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
