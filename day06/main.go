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

func turnRight(dir rune) rune {
	switch dir {
	case 'N':
		return 'E'
	case 'E':
		return 'S'
	case 'S':
		return 'W'
	case 'W':
		return 'N'
	}
	return 'N'
}

func part1(input string) int {
	var board [][]rune

	for _, line := range strings.Split(input, "\n") {
		board = append(board, []rune(line))
	}

	pos := []int{0, 0}

	dirMap := map[rune][]int{
		'N': {-1, 0},
		'S': {1, 0},
		'E': {0, 1},
		'W': {0, -1},
	}
	dir := 'N'

	for i, row := range board {
		for j, col := range row {
			if col == '^' {
				pos = []int{i, j}
			}
		}
	}
	visited := make(map[[2]int]bool)
	for !(pos[0] == len(board)-1 || pos[0] == 0 || pos[1] == len(board[0])-1 || pos[1] == 0) {
		nRow := pos[0] + dirMap[dir][0]
		nCol := pos[1] + dirMap[dir][1]
		if nRow < 0 || nCol < 0 || nRow >= len(board) || nCol >= len(board[0]) {
			continue
		}
		next := board[nRow][nCol]
		if next == '#' {
			dir = turnRight(dir)
		} else {
			pos[0] = nRow
			pos[1] = nCol
			visited[[2]int{nRow, nCol}] = true
		}
	}
	return len(visited) + 1
}

var dirMap = map[rune][]int{
	'N': {-1, 0},
	'S': {1, 0},
	'E': {0, 1},
	'W': {0, -1},
}

func checkIsStuckInLoop(pos [2]int, board [][]rune) bool {
	visited := make(map[[2]int]int)
	dir := 'N'

	for !(pos[0] == len(board)-1 || pos[0] == 0 || pos[1] == len(board[0])-1 || pos[1] == 0) {
		if visited[pos] > 4 {
			return true
		}
		nRow := pos[0] + dirMap[dir][0]
		nCol := pos[1] + dirMap[dir][1]
		if nRow < 0 || nCol < 0 || nRow >= len(board) || nCol >= len(board[0]) {
			continue
		}
		next := board[nRow][nCol]
		if next == '#' || next == 'O' {
			dir = turnRight(dir)
		} else {
			pos[0] = nRow
			pos[1] = nCol
			visited[[2]int{nRow, nCol}]++
		}
	}
	return false
}

func part2(input string) int {
	var board [][]rune

	for _, line := range strings.Split(input, "\n") {
		board = append(board, []rune(line))
	}
	pos := [2]int{0, 0}
	for i, row := range board {
		for j, col := range row {
			if col == '^' {
				pos = [2]int{i, j}
			}
		}
	}

	ct := 0
	for i, row := range board {
		for j, cell := range row {
			if cell == '.' {
				board[i][j] = 'O'
				if checkIsStuckInLoop(pos, board) {
					ct++
				}
				board[i][j] = '.'
			}
		}
	}
	return ct
}
