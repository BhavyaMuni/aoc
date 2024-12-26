package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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

	t := time.Now()
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
	fmt.Println("Time:", time.Since(t))
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
func adv(operand uint64) (uint64, uint64, uint64, int, uint64) {
	a /= uint64(math.Pow(2, float64(ComboOperand(operand))))
	return a, b, c, ip + 2, 0

}
func bxl(operand uint64) (uint64, uint64, uint64, int, uint64) {
	b = uint64(b) ^ uint64(operand)
	return a, b, c, ip + 2, 0

}
func bst(operand uint64) (uint64, uint64, uint64, int, uint64) {
	b = uint64(ComboOperand(operand)) % 8
	return a, b, c, ip + 2, 0

}
func jnz(operand uint64) (uint64, uint64, uint64, int, uint64) {
	if a == 0 {
		return a, b, c, ip + 2, 0

	}
	return a, b, c, int(operand), 0

}
func bxc(operand uint64) (uint64, uint64, uint64, int, uint64) {
	b = b ^ c
	return a, b, c, ip + 2, 0

}
func out(operand uint64) (uint64, uint64, uint64, int, uint64) {
	out := ComboOperand(operand) % 8
	return a, b, c, ip + 2, out

}
func bdv(operand uint64) (uint64, uint64, uint64, int, uint64) {
	b = a / uint64(math.Pow(2, float64(ComboOperand(operand))))
	return a, b, c, ip + 2, 0

}
func cdv(operand uint64) (uint64, uint64, uint64, int, uint64) {
	c = a / uint64(math.Pow(2, float64(ComboOperand(operand))))
	return a, b, c, ip + 2, 0
}

var opcodes = map[uint64]func(uint64) (uint64, uint64, uint64, int, uint64){
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
		var next int
		var out uint64
		_, _, _, next, out = opcode(operand)
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

type seenOut struct {
	a, b, c, out uint64
	next         int
}

func getOutinDec(inp []uint64) int64 {
	var out int64 = 0
	for i, n := range inp {
		out += int64(n) * int64(math.Pow(8, float64(i)))
	}
	return out
}

func part2(input string) int {
	program := []uint64{}
	programInp := strings.TrimPrefix(strings.Split(input, "\n\n")[1], "Program: ")
	for _, c := range strings.Split(programInp, ",") {
		n, _ := strconv.Atoi(c)
		program = append(program, uint64(n))
	}

	seen := make(map[[5]uint64]seenOut)

	l := 0
	r := math.MaxInt64
	out := math.MaxInt64
	for l < r {
		mid := l + (r-l)/2
		fmt.Println("Trying: ", mid)
		ip = 0
		a = uint64(mid)
		b = 0
		c = 0

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
			var out uint64
			if _, ok := seen[[5]uint64{a, b, c, op, operand}]; ok {
				seenOut := seen[[5]uint64{a, b, c, op, operand}]
				a = seenOut.a
				b = seenOut.b
				c = seenOut.c
				ip = seenOut.next
				out = seenOut.out
			} else {
				oldA, oldB, oldC := a, b, c
				a, b, c, ip, out = opcode(operand)
				seen[[5]uint64{oldA, oldB, oldC, op, operand}] = seenOut{a, b, c, out, ip}
			}
			if op == 5 {
				outs = append(outs, out)
			}
		}
		if getOutinDec(outs) == getOutinDec(program) {
			out = mid
			r = mid
			fmt.Println("Found: ", mid)
		} else if getOutinDec(outs) < getOutinDec(program) {
			l = mid + 1
		} else {
			r = mid
		}
	}

	fmt.Println("Out:", out)

	return 0
}
