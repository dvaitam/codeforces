package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// Item represents a table release time and its identifier
type Item struct {
	release int // time when table becomes free
	id      int // table id
}

// PriorityQueue implements heap.Interface based on Item.release (min-heap)
type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].release < pq[j].release
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Item))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N, M, D int
	fmt.Fscan(reader, &N, &M, &D)
	arr := make([]struct{ t, idx int }, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(reader, &arr[i].t)
		arr[i].idx = i
	}
	// sort by arrival time
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].t < arr[j].t
	})
	ans := make([]int, N)
	var pq PriorityQueue
	heap.Init(&pq)
	tables := 0
	for _, e := range arr {
		t, i := e.t, e.idx
		// if a table is free (release time <= arrival)
		if pq.Len() > 0 && pq[0].release <= t {
			item := heap.Pop(&pq).(Item)
			ans[i] = item.id
			// allocate same table with new release time
			heap.Push(&pq, Item{release: t + D + 1, id: item.id})
		} else {
			// need new table
			tables++
			ans[i] = tables
			heap.Push(&pq, Item{release: t + D + 1, id: tables})
		}
	}
	// output
	fmt.Fprintln(writer, tables)
	for i, v := range ans {
		if i+1 < N {
			fmt.Fprintf(writer, "%d ", v)
		} else {
			fmt.Fprintf(writer, "%d\n", v)
		}
	}
}
