package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type fenwick struct {
	n   int
	bit []int64
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, bit: make([]int64, n+1)}
}

func (f *fenwick) add(idx int, delta int64) {
	for idx <= f.n {
		f.bit[idx] += delta
		idx += idx & -idx
	}
}

func (f *fenwick) sum(idx int) int64 {
	var res int64
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

type heapItem struct {
	need int64
	id   int
}

type maxHeap []heapItem

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].need > h[j].need }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(heapItem)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	inQueue := make([]bool, n+1)
	inStack := make([]bool, n+1)
	shared := make([]bool, n+1)
	queueTime := make([]int, n+1)
	stackTime := make([]int, n+1)

	fq := newFenwick(n)
	fs := newFenwick(n)
	var totalStack int64

	h := &maxHeap{}
	heap.Init(h)

	computeNeed := func(id int) int64 {
		qt := queueTime[id]
		st := stackTime[id]
		qFinish := fq.sum(qt)
		stackAbove := totalStack - fs.sum(st)
		return qFinish - stackAbove
	}

	insertShared := func(id int) {
		if shared[id] {
			return
		}
		shared[id] = true
		need := computeNeed(id)
		heap.Push(h, heapItem{need: need, id: id})
	}

	for i := 1; i <= n; i++ {
		var t, x int
		fmt.Fscan(in, &t, &x)
		switch t {
		case 1:
			inQueue[x] = true
			queueTime[x] = i
			fq.add(i, a[x])
			if inStack[x] {
				insertShared(x)
			}
		case 2:
			inStack[x] = true
			stackTime[x] = i
			fs.add(i, a[x])
			totalStack += a[x]
			if inQueue[x] {
				insertShared(x)
			}
		case 3:
			// move from queue to stack; removes from queue
			fq.add(queueTime[x], -a[x])
			inQueue[x] = false
			inStack[x] = true
			stackTime[x] = queueTime[x]
			fs.add(stackTime[x], a[x])
			totalStack += a[x]
		}

		for h.Len() > 0 {
			top := (*h)[0]
			actual := computeNeed(top.id)
			if actual < top.need {
				heap.Pop(h)
				heap.Push(h, heapItem{need: actual, id: top.id})
			} else {
				break
			}
		}

		if h.Len() == 0 || (*h)[0].need <= 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
