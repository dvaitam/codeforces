package main

import (
	"bufio"
	"fmt"
	"os"
)

// UnionFind is a disjoint set union structure.
type UnionFind struct {
	parent []int
	size   []int
}

// NewUnionFind initializes a UnionFind for 1..n
func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n+1)
	size := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &UnionFind{parent: parent, size: size}
}

// Find returns the representative of x
func (u *UnionFind) Find(x int) int {
	if u.parent[x] != x {
		u.parent[x] = u.Find(u.parent[x])
	}
	return u.parent[x]
}

// Union merges the sets containing x and y
func (u *UnionFind) Union(x, y int) {
	fx := u.Find(x)
	fy := u.Find(y)
	if fx == fy {
		return
	}
	if u.size[fx] < u.size[fy] {
		fx, fy = fy, fx
	}
	u.parent[fy] = fx
	u.size[fx] += u.size[fy]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	type Edge struct{ u, v, w int }
	edges := make([]Edge, m)
	even := make([]bool, n+1)

	for i := 0; i < m; i++ {
		var a, b, w int
		fmt.Fscan(reader, &a, &b, &w)
		edges[i] = Edge{a, b, w}
		if w%2 == 0 {
			even[a] = true
			even[b] = true
		}
	}

	// Build DSU for each bit
	uf := make([]*UnionFind, 30)
	for j := 0; j < 30; j++ {
		uf[j] = NewUnionFind(n)
	}
	for _, e := range edges {
		for j := 0; j < 30; j++ {
			if (e.w>>j)&1 == 1 {
				uf[j].Union(e.u, e.v)
			}
		}
	}

	// good[j][root] indicates component reachable using bit j edges
	// that contains a vertex incident to an even-weight edge
	good := make([][]bool, 30)
	for j := 1; j < 30; j++ { // j=0 does not help for answer=1
		good[j] = make([]bool, n+1)
		for v := 1; v <= n; v++ {
			if even[v] {
				root := uf[j].Find(v)
				good[j][root] = true
			}
		}
	}

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		ans := 2
		for j := 0; j < 30; j++ {
			if uf[j].Find(u) == uf[j].Find(v) {
				ans = 0
				break
			}
		}
		if ans != 0 {
			for j := 1; j < 30 && ans > 1; j++ {
				if good[j][uf[j].Find(u)] {
					ans = 1
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
