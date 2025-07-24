package main

import (
	"bufio"
	"fmt"
	"os"
)

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	d := newDSU(n)
	deg := make([]int64, n+1)
	has := make([]bool, n+1)
	var loops int64

	edges := int64(m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		if u == v {
			loops++
			has[u] = true
		} else {
			deg[u]++
			deg[v]++
			d.union(u, v)
			has[u] = true
			has[v] = true
		}
	}

	root := -1
	for i := 1; i <= n; i++ {
		if has[i] {
			if root == -1 {
				root = d.find(i)
			} else if d.find(i) != root {
				fmt.Fprintln(writer, 0)
				return
			}
		}
	}

	if edges < 2 {
		fmt.Fprintln(writer, 0)
		return
	}

	ans := loops * (loops - 1) / 2
	ans += loops * (edges - loops)
	for i := 1; i <= n; i++ {
		if deg[i] >= 2 {
			ans += deg[i] * (deg[i] - 1) / 2
		}
	}

	fmt.Fprintln(writer, ans)
}
