package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Segment struct {
	c int
	l int
	r int
}

type DSU struct {
	parent []int
	size   []int
	maxR   [][2]int
}

func NewDSU(n int, segs []Segment) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n), maxR: make([][2]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
		if segs[i].c == 0 {
			d.maxR[i][0] = segs[i].r
			d.maxR[i][1] = -1
		} else {
			d.maxR[i][1] = segs[i].r
			d.maxR[i][0] = -1
		}
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) int {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return a
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
	if d.maxR[a][0] < d.maxR[b][0] {
		d.maxR[a][0] = d.maxR[b][0]
	}
	if d.maxR[a][1] < d.maxR[b][1] {
		d.maxR[a][1] = d.maxR[b][1]
	}
	return a
}

type Item struct {
	r  int
	id int
}

type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].r > h[j].r }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func clean(h *MaxHeap, d *DSU, color int, thresh int) {
	for h.Len() > 0 {
		it := (*h)[0]
		root := d.Find(it.id)
		val := d.maxR[root][color]
		if it.r != val {
			heap.Pop(h)
			continue
		}
		if val < thresh {
			heap.Pop(h)
			continue
		}
		break
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		segs := make([]Segment, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &segs[i].c, &segs[i].l, &segs[i].r)
		}
		order := make([]int, n)
		for i := range order {
			order[i] = i
		}
		sort.Slice(order, func(i, j int) bool { return segs[order[i]].l < segs[order[j]].l })

		d := NewDSU(n, segs)
		var h0, h1 MaxHeap
		heap.Init(&h0)
		heap.Init(&h1)

		for _, idx := range order {
			seg := segs[idx]
			l := seg.l
			r := seg.r
			c := seg.c
			clean(&h0, d, 0, l)
			clean(&h1, d, 1, l)

			root := d.Find(idx)
			oppHeap := &h1
			oppColor := 1
			ownHeap := &h0
			if c == 1 {
				oppHeap = &h0
				oppColor = 0
				ownHeap = &h1
			}
			for oppHeap.Len() > 0 {
				it := (*oppHeap)[0]
				rootOpp := d.Find(it.id)
				val := d.maxR[rootOpp][oppColor]
				if it.r != val {
					heap.Pop(oppHeap)
					continue
				}
				if val < l {
					break
				}
				heap.Pop(oppHeap)
				root = d.Union(root, rootOpp)
			}
			root = d.Find(root)
			if d.maxR[root][c] < r {
				d.maxR[root][c] = r
			}
			heap.Push(ownHeap, Item{r: d.maxR[root][c], id: root})
			if d.maxR[root][oppColor] >= l && d.maxR[root][oppColor] >= 0 {
				heap.Push(oppHeap, Item{r: d.maxR[root][oppColor], id: root})
			}
		}

		comps := make(map[int]struct{})
		for i := 0; i < n; i++ {
			comps[d.Find(i)] = struct{}{}
		}
		fmt.Fprintln(writer, len(comps))
	}
}
