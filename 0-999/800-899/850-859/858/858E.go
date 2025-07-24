package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Move struct {
	from string
	to   string
}

func parseNumber(s string, n int) (int, bool) {
	if len(s) == 0 {
		return 0, false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, false
		}
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	if v < 1 || v > n || strconv.Itoa(v) != s {
		return 0, false
	}
	return v, true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	name := make([]string, n)
	typeFile := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &name[i], &typeFile[i])
	}
	e := 0
	for _, t := range typeFile {
		if t == 1 {
			e++
		}
	}

	used := make(map[string]int)
	for i, nm := range name {
		used[nm] = i
	}

	numUsed := make([]bool, n+1)
	for _, nm := range name {
		if v, ok := parseNumber(nm, n); ok {
			numUsed[v] = true
		}
	}

	exAvail := []int{}
	regAvail := []int{}
	for i := 1; i <= e; i++ {
		if !numUsed[i] {
			exAvail = append(exAvail, i)
		}
	}
	for i := e + 1; i <= n; i++ {
		if !numUsed[i] {
			regAvail = append(regAvail, i)
		}
	}

	exWrongNum := []int{}
	regWrongNum := []int{}
	exWrongOther := []int{}
	regWrongOther := []int{}
	for i := 0; i < n; i++ {
		if v, ok := parseNumber(name[i], n); ok {
			if v <= e {
				if typeFile[i] == 1 {
					// correct
				} else {
					regWrongNum = append(regWrongNum, i)
				}
			} else {
				if typeFile[i] == 0 {
					// correct
				} else {
					exWrongNum = append(exWrongNum, i)
				}
			}
		} else {
			if typeFile[i] == 1 {
				exWrongOther = append(exWrongOther, i)
			} else {
				regWrongOther = append(regWrongOther, i)
			}
		}
	}

	moves := []Move{}
	pop := func(a *[]int) int {
		v := (*a)[len(*a)-1]
		*a = (*a)[:len(*a)-1]
		return v
	}

	moveFile := func(i int, newName string) {
		moves = append(moves, Move{from: name[i], to: newName})
		delete(used, name[i])
		used[newName] = i
		if v, ok := parseNumber(name[i], n); ok {
			numUsed[v] = false
			if v <= e {
				exAvail = append(exAvail, v)
			} else {
				regAvail = append(regAvail, v)
			}
		}
		name[i] = newName
		if v, ok := parseNumber(newName, n); ok {
			numUsed[v] = true
		}
	}

	findTemp := func() string {
		cand := "tmp"
		if _, ok := used[cand]; !ok {
			return cand
		}
		for i := 0; ; i++ {
			c := fmt.Sprintf("tmp%d", i)
			if len(c) > 6 {
				c = strings.Repeat("a", 6)
			}
			if _, ok := used[c]; !ok {
				return c
			}
		}
	}

	temp := findTemp()

	for len(exWrongNum) > 0 || len(regWrongNum) > 0 {
		if len(exAvail) == 0 && len(regAvail) == 0 {
			var idx int
			if len(exWrongNum) > 0 {
				idx = exWrongNum[len(exWrongNum)-1]
				exWrongNum = exWrongNum[:len(exWrongNum)-1]
			} else {
				idx = regWrongNum[len(regWrongNum)-1]
				regWrongNum = regWrongNum[:len(regWrongNum)-1]
			}
			oldVal, _ := parseNumber(name[idx], n)
			moveFile(idx, temp)
			if oldVal <= e {
				exAvail = append(exAvail, oldVal)
			} else {
				regAvail = append(regAvail, oldVal)
			}
			if typeFile[idx] == 1 {
				exWrongOther = append(exWrongOther, idx)
			} else {
				regWrongOther = append(regWrongOther, idx)
			}
			continue
		}

		for len(exWrongNum) > 0 && len(exAvail) > 0 {
			idx := exWrongNum[len(exWrongNum)-1]
			exWrongNum = exWrongNum[:len(exWrongNum)-1]
			val := pop(&exAvail)
			oldVal, _ := parseNumber(name[idx], n)
			moveFile(idx, strconv.Itoa(val))
			if oldVal <= e {
				exAvail = append(exAvail, oldVal)
			} else {
				regAvail = append(regAvail, oldVal)
			}
		}
		for len(regWrongNum) > 0 && len(regAvail) > 0 {
			idx := regWrongNum[len(regWrongNum)-1]
			regWrongNum = regWrongNum[:len(regWrongNum)-1]
			val := pop(&regAvail)
			oldVal, _ := parseNumber(name[idx], n)
			moveFile(idx, strconv.Itoa(val))
			if oldVal <= e {
				exAvail = append(exAvail, oldVal)
			} else {
				regAvail = append(regAvail, oldVal)
			}
		}
	}

	for len(exWrongOther) > 0 {
		idx := exWrongOther[len(exWrongOther)-1]
		exWrongOther = exWrongOther[:len(exWrongOther)-1]
		val := pop(&exAvail)
		moveFile(idx, strconv.Itoa(val))
	}
	for len(regWrongOther) > 0 {
		idx := regWrongOther[len(regWrongOther)-1]
		regWrongOther = regWrongOther[:len(regWrongOther)-1]
		val := pop(&regAvail)
		moveFile(idx, strconv.Itoa(val))
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, len(moves))
	for _, mv := range moves {
		fmt.Fprintf(writer, "move %s %s\n", mv.from, mv.to)
	}
}
