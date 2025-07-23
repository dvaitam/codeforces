package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	pos := map[byte][2]int{
		'1': {0, 0}, '2': {1, 0}, '3': {2, 0},
		'4': {0, 1}, '5': {1, 1}, '6': {2, 1},
		'7': {0, 2}, '8': {1, 2}, '9': {2, 2},
		'0': {1, 3},
	}

	board := [4][3]bool{}
	for d, p := range pos {
		board[p[1]][p[0]] = true
		_ = d
	}

	moves := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		a := pos[s[i-1]]
		b := pos[s[i]]
		moves[i-1] = [2]int{b[0] - a[0], b[1] - a[1]}
	}

	for d := byte('0'); d <= '9'; d++ {
		start, ok := pos[d]
		if !ok {
			continue
		}
		x, y := start[0], start[1]
		valid := true
		for _, mv := range moves {
			x += mv[0]
			y += mv[1]
			if y < 0 || y >= 4 || x < 0 || x >= 3 || !board[y][x] {
				valid = false
				break
			}
		}
		if valid && d != s[0] {
			fmt.Fprintln(out, "NO")
			return
		}
	}
	fmt.Fprintln(out, "YES")
}
