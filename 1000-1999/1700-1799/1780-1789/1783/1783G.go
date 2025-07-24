package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const INF int = 1 << 60
const NEG_INF int = -INF

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{}   { old := *h; x := old[len(old)-1]; *h = old[:len(old)-1]; return x }

type pair struct{ val, id int }
type PairHeap []pair

func (h PairHeap) Len() int { return len(h) }
func (h PairHeap) Less(i, j int) bool {
	if h[i].val == h[j].val {
		return h[i].id < h[j].id
	}
	return h[i].val > h[j].val
}
func (h PairHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PairHeap) Push(x interface{}) { *h = append(*h, x.(pair)) }
func (h *PairHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// multiset with max retrieval
type SubHeap struct {
	h   IntHeap
	cnt map[int]int
}

func NewSubHeap() *SubHeap { return &SubHeap{cnt: make(map[int]int)} }
func (s *SubHeap) push(v int) {
	heap.Push(&s.h, v)
	s.cnt[v]++
}
func (s *SubHeap) remove(v int) {
	if c, ok := s.cnt[v]; ok {
		if c == 1 {
			delete(s.cnt, v)
		} else {
			s.cnt[v] = c - 1
		}
	}
}
func (s *SubHeap) top() int {
	for s.h.Len() > 0 {
		v := s.h[0]
		if s.cnt[v] > 0 {
			return v
		}
		heap.Pop(&s.h)
	}
	return NEG_INF
}

type CentroidData struct {
	weight int
	subs   []*SubHeap
	subTop []int
	global PairHeap
	cand   int
}

func (cd *CentroidData) updateSub(id int, oldVal, newVal int) {
	sh := cd.subs[id]
	if oldVal != NEG_INF {
		sh.remove(oldVal)
	}
	if newVal != NEG_INF {
		sh.push(newVal)
	}
	top := sh.top()
	if top != cd.subTop[id] {
		cd.subTop[id] = top
		heap.Push(&cd.global, pair{top, id})
	}
}

func (cd *CentroidData) getTop() (int, int) {
	for cd.global.Len() > 0 {
		p := cd.global[0]
		if cd.subTop[p.id] == p.val {
			return p.val, p.id
		}
		heap.Pop(&cd.global)
	}
	return NEG_INF, -1
}

func (cd *CentroidData) recompute() {
	v1, id1 := cd.getTop()
	if id1 == -1 {
		cd.cand = 2 * cd.weight
		return
	}
	heap.Pop(&cd.global)
	v2, _ := cd.getTop()
	heap.Push(&cd.global, pair{v1, id1})
	cand := 2 * cd.weight
	if v1 != NEG_INF && cand < cd.weight+v1 {
		cand = cd.weight + v1
	}
	if v2 != NEG_INF && cand < v1+v2 {
		cand = v1 + v2
	}
	cd.cand = cand
}

var (
	n         int
	g         [][]int
	weight    []int
	nodePaths [][]info
	centroid  []*CentroidData
	candHeap  PairHeap
	candVal   []int
	used      []bool
	sz        []int
)

type info struct{ c, sub, dist int }

func dfsSize(v, p int) int {
	sz[v] = 1
	for _, to := range g[v] {
		if to == p || used[to] {
			continue
		}
		sz[v] += dfsSize(to, v)
	}
	return sz[v]
}

func dfsCentroid(v, p, total int) int {
	for _, to := range g[v] {
		if to != p && !used[to] && sz[to] > total/2 {
			return dfsCentroid(to, v, total)
		}
	}
	return v
}

func dfsCollect(v, p, dist, cent, sub int) {
	nodePaths[v] = append(nodePaths[v], info{cent, sub, dist})
	for _, to := range g[v] {
		if to == p || used[to] {
			continue
		}
		dfsCollect(to, v, dist+1, cent, sub)
	}
}

func build(v, p int) {
	total := dfsSize(v, -1)
	c := dfsCentroid(v, -1, total)
	used[c] = true
	parent := p
	_ = parent
	// count subtrees
	cnt := 0
	for _, to := range g[c] {
		if !used[to] {
			dfsCollect(to, c, 1, c, cnt)
			cnt++
		}
	}
	centroid[c] = &CentroidData{weight: weight[c], subs: make([]*SubHeap, cnt), subTop: make([]int, cnt)}
	for i := 0; i < cnt; i++ {
		centroid[c].subs[i] = NewSubHeap()
		centroid[c].subTop[i] = NEG_INF
	}
	nodePaths[c] = append(nodePaths[c], info{c, -1, 0})
	for _, to := range g[c] {
		if !used[to] {
			build(to, c)
		}
	}
}

func initStructures() {
	for v := 0; v < n; v++ {
		for _, it := range nodePaths[v] {
			c := it.c
			if it.sub != -1 {
				val := weight[v] + it.dist
				centroid[c].subs[it.sub].push(val)
			}
		}
	}
	candVal = make([]int, n)
	for c := 0; c < n; c++ {
		if centroid[c] == nil {
			continue
		}
		for i := range centroid[c].subs {
			top := centroid[c].subs[i].top()
			centroid[c].subTop[i] = top
			heap.Push(&centroid[c].global, pair{top, i})
		}
		centroid[c].recompute()
		candVal[c] = centroid[c].cand
		heap.Push(&candHeap, pair{candVal[c], c})
	}
}

func updateNode(v, newW int) {
	old := weight[v]
	if old == newW { // still need to recompute for centroid if weight changes? but if equal no change
	}
	weight[v] = newW
	for _, it := range nodePaths[v] {
		c := it.c
		if it.sub == -1 {
			centroid[c].weight = newW
		} else {
			oldVal := old + it.dist
			newVal := newW + it.dist
			centroid[c].updateSub(it.sub, oldVal, newVal)
		}
	}
	for _, it := range nodePaths[v] {
		c := it.c
		centroid[c].recompute()
		candVal[c] = centroid[c].cand
		heap.Push(&candHeap, pair{candVal[c], c})
	}
}

func globalD() int {
	for candHeap.Len() > 0 {
		p := candHeap[0]
		if candVal[p.id] == p.val {
			return p.val
		}
		heap.Pop(&candHeap)
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	weight = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &weight[i])
	}
	g = make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	nodePaths = make([][]info, n)
	centroid = make([]*CentroidData, n)
	used = make([]bool, n)
	sz = make([]int, n)
	build(0, -1)
	initStructures()

	var m int
	fmt.Fscan(in, &m)
	for ; m > 0; m-- {
		var v, x int
		fmt.Fscan(in, &v, &x)
		v--
		updateNode(v, x)
		ans := globalD()
		fmt.Fprintln(out, (ans+1)/2)
	}
}
