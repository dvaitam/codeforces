package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type edge struct {
	u, v int
	w    int64
	id   int
	inMS bool
}

type dsu struct {
	p, sz []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	sz := make([]int, n)
	for i := range p {
		p[i] = i
		sz[i] = 1
	}
	return &dsu{p: p, sz: sz}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) unite(a, b int) bool {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return false
	}
	if d.sz[ra] < d.sz[rb] {
		ra, rb = rb, ra
	}
	d.p[rb] = ra
	d.sz[ra] += d.sz[rb]
	return true
}

type pair struct {
	to int
	w  int64
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	cost := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &cost[i])
	}
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}

	l := make([]int, n)
	r := make([]int, n)
	type state struct {
		v, p  int
		stage int
	}
	stack := []state{{v: 0, p: -1, stage: 0}}
	leafIdx := 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur.stage == 0 {
			stack = append(stack, state{v: cur.v, p: cur.p, stage: 1})
			for _, to := range g[cur.v] {
				if to == cur.p {
					continue
				}
				stack = append(stack, state{v: to, p: cur.v, stage: 0})
			}
		} else {
			if cur.v != 0 && len(g[cur.v]) == 1 {
				leafIdx++
				l[cur.v] = leafIdx
				r[cur.v] = leafIdx
			} else {
				l[cur.v] = int(1e9)
				r[cur.v] = -1
				for _, to := range g[cur.v] {
					if to == cur.p {
						continue
					}
					if l[to] < l[cur.v] {
						l[cur.v] = l[to]
					}
					if r[to] > r[cur.v] {
						r[cur.v] = r[to]
					}
				}
			}
		}
	}
	m := leafIdx
	edges := make([]edge, n)
	for i := 0; i < n; i++ {
		edges[i] = edge{
			u:  l[i],
			v:  r[i] + 1,
			w:  cost[i],
			id: i,
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		if edges[i].w == edges[j].w {
			return edges[i].id < edges[j].id
		}
		return edges[i].w < edges[j].w
	})

	d := newDSU(m + 1 + 1) // vertices are 1..m+1
	mstAdj := make([][]pair, m+2)
	var mstWeight int64
	for i := range edges {
		u := edges[i].u
		v := edges[i].v
		if d.unite(u, v) {
			edges[i].inMS = true
			mstWeight += edges[i].w
			mstAdj[u] = append(mstAdj[u], pair{to: v, w: edges[i].w})
			mstAdj[v] = append(mstAdj[v], pair{to: u, w: edges[i].w})
		}
	}

	// binary lifting on MST
	log := 0
	for (1 << log) <= m+1 {
		log++
	}
	up := make([][]int, log)
	mx := make([][]int64, log)
	for i := 0; i < log; i++ {
		up[i] = make([]int, m+2)
		mx[i] = make([]int64, m+2)
	}
	depth := make([]int, m+2)
	stack2 := []int{1}
	up[0][1] = 0
	mx[0][1] = 0
	depth[1] = 0
	for len(stack2) > 0 {
		v := stack2[len(stack2)-1]
		stack2 = stack2[:len(stack2)-1]
		for _, e := range mstAdj[v] {
			if e.to == up[0][v] {
				continue
			}
			up[0][e.to] = v
			mx[0][e.to] = e.w
			depth[e.to] = depth[v] + 1
			stack2 = append(stack2, e.to)
		}
	}
	for k := 1; k < log; k++ {
		for v := 1; v <= m+1; v++ {
			ancestor := up[k-1][v]
			up[k][v] = up[k-1][ancestor]
			mx[k][v] = max64(mx[k-1][v], mx[k-1][ancestor])
		}
	}

	getMax := func(u, v int) int64 {
		res := int64(0)
		if depth[u] < depth[v] {
			u, v = v, u
		}
		diff := depth[u] - depth[v]
		for k := 0; k < log; k++ {
			if diff>>k&1 == 1 {
				if mx[k][u] > res {
					res = mx[k][u]
				}
				u = up[k][u]
			}
		}
		if u == v {
			return res
		}
		for k := log - 1; k >= 0; k-- {
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

	possible := make([]bool, n)
	for _, e := range edges {
		if e.w == 0 {
			possible[e.id] = true
		}
		if e.inMS {
			possible[e.id] = true
		}
	}
	for _, e := range edges {
		if e.inMS || e.w == 0 {
			continue
		}
		if getMax(e.u, e.v) == e.w {
			possible[e.id] = true
		}
	}

	ids := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if possible[i] {
			ids = append(ids, i+1)
		}
	}
	fmt.Fprintln(out, mstWeight, len(ids))
	for i, v := range ids {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	if len(ids) > 0 {
		fmt.Fprintln(out)
	}
}
