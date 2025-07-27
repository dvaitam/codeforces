package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Item struct {
	val int
	id  int
}

type MaxHeap []Item

func (h MaxHeap) Len() int { return len(h) }
func (h MaxHeap) Less(i, j int) bool {
	if h[i].val != h[j].val {
		return h[i].val > h[j].val
	}
	return h[i].id < h[j].id
}
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

var (
	n, q  int
	maxW  = 100000
	w     []int
	v     []int
	cnt   []int64
	heaps []MaxHeap
	seg   *SegTree
)

type SegTree struct {
	size int
	tree []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	return &SegTree{size: size, tree: make([]int, 2*size)}
}

func better(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	if v[a] != v[b] {
		if v[a] > v[b] {
			return a
		}
		return b
	}
	if w[a] != w[b] {
		if w[a] < w[b] {
			return a
		}
		return b
	}
	if a < b {
		return a
	}
	return b
}

func (t *SegTree) Update(pos int, id int) {
	p := pos + t.size - 1
	t.tree[p] = id
	for p > 1 {
		p >>= 1
		t.tree[p] = better(t.tree[p<<1], t.tree[p<<1|1])
	}
}

func (t *SegTree) Query(l, r int) int {
	if l > r {
		return 0
	}
	l += t.size - 1
	r += t.size - 1
	resL, resR := 0, 0
	for l <= r {
		if l&1 == 1 {
			resL = better(resL, t.tree[l])
			l++
		}
		if r&1 == 0 {
			resR = better(t.tree[r], resR)
			r--
		}
		l >>= 1
		r >>= 1
	}
	return better(resL, resR)
}

func updateWeight(wt int) {
	h := &heaps[wt]
	for h.Len() > 0 {
		top := (*h)[0]
		if cnt[top.id] > 0 {
			break
		}
		heap.Pop(h)
	}
	id := 0
	if h.Len() > 0 {
		id = (*h)[0].id
	}
	seg.Update(wt, id)
}

func addDiamonds(id int, delta int64) {
	cnt[id] += delta
	heap.Push(&heaps[w[id]], Item{val: v[id], id: id})
	updateWeight(w[id])
}

func queryBag(cap int64) int64 {
	var ans int64
	type change struct {
		id int
		d  int64
	}
	mods := make([]change, 0)
	for cap > 0 {
		up := int(cap)
		if up > maxW {
			up = maxW
		}
		id := seg.Query(1, up)
		if id == 0 {
			break
		}
		if cnt[id] == 0 || int64(w[id]) > cap {
			updateWeight(w[id])
			continue
		}
		take := cap / int64(w[id])
		if take > cnt[id] {
			take = cnt[id]
		}
		if take == 0 {
			updateWeight(w[id])
			continue
		}
		cnt[id] -= take
		mods = append(mods, change{id, take})
		heap.Push(&heaps[w[id]], Item{val: v[id], id: id})
		updateWeight(w[id])
		ans += take * int64(v[id])
		cap -= take * int64(w[id])
	}
	for _, c := range mods {
		cnt[c.id] += c.d
		heap.Push(&heaps[w[c.id]], Item{val: v[c.id], id: c.id})
		updateWeight(w[c.id])
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &q)
	w = make([]int, n+1)
	v = make([]int, n+1)
	cnt = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &cnt[i], &w[i], &v[i])
	}
	heaps = make([]MaxHeap, maxW+1)
	seg = NewSegTree(maxW + 2)
	for i := 1; i <= n; i++ {
		if cnt[i] > 0 {
			heap.Push(&heaps[w[i]], Item{val: v[i], id: i})
		}
	}
	for wt := 1; wt <= maxW; wt++ {
		updateWeight(wt)
	}
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var k, d int64
			fmt.Fscan(reader, &k, &d)
			addDiamonds(int(d), k)
		} else if t == 2 {
			var k, d int64
			fmt.Fscan(reader, &k, &d)
			addDiamonds(int(d), -k)
		} else {
			var c int64
			fmt.Fscan(reader, &c)
			fmt.Fprintln(writer, queryBag(c))
		}
	}
}
