package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Item struct {
	next int64
	id   int
}

type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].next < pq[j].next }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var h int64
		var n int
		fmt.Fscan(in, &h, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		c := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		pq := &PQ{}
		heap.Init(pq)
		for i := 0; i < n; i++ {
			heap.Push(pq, Item{next: 1, id: i})
		}
		var current int64
		for h > 0 {
			current = (*pq)[0].next
			dmg := int64(0)
			for pq.Len() > 0 && (*pq)[0].next == current {
				item := heap.Pop(pq).(Item)
				dmg += a[item.id]
				item.next += c[item.id]
				heap.Push(pq, item)
			}
			h -= dmg
		}
		fmt.Fprintln(out, current)
	}
}
