package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	id int
}

type DSU struct {
	parent []int
	rank   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n+1), rank: make([]int, n+1)}
	for i := 0; i <= n; i++ {
		d.parent[i] = i
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

func (d *DSU) union(a, b int) bool {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return false
	}
	if d.rank[a] < d.rank[b] {
		a, b = b, a
	}
	d.parent[b] = a
	if d.rank[a] == d.rank[b] {
		d.rank[a]++
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][2]int, m)
	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		edges[i] = [2]int{u, v}
		id := i + 1
		adj[u] = append(adj[u], Edge{v, id})
		adj[v] = append(adj[v], Edge{u, id})
	}

	dsu := NewDSU(n)
	mst := make(map[int]struct{})
	for i, e := range edges {
		if dsu.union(e[0], e[1]) {
			mst[i+1] = struct{}{}
		}
	}

	ans := make([]byte, n)
	for r := 1; r <= n; r++ {
		visited := make([]bool, n+1)
		selected := make(map[int]struct{})
		type item struct{ v, idx int }
		st := []item{{r, 0}}
		visited[r] = true
		for len(st) > 0 {
			top := &st[len(st)-1]
			v := top.v
			if top.idx >= len(adj[v]) {
				st = st[:len(st)-1]
				continue
			}
			e := adj[v][top.idx]
			top.idx++
			if !visited[e.to] {
				visited[e.to] = true
				selected[e.id] = struct{}{}
				st = append(st, item{e.to, 0})
			}
		}
		if len(selected) != n-1 {
			ans[r-1] = '0'
			continue
		}
		if len(selected) != len(mst) {
			ans[r-1] = '0'
			continue
		}
		ok := true
		for id := range selected {
			if _, ok2 := mst[id]; !ok2 {
				ok = false
				break
			}
		}
		if ok {
			ans[r-1] = '1'
		} else {
			ans[r-1] = '0'
		}
	}
	fmt.Fprintln(out, string(ans))
}
