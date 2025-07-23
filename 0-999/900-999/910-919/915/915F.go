package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) union(a, b int) int64 {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return 0
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	pairs := int64(d.size[ra]) * int64(d.size[rb])
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return pairs
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &vals[i])
	}
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}

	type node struct {
		val int
		idx int
	}
	order := make([]node, n)
	for i := 0; i < n; i++ {
		order[i] = node{val: vals[i], idx: i}
	}

	// contributions when value is maximum
	sort.Slice(order, func(i, j int) bool { return order[i].val < order[j].val })
	dsu := NewDSU(n)
	active := make([]bool, n)
	var maxSum int64
	for _, p := range order {
		v := p.idx
		active[v] = true
		for _, to := range g[v] {
			if active[to] {
				maxSum += int64(p.val) * dsu.union(v, to)
			}
		}
	}

	// contributions when value is minimum
	sort.Slice(order, func(i, j int) bool { return order[i].val > order[j].val })
	dsu = NewDSU(n)
	for i := 0; i < n; i++ {
		active[i] = false
	}
	var minSum int64
	for _, p := range order {
		v := p.idx
		active[v] = true
		for _, to := range g[v] {
			if active[to] {
				minSum += int64(p.val) * dsu.union(v, to)
			}
		}
	}

	fmt.Println(maxSum - minSum)
}
