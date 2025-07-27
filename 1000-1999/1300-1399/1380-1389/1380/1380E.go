package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) Union(a, b int) int {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return a
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	t := make([]int, n+1)
	items := make([][]int, m+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &t[i])
		items[t[i]] = append(items[t[i]], i)
	}

	dsu := NewDSU(m)
	for i := 1; i <= m; i++ {
		dsu.size[i] = len(items[i])
	}

	ans := 0
	for i := 1; i < n; i++ {
		if t[i] != t[i+1] {
			ans++
		}
	}

	res := make([]int, 0, m)
	res = append(res, ans)

	queries := make([][2]int, m-1)
	for i := 0; i < m-1; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
	}

	for _, q := range queries {
		a := dsu.Find(q[0])
		b := dsu.Find(q[1])
		if a != b {
			if dsu.size[a] < dsu.size[b] {
				a, b = b, a
			}
			for _, x := range items[b] {
				if x > 1 && dsu.Find(t[x-1]) == a {
					ans--
				}
				if x < n && dsu.Find(t[x+1]) == a {
					ans--
				}
			}
			dsu.parent[b] = a
			dsu.size[a] += dsu.size[b]
			items[a] = append(items[a], items[b]...)
			items[b] = nil
		}
		res = append(res, ans)
	}

	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
