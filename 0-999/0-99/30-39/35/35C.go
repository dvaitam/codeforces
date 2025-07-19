package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	fmt.Fscan(in, &k)
	total := n * m
	visited := make([]bool, total)
	q := make([]int, 0, total)
	var last int
	for i := 0; i < k; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		idx := x*m + y
		if !visited[idx] {
			visited[idx] = true
			q = append(q, idx)
		}
	}
	for head := 0; head < len(q); head++ {
		idx := q[head]
		last = idx
		r := idx / m
		c := idx % m
		// up
		if r > 0 {
			ni := (r-1)*m + c
			if !visited[ni] {
				visited[ni] = true
				q = append(q, ni)
			}
		}
		// down
		if r+1 < n {
			ni := (r+1)*m + c
			if !visited[ni] {
				visited[ni] = true
				q = append(q, ni)
			}
		}
		// left
		if c > 0 {
			ni := r*m + (c - 1)
			if !visited[ni] {
				visited[ni] = true
				q = append(q, ni)
			}
		}
		// right
		if c+1 < m {
			ni := r*m + (c + 1)
			if !visited[ni] {
				visited[ni] = true
				q = append(q, ni)
			}
		}
	}
	rf := last/m + 1
	cf := last%m + 1
   fmt.Printf("%d %d\n", rf, cf)
}
