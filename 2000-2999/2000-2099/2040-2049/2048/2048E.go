package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	u, v  int
	right int
	color int
}

type dsu struct {
	p []int
	s []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &dsu{p: p, s: s}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) bool {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return false
	}
	if d.s[a] < d.s[b] {
		a, b = b, a
	}
	d.p[b] = a
	d.s[a] += d.s[b]
	return true
}

func buildMatchings(n, m int) [][][2]int {
	// Round-robin factorization for 2n vertices, keep first m matchings.
	rounds := make([][][2]int, 0, m)
	total := 2*n - 1
	for r := 0; r < m; r++ {
		pairs := make([][2]int, 0, n)
		pairs = append(pairs, [2]int{2*n - 1, r})
		for k := 1; k < n; k++ {
			a := (r + k) % total
			b := (r - k + total) % total
			pairs = append(pairs, [2]int{a, b})
		}
		rounds = append(rounds, pairs)
	}
	return rounds
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		if m > 2*n-1 {
			fmt.Fprintln(out, "NO")
			continue
		}

		matchings := buildMatchings(n, m)
		edges := make([]edge, 0, n*m)
		for r, pairs := range matchings {
			for _, p := range pairs {
				edges = append(edges, edge{u: p[0], v: p[1], right: r})
			}
		}

		remaining := make([]int, len(edges))
		for i := range remaining {
			remaining[i] = i
		}

		color := 0
		for len(remaining) > 0 && color < n {
			color++
			d := newDSU(2 * n)
			used := make([]bool, len(remaining))
			for i, idx := range remaining {
				e := &edges[idx]
				if d.union(e.u, e.v) {
					e.color = color
					used[i] = true
				}
			}
			newRemaining := make([]int, 0, len(remaining))
			for i, idx := range remaining {
				if !used[i] {
					newRemaining = append(newRemaining, idx)
				}
			}
			remaining = newRemaining
		}

		if len(remaining) > 0 {
			fmt.Fprintln(out, "NO")
			continue
		}

		ans := make([][]int, 2*n)
		for i := range ans {
			ans[i] = make([]int, m)
		}
		for _, e := range edges {
			ans[e.u][e.right] = e.color
			ans[e.v][e.right] = e.color
		}

		fmt.Fprintln(out, "YES")
		for i := 0; i < 2*n; i++ {
			row := ans[i]
			for j, v := range row {
				if j > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			fmt.Fprintln(out)
		}
	}
}
