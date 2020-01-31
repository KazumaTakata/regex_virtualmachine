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

func is_quantifier(ch rune) bool {
	if ch == '+' || ch == '*' || ch == '?' {
		return true
	}
	return false
}
func main() {

	operators := []shunting.Operator{}
	operators = append(operators, shunting.Operator{Value: '|', Precedence: 0, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: ',', Precedence: 1, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: '+', Precedence: 2, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: '*', Precedence: 2, IsLeftAssociative: true})
	operators = append(operators, shunting.Operator{Value: '?', Precedence: 2, IsLeftAssociative: true})

	i2p := shunting.NewIn2Post(operators, true)

	input_regex := "a+(b+)"

	preprocessed := Preprocess(input_regex)
	fmt.Printf("%s\n", preprocessed)

	postfix := i2p.Parse(preprocessed)
	postfix = []byte(postfix)
	fmt.Printf("%s\n", postfix)
	insts := compileToBytecode(postfix)

	for i, ins := range insts {
		fmt.Printf("%d: %+v\n", i, ins)
	}

	saved := make([]int, 100)

	if Execute(insts, "aabbbbbb", 0, 0, saved) {
		fmt.Printf("matched")
		fmt.Printf("%v", saved)
	} else {
		fmt.Printf("not matched")
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

func compileToBytecode(postfix []byte) []Inst {

	inst_stack := InstStack{}
	paren_counter := 2

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
				if !inst_stack.empty() {
					prev_inst := inst_stack.pop()
					new_inst := append(prev_inst, Inst{opcode: Save, save_id: paren_counter})
					inst_stack.push(new_inst)
				} else {
					inst_stack.push([]Inst{Inst{opcode: Save, save_id: paren_counter}})
				}
			}
		case ')':
			{
				prev_inst := inst_stack.pop()
				new_inst := append(prev_inst, Inst{opcode: Save, save_id: paren_counter + 1})
				inst_stack.push(new_inst)
				paren_counter++
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

	return instructios
}
