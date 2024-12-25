package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"reflect"
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

type Node struct {
	x, y, dx, dy, pts int
	path              [][2]int
	seen              map[[2]int]bool
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

	startNode := Node{startPos[0], startPos[1], 0, 1, 0, [][2]int{}, map[[2]int]bool{{startPos[0], startPos[1]}: true}}
	queue := []Node{startNode}
	min_pts := math.MaxInt64
	for len(queue) > 0 {
		// fmt.Println(queue)
		// pop the minimum point from the queue
		m := math.MaxInt64
		mi := 0
		for i := 0; i < len(queue); i++ {
			if queue[i].pts < m {
				m = queue[i].pts
				mi = i
			}
		}

		x, y, dx, dy, pts := queue[mi].x, queue[mi].y, queue[mi].dx, queue[mi].dy, queue[mi].pts
		queue = append(queue[:mi], queue[mi+1:]...)

		if x < 0 || x >= len(board) || y < 0 || y >= len(board[0]) {
			continue
		}
		if board[x][y] == 'E' {
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
			if board[newX][newY] == '#' || seen[[4]int{newX, newY}] {
				continue
			}
			if board[newX][newY] != 'E' {
				seen[[4]int{newX, newY}] = true
			}
			prev[[2]int{newX, newY}] = [2]int{x, y}
			newNode := Node{newX, newY, dir[0], dir[1], pts + 1, [][2]int{}, map[[2]int]bool{{newX, newY}: true}}
			if dir != [2]int{dx, dy} {
				newNode.pts += 1000
			}

			queue = append(queue, newNode)
		}
	}

	return min_pts
}

func checkIfPathExists(path [][2]int, seen [][][2]int) bool {
	for i := 0; i < len(seen); i++ {
		if reflect.DeepEqual(seen[i], path) {
			return true
		}
	}
	return false
}

func appendPathToCounter(path [][2]int, ct map[[2]int]struct{}) map[[2]int]struct{} {
	for i := 0; i < len(path); i++ {
		ct[[2]int{path[i][0], path[i][1]}] = struct{}{}
	}
	return ct
}

func isNotVisited(x, y int, seen map[[2]int]bool) bool {
	_, ok := seen[[2]int{x, y}]
	return !ok
}

func part2(input string) int {

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

	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	counter := make(map[[2]int]struct{})

	startNode := Node{startPos[0], startPos[1], 0, 1, 0, [][2]int{{startPos[0], startPos[1]}}, map[[2]int]bool{{startPos[0], startPos[1]}: true}}
	queue := [][]Node{{startNode}}
	min_pts := math.MaxInt64
	for len(queue) > 0 {
		// fmt.Println(queue)
		// pop the minimum point from the queue
		m := math.MaxInt64
		mi := 0
		for i := 0; i < len(queue); i++ {
			if queue[i][len(queue[i])-1].pts < m {
				m = queue[i][len(queue[i])-1].pts
				mi = i
			}
		}
		selectedPath := queue[mi]
		lastNode := selectedPath[len(selectedPath)-1]
		x, y, dx, dy, pts, currPath, curSeen := lastNode.x, lastNode.y, lastNode.dx, lastNode.dy, lastNode.pts, lastNode.path, lastNode.seen
		queue = append(queue[:mi], queue[mi+1:]...)

		if x < 0 || x >= len(board) || y < 0 || y >= len(board[0]) {
			continue
		}
		if board[x][y] == 'E' {
			if pts > min_pts {
				break
			}
			fmt.Println("Found path with pts: ", pts)
			min_pts = min(min_pts, pts)
			continue
		}
		fmt.Println("x, y, pts", x, y, pts)
		for _, dir := range dirs {
			newX, newY := x+dir[0], y+dir[1]
			newPath := [][2]int{}
			newPath = append(newPath, currPath...)
			newPath = append(newPath, [2]int{newX, newY})
			newSeen := make(map[[2]int]bool)
			for k, v := range curSeen {
				newSeen[k] = v
			}
			if dir[0]+dx == 0 && dir[1]+dy == 0 {
				continue
			}
			if newX < 0 || newX >= len(board) || newY < 0 || newY >= len(board[0]) {
				continue
			}
			if board[newX][newY] == '#' || curSeen[[2]int{newX, newY}] {
				continue
			}
			if board[newX][newY] != 'E' {
				newSeen[[2]int{newX, newY}] = true
			}
			newNode := Node{newX, newY, dir[0], dir[1], pts + 1, newPath, curSeen}
			if dir != [2]int{dx, dy} {
				newNode.pts += 1000
			}
			nPath := []Node{}
			nPath = append(nPath, selectedPath...)
			nPath = append(nPath, newNode)
			queue = append(queue, nPath)
		}
	}

	return len(counter)
}
