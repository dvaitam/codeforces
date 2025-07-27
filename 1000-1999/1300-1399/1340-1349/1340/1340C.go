package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Item struct {
	dist int64
	idx  int
	rem  int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	d := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &d[i])
	}
	sort.Ints(d)

	var g, r int
	fmt.Fscan(in, &g, &r)

	const INF int64 = 1 << 60
	dist := make([]int64, m*g)
	for i := range dist {
		dist[i] = INF
	}
	pq := &PriorityQueue{}
	heap.Init(pq)
	dist[0] = 0
	heap.Push(pq, Item{0, 0, 0})

	ans := int64(-1)
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		id := cur.idx*g + cur.rem
		if cur.dist != dist[id] {
			continue
		}
		if ans != -1 && cur.dist >= ans {
			break
		}
		for _, next := range []int{cur.idx - 1, cur.idx + 1} {
			if next < 0 || next >= m {
				continue
			}
			delta := abs(d[next] - d[cur.idx])
			nt := cur.rem + delta
			if nt > g {
				continue
			}
			nd := cur.dist + int64(delta)
			if next == m-1 {
				if ans == -1 || nd < ans {
					ans = nd
				}
				continue
			}
			if nt == g {
				nd += int64(r)
				nt = 0
			}
			nid := next*g + nt
			if nd < dist[nid] {
				dist[nid] = nd
				heap.Push(pq, Item{nd, next, nt})
			}
		}
	}

	fmt.Fprintln(out, ans)
}
