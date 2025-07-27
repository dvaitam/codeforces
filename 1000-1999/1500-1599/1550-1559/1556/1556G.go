package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This solution implements a compressed graph for the n-dimensional hypercube
// built from disjoint blocked segments. Vertices of the compressed graph
// correspond to power-of-two aligned segments. Queries are processed in reverse
// order using DSU.

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	type Query struct {
		t    int // 0 ask, 1 block
		a, b uint64
	}
	queries := make([]Query, m)
	blocks := []Query{}
	for i := 0; i < m; i++ {
		var s string
		fmt.Fscan(in, &s)
		if s == "ask" {
			var a, b uint64
			fmt.Fscan(in, &a, &b)
			queries[i] = Query{t: 0, a: a, b: b}
		} else {
			var l, r uint64
			fmt.Fscan(in, &l, &r)
			queries[i] = Query{t: 1, a: l, b: r}
			blocks = append(blocks, Query{t: i + 1, a: l, b: r})
		}
	}

	sort.Slice(blocks, func(i, j int) bool { return blocks[i].a < blocks[j].a })

	// DSU
	parent := []int{}
	size := []int{}
	active := []bool{}
	adj := [][]int{}
	actAt := make([][]int, m+1)

	var vertexKey = map[[2]uint64]int{}
	var addVertex func(uint64, int, int) int
	addVertex = func(start uint64, k int, t int) int {
		id := len(parent)
		parent = append(parent, id)
		size = append(size, 1)
		active = append(active, false)
		adj = append(adj, nil)
		vertexKey[[2]uint64{start, uint64(k)}] = id
		if t >= 0 && t <= m {
			actAt[t] = append(actAt[t], id)
		}
		return id
	}

	var addEdge = func(a, b int) {
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	// build tree recursively
	var build func(uint64, int, int, int) []struct {
		start, end uint64
		id         int
	}
	build = func(start uint64, k int, li, ri int) []struct {
		start, end uint64
		id         int
	} {
		end := start + (uint64(1) << k)
		// if no blocks in range
		if li > ri {
			id := addVertex(start, k, 0)
			return []struct {
				start, end uint64
				id         int
			}{{start, end, id}}
		}
		if li == ri && blocks[li].a == start && blocks[li].b+1 == end {
			id := addVertex(start, k, int(blocks[li].t))
			return []struct {
				start, end uint64
				id         int
			}{{start, end, id}}
		}
		mid := start + (uint64(1) << (k - 1))
		mi := li - 1
		for mi+1 <= ri && blocks[mi+1].a < mid {
			mi++
		}
		left := build(start, k-1, li, mi)
		right := build(mid, k-1, mi+1, ri)
		// merge edges between left and right
		i, j := 0, 0
		pos := start
		for pos < mid {
			L := left[i]
			R := right[j]
			addEdge(L.id, R.id)
			stepL := L.end - pos
			stepR := R.end - (mid + (pos - start))
			step := stepL
			if stepR < step {
				step = stepR
			}
			pos += step
			if pos >= L.end {
				i++
			}
			if pos-start >= R.end-mid {
				j++
			}
		}
		return append(left, right...)
	}

	build(0, n, 0, len(blocks)-1)

	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(a, b int) {
		a = find(a)
		b = find(b)
		if a == b {
			return
		}
		if size[a] < size[b] {
			a, b = b, a
		}
		parent[b] = a
		size[a] += size[b]
	}
	activateVert := func(v int) {
		if active[v] {
			return
		}
		active[v] = true
		for _, nb := range adj[v] {
			if active[nb] {
				union(v, nb)
			}
		}
	}

	// activate time 0 vertices
	for _, v := range actAt[0] {
		activateVert(v)
	}

	// map number to vertex
	var rootStart uint64 = 0
	var search func(uint64, uint64, int) int
	search = func(x uint64, start uint64, k int) int {
		if id, ok := vertexKey[[2]uint64{start, uint64(k)}]; ok {
			if active[id] || k == 0 {
				if x >= start && x < start+(uint64(1)<<k) {
					return id
				}
			}
		}
		mid := start + (uint64(1) << (k - 1))
		if x < mid {
			return search(x, start, k-1)
		}
		return search(x, mid, k-1)
	}

	out := bufio.NewWriter(os.Stdout)
	ans := make([]byte, 0, m)
	for i := m - 1; i >= 0; i-- {
		q := queries[i]
		if q.t == 0 {
			a := search(q.a, rootStart, n)
			b := search(q.b, rootStart, n)
			if active[a] && active[b] && find(a) == find(b) {
				ans = append(ans, '1', '\n')
			} else {
				ans = append(ans, '0', '\n')
			}
		} else {
			for _, v := range actAt[i+1] {
				activateVert(v)
			}
		}
	}
	out.Write(ans)
	out.Flush()
}
