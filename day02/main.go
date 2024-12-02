package main

import (
	_ "embed"
	"flag"
	"fmt"
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

func checkIncreasing(val, prev int) bool {
	return val-prev <= 3 && val-prev >= 1
}

func checkDecreasing(val, prev int) bool {
	return prev-val <= 3 && prev-val >= 1
}

func isSafe(values []int) bool {
	isDecreasing, isIncreasing := true, true
	for i, val := range values {
		if i == 0 {
			continue
		}
		isDecreasing = checkDecreasing(val, values[i-1]) && isDecreasing
		isIncreasing = checkIncreasing(val, values[i-1]) && isIncreasing
		if !isDecreasing && !isIncreasing {
			break
		}
	}
	return isDecreasing || isIncreasing
}

func part1(input string) int {
	out := 0
	for _, line := range strings.Split(input, "\n") {
		reports := []int{}

		for _, c := range strings.Split(line, " ") {
			val, err := strconv.Atoi(c)
			if err != nil {
				fmt.Println(err.Error())
			}
			reports = append(reports, val)
		}

		if isSafe(reports) {
			out++
		}
	}
	return out
}

func part2(input string) int {
	out := 0
	for _, line := range strings.Split(input, "\n") {
		reports := []int{}

		for _, c := range strings.Split(line, " ") {
			val, err := strconv.Atoi(c)
			if err != nil {
				fmt.Println(err.Error())
			}
			reports = append(reports, val)
		}

		if isSafe(reports) {
			out++
		} else {
			for i := 0; i < len(reports); i++ {
				reportsWithout := append([]int{}, reports...)
				reportsWithout = append(reportsWithout[:i], reportsWithout[i+1:]...)
				if isSafe(reportsWithout) {
					out++
					break
				}
			}
		}
	}
	return out
}
