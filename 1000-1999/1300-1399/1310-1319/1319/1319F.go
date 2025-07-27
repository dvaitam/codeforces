package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	to int
	w  int
}

const INF = int(1 << 60)

var (
	n     int
	g     [][]int
	LOG   int
	up    [][]int
	depth []int
	tin   []int
	tout  []int
	timer int
)

func buildLCA(root int) {
	LOG = 1
	for (1 << LOG) <= n {
		LOG++
	}
	up = make([][]int, LOG)
	for i := range up {
		up[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	tin = make([]int, n+1)
	tout = make([]int, n+1)

	parent := make([]int, n+1)
	parent[root] = root
	stack := []int{root}
	order := []int{root}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
			order = append(order, to)
		}
	}
	timer = 0
	type Node struct {
		v   int
		idx int
	}
	st := []Node{{root, 0}}
	for len(st) > 0 {
		cur := &st[len(st)-1]
		v := cur.v
		if cur.idx == 0 {
			tin[v] = timer
			timer++
		}
		if cur.idx < len(g[v]) {
			to := g[v][cur.idx]
			cur.idx++
			if to == parent[v] {
				continue
			}
			st = append(st, Node{to, 0})
		} else {
			tout[v] = timer
			timer++
			st = st[:len(st)-1]
		}
	}
	for i := 1; i <= n; i++ {
		up[0][i] = parent[i]
	}
	for k := 1; k < LOG; k++ {
		for i := 1; i <= n; i++ {
			up[k][i] = up[k-1][up[k-1][i]]
		}
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for k := 0; diff > 0; k++ {
		if diff&1 == 1 {
			a = up[k][a]
		}
		diff >>= 1
	}
	if a == b {
		return a
	}
	for k := LOG - 1; k >= 0; k-- {
		if up[k][a] != up[k][b] {
			a = up[k][a]
			b = up[k][b]
		}
	}
	return up[0][a]
}

func dist(a, b int) int {
	l := lca(a, b)
	return depth[a] + depth[b] - 2*depth[l]
}

func buildVirtualTree(nodes []int) ([]int, map[int][]Edge) {
	if len(nodes) == 0 {
		return nil, nil
	}
	sort.Slice(nodes, func(i, j int) bool { return tin[nodes[i]] < tin[nodes[j]] })
	// add LCAs
	m := len(nodes)
	for i := 0; i < m-1; i++ {
		l := lca(nodes[i], nodes[i+1])
		nodes = append(nodes, l)
	}
	sort.Slice(nodes, func(i, j int) bool { return tin[nodes[i]] < tin[nodes[j]] })
	uniq := nodes[:0]
	last := -1
	for _, v := range nodes {
		if last != v {
			uniq = append(uniq, v)
			last = v
		}
	}
	nodes = uniq
	adj := make(map[int][]Edge, len(nodes)*2)
	stack := []int{}
	stack = append(stack, nodes[0])
	for i := 1; i < len(nodes); i++ {
		v := nodes[i]
		l := lca(v, stack[len(stack)-1])
		for len(stack) >= 2 && depth[stack[len(stack)-2]] >= depth[l] {
			u := stack[len(stack)-1]
			p := stack[len(stack)-2]
			w := dist(u, p)
			adj[u] = append(adj[u], Edge{p, w})
			adj[p] = append(adj[p], Edge{u, w})
			stack = stack[:len(stack)-1]
		}
		if stack[len(stack)-1] != l {
			u := stack[len(stack)-1]
			w := dist(u, l)
			adj[u] = append(adj[u], Edge{l, w})
			adj[l] = append(adj[l], Edge{u, w})
			stack[len(stack)-1] = l
		}
		stack = append(stack, v)
	}
	for len(stack) > 1 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		p := stack[len(stack)-1]
		w := dist(u, p)
		adj[u] = append(adj[u], Edge{p, w})
		adj[p] = append(adj[p], Edge{u, w})
	}
	return nodes, adj
}

type Virus struct {
	v int
	s int
}

type Item struct {
	cycles int
	id     int
	node   int
	dist   int
}

// priority queue
type PQ []Item

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	if pq[i].cycles == pq[j].cycles {
		return pq[i].id < pq[j].id
	}
	return pq[i].cycles < pq[j].cycles
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func solveScenario(viruses []Virus, queries []int) []int {
	nodes := make([]int, 0, len(viruses)+len(queries))
	for _, v := range viruses {
		nodes = append(nodes, v.v)
	}
	for _, u := range queries {
		nodes = append(nodes, u)
	}
	nodes, adj := buildVirtualTree(nodes)
	idx := make(map[int]int, len(nodes))
	for i, v := range nodes {
		idx[v] = i
	}
	// build adjacency list using indices
	g2 := make([][]Edge, len(nodes))
	for v, es := range adj {
		i := idx[v]
		for _, e := range es {
			g2[i] = append(g2[i], Edge{idx[e.to], e.w})
		}
	}
	owner := make([]int, len(nodes))
	cycles := make([]int, len(nodes))
	distArr := make([]int, len(nodes))
	for i := range cycles {
		cycles[i] = INF
	}
	pq := &PQ{}
	heap.Init(pq)
	for id, vr := range viruses {
		pos := idx[vr.v]
		if cycles[pos] > 0 || (cycles[pos] == 0 && owner[pos] == 0) {
			owner[pos] = id + 1
			cycles[pos] = 0
			distArr[pos] = 0
			heap.Push(pq, Item{0, id + 1, pos, 0})
		}
	}
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.cycles != cycles[it.node] || it.id != owner[it.node] || it.dist != distArr[it.node] {
			continue
		}
		speed := viruses[it.id-1].s
		for _, e := range g2[it.node] {
			nd := it.dist + e.w
			nc := (nd + speed - 1) / speed
			if cycles[e.to] > nc || (cycles[e.to] == nc && owner[e.to] > it.id) {
				cycles[e.to] = nc
				owner[e.to] = it.id
				distArr[e.to] = nd
				heap.Push(pq, Item{nc, it.id, e.to, nd})
			}
		}
	}
	ans := make([]int, len(queries))
	for i, u := range queries {
		ans[i] = owner[idx[u]]
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fscan(in, &n)
	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	buildLCA(1)
	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var k, m int
		fmt.Fscan(in, &k, &m)
		viruses := make([]Virus, k)
		nodes := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &viruses[i].v, &viruses[i].s)
			nodes[i] = viruses[i].v
		}
		queries := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &queries[i])
		}
		ans := solveScenario(viruses, queries)
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
