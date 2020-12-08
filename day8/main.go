package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

type instruction struct {
	op  string
	arg int
}

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	instructions := []*instruction{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		instructions = append(instructions, parseInstruction(scanner.Text()))
	}

	acc, hasLoop := executeInstructionSet(instructions)
	if !hasLoop {
		panic("expected a loop in instructions, did not encounter one")
	}
	fmt.Println("accumulator value before first repeated instruction:", acc)

	fixInstructionSet(instructions)
	acc, _ = executeInstructionSet(instructions)
	fmt.Println("accumulator value after executing fixed instructions:", acc)
}

// destructively fixes the instruction set, ensuring there are no loops
func fixInstructionSet(instructions []*instruction) {
	for i, inst := range instructions {
		originalOp := inst.op
		// swap the jmp/nop operation if we've found one
		switch originalOp {
		case "acc":
			continue
		case "jmp":
			inst.op = "nop"
		case "nop":
			inst.op = "jmp"
		}

		// If the swapped operation causes no loop, we've fixed the instructions
		_, hasLoop := executeInstructionSet(instructions)
		if !hasLoop {
			fmt.Println("fixed instruction:", i)
			return
		}

		// reset the instruction and try the next one
		inst.op = originalOp
	}

	panic("unable to fix instructions")
}

// returns the accumulator value after executing the given instructions
// and whether any instructions have repeated.
// Since no args are dynamic, we can guarantee that any repeated instruction
// causes an infinite loop
func executeInstructionSet(instructions []*instruction) (int, bool) {
	instructionCount := make(map[int]struct{})
	acc, ptr := 0, 0
	for ptr < len(instructions) {
		_, ok := instructionCount[ptr]
		if ok {
			return acc, true
		}

		i := instructions[ptr]
		instructionCount[ptr] = struct{}{}
		i.execute(&acc, &ptr)
	}

	return acc, false
}

func (i instruction) execute(acc *int, ptr *int) {
	switch i.op {
	case "jmp":
		*ptr += i.arg
	case "acc":
		*acc += i.arg
		fallthrough
	case "nop":
		*ptr++
	}

}

func parseInstruction(line string) *instruction {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		panic("instruction was invalid")
	}
	op := parts[0]
	arg, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	return &instruction{
		op:  op,
		arg: arg,
	}
}
