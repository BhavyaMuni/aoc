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

const BOARD_SIZE = 71
const NUM_BYTES = 1024

func PrintBoard(board [][]rune) {
	for _, line := range board {
		fmt.Println(string(line))
	}
}

func FindExit(board [][]rune) (steps int, possible bool) {
	start := [3]int{0, 0, 0}
	queue := [][3]int{start}
	visited := make(map[[3]int]bool)
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if visited[curr] {
			continue
		}
		visited[curr] = true
		i, j, currSteps := curr[0], curr[1], curr[2]
		if i < 0 || i >= BOARD_SIZE || j < 0 || j >= BOARD_SIZE || board[j][i] == '#' {
			continue
		}
		if i == BOARD_SIZE-1 && j == BOARD_SIZE-1 {
			steps = currSteps
			possible = true
			return
		}
		queue = append(queue, [3]int{i + 1, j, steps + 1})
		queue = append(queue, [3]int{i - 1, j, steps + 1})
		queue = append(queue, [3]int{i, j + 1, steps + 1})
		queue = append(queue, [3]int{i, j - 1, steps + 1})
	}
	steps = -1
	possible = false
	return
}

func part1(input string) int {
	board := make([][]rune, BOARD_SIZE, BOARD_SIZE)
	for i := 0; i < BOARD_SIZE; i++ {
		board[i] = []rune(strings.Repeat(".", BOARD_SIZE))
	}
	for n, line := range strings.Split(input, "\n") {
		if n >= NUM_BYTES {
			break
		}
		var i, j int
		fmt.Sscanf(line, "%d,%d", &i, &j)
		board[j][i] = '#'
	}
	PrintBoard(board)
	steps, poss := FindExit(board)
	if poss {
		return steps
	}
	return -1
}

func part2(input string) int {
	board := make([][]rune, BOARD_SIZE, BOARD_SIZE)
	for i := 0; i < BOARD_SIZE; i++ {
		board[i] = []rune(strings.Repeat(".", BOARD_SIZE))
	}

	for newNumBytes := 1; newNumBytes <= 3450; newNumBytes++ {
		for n, line := range strings.Split(input, "\n") {
			if n >= newNumBytes {
				break
			}
			var i, j int
			fmt.Sscanf(line, "%d,%d", &i, &j)
			board[j][i] = '#'
		}
		_, poss := FindExit(board)
		if !poss {
			fmt.Println("Not possible with line number: ", newNumBytes)
			return newNumBytes
		}
	}
	return 0
}
