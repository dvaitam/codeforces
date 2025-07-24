package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// MaxHeap implements a max-heap for int64 values.
type MaxHeap []int64

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	v := old[len(old)-1]
	*h = old[:len(old)-1]
	return v
}

type month struct {
	val int64
	idx int
}

type task struct {
	dl   int
	cost int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int64, n)
	prefix := make([]int64, n+1)
	months := make([]month, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		prefix[i+1] = prefix[i] + a[i]
		months[i] = month{val: a[i], idx: i + 1}
	}

	b := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &b[i])
	}
	p := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &p[i])
	}

	sort.Slice(months, func(i, j int) bool { return months[i].val < months[j].val })
	suf := make([]int, n)
	maxIdx := 0
	for i := n - 1; i >= 0; i-- {
		if months[i].idx > maxIdx {
			maxIdx = months[i].idx
		}
		suf[i] = maxIdx
	}

	maxA := months[n-1].val
	tasks := make([]task, 0, m)
	for i := 0; i < m; i++ {
		if b[i] > maxA {
			continue
		}
		pos := sort.Search(n, func(j int) bool { return months[j].val >= b[i] })
		dl := suf[pos]
		tasks = append(tasks, task{dl: dl, cost: p[i]})
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].dl == tasks[j].dl {
			return tasks[i].cost < tasks[j].cost
		}
		return tasks[i].dl < tasks[j].dl
	})

	mh := &MaxHeap{}
	heap.Init(mh)
	var total int64
	for _, t := range tasks {
		heap.Push(mh, t.cost)
		total += t.cost
		for total > prefix[t.dl] {
			v := heap.Pop(mh).(int64)
			total -= v
		}
	}
	fmt.Fprintln(out, mh.Len())
}
