package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Movie struct {
	a, b int64
	id   int
}

type Item struct {
	deadline int64
	id       int
}

type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].deadline < pq[j].deadline }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func schedule(movies []Movie, m int, D int) (bool, []int64) {
	n := len(movies)
	res := make([]int64, n)
	var pq PQ
	idx := 0
	day := int64(0)
	for idx < n || len(pq) > 0 {
		if len(pq) == 0 {
			day = max64(day, movies[idx].a)
		}
		for idx < n && movies[idx].a <= day {
			heap.Push(&pq, Item{deadline: movies[idx].b + int64(D), id: movies[idx].id})
			idx++
		}
		for i := 0; i < m && len(pq) > 0; i++ {
			it := heap.Pop(&pq).(Item)
			if it.deadline < day {
				return false, nil
			}
			res[it.id] = day
		}
		day++
	}
	return true, res
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		movies := make([]Movie, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &movies[i].a, &movies[i].b)
			movies[i].id = i
		}
		sorted := make([]Movie, n)
		copy(sorted, movies)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].a < sorted[j].a })

		lo, hi := 0, n
		for lo < hi {
			mid := (lo + hi) / 2
			ok, sched := schedule(sorted, m, mid)
			if ok {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		_, final := schedule(sorted, m, lo)
		fmt.Fprintln(out, lo)
		for i := 0; i < n; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, final[movies[i].id])
		}
		out.WriteByte('\n')
	}
}
