package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	u, v int
	w    int
}

type DSU struct {
	p, r []int
}

func NewDSU(n int) *DSU {
	d := &DSU{p: make([]int, n), r: make([]int, n)}
	for i := 0; i < n; i++ {
		d.p[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) union(x, y int) bool {
	fx := d.find(x)
	fy := d.find(y)
	if fx == fy {
		return false
	}
	if d.r[fx] < d.r[fy] {
		fx, fy = fy, fx
	}
	d.p[fy] = fx
	if d.r[fx] == d.r[fy] {
		d.r[fx]++
	}
	return true
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func mstEdges(edges []Edge, n int, x int) []int {
	es := make([]Edge, len(edges))
	copy(es, edges)
	sort.Slice(es, func(i, j int) bool {
		d1 := abs(es[i].w - x)
		d2 := abs(es[j].w - x)
		if d1 == d2 {
			return es[i].w < es[j].w
		}
		return d1 < d2
	})
	dsu := NewDSU(n + 1)
	ws := make([]int, 0, n-1)
	for _, e := range es {
		if dsu.union(e.u, e.v) {
			ws = append(ws, e.w)
			if len(ws) == n-1 {
				break
			}
		}
	}
	sort.Ints(ws)
	return ws
}

type MSTInfo struct {
	weights []int
	prefix  []int64
}

func buildInfo(ws []int) MSTInfo {
	prefix := make([]int64, len(ws)+1)
	for i, w := range ws {
		prefix[i+1] = prefix[i] + int64(w)
	}
	return MSTInfo{weights: ws, prefix: prefix}
}

func (info MSTInfo) cost(x int) int64 {
	w := info.weights
	idx := sort.Search(len(w), func(i int) bool { return w[i] > x })
	sumLess := info.prefix[idx]
	total := info.prefix[len(w)]
	return int64(x)*int64(idx) - sumLess + (total - sumLess) - int64(x)*int64(len(w)-idx)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	edges := make([]Edge, m)
	weights := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
		weights[i] = edges[i].w
	}

	var p, k int
	var a, b, c int
	fmt.Fscan(in, &p, &k, &a, &b, &c)
	q := make([]int, k)
	for i := 0; i < p; i++ {
		fmt.Fscan(in, &q[i])
	}
	for i := p; i < k; i++ {
		q[i] = (q[i-1]*a + b) % c
	}

	// build boundaries
	boundariesMap := make(map[int]struct{})
	boundariesMap[0] = struct{}{}
	boundariesMap[c] = struct{}{}
	for i := 0; i < m; i++ {
		for j := i; j < m; j++ {
			v := (weights[i] + weights[j] + 1) / 2
			if v < 0 {
				v = 0
			}
			if v > c {
				v = c
			}
			boundariesMap[v] = struct{}{}
		}
	}
	boundaries := make([]int, 0, len(boundariesMap))
	for v := range boundariesMap {
		boundaries = append(boundaries, v)
	}
	sort.Ints(boundaries)

	starts := make([]int, 0)
	for _, v := range boundaries {
		if v < c {
			starts = append(starts, v)
		}
	}
	// ensure last interval up to c-1
	starts = append(starts, c)

	infos := make([]MSTInfo, len(starts))
	for i, st := range starts {
		if st == c {
			// dummy for last; won't be used
			infos[i] = MSTInfo{}
		} else {
			ws := mstEdges(edges, n, st)
			infos[i] = buildInfo(ws)
		}
	}

	// answer queries
	var xorResult int64
	for _, x := range q {
		idx := sort.Search(len(starts), func(i int) bool { return starts[i] > x }) - 1
		if idx < 0 {
			idx = 0
		}
		ans := infos[idx].cost(x)
		xorResult ^= ans
	}
	fmt.Fprintln(out, xorResult)
}
