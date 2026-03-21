package main

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strconv"
)

type DSU struct {
	p  []int
	sz []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	sz := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{p: p, sz: sz}
}

func (d *DSU) Find(x int) int {
	for d.p[x] != x {
		d.p[x] = d.p[d.p[x]]
		x = d.p[x]
	}
	return x
}

func (d *DSU) Union(a, b int) bool {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return false
	}
	if d.sz[a] < d.sz[b] {
		a, b = b, a
	}
	d.p[b] = a
	d.sz[a] += d.sz[b]
	return true
}

type FixedEdge struct {
	u, v int
	w    int64
	idx  int
}

type TreeEdge struct {
	to   int
	free bool
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	ptr := 0
	nextInt := func() int {
		for ptr < len(data) && (data[ptr] < '0' || data[ptr] > '9') {
			ptr++
		}
		val := 0
		for ptr < len(data) && data[ptr] >= '0' && data[ptr] <= '9' {
			val = val*10 + int(data[ptr]-'0')
			ptr++
		}
		return val
	}

	n := nextInt()
	m := nextInt()

	adj := make([][]int, n+1)
	fixed := make([]FixedEdge, m)
	var x int64

	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		w := int64(nextInt())
		fixed[i] = FixedEdge{u: u, v: v, w: w, idx: i}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		x ^= w
	}

	next := make([]int, n+1)
	prev := make([]int, n+1)
	inList := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		prev[i] = i - 1
		next[i] = i + 1
		inList[i] = true
	}
	next[n] = 0
	head := 1

	remove := func(v int) {
		if !inList[v] {
			return
		}
		p := prev[v]
		nx := next[v]
		if p == 0 {
			head = nx
		} else {
			next[p] = nx
		}
		if nx != 0 {
			prev[nx] = p
		}
		inList[v] = false
	}

	tag := make([]int, n+1)
	timer := 0
	compCnt := 0
	freeTree := make([][2]int, 0, n-1)
	queue := make([]int, 0, n)

	for s := 1; s <= n; s++ {
		if !inList[s] {
			continue
		}
		compCnt++
		remove(s)
		queue = queue[:0]
		queue = append(queue, s)

		for qi := 0; qi < len(queue); qi++ {
			v := queue[qi]
			timer++
			for _, to := range adj[v] {
				tag[to] = timer
			}
			for cur := head; cur != 0; {
				nx := next[cur]
				if tag[cur] != timer {
					remove(cur)
					freeTree = append(freeTree, [2]int{v, cur})
					queue = append(queue, cur)
				}
				cur = nx
			}
		}
	}

	totalFree := int64(n)*(int64(n)-1)/2 - int64(m)
	hasCycleInFree := totalFree > int64(len(freeTree))

	tree := make([][]TreeEdge, n+1)
	ds := NewDSU(n)

	for _, e := range freeTree {
		u, v := e[0], e[1]
		ds.Union(u, v)
		tree[u] = append(tree[u], TreeEdge{to: v, free: true})
		tree[v] = append(tree[v], TreeEdge{to: u, free: true})
	}

	sort.Slice(fixed, func(i, j int) bool {
		if fixed[i].w == fixed[j].w {
			return fixed[i].idx < fixed[j].idx
		}
		return fixed[i].w < fixed[j].w
	})

	selected := make([]bool, m)
	var base int64

	for _, e := range fixed {
		if ds.Union(e.u, e.v) {
			selected[e.idx] = true
			base += e.w
			tree[e.u] = append(tree[e.u], TreeEdge{to: e.v, free: false})
			tree[e.v] = append(tree[e.v], TreeEdge{to: e.u, free: false})
		}
	}

	if hasCycleInFree || x == 0 {
		out := bufio.NewWriterSize(os.Stdout, 1<<20)
		out.WriteString(strconv.FormatInt(base, 10))
		out.WriteByte('\n')
		out.Flush()
		return
	}

	const LOG = 20
	up := make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, n+1)
	}
	depth := make([]int, n+1)
	prefFree := make([]int, n+1)

	order := make([]int, 0, n)
	order = append(order, 1)
	for i := 0; i < len(order); i++ {
		v := order[i]
		for _, e := range tree[v] {
			if e.to == up[0][v] {
				continue
			}
			up[0][e.to] = v
			depth[e.to] = depth[v] + 1
			prefFree[e.to] = prefFree[v]
			if e.free {
				prefFree[e.to]++
			}
			order = append(order, e.to)
		}
	}

	for j := 1; j < LOG; j++ {
		for v := 1; v <= n; v++ {
			up[j][v] = up[j-1][up[j-1][v]]
		}
	}

	lca := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		bit := 0
		for diff > 0 {
			if diff&1 == 1 {
				a = up[bit][a]
			}
			diff >>= 1
			bit++
		}
		if a == b {
			return a
		}
		for j := LOG - 1; j >= 0; j-- {
			if up[j][a] != up[j][b] {
				a = up[j][a]
				b = up[j][b]
			}
		}
		return up[0][a]
	}

	const INF int64 = 1 << 60
	best := INF
	for _, e := range fixed {
		if selected[e.idx] {
			continue
		}
		l := lca(e.u, e.v)
		if prefFree[e.u]+prefFree[e.v]-2*prefFree[l] > 0 && e.w < best {
			best = e.w
		}
	}

	if best > x {
		best = x
	}

	ans := base + best

	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	out.WriteString(strconv.FormatInt(ans, 10))
	out.WriteByte('\n')
	out.Flush()
}
