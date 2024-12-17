package main

import (
	_ "embed"
	"flag"
	"fmt"
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

type Equation struct {
	A float64
	B float64
	C float64
}

func part1(input string) int {
	totalCost := 0

	for _, prizes := range strings.Split(input, "\n\n") {
		lines := strings.Split(prizes, "\n")
		var eqA, eqB Equation
		fmt.Sscanf(lines[0], "Button A: X+%f, Y+%f", &eqA.A, &eqB.A)
		fmt.Sscanf(lines[1], "Button B: X+%f, Y+%f", &eqA.B, &eqB.B)
		fmt.Sscanf(lines[2], "Prize: X=%f, Y=%f", &eqA.C, &eqB.C)

		eqB.B *= eqA.A
		eqB.C *= eqA.A
		eqA.B *= eqB.A
		eqA.C *= eqB.A

		b2 := (eqA.C - eqB.C) / (eqA.B - eqB.B)
		b1 := (eqA.C - (eqA.B * b2)) / (eqA.A * eqB.A)

		if b2 == float64(int(b2)) && b1 == float64(int(b1)) {
			cost := b1*3 + b2*1
			totalCost += int(cost)
		}

	}
	return totalCost
}

func part2(input string) int {
	totalCost := 0

	for _, prizes := range strings.Split(input, "\n\n") {
		lines := strings.Split(prizes, "\n")
		var eqA, eqB Equation
		fmt.Sscanf(lines[0], "Button A: X+%f, Y+%f", &eqA.A, &eqB.A)
		fmt.Sscanf(lines[1], "Button B: X+%f, Y+%f", &eqA.B, &eqB.B)
		fmt.Sscanf(lines[2], "Prize: X=%f, Y=%f", &eqA.C, &eqB.C)

		eqA.C += 10000000000000
		eqB.C += 10000000000000

		eqB.B *= eqA.A
		eqB.C *= eqA.A
		eqA.B *= eqB.A
		eqA.C *= eqB.A

		b2 := (eqA.C - eqB.C) / (eqA.B - eqB.B)
		b1 := (eqA.C - (eqA.B * b2)) / (eqA.A * eqB.A)

		if b2 == float64(int(b2)) && b1 == float64(int(b1)) {
			cost := b1*3 + b2*1
			totalCost += int(cost)
		}

	}
	return totalCost
}
