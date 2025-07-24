package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Edge represents a directed edge with a letter label
type Edge struct {
	to int
	ch byte
}

type pair struct {
	ch byte
	id int
}

type key struct {
	ch byte
	id int
}

type hash struct {
	a uint64
	b uint64
}

const (
	base1 uint64 = 911382323
	base2 uint64 = 972663749
)

var (
	n            int
	g            [][]Edge
	sz           []int
	depth        []int
	parent       []int
	subHash      []hash
	subID        []int
	idMap        map[hash]int
	maxDepth     int
	nodesAtDepth [][]int
	cntLevel     []int
)

func dfs(v, d int) {
	depth[v] = d
	if d > maxDepth {
		maxDepth = d
	}
	for _, e := range g[v] {
		parent[e.to] = v
		dfs(e.to, d+1)
	}
	sz[v] = 1
	pairs := make([]pair, len(g[v]))
	for i, e := range g[v] {
		pairs[i] = pair{ch: e.ch, id: subID[e.to]}
		sz[v] += sz[e.to]
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].ch == pairs[j].ch {
			return pairs[i].id < pairs[j].id
		}
		return pairs[i].ch < pairs[j].ch
	})
	h := hash{a: 1, b: 1}
	for _, p := range pairs {
		h.a = h.a*base1 + uint64(p.ch) + 1
		h.a = h.a*base1 + uint64(p.id)
		h.b = h.b*base2 + uint64(p.ch) + 1
		h.b = h.b*base2 + uint64(p.id)
	}
	h.a = h.a*base1 + 7
	h.b = h.b*base2 + 7
	subHash[v] = h
	if id, ok := idMap[h]; ok {
		subID[v] = id
	} else {
		id := len(idMap) + 1
		idMap[h] = id
		subID[v] = id
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fscan(in, &n)
	g = make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		var s string
		fmt.Fscan(in, &u, &v, &s)
		g[u] = append(g[u], Edge{to: v, ch: s[0]})
	}

	sz = make([]int, n+1)
	depth = make([]int, n+1)
	parent = make([]int, n+1)
	subHash = make([]hash, n+1)
	subID = make([]int, n+1)
	idMap = make(map[hash]int)

	dfs(1, 0)

	nodesAtDepth = make([][]int, maxDepth+1)
	cntLevel = make([]int, maxDepth+1)
	for v := 1; v <= n; v++ {
		d := depth[v]
		nodesAtDepth[d] = append(nodesAtDepth[d], v)
		cntLevel[d]++
	}

	dupSize := make([]int, maxDepth+1)
	for p := 1; p <= maxDepth; p++ {
		for _, u := range nodesAtDepth[p-1] {
			if len(g[u]) == 0 {
				continue
			}
			m := make(map[key]bool)
			for _, e1 := range g[u] {
				v := e1.to
				if depth[v] != p {
					continue
				}
				for _, e2 := range g[v] {
					k := key{ch: e2.ch, id: subID[e2.to]}
					if _, ok := m[k]; ok {
						dupSize[p] += sz[e2.to]
					} else {
						m[k] = true
					}
				}
			}
		}
	}

	bestSize := n
	bestP := 1
	for p := 1; p <= maxDepth; p++ {
		if cntLevel[p] == 0 {
			continue
		}
		cur := n - cntLevel[p] - dupSize[p]
		if cur < bestSize {
			bestSize = cur
			bestP = p
		}
	}

	fmt.Fprintln(out, bestSize)
	fmt.Fprintln(out, bestP)
}
