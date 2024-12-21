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

func PrintBoard(board [][]rune) {
	for _, row := range board {
		fmt.Println(string(row))
	}
}

func MoveInDirection(board [][]rune, x, y int, dir [2]int) ([][]rune, [2]int) {
	ogX, ogY := x, y
	stones := map[[2]int]rune{}
	for ; x >= 0 && x < len(board) && y >= 0 && y < len(board[0]); x, y = x+dir[0], y+dir[1] {
		if board[x][y] == '#' {
			// hit a wall
			return board, [2]int{ogX, ogY}
		} else if board[x][y] == '.' {
			// move all stones one down
			board[ogX][ogY] = '.'

			for stone, char := range stones {
				board[stone[0]+dir[0]][stone[1]+dir[1]] = char
			}
			return board, [2]int{ogX + dir[0], ogY + dir[1]}
		}
		stones[[2]int{x, y}] = board[x][y]
	}

	return board, [2]int{ogX, ogY}
}

func part1(input string) int {

	board := [][]rune{}
	lines := strings.Split(input, "\n\n")[0]
	movements := strings.ReplaceAll(strings.Split(input, "\n\n")[1], "\n", "")
	for _, line := range strings.Split(lines, "\n") {
		board = append(board, []rune(line))
	}

	playerPos := [2]int{}
	for i, row := range board {
		for j, cell := range row {
			if cell == '@' {
				playerPos = [2]int{i, j}
				break
			}
		}
	}

	for _, move := range movements {
		switch move {
		case '^':
			// move north
			board, playerPos = MoveInDirection(board, playerPos[0], playerPos[1], [2]int{-1, 0})
		case 'v':
			// move south
			board, playerPos = MoveInDirection(board, playerPos[0], playerPos[1], [2]int{1, 0})
		case '>':
			// move east
			board, playerPos = MoveInDirection(board, playerPos[0], playerPos[1], [2]int{0, 1})
		case '<':
			// move west
			board, playerPos = MoveInDirection(board, playerPos[0], playerPos[1], [2]int{0, -1})
		}
	}

	total := 0
	for i, row := range board {
		for j, cell := range row {
			if cell == 'O' {
				total += 100*i + j
			}
		}
	}

	return total

}

type Box struct {
	x, yL, yR int
}

func CheckNextBox(board [][]rune, b Box, dir [2]int) bool {
	if dir[1] == 1 {
		if board[b.x][b.yR+dir[1]] != '.' {
			return false
		}
	} else if dir[1] == -1 {
		if board[b.x][b.yL+dir[1]] != '.' {
			return false
		}
	} else {
		if board[b.x+dir[0]][b.yL] != '.' || board[b.x+dir[0]][b.yR] != '.' {
			return false
		}
	}
	return true
}

func MoveBox(board [][]rune, b Box, dir [2]int) ([][]rune, Box) {
	if dir[0] == 0 {
		if dir[1] == 1 {
			board[b.x][b.yR+dir[1]], board[b.x][b.yR] = board[b.x][b.yR], board[b.x][b.yR+dir[1]]
			board[b.x][b.yL+dir[1]], board[b.x][b.yL] = board[b.x][b.yL], board[b.x][b.yL+dir[1]]
		} else {
			board[b.x][b.yL+dir[1]], board[b.x][b.yL] = board[b.x][b.yL], board[b.x][b.yL+dir[1]]
			board[b.x][b.yR+dir[1]], board[b.x][b.yR] = board[b.x][b.yR], board[b.x][b.yR+dir[1]]
		}
		return board, Box{b.x + dir[0], b.yL + dir[1], b.yR + dir[1]}
	}
	board[b.x][b.yL], board[b.x+dir[0]][b.yL+dir[1]] = board[b.x+dir[0]][b.yL+dir[1]], board[b.x][b.yL]
	board[b.x][b.yR], board[b.x+dir[0]][b.yR+dir[1]] = board[b.x+dir[0]][b.yR+dir[1]], board[b.x][b.yR]
	return board, Box{b.x + dir[0], b.yL + dir[1], b.yR + dir[1]}
}

