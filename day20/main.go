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

func get_dists(board [][]rune, start [2]int) map[[2]int]int {
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
			visited[curr] = dist
			for _, dir := range directions {
				newPos := [3]int{curr[0] + dir[0], curr[1] + dir[1], dist + 1}
				queue = append(queue, newPos)
			}
		}
	}
	return visited
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
	dists := get_dists(board, start)
	// baseDist = dists[start]
	saves := make(map[int]int)
	dirs := [][2]int{{0, 2}, {0, -2}, {2, 0}, {-2, 0}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for pos := range dists {
		for _, dir := range dirs {
			newPos := [2]int{pos[0] + dir[0], pos[1] + dir[1]}
			if newPos[0] < 0 || newPos[0] >= len(board) || newPos[1] < 0 || newPos[1] >= len(board[0]) {
				continue
			}
			if board[newPos[0]][newPos[1]] == '#' {
				continue
			}
			if dists[newPos]-dists[pos]-2 > 0 {
				saves[dists[newPos]-dists[pos]-2]++
			}
		}
	}
	for k, v := range saves {
		if k >= 100 {
			counter += v
		}
	}
	return counter
}

func GenerateDirections(dist int) [][2]int {
	dirs := make(map[[2]int]struct{})
	for j := 1; j <= dist; j++ {
		for i := 0; i <= j; i++ {
			inv := j - i
			dirs[[2]int{i, inv}] = struct{}{}
			dirs[[2]int{inv, -i}] = struct{}{}
			dirs[[2]int{-i, -inv}] = struct{}{}
			dirs[[2]int{-inv, i}] = struct{}{}
		}
	}
	out := make([][2]int, 0, len(dirs))
	for k := range dirs {
		out = append(out, k)
	}
	return out
}

func GetManhattanDist(a, b [2]int) int {
	return int(math.Abs(float64(a[0]-b[0])) + math.Abs(float64(a[1]-b[1])))
}

func part2(input string) int {
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
	dists := get_dists(board, start)
	// baseDist = dists[start]
	saves := make(map[int]int)
	dirs := GenerateDirections(20)
	cheats := map[[2][2]int]struct{}{}
	for pos := range dists {
		for _, dir := range dirs {
			newPos := [2]int{pos[0] + dir[0], pos[1] + dir[1]}
			if newPos[0] < 0 || newPos[0] >= len(board) || newPos[1] < 0 || newPos[1] >= len(board[0]) {
				continue
			}
			if board[newPos[0]][newPos[1]] == '#' {
				continue
			}
			mDist := GetManhattanDist(pos, newPos)
			if dists[newPos]-dists[pos]-mDist >= 100 {
				saves[dists[newPos]-dists[pos]-mDist]++
				cheats[[2][2]int{pos, newPos}] = struct{}{}
			}
		}
	}
	for _, v := range saves {
		counter += v
	}
	return counter
}
