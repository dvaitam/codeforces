package main

import (
	"bufio"
	"fmt"
	"os"
)

// Heavy-Light Decomposition
type HLD struct {
	n        int
	t        int
	heavy    []int
	par      []int
	sub      []int
	depth    []int
	in       []int
	head     []int
	tour     []int
	adj      [][]int
	children [][]int
}

func NewHLD(n int) *HLD {
	h := &HLD{
		n:        n,
		t:        0,
		heavy:    make([]int, n),
		par:      make([]int, n),
		sub:      make([]int, n),
		depth:    make([]int, n),
		in:       make([]int, n),
		head:     make([]int, n),
		tour:     make([]int, n),
		adj:      make([][]int, n),
		children: make([][]int, n),
	}
	return h
}

func (h *HLD) AddEdge(u, v int) {
	h.adj[u] = append(h.adj[u], v)
	h.adj[v] = append(h.adj[v], u)
}

func (h *HLD) dfsPrep(u, p int) {
	h.sub[u] = 1
	h.heavy[u] = -1
	h.par[u] = p
	h.depth[u] = h.depth[p] + 1
	for _, v := range h.adj[u] {
		if v == p {
			continue
		}
		h.dfsPrep(v, u)
		h.sub[u] += h.sub[v]
		if h.heavy[u] == -1 || h.sub[v] > h.sub[h.heavy[u]] {
			h.heavy[u] = v
		}
	}
}

func (h *HLD) dfsHLD(u, p int) {
	h.in[u] = h.t
	h.tour[h.t] = u
	h.t++
	// head: new chain if u is not heavy child of p
	if u != h.heavy[p] {
		h.head[u] = u
	} else {
		h.head[u] = h.head[p]
	}
	// recurse heavy child first
	if h.heavy[u] != -1 {
		h.dfsHLD(h.heavy[u], u)
		h.children[u] = append(h.children[u], h.heavy[u])
	}
	// recurse light children
	for _, v := range h.adj[u] {
		if v == p || v == h.heavy[u] {
			continue
		}
		h.children[u] = append(h.children[u], v)
		h.dfsHLD(v, u)
	}
}

func (h *HLD) Init() {
	h.dfsPrep(0, 0)
	h.dfsHLD(0, 0)
}

// find first child of anc on path to dec
func firstAncestor(h *HLD, anc, dec int) int {
	ch := h.children[anc]
	l, r := 1, len(ch)
	for l < r {
		m := (l + r) / 2
		if m == len(ch) || h.in[ch[m]] > h.in[dec] {
			r = m
		} else {
			l = m + 1
		}
	}
	return ch[l-1]
}

var (
	off          int
	adj          [][]int
	leafToTime   []int
	leaf1, leaf2 []int
)

func rec1(u, l, r int) {
	if u != 0 {
		p := (u - 1) / 2
		adj[u] = append(adj[u], p)
	}
	if l == r {
		leafToTime[u] = l
		leaf1[l] = u
	} else {
		m := (l + r) >> 1
		rec1(2*u+1, l, m)
		rec1(2*u+2, m+1, r)
	}
}

func rec2(u, l, r int) {
	idx := u + off
	if u != 0 {
		p := (u - 1) / 2
		adj[p+off] = append(adj[p+off], idx)
	}
	if l == r {
		leaf2[l] = idx
	} else {
		m := (l + r) >> 1
		rec2(2*u+1, l, m)
		rec2(2*u+2, m+1, r)
	}
}

func add1(u, l, r, s, e, x int) {
	if r < s || e < l {
		return
	}
	if s <= l && r <= e {
		adj[u] = append(adj[u], x)
		return
	}
	m := (l + r) >> 1
	add1(2*u+1, l, m, s, e, x)
	add1(2*u+2, m+1, r, s, e, x)
}

func add2(u, l, r, s, e, x int) {
	if r < s || e < l {
		return
	}
	if s <= l && r <= e {
		adj[x] = append(adj[x], u+off)
		return
	}
	m := (l + r) >> 1
	add2(2*u+1, l, m, s, e, x)
	add2(2*u+2, m+1, r, s, e, x)
}

