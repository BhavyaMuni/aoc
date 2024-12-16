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

func part1(input string) int {
	board := make([][]rune, 0)
	for _, line := range strings.Split(input, "\n") {
		board = append(board, []rune(line))
	}
	var visited = make(map[[2]int]bool)
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	var area func(x, y int, char rune) (int, int)
	area = func(x, y int, char rune) (int, int) {
		if x < 0 || x >= len(board) || y < 0 || y >= len(board[0]) || board[x][y] != char {
			return 0, 1
		}
		if visited[[2]int{x, y}] {
			return 0, 0
		}
		visited[[2]int{x, y}] = true
		as, ps := 1, 0
		for _, d := range dirs {
			a, p := area(x+d[0], y+d[1], char)
			as += a
			ps += p
		}

		return as, ps
	}

	ct := 0

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if visited[[2]int{i, j}] {
				continue
			}
			ar, p := area(i, j, board[i][j])
			ct += ar * p
		}
	}

	return ct
}

func isSame(x, y int, z rune, grid [][]rune) bool {
	if x < 0 || x >= len(grid) || y < 0 || y >= len(grid[0]) {
		return false
	}
	return grid[x][y] == z
}

func numCorners(areas map[string]bool) int {
	corners := 0
	if areas["N"] && areas["W"] && !areas["NW"] {
		corners++
	}
	if areas["N"] && areas["E"] && !areas["NE"] {
		corners++
	}
	if areas["S"] && areas["W"] && !areas["SW"] {
		corners++
	}
	if areas["S"] && areas["E"] && !areas["SE"] {
		corners++
	}

	if !(areas["N"] || areas["W"]) {
		corners++
	}
	if !(areas["N"] || areas["E"]) {
		corners++
	}
	if !(areas["S"] || areas["W"]) {
		corners++
	}
	if !(areas["S"] || areas["E"]) {
		corners++
	}
	return corners
}

func part2(input string) int {
	board := make([][]rune, 0)
	for _, line := range strings.Split(input, "\n") {
		board = append(board, []rune(line))
	}
	var visited = make(map[[2]int]bool)
	dirs := map[string][2]int{"N": {1, 0}, "S": {-1, 0}, "E": {0, 1}, "W": {0, -1}}
	diagonals := map[string][2]int{"NW": {1, -1}, "NE": {1, 1}, "SW": {-1, -1}, "SE": {-1, 1}}
	corners := 0
	var area func(x, y int, char rune) (int, int)
	area = func(x, y int, char rune) (int, int) {
		sames := map[string]bool{}
		if x < 0 || x >= len(board) || y < 0 || y >= len(board[0]) || board[x][y] != char {
			return 0, 1
		}
		if visited[[2]int{x, y}] {
			return 0, 0
		}
		visited[[2]int{x, y}] = true
		as, ps := 1, 0
		for dir, d := range dirs {
			a, p := area(x+d[0], y+d[1], char)
			sames[dir] = isSame(x+d[0], y+d[1], char, board)
			as += a
			ps += p
		}
		for dir, d := range diagonals {
			sames[dir] = isSame(x+d[0], y+d[1], char, board)
		}
		corners += numCorners(sames)
		return as, ps
	}

	ct := 0

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if visited[[2]int{i, j}] {
				continue
			}
			ar, _ := area(i, j, board[i][j])
			ct += ar * corners
			corners = 0
		}
	}

	return ct
}
