package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
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

func part1(input string) int {
	rulesList := strings.Split(strings.Split(input, "\n\n")[0], "\n")
	pagesList := strings.Split(strings.Split(input, "\n\n")[1], "\n")

	rules := make(map[int][]int)
	for _, rule := range rulesList {
		var key, val int
		fmt.Sscanf(rule, "%d|%d", &key, &val)
		rules[key] = append(rules[key], val)
	}

	ct := 0
	nu := 0
	for _, line := range pagesList {
		found := true
		seen := make([]int, 0)
		pages := strings.Split(line, ",")
		for i := 0; i < len(pages); i++ {
			page, _ := strconv.Atoi(pages[i])
			if _, ok := rules[page]; ok {
				for _, rule := range rules[page] {
					if slices.Contains(seen, rule) {
						found = false
						break
					}
				}
			}
			seen = append(seen, page)
		}
		if found {
			mid, _ := strconv.Atoi(pages[len(pages)/2])
			ct += mid
			nu++
		}
	}
	fmt.Println("NU", nu)

	return ct
}

func part2(input string) int {
	rulesList := strings.Split(strings.Split(input, "\n\n")[0], "\n")
	pagesList := strings.Split(strings.Split(input, "\n\n")[1], "\n")

	rules := make(map[int][]int)
	for _, rule := range rulesList {
		var key, val int
		fmt.Sscanf(rule, "%d|%d", &key, &val)
		rules[key] = append(rules[key], val)
	}

	ct := 0
	nu := 0

	for _, line := range pagesList {

		pagesStr := strings.Split(line, ",")
		pages := make([]int, len(pagesStr))
		for i, line := range pagesStr {
			pages[i], _ = strconv.Atoi(line)
		}

		found := true
		for i := 0; i < len(pages); i++ {
			page := pages[i]

			for j := 0; j < len(rules[page]); j++ {
				for _, val := range rules[page] {
					// rule := rules[pages[i]][j]
					ruleIndex := slices.Index(pages, val)
					pageIndex := slices.Index(pages, page)

					if ruleIndex < pageIndex && ruleIndex > -1 && pageIndex > -1 {
						found = false
						pages[pageIndex], pages[ruleIndex] = val, page
						j = 0
					}

				}
			}
		}
		if !found {
			ct += pages[len(pages)/2]
			nu++
		}

	}
	return ct
}
