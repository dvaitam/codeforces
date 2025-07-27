package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ x, y int }

var (
	n, m  int
	q     int
	grid  [][]bool
	rlen  [][]int
	dlen  [][]int
	total int64
)

func recompute(i, j int, prev bool) {
	oldR := rlen[i][j]
	oldD := dlen[i][j]
	var oldC int64
	if prev {
		oldC = int64(oldR + oldD - 1)
	}
	var newR, newD int
	if grid[i][j] {
		newR = 1
		if j < m && grid[i][j+1] {
			newR += dlen[i][j+1]
		}
		newD = 1
		if i < n && grid[i+1][j] {
			newD += rlen[i+1][j]
		}
	}
	rlen[i][j] = newR
	dlen[i][j] = newD
	var newC int64
	if grid[i][j] {
		newC = int64(newR + newD - 1)
	}
	if oldC != newC {
		total += newC - oldC
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m, &q)
	grid = make([][]bool, n+2)
	rlen = make([][]int, n+2)
	dlen = make([][]int, n+2)
	for i := 0; i <= n+1; i++ {
		grid[i] = make([]bool, m+2)
		rlen[i] = make([]int, m+2)
		dlen[i] = make([]int, m+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			grid[i][j] = true
		}
	}
	for i := n; i >= 1; i-- {
		for j := m; j >= 1; j-- {
			r := 1
			if j < m {
				r += dlen[i][j+1]
			}
			d := 1
			if i < n {
				d += rlen[i+1][j]
			}
			rlen[i][j] = r
			dlen[i][j] = d
			total += int64(r + d - 1)
		}
	}
	out := bufio.NewWriter(os.Stdout)
	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		prev := grid[x][y]
		grid[x][y] = !grid[x][y]
		queue := []pair{{x, y}}
		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]
			oldR := rlen[p.x][p.y]
			oldD := dlen[p.x][p.y]
			pr := grid[p.x][p.y]
			if p.x == x && p.y == y {
				pr = prev
			}
			recompute(p.x, p.y, pr)
			if rlen[p.x][p.y] != oldR || dlen[p.x][p.y] != oldD {
				if p.x > 1 {
					queue = append(queue, pair{p.x - 1, p.y})
				}
				if p.y > 1 {
					queue = append(queue, pair{p.x, p.y - 1})
				}
			}
		}
		fmt.Fprintln(out, total)
	}
	out.Flush()
}
