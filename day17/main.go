package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.design/x/clipboard"
)

//go:embed input.txt
var input string

func init() {
	// do this in init (not main) so test file has same input
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty input.txt file")
	}
}

func main() {
	var part int
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.Parse()
	fmt.Println("Running part", part)
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	if part == 1 {
		ans := part1(input)
		clipboard.Write(clipboard.FmtText, []byte(fmt.Sprintf("%v", ans)))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		clipboard.Write(clipboard.FmtText, []byte(fmt.Sprintf("%v", ans)))
		fmt.Println("Output:", ans)
	}
}

var a uint64 = 0
var b uint64 = 0
var c uint64 = 0
var ip = 0

func ComboOperand(input uint64) uint64 {
	switch input {
	case 4:
		return a
	case 5:
		return b
	case 6:
		return c
	case 7:
		return 0
	default:
		return input
	}
}
func adv(operand uint64) (uint64, int) {
	a /= uint64(math.Pow(2, float64(ComboOperand(operand))))
	return a, ip + 2
}
func bxl(operand uint64) (uint64, int) {
	b = uint64(b) ^ uint64(operand)
	return b, ip + 2
}
func bst(operand uint64) (uint64, int) {
	b = uint64(ComboOperand(operand)) % 8
	return b, ip + 2
}
func jnz(operand uint64) (uint64, int) {
	if a == 0 {
		return operand, ip + 2
	}
	return operand, int(operand)
}
func bxc(operand uint64) (uint64, int) {
	b = b ^ c
	return b, ip + 2
}
func out(operand uint64) (uint64, int) {
	out := ComboOperand(operand) % 8
	return out, ip + 2
}
func bdv(operand uint64) (uint64, int) {
	b = a / uint64(math.Pow(2, float64(ComboOperand(operand))))
	return b, ip + 2
}
func cdv(operand uint64) (uint64, int) {
	c = a / uint64(math.Pow(2, float64(ComboOperand(operand))))
	return c, ip + 2
}

var opcodes = map[uint64]func(uint64) (uint64, int){
	0: adv,
	1: bxl,
	2: bst,
	3: jnz,
	4: bxc,
	5: out,
	6: bdv,
	7: cdv,
}

func part1(input string) int {
	program := []uint64{}
	programInp := strings.TrimPrefix(strings.Split(input, "\n\n")[1], "Program: ")
	for _, c := range strings.Split(programInp, ",") {
		n, _ := strconv.Atoi(c)
		program = append(program, uint64(n))
	}

	ip = 0
	a = 30344604
	outs := []uint64{}
	for ip < len(program) {
		op := program[ip]
		operand := program[ip+1]
		if op == 8 {
			break
		}
		opcode, ok := opcodes[op]
		if !ok {
			panic("unknown opcode")
		}
		out, next := opcode(operand)
		ip = next
		if op == 5 {
			outs = append(outs, out)
		}
	}
	for _, o := range outs {
		fmt.Print(o, ",")
	}
	fmt.Println()

	return 0

}

func part2(input string) int {
	panic("not implemented")
}
