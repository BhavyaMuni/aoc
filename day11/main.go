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

func part1(input string) int {
	// brute force
	stones := make([]int, 0)
	for _, line := range strings.Split(input, " ") {
		n, _ := strconv.Atoi(line)
		stones = append(stones, n)
	}

	for i := 0; i < 25; i++ {
		for s := 0; s < len(stones); s++ {
			stone := stones[s]
			if stone == 0 {
				stones[s] = 1
			} else if len(strconv.Itoa(stone))%2 == 0 {
				stoneStr := strconv.Itoa(stone)
				h1 := stoneStr[:len(stoneStr)/2]
				h2 := stoneStr[len(stoneStr)/2:]
				h1Int, _ := strconv.Atoi(h1)
				h2Int, _ := strconv.Atoi(h2)
				newStones := append([]int{}, stones[:s]...)
				newStones = append(newStones, h1Int)
				newStones = append(newStones, h2Int)
				newStones = append(newStones, stones[s+1:]...)
				stones = newStones
				s++
			} else {
				stones[s] *= 2024
			}
		}
	}

	return len(stones)
}

func part2(input string) int {
	// dp
	stones := make([]int, 0)
	for _, line := range strings.Split(input, " ") {
		n, _ := strconv.Atoi(line)
		stones = append(stones, n)
	}
	ct := 0
	memo := make(map[[2]int]int)
	var dp func(n, l int) int
	dp = func(n, l int) int {
		if val, ok := memo[[2]int{n, l}]; ok {
			return val
		}

		if l == 75 {
			return 1
		}

		if n == 0 {
			memo[[2]int{n, l}] = dp(1, l+1)
			return memo[[2]int{n, l}]
		} else if len(strconv.Itoa(n))%2 == 0 {
			ns := strconv.Itoa(n)
			h1 := ns[:len(ns)/2]
			h2 := ns[len(ns)/2:]
			h1Int, _ := strconv.Atoi(h1)
			h2Int, _ := strconv.Atoi(h2)
			memo[[2]int{n, l}] = dp(h1Int, l+1) + dp(h2Int, l+1)
			return memo[[2]int{n, l}]
		} else {
			memo[[2]int{n, l}] = dp(n*2024, l+1)
			return memo[[2]int{n, l}]
		}
	}

	for i := 0; i < len(stones); i++ {
		ct += dp(stones[i], 0)
	}
	return ct
}
