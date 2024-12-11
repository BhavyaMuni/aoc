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
	trails := [][]int{}

	for _, line := range strings.Split(input, "\n") {
		trail := []int{}
		for _, c := range line {
			n, _ := strconv.Atoi(string(c))
			trail = append(trail, n)
		}
		trails = append(trails, trail)
	}

	visited := map[[2]int]bool{}

	// write a dfs to find path from 0 to 9 in the trails
	var dfs func(x, y, n int)

	dfs = func(x, y, n int) {
		if x < 0 || y < 0 || x >= len(trails) || y >= len(trails[x]) {
			return
		}
		if trails[x][y] == 9 && n == 9 {
			visited[[2]int{x, y}] = true
			return
		}
		if trails[x][y] != n {
			return
		}
		dfs(x+1, y, n+1)
		dfs(x-1, y, n+1)
		dfs(x, y+1, n+1)
		dfs(x, y-1, n+1)
	}

	ct := 0
	for x, row := range trails {
		for y, n := range row {
			if n == 0 {
				dfs(x, y, 0)
				ct += len(visited)
				visited = make(map[[2]int]bool)
			}
		}
	}
	return ct
}

func part2(input string) int {
	trails := [][]int{}

	for _, line := range strings.Split(input, "\n") {
		trail := []int{}
		for _, c := range line {
			n, _ := strconv.Atoi(string(c))
			trail = append(trail, n)
		}
		trails = append(trails, trail)
	}

	// write a dfs to find path from 0 to 9 in the trails
	var dfs func(x, y, n int) int

	dfs = func(x, y, n int) int {
		if x < 0 || y < 0 || x >= len(trails) || y >= len(trails[x]) {
			return 0
		}
		if trails[x][y] == 9 && n == 9 {
			return 1
		}
		if trails[x][y] != n {
			return 0
		}
		return dfs(x+1, y, n+1) + dfs(x-1, y, n+1) + dfs(x, y+1, n+1) + dfs(x, y-1, n+1)
	}

	ct := 0
	for x, row := range trails {
		for y, n := range row {
			if n == 0 {
				ct += dfs(x, y, 0)
			}
		}
	}
	return ct
}