func addUpperHelper(h *HLD, u, v, c int) {
	z := leaf2[h.in[c]]
	for h.head[u] != h.head[v] {
		if h.depth[h.head[u]] < h.depth[h.head[v]] {
			u, v = v, u
		}
		add1(0, 0, h.n-1, h.in[h.head[u]], h.in[u], z)
		u = h.par[h.head[u]]
	}
	if h.depth[u] > h.depth[v] {
		u, v = v, u
	}
	add1(0, 0, h.n-1, h.in[u], h.in[v], z)
}

func addLowerHelper(h *HLD, u, v, c int) {
	z := leaf1[h.in[c]]
	for h.head[u] != h.head[v] {
		if h.depth[h.head[u]] < h.depth[h.head[v]] {
			u, v = v, u
		}
		add2(0, 0, h.n-1, h.in[h.head[u]], h.in[u], z)
		u = h.par[h.head[u]]
	}
	if h.depth[u] > h.depth[v] {
		u, v = v, u
	}
	add2(0, 0, h.n-1, h.in[u], h.in[v], z)
}

func addUpper(h *HLD, a, b, c int) {
	for _, u := range []int{a, b} {
		if u == c {
			continue
		}
		if h.in[c] <= h.in[u] && h.in[c]+h.sub[c] >= h.in[u]+h.sub[u] {
			v := firstAncestor(h, c, u)
			addUpperHelper(h, u, v, c)
		} else {
			addUpperHelper(h, u, h.par[c], c)
		}
	}
}

func addLower(h *HLD, a, b, c int) {
	for _, u := range []int{a, b} {
		if u == c {
			continue
		}
		if h.in[c] <= h.in[u] && h.in[c]+h.sub[c] >= h.in[u]+h.sub[u] {
			v := firstAncestor(h, c, u)
			addLowerHelper(h, u, v, c)
		} else {
			addLowerHelper(h, u, h.par[c], c)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(reader, &n, &m)
	hoc := NewHLD(n)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		hoc.AddEdge(u, v)
	}
	hoc.Init()
	// build segment trees
	size := 1
	for size < n {
		size <<= 1
	}
	total1 := 2*size - 1
	off = total1
	total := total1 * 2
	adj = make([][]int, total)
	leafToTime = make([]int, total)
	leaf1 = make([]int, n)
	leaf2 = make([]int, n)
	for i := range leafToTime {
		leafToTime[i] = -1
	}
	rec1(0, 0, n-1)
	rec2(0, 0, n-1)
	// connect leaf2->leaf1
	for t := 0; t < n; t++ {
		adj[leaf2[t]] = append(adj[leaf2[t]], leaf1[t])
	}
	// process queries
	for i := 0; i < m; i++ {
		var t, a, b, c int
		fmt.Fscan(reader, &t, &a, &b, &c)
		a--
		b--
		c--
		if t == 2 {
			addUpper(hoc, a, b, c)
		} else {
			addLower(hoc, a, b, c)
		}
	}
	// topological sort
	inDeg := make([]int, total)
	for u := 0; u < total; u++ {
		for _, v := range adj[u] {
			inDeg[v]++
		}
	}
	queue := make([]int, 0, total)
	for u := 0; u < total; u++ {
		if inDeg[u] == 0 {
			queue = append(queue, u)
		}
	}
	ord := make([]int, 0, total)
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		ord = append(ord, u)
		for _, v := range adj[u] {
			inDeg[v]--
			if inDeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	if len(ord) != total {
		fmt.Println(-1)
		return
	}
	// assign answers
	ans := make([]int, n)
	timeToNode := make([]int, n)
	for i := 0; i < n; i++ {
		timeToNode[hoc.in[i]] = i
	}
	ptr := 0
	for _, u := range ord {
		if leafToTime[u] != -1 {
			ptr++
			node := timeToNode[leafToTime[u]]
			ans[node] = ptr
		}
	}
	// output
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteString(" ")
		}
		writer.WriteString(fmt.Sprintf("%d", ans[i]))
	}
	writer.WriteByte('\n')
}
