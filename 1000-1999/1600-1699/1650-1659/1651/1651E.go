package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	left   []int
	right  []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), left: make([]int, n), right: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	d.parent[rb] = ra
	d.left[ra] += d.left[rb]
	d.right[ra] += d.right[rb]
}

func calc(n int, adj [][]int, l, r, L, R int) int {
	size := 2 * n
	active := make([]bool, size)
	d := NewDSU(size)

	for i := l - 1; i <= r-1; i++ {
		active[i] = true
		d.left[i] = 1
	}
	for j := L - 1; j <= R-1; j++ {
		active[j] = true
		d.right[j] = 1
	}

	for v := 0; v < size; v++ {
		if !active[v] {
			continue
		}
		for _, u := range adj[v] {
			if active[u] {
				d.union(v, u)
			}
		}
	}

	visited := make([]bool, size)
	mm := 0
	for v := 0; v < size; v++ {
		if !active[v] {
			continue
		}
		root := d.find(v)
		if visited[root] {
			continue
		}
		visited[root] = true
		if d.left[root] < d.right[root] {
			mm += d.left[root]
		} else {
			mm += d.right[root]
		}
	}
	return mm
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	size := 2 * n
	adj := make([][]int, size)
	for i := 0; i < 2*n; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		x--
		y--
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	var ans int64 = 0
	for l := 1; l <= n; l++ {
		for r := l; r <= n; r++ {
			for L := n + 1; L <= 2*n; L++ {
				for R := L; R <= 2*n; R++ {
					mm := calc(n, adj, l, r, L, R)
					ans += int64(mm)
				}
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
