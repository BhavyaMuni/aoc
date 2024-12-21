package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
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

func getPath(prev map[[2]int][2]int, start [2]int, end [2]int) [][2]int {
	path := [][2]int{}
	curr := end
	for _, ok := prev[curr]; ok; _, ok = prev[curr] {
		path = append(path, curr)
		curr = prev[curr]
	}
	return path
}

func part1(input string) int {

	board := [][]rune{}

	for _, line := range strings.Split(input, "\n") {
		board = append(board, []rune(line))
	}

	var startPos [2]int
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			if board[i][j] == 'S' {
				startPos = [2]int{i, j}
				break
			}
		}
	}
	seen := make(map[[4]int]bool)
	prev := make(map[[2]int][2]int)
	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	// rotation := [2]int{0, 1}

	queue := [][5]int{{startPos[0], startPos[1], 0, 1, 0}}
	min_pts := math.MaxInt64
	for len(queue) > 0 {
		// fmt.Println(queue)
		// pop the minimum point from the queue
		m := math.MaxInt64
		mi := 0
		for i := 0; i < len(queue); i++ {
			if queue[i][4] < m {
				m = queue[i][4]
				mi = i
			}
		}

		x, y, dx, dy, pts := queue[mi][0], queue[mi][1], queue[mi][2], queue[mi][3], queue[mi][4]
		queue = append(queue[:mi], queue[mi+1:]...)

		if x < 0 || x >= len(board) || y < 0 || y >= len(board[0]) {
			continue
		}
		if board[x][y] == 'E' {
			fmt.Println(pts)
			fmt.Println(getPath(prev, startPos, [2]int{x, y}))
			min_pts = min(min_pts, pts)
			continue
		}

		for _, dir := range dirs {
			newX, newY := x+dir[0], y+dir[1]
			if dir[0]+dx == 0 && dir[1]+dy == 0 {
				continue
			}
			if newX < 0 || newX >= len(board) || newY < 0 || newY >= len(board[0]) {
				continue
			}
			if board[newX][newY] == '#' || seen[[4]int{newX, newY, dx, dy}] {
				continue
			}
			if board[newX][newY] != 'E' {
				seen[[4]int{newX, newY, dx, dy}] = true
			}
			prev[[2]int{newX, newY}] = [2]int{x, y}
			if dir != [2]int{dx, dy} {
				queue = append(queue, [5]int{newX, newY, dir[0], dir[1], pts + 1001})
			} else {
				queue = append(queue, [5]int{newX, newY, dir[0], dir[1], pts + 1})
			}
		}
	}

	return min_pts
}

func part2(input string) int {
	panic("not implemented")
}
