package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Item struct {
	deadline int
	idx      int
}

// PQ implements a min-heap ordered by deadline, then index.
type PQ []Item

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	if pq[i].deadline == pq[j].deadline {
		return pq[i].idx < pq[j].idx
	}
	return pq[i].deadline < pq[j].deadline
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func canPlace(target, pos int, n int, p []int, g [][]int, indeg0 []int) bool {
	indeg := make([]int, n)
	copy(indeg, indeg0)
	pq := &PQ{}
	heap.Init(pq)
	for i := 0; i < n; i++ {
		if i == target {
			continue
		}
		if indeg[i] == 0 {
			heap.Push(pq, Item{deadline: p[i], idx: i})
		}
	}
	for t := 1; t <= n; t++ {
		if t == pos {
			if indeg[target] != 0 || p[target] < t {
				return false
			}
			if pq.Len() > 0 && (*pq)[0].deadline < t {
				return false
			}
			for _, to := range g[target] {
				indeg[to]--
				if indeg[to] == 0 && to != target {
					heap.Push(pq, Item{deadline: p[to], idx: to})
				}
			}
			continue
		}
		if pq.Len() == 0 {
			return false
		}
		if (*pq)[0].deadline < t {
			return false
		}
		item := heap.Pop(pq).(Item)
		for _, to := range g[item.idx] {
			indeg[to]--
			if indeg[to] == 0 && to != target {
				heap.Push(pq, Item{deadline: p[to], idx: to})
			}
		}
	}
	return true
}

func earliestPosition(i int, n int, p []int, g [][]int, indeg []int) int {
	l, r := 1, p[i]
	ans := p[i]
	for l <= r {
		mid := (l + r) / 2
		if canPlace(i, mid, n, p, g, indeg) {
			ans = mid
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	g := make([][]int, n)
	indeg := make([]int, n)
	for j := 0; j < m; j++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		a--
		b--
		g[a] = append(g[a], b)
		indeg[b]++
	}
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = earliestPosition(i, n, p, g, indeg)
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, res[i])
	}
	fmt.Fprintln(writer)
}
