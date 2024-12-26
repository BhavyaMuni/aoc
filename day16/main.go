package main

import (
	"container/heap"
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

type Node1 struct {
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

	startNode := Node1{startPos[0], startPos[1], 0, 1, 0, [][2]int{}, map[[2]int]bool{{startPos[0], startPos[1]}: true}}
	queue := []Node1{startNode}
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
			newNode := Node1{newX, newY, dir[0], dir[1], pts + 1, [][2]int{}, map[[2]int]bool{{newX, newY}: true}}
			if dir != [2]int{dx, dy} {
				newNode.pts += 1000
			}

			queue = append(queue, newNode)
		}
	}

	return min_pts
}

type Position struct {
	x, y, dx, dy int
}

type Path struct {
	pts  int
	pos  Position
	path [][2]int
}

type PriorityQueue []Path

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].pts < pq[j].pts
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Path))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
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

	startNode := Position{startPos[0], startPos[1], 0, 1}
	queue := &PriorityQueue{}
	heap.Init(queue)
	min_pts := math.MaxInt64
	seen := make(map[Position]int)
	heap.Push(queue, Path{0, startNode, [][2]int{startPos}})
	for queue.Len() > 0 {
		state := heap.Pop(queue).(Path)
		pos := state.pos

		if state.pts > min_pts {
			break
		}
		if board[pos.x][pos.y] == 'E' {
			for _, p := range state.path {
				counter[[2]int{p[0], p[1]}] = struct{}{}
			}
			continue
		}
		if s, ok := seen[pos]; ok && s < state.pts {
			continue
		}
		seen[pos] = state.pts
		x, y, dx, dy := pos.x, pos.y, pos.dx, pos.dy
		for _, dir := range dirs {
			newX, newY := x+dir[0], y+dir[1]
			if newX < 0 || newX >= len(board) || newY < 0 || newY >= len(board[0]) {
				continue
			}
			if board[newX][newY] == '#' {
				continue
			}
			newPath := [][2]int{}
			newPath = append(newPath, state.path...)
			if dir != [2]int{dx, dy} {
				heap.Push(queue, Path{state.pts + 1001, Position{newX, newY, dir[0], dir[1]}, append(newPath, [2]int{newX, newY})})
			} else {
				heap.Push(queue, Path{state.pts + 1, Position{newX, newY, dir[0], dir[1]}, append(newPath, [2]int{newX, newY})})
			}
		}
	}

	return len(counter)
}
