package main

import (
	_ "embed"
	"flag"
	"fmt"
	"regexp"
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

func part1(input string) int {
	re, error := regexp.Compile(`mul\([0-9]{1,3},[0-9]{1,3}\)`)
	if error != nil {
		fmt.Println("Error compiling regex")
	}

	matches := re.FindAll([]byte(input), -1)
	out := 0
	for _, match := range matches {
		var n1, n2 int
		fmt.Sscanf(string(match), "mul(%d,%d)", &n1, &n2)
		out += n1 * n2
	}

	return out
}

func part2(input string) int {
	// match either mul() or do() don't()
	re, error := regexp.Compile(`mul\([0-9]{1,3},[0-9]{1,3}\)|do\(\)|don't\(\)`)
	if error != nil {
		fmt.Println("Error compiling regex")
	}

	matches := re.FindAll([]byte(input), -1)
	out := 0
	enabled := true
	for _, match := range matches {
		switch string(match) {
		case "do()":
			enabled = true
		case "don't()":
			enabled = false
		default:
			if enabled {
				var n1, n2 int
				fmt.Sscanf(string(match), "mul(%d,%d)", &n1, &n2)
				out += n1 * n2
			}
		}
	}

	return out
}
