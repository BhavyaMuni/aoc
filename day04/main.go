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

var directions = [][]int{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

func checkSurrounding(board [][]rune, row int, col int) int {
	//for each direction check if M is first one A is second one and S is third one
	toCheck := []rune{'M', 'A', 'S'}
	numFound := 0
	for _, dir := range directions {
		found := true
		for i, letterToCheck := range toCheck {
			newRow := row + dir[0]*(i+1)
			newCol := col + dir[1]*(i+1)
			if newRow >= 0 && newRow < len(board) && newCol >= 0 && newCol < len(board[0]) {
				if board[newRow][newCol] != letterToCheck {
					found = false
				}
			} else {
				found = false
			}
		}
		if found {
			numFound++
		}
	}
	return numFound
}

func checkDiagonal(board [][]rune, row int, col int) bool {
	diagonals := [][]int{
		{-1, -1},
		{-1, 1},
		{1, -1},
		{1, 1},
	}

	// 2 diagonals
	counts := map[int]string{
		1: "",
		2: "",
	}

	for _, dir := range diagonals {
		newRow := row + dir[0]
		newCol := col + dir[1]
		if newRow >= 0 && newRow < len(board) && newCol >= 0 && newCol < len(board[0]) {
			if board[newRow][newCol] == 'S' || board[newRow][newCol] == 'M' {
				if dir[0]*dir[1] == 1 {
					counts[1] += string(board[newRow][newCol])
				} else {
					counts[2] += string(board[newRow][newCol])
				}
			}
		}
	}
	valid := (counts[1] == "SM" || counts[1] == "MS") && (counts[2] == "SM" || counts[2] == "MS")
	return valid
}

func part1(input string) int {
	var board [][]rune

	ct := 0
	for _, line := range strings.Split(input, "\n") {
		row := []rune(line)
		board = append(board, row)
	}

	for row_ix, row := range board {
		for cell_ix, cell := range row {

			if cell == 'X' {
				ct += checkSurrounding(board, row_ix, cell_ix)
			}
		}
	}
	return ct
}

func part2(input string) int {
	var board [][]rune

	ct := 0
	for _, line := range strings.Split(input, "\n") {
		row := []rune(line)
		board = append(board, row)
	}

	for row_ix, row := range board {
		for cell_ix, cell := range row {
			if cell == 'A' && checkDiagonal(board, row_ix, cell_ix) {
				ct++
			}
		}
	}

	return ct
}