func MoveBoxes(board [][]rune, x, y int, dir [2]int) ([][]rune, [2]int) {
	ogX, ogY := x, y
	boxes := map[Box]bool{}

	tilesToCheck := [][2]int{{x + dir[0], y + dir[1]}}

	isWall := false
	for ; len(tilesToCheck) > 0; tilesToCheck = tilesToCheck[1:] {
		x, y := tilesToCheck[0][0], tilesToCheck[0][1]
		if x < 0 || x >= len(board) || y < 0 || y >= len(board[0]) {
			continue
		}
		if board[x][y] == '#' {
			isWall = true
			break
		}
		if dir[0] != 0 {
			if board[x][y] == '[' {
				boxes[Box{x, y, y + 1}] = true
				tilesToCheck = append(tilesToCheck, [2]int{x + dir[0], y + dir[1]}, [2]int{x + dir[0], y + dir[1] + 1})
			} else if board[x][y] == ']' {
				boxes[Box{x, y - 1, y}] = true
				tilesToCheck = append(tilesToCheck, [2]int{x + dir[0], y + dir[1]}, [2]int{x + dir[0], y + dir[1] - 1})
			}
		} else {
			if board[x][y] == '[' {
				boxes[Box{x, y, y + 1}] = true
				tilesToCheck = append(tilesToCheck, [2]int{x, y + dir[1]})
			} else if board[x][y] == ']' {
				boxes[Box{x, y - 1, y}] = true
				tilesToCheck = append(tilesToCheck, [2]int{x, y + dir[1]})
			}

		}
	}

	if isWall {
		return board, [2]int{ogX, ogY}
	} else if len(boxes) == 0 {
		board[ogX][ogY] = '.'
		board[ogX+dir[0]][ogY+dir[1]] = '@'
	} else {
		for len(boxes) > 0 {
			box := Box{}
			for b := range boxes {
				box = b
				break
			}
			if CheckNextBox(board, box, dir) {
				board, _ = MoveBox(board, box, dir)
				delete(boxes, box)
			}
		}
		board[ogX][ogY] = '.'
		board[ogX+dir[0]][ogY+dir[1]] = '@'
	}

	return board, [2]int{ogX + dir[0], ogY + dir[1]}
}

func part2(input string) int {
	board := [][]rune{}
	lines := strings.Split(input, "\n\n")[0]
	movements := strings.ReplaceAll(strings.Split(input, "\n\n")[1], "\n", "")
	for _, line := range strings.Split(lines, "\n") {
		newLine := []rune{}
		for _, char := range line {
			switch char {
			case '#':
				newLine = append(newLine, '#', '#')
			case '.':
				newLine = append(newLine, '.', '.')
			case 'O':
				newLine = append(newLine, '[', ']')
			case '@':
				newLine = append(newLine, '@', '.')
			}
		}
		board = append(board, newLine)
	}

	PrintBoard(board)

	playerPos := [2]int{}
	for i, row := range board {
		for j, cell := range row {
			if cell == '@' {
				playerPos = [2]int{i, j}
				break
			}
		}
	}

	for _, move := range movements {
		switch move {
		case '^':
			// move north
			board, playerPos = MoveBoxes(board, playerPos[0], playerPos[1], [2]int{-1, 0})
		case 'v':
			// move south
			board, playerPos = MoveBoxes(board, playerPos[0], playerPos[1], [2]int{1, 0})
		case '>':
			// move east
			board, playerPos = MoveBoxes(board, playerPos[0], playerPos[1], [2]int{0, 1})
		case '<':
			// move west
			board, playerPos = MoveBoxes(board, playerPos[0], playerPos[1], [2]int{0, -1})
		}
	}

	total := 0
	for i, row := range board {
		for j, cell := range row {
			if cell == '[' {
				total += 100*i + j
			}
		}
	}
	PrintBoard(board)
	return total
}
