package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const INF = int64(1) << 60

// IntHeap is a max-heap of int64
type IntHeap []int64

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// MultiSet supports insert, delete by value, and get max
type MultiSet struct {
	maxH *IntHeap
	delH *IntHeap
}

func NewMultiSet() *MultiSet {
	mh := &IntHeap{}
	dh := &IntHeap{}
	heap.Init(mh)
	heap.Init(dh)
	return &MultiSet{mh, dh}
}
func (ms *MultiSet) Insert(v int64) {
	heap.Push(ms.maxH, v)
}
func (ms *MultiSet) Delete(v int64) {
	heap.Push(ms.delH, v)
}
func (ms *MultiSet) Top() int64 {
	for ms.maxH.Len() > 0 && ms.delH.Len() > 0 && (*ms.maxH)[0] == (*ms.delH)[0] {
		heap.Pop(ms.maxH)
		heap.Pop(ms.delH)
	}
	if ms.maxH.Len() == 0 {
		return -INF
	}
	return (*ms.maxH)[0]
}

// SegmentTree for range max
type SegmentTree struct {
	n int
	t []int64
}

func NewSegmentTree(sz int) *SegmentTree {
	n := 1
	for n < sz {
		n <<= 1
	}
	t := make([]int64, 2*n)
	for i := range t {
		t[i] = -INF
	}
	return &SegmentTree{n, t}
}
func (st *SegmentTree) Update(pos int, v int64) {
	i := pos + st.n
	st.t[i] = v
	for i >>= 1; i > 0; i >>= 1 {
		if st.t[2*i] > st.t[2*i+1] {
			st.t[i] = st.t[2*i]
		} else {
			st.t[i] = st.t[2*i+1]
		}
	}
}

// Query max in [l,r]
func (st *SegmentTree) Query(l, r int) int64 {
	res := -INF
	l += st.n
	r += st.n
	for l <= r {
		if (l & 1) == 1 {
			if st.t[l] > res {
				res = st.t[l]
			}
			l++
		}
		if (r & 1) == 0 {
			if st.t[r] > res {
				res = st.t[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, m int
	fmt.Fscan(reader, &n, &m)
	names := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &names[i])
	}
	// Build trie
	alpha := 26
	trie := [][]int{{}}
	trie[0] = make([]int, alpha)
	for i := range trie[0] {
		trie[0][i] = -1
	}
	patNode := make([]int, n)
	for i, s := range names {
		u := 0
		for _, ch := range s {
			c := int(ch - 'a')
			if trie[u][c] == -1 {
				trie = append(trie, make([]int, alpha))
				for j := range trie[len(trie)-1] {
					trie[len(trie)-1][j] = -1
				}
				trie[u][c] = len(trie) - 1
			}
			u = trie[u][c]
		}
		patNode[i] = u
	}
	sz := len(trie)
	// build failure links and go
	fail := make([]int, sz)
	goTo := make([]int, sz*alpha)
	for i := 0; i < alpha; i++ {
		if trie[0][i] != -1 {
			fail[trie[0][i]] = 0
			goTo[i] = trie[0][i]
		} else {
			goTo[i] = 0
		}
	}
	q := make([]int, 0, sz)
	for i := 0; i < alpha; i++ {
		if trie[0][i] != -1 {
			q = append(q, trie[0][i])
		}
	}
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for c := 0; c < alpha; c++ {
			v := trie[u][c]
			if v != -1 {
				fail[v] = goTo[fail[u]*alpha+c]
				goTo[u*alpha+c] = v
				q = append(q, v)
			} else {
				goTo[u*alpha+c] = goTo[fail[u]*alpha+c]
			}
		}
	}
	// failure tree
	adj := make([][]int, sz)
	for v := 1; v < sz; v++ {
		u := fail[v]
		adj[u] = append(adj[u], v)
	}
	// HLD prep: size & heavy
	size := make([]int, sz)
	heavy := make([]int, sz)
	for i := range heavy {
		heavy[i] = -1
	}
	// post-order
	order := make([]int, 0, sz)
	type stItem struct{ u, idx int }
	stack := []stItem{{0, 0}}
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		u := top.u
		if top.idx < len(adj[u]) {
			v := adj[u][top.idx]
			top.idx++
			stack = append(stack, stItem{v, 0})
		} else {
			order = append(order, u)
			stack = stack[:len(stack)-1]
		}
	}
	for _, u := range order {
		size[u] = 1
		maxSz := 0
		for _, v := range adj[u] {
			size[u] += size[v]
			if size[v] > maxSz {
				maxSz = size[v]
				heavy[u] = v
			}
		}
	}
	// dfs2 for pos and head
	head := make([]int, sz)
	pos := make([]int, sz)
	curPos := 0
	type hp struct{ u, h int }
	stk := []hp{{0, 0}}
	for len(stk) > 0 {
		x := stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		u, h := x.u, x.h
		head[u] = h
		pos[u] = curPos
		curPos++
		// push light children
		for i := len(adj[u]) - 1; i >= 0; i-- {
			v := adj[u][i]
			if v != heavy[u] {
				stk = append(stk, hp{v, v})
			}
		}
		// then heavy child
		if heavy[u] != -1 {
			stk = append(stk, hp{heavy[u], h})
		}
	}
	// multiset per node
	counts := make([]int, sz)
	for _, u := range patNode {
		counts[u]++
	}
	msArr := make([]*MultiSet, sz)
	nodeVal := make([]int64, sz)
	for u := 0; u < sz; u++ {
		if counts[u] > 0 {
			ms := NewMultiSet()
			for i := 0; i < counts[u]; i++ {
				ms.Insert(0)
			}
			msArr[u] = ms
			nodeVal[u] = 0
		} else {
			nodeVal[u] = -INF
		}
	}
	// segment tree
	st := NewSegmentTree(sz)
	for u := 0; u < sz; u++ {
		st.Update(pos[u], nodeVal[u])
	}
	// pattern values
	patVal := make([]int64, n)
	for i := range patVal {
		patVal[i] = 0
	}
	// process queries
	for qi := 0; qi < m; qi++ {
		var tp int
		fmt.Fscan(reader, &tp)
		if tp == 1 {
			var idx int
			var x int64
			fmt.Fscan(reader, &idx, &x)
			idx--
			u := patNode[idx]
			old := patVal[idx]
			patVal[idx] = x
			ms := msArr[u]
			ms.Delete(old)
			ms.Insert(x)
			top := ms.Top()
			if top != nodeVal[u] {
				nodeVal[u] = top
				st.Update(pos[u], top)
			}
		} else {
			var qstr string
			fmt.Fscan(reader, &qstr)
			cur := 0
			ans := -INF
			for _, ch := range qstr {
				c := int(ch - 'a')
				cur = goTo[cur*alpha+c]
				// query path from cur to root
				u := cur
				for u != 0 {
					h := head[u]
					res := st.Query(pos[h], pos[u])
					if res > ans {
						ans = res
					}
					u = fail[h]
				}
				// include root
				if nodeVal[0] > ans {
					ans = nodeVal[0]
				}
			}
			if ans < 0 {
				fmt.Fprintln(writer, -1)
			} else {
				fmt.Fprintln(writer, ans)
			}
		}
	}
}
