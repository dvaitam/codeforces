package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type TaskHeap []Task

type Task struct {
	val int64
	id  int
	idx int
}

func (h TaskHeap) Len() int            { return len(h) }
func (h TaskHeap) Less(i, j int) bool  { return h[i].val < h[j].val }
func (h TaskHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *TaskHeap) Push(x interface{}) { *h = append(*h, x.(Task)) }
func (h *TaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	seqs := make([][]int64, n)
	total := 0

	for i := 0; i < n; i++ {
		var k int
		var a1, x, y, m int64
		fmt.Fscan(in, &k, &a1, &x, &y, &m)
		total += k
		seq := make([]int64, k)
		if k > 0 {
			seq[0] = a1
			for j := 1; j < k; j++ {
				seq[j] = (seq[j-1]*x + y) % m
			}
		}
		seqs[i] = seq
	}

	const limit = 200000
	var order []struct {
		val int64
		id  int
	}
	if total <= limit {
		order = make([]struct {
			val int64
			id  int
		}, 0, total)
	}

	h := &TaskHeap{}
	heap.Init(h)
	for i := 0; i < n; i++ {
		if len(seqs[i]) > 0 {
			heap.Push(h, Task{val: seqs[i][0], id: i, idx: 0})
		}
	}

	var prev int64
	first := true
	ans := int64(0)

	for h.Len() > 0 {
		node := heap.Pop(h).(Task)
		val := node.val
		if !first && val < prev {
			ans++
		}
		prev = val
		first = false

		if total <= limit {
			order = append(order, struct {
				val int64
				id  int
			}{val: val, id: node.id + 1})
		}

		nextIdx := node.idx + 1
		if nextIdx < len(seqs[node.id]) {
			heap.Push(h, Task{val: seqs[node.id][nextIdx], id: node.id, idx: nextIdx})
		}
	}

	fmt.Fprintln(out, ans)
	if total <= limit {
		for _, item := range order {
			fmt.Fprintf(out, "%d %d\n", item.val, item.id)
		}
	}
}
