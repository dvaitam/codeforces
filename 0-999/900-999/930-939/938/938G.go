package main

import (
	"bufio"
	"fmt"
	"os"
)

// Edge represents an undirected weighted edge
type Edge struct {
	u, v int
	w    int
}

// linear basis operation record for rollback
type basisOp struct {
	root int
	idx  int
}

// union operation record for rollback
type unionOp struct {
	joined   bool
	rv, ru   int
	sizeRu   int
	basisLen int
}

// Basis implements xor linear basis of up to 31 bits
type Basis struct {
	b [31]int
}

func (bs *Basis) insert(x int, root int, hist *[]basisOp) {
	v := x
	for i := 30; i >= 0; i-- {
		if (v>>i)&1 == 0 {
			continue
		}
		if bs.b[i] == 0 {
			bs.b[i] = v
			*hist = append(*hist, basisOp{root, i})
			return
		}
		if v^bs.b[i] < v {
			v ^= bs.b[i]
		}
	}
}

func (bs *Basis) minimize(x int) int {
	v := x
	for i := 30; i >= 0; i-- {
		if bs.b[i] != 0 && (v^bs.b[i]) < v {
			v ^= bs.b[i]
		}
	}
	return v
}

// DSU with xor weights and rollback
type DSU struct {
	parent []int
	size   []int
	xor    []int
	basis  []Basis
	hist   []unionOp
	bHist  []basisOp
}

var (
	seg     [][]Edge
	queries []Query
	answers []int
	n, m, q int
)

func NewDSU(n int) *DSU {
	d := &DSU{}
	d.parent = make([]int, n+1)
	d.size = make([]int, n+1)
	d.xor = make([]int, n+1)
	d.basis = make([]Basis, n+1)
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) (int, int) {
	xr := 0
	for x != d.parent[x] {
		xr ^= d.xor[x]
		x = d.parent[x]
	}
	return x, xr
}

func (d *DSU) union(x, y, w int) {
	ru, xu := d.find(x)
	rv, xv := d.find(y)
	val := xu ^ xv ^ w
	before := len(d.bHist)
	if ru == rv {
		d.basis[ru].insert(val, ru, &d.bHist)
		d.hist = append(d.hist, unionOp{joined: false, basisLen: before})
		return
	}
	if d.size[ru] < d.size[rv] {
		ru, rv = rv, ru
	}
	d.parent[rv] = ru
	d.xor[rv] = val
	d.size[ru] += d.size[rv]
	for i := 30; i >= 0; i-- {
		if d.basis[rv].b[i] != 0 {
			d.basis[ru].insert(d.basis[rv].b[i], ru, &d.bHist)
		}
	}
	d.hist = append(d.hist, unionOp{joined: true, rv: rv, ru: ru, sizeRu: d.size[ru] - d.size[rv], basisLen: before})
}

func (d *DSU) rollback(to int) {
	for len(d.hist) > to {
		op := d.hist[len(d.hist)-1]
		d.hist = d.hist[:len(d.hist)-1]
		for len(d.bHist) > op.basisLen {
			bh := d.bHist[len(d.bHist)-1]
			d.bHist = d.bHist[:len(d.bHist)-1]
			d.basis[bh.root].b[bh.idx] = 0
		}
		if op.joined {
			d.parent[op.rv] = op.rv
			d.xor[op.rv] = 0
			d.size[op.ru] = op.sizeRu
		}
	}
}

func (d *DSU) query(x, y int) int {
	ru, xu := d.find(x)
	rv, xv := d.find(y)
	val := xu ^ xv
	if ru != rv {
		// graph should remain connected, so this shouldn't happen
	}
	return d.basis[ru].minimize(val)
}

// segment tree add helper
func addEdge(node, l, r, ql, qr int, e Edge) {
	if ql > qr || ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		seg[node] = append(seg[node], e)
		return
	}
	mid := (l + r) / 2
	if ql <= mid {
		addEdge(node*2, l, mid, ql, qr, e)
	}
	if qr > mid {
		addEdge(node*2+1, mid+1, r, ql, qr, e)
	}
}

func dfs(node, l, r int, dsu *DSU) {
	before := len(dsu.hist)
	for _, e := range seg[node] {
		dsu.union(e.u, e.v, e.w)
	}
	if l == r {
		if queries[l].typ == 3 {
			ans := dsu.query(queries[l].x, queries[l].y)
			answers = append(answers, ans)
		}
	} else {
		mid := (l + r) / 2
		dfs(node*2, l, mid, dsu)
		dfs(node*2+1, mid+1, r, dsu)
	}
	dsu.rollback(before)
}

// Query represents a single query
type Query struct {
	typ int
	x   int
	y   int
	w   int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &m)
	edgeStart := make(map[[2]int]int)
	edgeWeight := make(map[[2]int]int)

	for i := 0; i < m; i++ {
		var x, y, w int
		fmt.Fscan(reader, &x, &y, &w)
		if x > y {
			x, y = y, x
		}
		key := [2]int{x, y}
		edgeStart[key] = 0
		edgeWeight[key] = w
	}

	fmt.Fscan(reader, &q)
	queries = make([]Query, q+1)
	seg = make([][]Edge, 4*q+4)
	for i := 1; i <= q; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var x, y, w int
			fmt.Fscan(reader, &x, &y, &w)
			if x > y {
				x, y = y, x
			}
			key := [2]int{x, y}
			edgeStart[key] = i
			edgeWeight[key] = w
			queries[i] = Query{typ: 1, x: x, y: y, w: w}
		} else if t == 2 {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			if x > y {
				x, y = y, x
			}
			key := [2]int{x, y}
			start := edgeStart[key]
			w := edgeWeight[key]
			addEdge(1, 1, q, start+1, i-1, Edge{u: x, v: y, w: w})
			delete(edgeStart, key)
			delete(edgeWeight, key)
			queries[i] = Query{typ: 2, x: x, y: y}
		} else {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			if x > y {
				x, y = y, x
			}
			queries[i] = Query{typ: 3, x: x, y: y}
		}
	}

	for key, st := range edgeStart {
		w := edgeWeight[key]
		addEdge(1, 1, q, st+1, q, Edge{u: key[0], v: key[1], w: w})
	}

	dsu := NewDSU(n)
	dfs(1, 1, q, dsu)

	for i, v := range answers {
		if i > 0 {
			fmt.Fprint(writer, "\n")
		}
		fmt.Fprint(writer, v)
	}
	if len(answers) > 0 {
		fmt.Fprint(writer, "\n")
	}
}
