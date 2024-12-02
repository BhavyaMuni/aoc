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

func part1(input string) int {
	out := 0
	for _, line := range strings.Split(input, "\n") {
		isIncreasing, isDecreasing, isSafe := true, true, true
		reports := []int{}

		for _, c := range strings.Split(line, " ") {
			val, err := strconv.Atoi(c)
			if err != nil {
				fmt.Println(err.Error())
			}
			reports = append(reports, val)
		}

		for i, val := range reports {
			if i == 0 {
				continue
			}

			isDecreasing = val > reports[i-1] && isDecreasing
			isIncreasing = val < reports[i-1] && isIncreasing
			diff := math.Abs(float64(val - reports[i-1]))
			isSafe = diff >= 1 && diff <= 3 && isSafe && (isDecreasing || isIncreasing)
		}

		if isSafe {
			out++
		}
	}
	return out
}

func part2(input string) int {
	return 0
}
