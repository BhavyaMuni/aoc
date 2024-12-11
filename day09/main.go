package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
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

func AppendCharNTimes(st []string, char string, n int) []string {
	for i := 0; i < n; i++ {
		st = append(st, char)
	}
	return st
}

func CalculateChecksum(fs []string) int {
	checksum := 0
	for i := 0; i < len(fs); i++ {
		if fs[i] == "." {
			continue
		}
		n, _ := strconv.Atoi(fs[i])
		checksum += i * n
	}
	return checksum
}

func part1(input string) int {
	fs := make([]string, 0)
	id := 0
	for ix, s := range input {
		size, _ := strconv.Atoi(string(s))
		if ix%2 == 0 {
			fs = AppendCharNTimes(fs, strconv.Itoa(id), size)
			id += 1
		} else {
			fs = AppendCharNTimes(fs, ".", size)
		}
	}

	l := 0
	r := len(fs) - 1
	for r > l {
		if fs[r] == "." {
			r--
			continue
		}
		if fs[l] != "." {
			l++
			continue
		}
		fs[l], fs[r] = fs[r], fs[l]
		l++
		r--
	}
	return CalculateChecksum(fs)
}

func part2(input string) int {
	fs := make([]string, 0)
	id := 0
	for ix, s := range input {
		size, _ := strconv.Atoi(string(s))
		if ix%2 == 0 {
			fs = AppendCharNTimes(fs, strconv.Itoa(id), size)
			id += 1
		} else {
			fs = AppendCharNTimes(fs, ".", size)
		}
	}

	freeSpaces := map[int]int{}
	curr := 0
	j := 0
	for i := 0; i < len(fs); i++ {
		if fs[i] == "." {
			curr++
		} else {
			if curr > 0 {
				freeSpaces[j] = curr
			}
			curr = 0
			j = i + 1
		}
	}

	r := len(fs) - 1
	l := r
	for r > 0 {
		if fs[r] == "." {
			r--
			continue
		}
		l = r - 1
		for l >= 0 && fs[l] == fs[r] {
			l--
		}

		ln := r - l

		for i := 0; i < r; i++ {
			if indexes := freeSpaces[i]; indexes >= ln {
				temp := append([]string{}, fs[i:i+ln]...)
				copy(fs[i:i+ln], fs[l+1:r+1])
				copy(fs[l+1:r+1], temp)
				delete(freeSpaces, i)
				if indexes-ln > 0 {
					freeSpaces[i+ln] = indexes - ln
				}
				break
			}
		}

		r = l
	}

	return CalculateChecksum(fs)
}
