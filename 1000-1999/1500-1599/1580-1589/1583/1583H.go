package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const LOG = 19

type edge struct {
	u, v int
	c    int
	t    int
}

type query struct {
	v   int
	x   int
	idx int
}

var (
	g [][]struct {
		to   int
		toll int
	}
	up    [][]int
	mx    [][]int
	depth []int
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func buildLCA(n int) {
	up = make([][]int, LOG)
	mx = make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, n+1)
		mx[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	// BFS from node 1
	q := make([]int, 0, n)
	q = append(q, 1)
	up[0][1] = 0
	mx[0][1] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range g[v] {
			if e.to == up[0][v] {
				continue
			}
			up[0][e.to] = v
			mx[0][e.to] = e.toll
			depth[e.to] = depth[v] + 1
			q = append(q, e.to)
		}
	}
	for k := 1; k < LOG; k++ {
		for v := 1; v <= n; v++ {
			anc := up[k-1][v]
			up[k][v] = up[k-1][anc]
			if anc != 0 {
				mx[k][v] = maxInt(mx[k-1][v], mx[k-1][anc])
			} else {
				mx[k][v] = mx[k-1][v]
			}
		}
	}
}

func maxEdge(u, v int) int {
	if u == v {
		return 0
	}
	res := 0
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for k := LOG - 1; k >= 0; k-- {
		if diff&(1<<k) != 0 {
			if mx[k][u] > res {
				res = mx[k][u]
			}
			u = up[k][u]
		}
	}
	if u == v {
		return res
	}
	for k := LOG - 1; k >= 0; k-- {
		if up[k][u] != up[k][v] {
			if mx[k][u] > res {
				res = mx[k][u]
			}
			if mx[k][v] > res {
				res = mx[k][v]
			}
			u = up[k][u]
			v = up[k][v]
		}
	}
	if mx[0][u] > res {
		res = mx[0][u]
	}
	if mx[0][v] > res {
		res = mx[0][v]
	}
	return res
}

// DSU structure maintaining best enjoyment and diameter among nodes with that enjoyment

type DSU struct {
	parent []int
	size   []int
	emax   []int
	a      []int
	b      []int
}

func newDSU(n int, e []int) *DSU {
	parent := make([]int, n+1)
	size := make([]int, n+1)
	emax := make([]int, n+1)
	a := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
		emax[i] = e[i]
		a[i] = i
		b[i] = i
	}
	return &DSU{parent, size, emax, a, b}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(u, v int) {
	ru := d.find(u)
	rv := d.find(v)
	if ru == rv {
		return
	}
	if d.size[ru] < d.size[rv] {
		ru, rv = rv, ru
	}
	d.parent[rv] = ru
	d.size[ru] += d.size[rv]

	if d.emax[ru] < d.emax[rv] {
		d.emax[ru] = d.emax[rv]
		d.a[ru] = d.a[rv]
		d.b[ru] = d.b[rv]
	} else if d.emax[ru] == d.emax[rv] {
		nodes := []int{d.a[ru], d.b[ru], d.a[rv], d.b[rv]}
		bestA := d.a[ru]
		bestB := d.b[ru]
		bestD := maxEdge(bestA, bestB)
		for i := 0; i < len(nodes); i++ {
			for j := i + 1; j < len(nodes); j++ {
				dtmp := maxEdge(nodes[i], nodes[j])
				if dtmp > bestD {
					bestD = dtmp
					bestA = nodes[i]
					bestB = nodes[j]
				}
			}
		}
		d.a[ru] = bestA
		d.b[ru] = bestB
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	e := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &e[i])
	}
	edges := make([]edge, n-1)
	g = make([][]struct {
		to   int
		toll int
	}, n+1)
	for i := 0; i < n-1; i++ {
		var a, b, c, t int
		fmt.Fscan(reader, &a, &b, &c, &t)
		edges[i] = edge{a, b, c, t}
		g[a] = append(g[a], struct {
			to   int
			toll int
		}{b, t})
		g[b] = append(g[b], struct {
			to   int
			toll int
		}{a, t})
	}
	buildLCA(n)

	qs := make([]query, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &qs[i].v, &qs[i].x)
		qs[i].idx = i
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].c > edges[j].c })
	sort.Slice(qs, func(i, j int) bool { return qs[i].v > qs[j].v })

	d := newDSU(n, e)
	ansE := make([]int, q)
	ansT := make([]int, q)
	ei := 0
	for _, qu := range qs {
		for ei < len(edges) && edges[ei].c >= qu.v {
			d.union(edges[ei].u, edges[ei].v)
			ei++
		}
		root := d.find(qu.x)
		ansE[qu.idx] = d.emax[root]
		tollA := maxEdge(qu.x, d.a[root])
		tollB := maxEdge(qu.x, d.b[root])
		if tollA > tollB {
			ansT[qu.idx] = tollA
		} else {
			ansT[qu.idx] = tollB
		}
	}

	for i := 0; i < q; i++ {
		fmt.Fprintf(writer, "%d %d\n", ansE[i], ansT[i])
	}
}
