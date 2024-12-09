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

func getAntinodes(c1, c2 [2]int, boardX, boardY int) (an1, an2 [2]int) {
	xDiff1, xDiff2 := c1[0]-c2[0], c2[0]-c1[0]
	yDiff1, yDiff2 := c1[1]-c2[1], c2[1]-c1[1]

	an1 = [2]int{c1[0] + xDiff1, c1[1] + yDiff1}
	an2 = [2]int{c2[0] + xDiff2, c2[1] + yDiff2}

	//check if an1 and an2 are in the board
	if an1[0] < 0 || an1[0] >= boardX || an1[1] < 0 || an1[1] >= boardY {
		an1 = [2]int{}
	}
	if an2[0] < 0 || an2[0] >= boardX || an2[1] < 0 || an2[1] >= boardY {
		an2 = [2]int{}
	}

	return
}

func getMultipleAntinodes(c1, c2 [2]int, boardX, boardY int) [][2]int {
	var ans [][2]int

	xDiff1, xDiff2 := c1[0]-c2[0], c2[0]-c1[0]
	yDiff1, yDiff2 := c1[1]-c2[1], c2[1]-c1[1]

	for i := 1; ; i++ {
		an1 := [2]int{c1[0] + xDiff1*i, c1[1] + yDiff1*i}

		if an1[0] < 0 || an1[0] >= boardX || an1[1] < 0 || an1[1] >= boardY {
			break
		}
		ans = append(ans, an1)
	}

	for i := 1; ; i++ {
		an2 := [2]int{c2[0] + xDiff2*i, c2[1] + yDiff2*i}

		if an2[0] < 0 || an2[0] >= boardX || an2[1] < 0 || an2[1] >= boardY {
			break
		}
		ans = append(ans, an2)
	}

	return ans
}

func part1(input string) int {
	var board [][]rune

	for _, line := range strings.Split(input, "\n") {
		board = append(board, []rune(line))
	}

	coords := make(map[rune][][2]int)

	for x, row := range board {
		for y, cell := range row {
			if cell != '.' {
				coords[cell] = append(coords[cell], [2]int{x, y})
			}
		}
	}

	out := map[[2]int]string{}

	for _, v := range coords {
		// get every pair of coords
		for i := 0; i < len(v); i++ {
			for j := i + 1; j < len(v); j++ {
				an1, an2 := getAntinodes(v[i], v[j], len(board), len(board[0]))
				if an1 != [2]int{} {
					out[an1] = "#"
				}
				if an2 != [2]int{} {
					out[an2] = "#"
				}
			}
		}
	}

	return len(out)
}

func part2(input string) int {
	var board [][]rune

	for _, line := range strings.Split(input, "\n") {
		board = append(board, []rune(line))
	}

	coords := make(map[rune][][2]int)

	for x, row := range board {
		for y, cell := range row {
			if cell != '.' {
				coords[cell] = append(coords[cell], [2]int{x, y})
			}
		}
	}

	out := map[[2]int]string{}

	for _, v := range coords {
		// get every pair of coords
		for i := 0; i < len(v); i++ {
			for j := i + 1; j < len(v); j++ {
				ans := getMultipleAntinodes(v[i], v[j], len(board), len(board[0]))
				for _, an := range ans {
					if board[an[0]][an[1]] == '.' {
						out[an] = "#"
						board[an[0]][an[1]] = '#'
					}
				}
			}
		}
	}
	for _, row := range board {
		fmt.Println(string(row))
	}

	numAntennas := 0
	for _, row := range board {
		for _, cell := range row {
			if cell != '#' && cell != '.' {
				numAntennas++
			}
		}
	}
	return len(out) + numAntennas
}
