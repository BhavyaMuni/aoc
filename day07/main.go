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

func calc(nums []int, ops []int) int {
	sum := nums[0]
	for i := 0; i < len(ops); i++ {
		switch ops[i] {
		case 0:
			sum += nums[i+1]
		case 1:
			sum *= nums[i+1]
		case 2:
			sum, _ = strconv.Atoi(strconv.Itoa(sum) + strconv.Itoa(nums[i+1]))
		}
	}
	return sum
}

var permOut [][]int

func perm(x []int, n int) [][]int {
	if len(x) == n {
		permOut = append(permOut, x)
		return permOut
	}
	for _, op := range [3]int{0, 1, 2} {
		cp := append([]int{}, x...)
		cp = append(cp, op)
		perm(cp, n)
	}
	return permOut
}

func part1(input string) int {
	out := 0
	for _, line := range strings.Split(input, "\n") {
		test, _ := strconv.Atoi(strings.Split(line, ": ")[0])
		allNums := strings.Split(line, ": ")[1]
		nums := make([]int, 0)

		for _, num := range strings.Split(allNums, " ") {
			n, _ := strconv.Atoi(num)
			nums = append(nums, n)
		}

		ops := perm([]int{}, len(nums)-1)

		for _, op := range ops {
			if calc(nums, op) == test {
				out += test
				break
			}
		}
		permOut = [][]int{}
	}

	return out
}

func part2(input string) int {
	panic("not implemented")
}
