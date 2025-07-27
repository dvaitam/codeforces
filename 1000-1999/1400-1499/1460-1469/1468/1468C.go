package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// customer holds the estimation and arrival id
type customer struct {
	money int
	id    int
}

// maxHeap orders customers by money descending, then by id ascending
// (earlier arrivals preferred on tie)
type maxHeap []customer

func (h maxHeap) Len() int { return len(h) }
func (h maxHeap) Less(i, j int) bool {
	if h[i].money == h[j].money {
		return h[i].id < h[j].id
	}
	return h[i].money > h[j].money
}
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(customer)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	served := make([]bool, q+5)
	queue := make([]int, 0, q)
	front := 0
	pq := &maxHeap{}
	heap.Init(pq)
	id := 0
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(reader, &t)
		switch t {
		case 1:
			var m int
			fmt.Fscan(reader, &m)
			id++
			queue = append(queue, id)
			heap.Push(pq, customer{money: m, id: id})
		case 2:
			for served[queue[front]] {
				front++
			}
			cid := queue[front]
			front++
			served[cid] = true
			fmt.Fprintln(writer, cid)
		case 3:
			for pq.Len() > 0 && served[(*pq)[0].id] {
				heap.Pop(pq)
			}
			if pq.Len() > 0 {
				c := heap.Pop(pq).(customer)
				served[c.id] = true
				fmt.Fprintln(writer, c.id)
			}
		}
	}
}
