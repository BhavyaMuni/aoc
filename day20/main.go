package main

import (
	_ "embed"
	"flag"
	"fmt"
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

func minPath(board [][]rune, start [2]int) int {
	directions := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	queue := [][3]int{{start[0], start[1], 0}}
	visited := make(map[[2]int]int)
	// steps := 0
	for len(queue) > 0 {
		for i := 0; i < len(queue); i++ {
			f := queue[0]
			queue = queue[1:]
			curr := [2]int{f[0], f[1]}
			dist := f[2]
			if curr[0] < 0 || curr[0] >= len(board) || curr[1] < 0 || curr[1] >= len(board[0]) {
				continue
			}
			if board[curr[0]][curr[1]] == '#' {
				continue
			}
			if _, ok := visited[curr]; ok {
				continue
			}
			if board[curr[0]][curr[1]] == 'E' {
				return dist
			}
			visited[curr] = dist
			for _, dir := range directions {
				newPos := [3]int{curr[0] + dir[0], curr[1] + dir[1], dist + 1}
				queue = append(queue, newPos)
			}
		}
	}
	return -1
}
func part1(input string) int {
	board := [][]rune{}
	for _, line := range strings.Split(input, "\n") {
		board = append(board, []rune(line))
	}

	start := [2]int{3, 1}
	for i, row := range board {
		for j, cell := range row {
			if cell == 'S' {
				start = [2]int{i, j}
				break
			}
		}
	}
	counter := 0
	basePath := minPath(board, start)
	for i, row := range board {
		for j, cell := range row {
			if cell == '#' {
				a := [2]int{i, j}
				board[a[0]][a[1]] = '.'
				m := minPath(board, start)
				if m > 0 && basePath-m >= 100 {
					counter++
				}
				board[a[0]][a[1]] = '#'
			}
		}
	}

	return counter
}

func part2(input string) int {
	return 0
}
