package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var (
	R, C   int
	grid   []string
	totalM int
)

var dr = []int{-1, 0, 1, 0}
var dc = []int{0, 1, 0, -1}

func reflect(dir int, ch byte) int {
	if ch == '/' {
		switch dir {
		case 0:
			return 3
		case 1:
			return 2
		case 2:
			return 1
		case 3:
			return 0
		}
	} else if ch == '\\' {
		switch dir {
		case 0:
			return 1
		case 1:
			return 0
		case 2:
			return 3
		case 3:
			return 2
		}
	}
	return dir
}

func simulate(startR, startC, dir int) bool {
	size := R * C
	visitedMirror := make([]bool, size)
	seen := make([]bool, size*4)
	count := 0

	r := startR
	c := startC

	for {
		r += dr[dir]
		c += dc[dir]
		if r < 0 || r >= R || c < 0 || c >= C {
			break
		}
		idx := r*C + c

		key := idx*4 + dir
		if seen[key] {
			break
		}
		seen[key] = true

		ch := grid[r][c]
		if ch == '/' || ch == '\\' {
			if !visitedMirror[idx] {
				visitedMirror[idx] = true
				count++
			}
			dir = reflect(dir, ch)
		}
	}

	return count == totalM
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &R, &C)
	grid = make([]string, R)
	totalM = 0
	readRow := func() string {
		buf := make([]byte, 0, C)
		for len(buf) < C {
			ch, err := in.ReadByte()
			if err != nil {
				buf = append(buf, '.')
				continue
			}
			if ch == '\n' || ch == '\r' {
				continue
			}
			buf = append(buf, ch)
		}
		return string(buf)
	}
	for i := 0; i < R; i++ {
		grid[i] = readRow()
		for j := 0; j < C; j++ {
			if grid[i][j] == '/' || grid[i][j] == '\\' {
				totalM++
			}
		}
	}

	results := make([]string, 0)

	// North
	for c := 0; c < C; c++ {
		if simulate(-1, c, 2) {
			results = append(results, "N"+strconv.Itoa(c+1))
		}
	}
	// South
	for c := 0; c < C; c++ {
		if simulate(R, c, 0) {
			results = append(results, "S"+strconv.Itoa(c+1))
		}
	}
	// East
	for r := 0; r < R; r++ {
		if simulate(r, C, 3) {
			results = append(results, "E"+strconv.Itoa(r+1))
		}
	}
	// West
	for r := 0; r < R; r++ {
		if simulate(r, -1, 1) {
			results = append(results, "W"+strconv.Itoa(r+1))
		}
	}

	fmt.Println(len(results))
	if len(results) > 0 {
		for i, s := range results {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(s)
		}
		fmt.Println()
	}
}
