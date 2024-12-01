package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"sort"
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
	var l1, l2 []int
	for _, line := range strings.Split(input, "\n") {
		var n1, n2 int
		fmt.Sscanf(line, "%d%d", &n1, &n2)
		l1 = append(l1, n1)
		l2 = append(l2, n2)
	}
	sort.Ints(l1)
	sort.Ints(l2)

	out := 0
	for i := 0; i < len(l1); i++ {
		out += int(math.Abs(float64(l1[i] - l2[i])))
	}
	return out
}

func part2(input string) int {
	var l2 []int
	l1cts := make(map[int]int)
	for _, line := range strings.Split(input, "\n") {
		var n1, n2 int
		fmt.Sscanf(line, "%d%d", &n1, &n2)
		l2 = append(l2, n2)
		l1cts[n1]++
	}

	out := 0
	for i := 0; i < len(l2); i++ {
		out += (l1cts[l2[i]] * l2[i])
	}
	return out
}
