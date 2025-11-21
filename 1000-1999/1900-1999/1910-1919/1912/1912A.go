package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type task struct {
	cost int64
	gain int64
}

type maxHeap []int64

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}
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

	var x int64
	var k int
	fmt.Fscan(in, &x, &k)

	tasks := make([]task, 0)

	for i := 0; i < k; i++ {
		var l int
		fmt.Fscan(in, &l)
		sum := int64(0)
		minPref := int64(0)
		for j := 0; j < l; j++ {
			var val int64
			fmt.Fscan(in, &val)
			sum += val
			if sum > 0 {
				cost := -minPref
				if cost < 0 {
					cost = 0
				}
				tasks = append(tasks, task{cost: cost, gain: sum})
				sum = 0
				minPref = 0
			} else if sum < minPref {
				minPref = sum
			}
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].cost == tasks[j].cost {
			return tasks[i].gain > tasks[j].gain
		}
		return tasks[i].cost < tasks[j].cost
	})

	h := &maxHeap{}
	heap.Init(h)

	idx := 0
	for idx < len(tasks) && tasks[idx].cost <= x {
		heap.Push(h, tasks[idx].gain)
		idx++
	}

	for h.Len() > 0 {
		x += heap.Pop(h).(int64)
		for idx < len(tasks) && tasks[idx].cost <= x {
			heap.Push(h, tasks[idx].gain)
			idx++
		}
	}

	fmt.Fprintln(out, x)
}
