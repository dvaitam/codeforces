package main

import (
	"bufio"
	"fmt"
	"os"
)

func computeLR(p []int) ([]int, []int) {
	n := len(p)
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = i
		for j := i - 1; j >= 0; j-- {
			if p[j] > p[i] {
				l[i] = j
				break
			}
		}
		r[i] = i
		for j := i + 1; j < n; j++ {
			if p[j] > p[i] {
				r[i] = j
				break
			}
		}
	}
	return l, r
}

func diameter(n int, edges [][2]int) int {
	g := make([][]int, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	bfs := func(start int) (int, int) {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		q := []int{start}
		dist[start] = 0
		idx := 0
		for idx < len(q) {
			v := q[idx]
			idx++
			for _, u := range g[v] {
				if dist[u] == -1 {
					dist[u] = dist[v] + 1
					q = append(q, u)
				}
			}
		}
		far := start
		for i, d := range dist {
			if d > dist[far] {
				far = i
			}
			if d == -1 {
				return -1, -1
			}
		}
		return far, dist[far]
	}
	f, _ := bfs(0)
	if f == -1 {
		return -1
	}
	_, d := bfs(f)
	return d
}

func solve(n int, p []int, s string) int {
	l, r := computeLR(p)
	var qIdx []int
	for i, ch := range s {
		if ch == '?' {
			qIdx = append(qIdx, i)
		}
	}
	best := -1
	total := 1 << len(qIdx)
	for mask := 0; mask < total; mask++ {
		edges := make([][2]int, 0, n)
		connected := true
		for i := 0; i < n; i++ {
			c := s[i]
			if c == '?' {
				bit := 0
				for j, idx := range qIdx {
					if idx == i {
						bit = (mask >> j) & 1
						break
					}
				}
				if bit == 0 {
					c = 'L'
				} else {
					c = 'R'
				}
			}
			var to int
			if c == 'L' {
				to = l[i]
			} else {
				to = r[i]
			}
			if to == i {
				connected = false
				break
			}
			edges = append(edges, [2]int{i, to})
		}
		if !connected {
			continue
		}
		d := diameter(n, edges)
		if d > best {
			best = d
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		var s string
		fmt.Fscan(reader, &s)
		fmt.Fprintln(writer, solve(n, p, s))
	}
}
