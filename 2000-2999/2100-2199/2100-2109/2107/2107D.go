package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type component struct {
	d    int
	u    int
	v    int
	path []int
	idx  int
}

type componentHeap []*component

func (h componentHeap) Len() int { return len(h) }

func (h componentHeap) Less(i, j int) bool {
	if h[i].d != h[j].d {
		return h[i].d > h[j].d
	}
	if h[i].u != h[j].u {
		return h[i].u > h[j].u
	}
	if h[i].v != h[j].v {
		return h[i].v > h[j].v
	}
	return h[i].idx < h[j].idx
}

func (h componentHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *componentHeap) Push(x interface{}) {
	*h = append(*h, x.(*component))
}

func (h *componentHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

type solver struct {
	n       int
	adj     [][]int
	alive   []bool
	visit   []int
	depth   []int
	parent  []int
	mark    int
	nextIdx int
}

func newSolver(n int, adj [][]int) *solver {
	al := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		al[i] = true
	}
	return &solver{
		n:      n,
		adj:    adj,
		alive:  al,
		visit:  make([]int, n+1),
		depth:  make([]int, n+1),
		parent: make([]int, n+1),
	}
}

func (s *solver) farthestNode(start int) int {
	s.mark++
	q := []int{start}
	s.visit[start] = s.mark
	s.depth[start] = 0
	best := start
	for head := 0; head < len(q); head++ {
		v := q[head]
		if s.depth[v] > s.depth[best] || (s.depth[v] == s.depth[best] && v > best) {
			best = v
		}
		for _, to := range s.adj[v] {
			if !s.alive[to] || s.visit[to] == s.mark {
				continue
			}
			s.visit[to] = s.mark
			s.depth[to] = s.depth[v] + 1
			q = append(q, to)
		}
	}
	return best
}

func (s *solver) farthestNodeAndPath(start int) (int, []int) {
	s.mark++
	q := []int{start}
	s.visit[start] = s.mark
	s.depth[start] = 0
	s.parent[start] = 0
	best := start
	for head := 0; head < len(q); head++ {
		v := q[head]
		if s.depth[v] > s.depth[best] || (s.depth[v] == s.depth[best] && v > best) {
			best = v
		}
		for _, to := range s.adj[v] {
			if !s.alive[to] || s.visit[to] == s.mark {
				continue
			}
			s.visit[to] = s.mark
			s.depth[to] = s.depth[v] + 1
			s.parent[to] = v
			q = append(q, to)
		}
	}
	path := []int{}
	cur := best
	for cur != 0 {
		path = append(path, cur)
		if cur == start {
			break
		}
		cur = s.parent[cur]
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return best, path
}

func (s *solver) buildComponent(start int) *component {
	if !s.alive[start] {
		return nil
	}
	endA := s.farthestNode(start)
	endB, path := s.farthestNodeAndPath(endA)
	u, v := endA, endB
	if v > u {
		u, v = v, u
	}
	s.nextIdx++
	return &component{
		d:    len(path),
		u:    u,
		v:    v,
		path: path,
		idx:  s.nextIdx,
	}
}

func solveCase(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	s := newSolver(n, adj)
	pq := &componentHeap{}
	heap.Init(pq)
	initial := s.buildComponent(1)
	heap.Push(pq, initial)
	ans := make([]int, 0, n*3)
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*component)
		ans = append(ans, cur.d, cur.u, cur.v)
		for _, node := range cur.path {
			s.alive[node] = false
		}
		for _, node := range cur.path {
			for _, to := range s.adj[node] {
				if !s.alive[to] {
					continue
				}
				if comp := s.buildComponent(to); comp != nil {
					heap.Push(pq, comp)
				}
			}
		}
	}
	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solveCase(reader, writer)
	}
}
