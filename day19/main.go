package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
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
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	t := time.Now()
	if part == 1 {
		ans := part1(input)
		clipboard.Write(clipboard.FmtText, []byte(fmt.Sprintf("%v", ans)))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		clipboard.Write(clipboard.FmtText, []byte(fmt.Sprintf("%v", ans)))
		fmt.Println("Output:", ans)
	}
	fmt.Println("Time", time.Since(t))
}

const BOARD_SIZE = 71
const NUM_BYTES = 1024

func part1(input string) int {

	available := map[string]struct{}{}
	maxLen := math.MinInt32
	for _, v := range strings.Split(strings.Split(input, "\n\n")[0], ", ") {
		available[v] = struct{}{}
		maxLen = max(maxLen, len(v))
	}
	seen := map[string]bool{}
	var dp func(words string) bool
	dp = func(words string) bool {
		if len(words) == 0 {
			return true
		}
		if v, ok := seen[words]; ok {
			return v
		}
		for i := 0; i < min(len(words), maxLen); i++ {
			if _, ok := available[words[:i+1]]; ok && dp(words[i+1:]) {
				seen[words] = true
				return true
			}
		}
		seen[words] = false
		return false
	}

	count := 0
	for _, line := range strings.Split(strings.Split(input, "\n\n")[1], "\n") {
		if dp(line) {
			count++
		}
	}

	return count
}

func part2(input string) int {

	available := map[string]struct{}{}
	maxLen := math.MinInt32
	for _, v := range strings.Split(strings.Split(input, "\n\n")[0], ", ") {
		available[v] = struct{}{}
		maxLen = max(maxLen, len(v))
	}
	seen := map[string]int{}
	var dp func(words string) int
	dp = func(words string) int {
		if len(words) == 0 {
			return 1
		}
		if v, ok := seen[words]; ok {
			return v
		}
		ct := 0
		for i := 0; i < min(len(words), maxLen); i++ {
			if _, ok := available[words[:i+1]]; ok {
				ct += dp(words[i+1:])
				seen[words] = ct
			}
		}
		seen[words] = ct
		return ct
	}

	count := 0
	for _, line := range strings.Split(strings.Split(input, "\n\n")[1], "\n") {
		dp(line)
		count += seen[line]
	}

	return count
}
