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

type Robot struct {
	pos [2]int
	vel [2]int
}

const BOARD_LEN = 101
const BOARD_WID = 103

func part1(input string) int {
	robots := make([]Robot, 0)

	for _, line := range strings.Split(input, "\n") {
		var r Robot
		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &r.pos[0], &r.pos[1], &r.vel[0], &r.vel[1])
		robots = append(robots, r)
	}

	halfX := BOARD_LEN / 2
	halfY := BOARD_WID / 2
	quads := make([]int, 4)
	// move all robots to 100 seconds
	for r, robot := range robots {
		robots[r].pos[0] += 100 * robot.vel[0]
		robots[r].pos[1] += 100 * robot.vel[1]

		robots[r].pos[0] %= BOARD_LEN
		robots[r].pos[1] %= BOARD_WID

		if robots[r].pos[0] < 0 {
			robots[r].pos[0] += BOARD_LEN
		}
		if robots[r].pos[1] < 0 {
			robots[r].pos[1] += BOARD_WID
		}

		switch {
		case robots[r].pos[0] == halfX && robots[r].pos[1] == halfY:
			continue
		case robots[r].pos[0] < halfX && robots[r].pos[1] < halfY:
			quads[0]++
		case robots[r].pos[0] > halfX && robots[r].pos[1] < halfY:
			quads[1]++
		case robots[r].pos[0] < halfX && robots[r].pos[1] > halfY:
			quads[2]++
		case robots[r].pos[0] > halfX && robots[r].pos[1] > halfY:
			quads[3]++
		}
	}

	return quads[0] * quads[1] * quads[2] * quads[3]
}

func PrintBoard(robots []Robot) {
	board := make([][]byte, BOARD_LEN)
	for i := 0; i < BOARD_LEN; i++ {
		board[i] = make([]byte, BOARD_WID)
		for j := 0; j < BOARD_WID; j++ {
			board[i][j] = '.'
		}
	}

	for _, robot := range robots {
		board[robot.pos[0]][robot.pos[1]] = '#'
	}

	for i := 0; i < BOARD_LEN; i++ {
		fmt.Println(string(board[i]))
	}
}

const DIST = 20

func DensityAroundRobot(robots []Robot, robot Robot) int {
	count := 0
	for _, r := range robots {
		if r == robot {
			continue
		}
		if math.Abs(float64(r.pos[0]-robot.pos[0])) <= DIST && math.Abs(float64(r.pos[1]-robot.pos[1])) <= DIST {
			count++
		}
	}
	return count
}

func part2(input string) int {
	robots := make([]Robot, 0)

	for _, line := range strings.Split(input, "\n") {
		var r Robot
		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &r.pos[0], &r.pos[1], &r.vel[0], &r.vel[1])
		robots = append(robots, r)
	}

	// time := 16128
	max, maxT := 0, 0

	for i := 0; i < 1; i++ {
		// move all robots to 100 seconds
		for r, robot := range robots {
			robots[r].pos[0] += 8270 * robot.vel[0]
			robots[r].pos[1] += 8270 * robot.vel[1]

			robots[r].pos[0] %= BOARD_LEN
			robots[r].pos[1] %= BOARD_WID

			if robots[r].pos[0] < 0 {
				robots[r].pos[0] += BOARD_LEN
			}
			if robots[r].pos[1] < 0 {
				robots[r].pos[1] += BOARD_WID
			}

			// // check if robot surrounging is dense
			// d := DensityAroundRobot(robots, robots[r])
			// if d > max {
			// 	max = d
			// 	maxT = i
			// }
		}
		PrintBoard(robots)
	}
	fmt.Println(max, maxT)
	return 0
}
