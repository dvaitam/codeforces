package main

import (
	"bufio"
	"fmt"
	"os"
)

func canEscape(board [][]byte, n int, start int) bool {
	extra := n + 6
	for i := 0; i < 3; i++ {
		board[i] = append(board[i], make([]byte, extra)...)
		for j := n; j < n+extra; j++ {
			board[i][j] = '.'
		}
	}
	type node struct{ r, c int }
	q := []node{{start, 0}}
	visited := make([][]bool, 3)
	for i := 0; i < 3; i++ {
		visited[i] = make([]bool, n+extra+1)
	}
	visited[start][0] = true
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		r, c := cur.r, cur.c
		if c >= n {
			return true
		}
		if board[r][c] != '.' {
			continue
		}
		for dr := -1; dr <= 1; dr++ {
			nr := r + dr
			if nr < 0 || nr >= 3 {
				continue
			}
			if board[r][c+1] != '.' {
				continue
			}
			if board[nr][c+1] != '.' {
				continue
			}
			if board[nr][c+2] != '.' {
				continue
			}
			if board[nr][c+3] != '.' {
				continue
			}
			nc := c + 3
			if nc >= n+extra {
				nc = n + extra - 1
			}
			if !visited[nr][nc] {
				visited[nr][nc] = true
				q = append(q, node{nr, nc})
			}
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		board := make([][]byte, 3)
		startRow := 0
		for i := 0; i < 3; i++ {
			var line string
			fmt.Fscan(in, &line)
			for len(line) < n {
				line += " "
			}
			for j := 0; j < n; j++ {
				if line[j] == 's' {
					startRow = i
					line = line[:j] + "." + line[j+1:]
					break
				}
			}
			board[i] = []byte(line)
		}
		if canEscape(board, n, startRow) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
