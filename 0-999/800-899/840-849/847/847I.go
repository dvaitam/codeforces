package main

import (
	"bufio"
	"fmt"
	"os"
)

type point struct{ r, c int }

type node struct {
	point
	val int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	var q, p int
	if _, err := fmt.Fscan(in, &n, &m, &q, &p); err != nil {
		return
	}
	grid := make([][]byte, n)
	letters := make([][]point, 26)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		row := []byte(s)
		grid[i] = row
		for j := 0; j < m; j++ {
			ch := row[j]
			if ch >= 'A' && ch <= 'Z' {
				letters[ch-'A'] = append(letters[ch-'A'], point{i, j})
			}
		}
	}
	noise := make([][]int, n)
	for i := 0; i < n; i++ {
		noise[i] = make([]int, m)
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for idx, posList := range letters {
		if len(posList) == 0 {
			continue
		}
		init := q * (idx + 1)
		if init == 0 {
			continue
		}
		vis := make([][]bool, n)
		for i := 0; i < n; i++ {
			vis[i] = make([]bool, m)
		}
		qnodes := make([]node, len(posList))
		for i, pnt := range posList {
			qnodes[i] = node{point: pnt, val: init}
			vis[pnt.r][pnt.c] = true
		}
		for head := 0; head < len(qnodes); head++ {
			cur := qnodes[head]
			noise[cur.r][cur.c] += cur.val
			nextVal := cur.val / 2
			if nextVal == 0 {
				continue
			}
			for _, d := range dirs {
				nr := cur.r + d[0]
				nc := cur.c + d[1]
				if nr < 0 || nr >= n || nc < 0 || nc >= m {
					continue
				}
				if grid[nr][nc] == '*' || vis[nr][nc] {
					continue
				}
				vis[nr][nc] = true
				qnodes = append(qnodes, node{point: point{nr, nc}, val: nextVal})
			}
		}
	}
	ans := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if noise[i][j] > p {
				ans++
			}
		}
	}
	fmt.Fprintln(out, ans)
}
