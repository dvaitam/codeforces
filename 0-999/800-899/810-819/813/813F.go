package main

import (
	"bufio"
	"fmt"
	"os"
)

// Edge represents an undirected edge between u and v
type Edge struct {
	u int
	v int
}

// operation record for rollback
type op struct {
	joined   bool
	rv       int
	ru       int
	sizeRu   int
	parityRV int
	badDelta int
}

// DSU with parity and rollback capability
type DSU struct {
	parent []int
	size   []int
	color  []int
	hist   []op
	bad    int
}

func NewDSU(n int) *DSU {
	d := &DSU{}
	d.parent = make([]int, n+1)
	d.size = make([]int, n+1)
	d.color = make([]int, n+1)
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) (int, int) {
	parity := 0
	for x != d.parent[x] {
		parity ^= d.color[x]
		x = d.parent[x]
	}
	return x, parity
}

func (d *DSU) union(x, y int) {
	ru, pu := d.find(x)
	rv, pv := d.find(y)
	if ru == rv {
		if pu^pv != 1 {
			d.bad++
			d.hist = append(d.hist, op{joined: false, badDelta: 1})
		} else {
			d.hist = append(d.hist, op{joined: false, badDelta: 0})
		}
		return
	}
	if d.size[ru] < d.size[rv] {
		ru, rv = rv, ru
		pu, pv = pv, pu
	}
	d.hist = append(d.hist, op{joined: true, rv: rv, ru: ru, sizeRu: d.size[ru], parityRV: d.color[rv]})
	d.parent[rv] = ru
	d.color[rv] = pu ^ pv ^ 1
	d.size[ru] += d.size[rv]
}

func (d *DSU) rollback(to int) {
	for len(d.hist) > to {
		op := d.hist[len(d.hist)-1]
		d.hist = d.hist[:len(d.hist)-1]
		d.bad -= op.badDelta
		if op.joined {
			d.parent[op.rv] = op.rv
			d.color[op.rv] = op.parityRV
			d.size[op.ru] = op.sizeRu
		}
	}
}

var seg [][]Edge
var ans []bool
var dsu *DSU

func addEdge(node, l, r, ql, qr int, e Edge) {
	if ql > r || qr < l || ql > qr {
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

func dfs(node, l, r int) {
	before := len(dsu.hist)
	for _, e := range seg[node] {
		dsu.union(e.u, e.v)
	}
	if l == r {
		ans[l] = dsu.bad == 0
	} else {
		mid := (l + r) / 2
		dfs(node*2, l, mid)
		dfs(node*2+1, mid+1, r)
	}
	dsu.rollback(before)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	seg = make([][]Edge, 4*q+5)
	ans = make([]bool, q+1)
	edgeStart := make(map[[2]int]int)

	for i := 1; i <= q; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		if x > y {
			x, y = y, x
		}
		key := [2]int{x, y}
		if st, ok := edgeStart[key]; ok {
			addEdge(1, 1, q, st, i-1, Edge{x, y})
			delete(edgeStart, key)
		} else {
			edgeStart[key] = i
		}
	}
	for key, st := range edgeStart {
		addEdge(1, 1, q, st, q, Edge{key[0], key[1]})
	}

	dsu = NewDSU(n)
	dfs(1, 1, q)

	for i := 1; i <= q; i++ {
		if ans[i] {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
