package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Flight represents a flight with its delay cost and index
type Flight struct {
	cost int64
	idx  int
}

type PQ []Flight

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].cost > pq[j].cost }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Flight)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	costs := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &costs[i])
	}

	pq := &PQ{}
	heap.Init(pq)
	ans := make([]int, n+1)
	var total int64
	idx := 1
	for t := k + 1; t <= k+n; t++ {
		for idx <= n && idx <= t {
			heap.Push(pq, Flight{costs[idx], idx})
			idx++
		}
		if pq.Len() == 0 {
			continue
		}
		f := heap.Pop(pq).(Flight)
		ans[f.idx] = t
		total += int64(t-f.idx) * f.cost
	}

	fmt.Fprintln(writer, total)
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
}
