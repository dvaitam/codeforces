package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// exam represents an exam available for preparation.
type exam struct {
	d   int // exam day
	idx int // exam index (1-based)
}

type examPQ []exam

func (pq examPQ) Len() int            { return len(pq) }
func (pq examPQ) Less(i, j int) bool  { return pq[i].d < pq[j].d }
func (pq examPQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *examPQ) Push(x interface{}) { *pq = append(*pq, x.(exam)) }
func (pq *examPQ) Pop() interface{} {
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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	type info struct{ s, d, c int }
	exams := make([]info, m+1) // 1-based
	examDay := make([]int, n+1)
	start := make([][]int, n+1)
	left := make([]int, m+1)

	for i := 1; i <= m; i++ {
		var s, d, c int
		fmt.Fscan(in, &s, &d, &c)
		exams[i] = info{s, d, c}
		examDay[d] = i
		start[s] = append(start[s], i)
		left[i] = c
	}

	res := make([]int, n+1)
	pq := &examPQ{}
	heap.Init(pq)

	for day := 1; day <= n; day++ {
		for _, idx := range start[day] {
			heap.Push(pq, exam{d: exams[idx].d, idx: idx})
		}
		if examDay[day] != 0 {
			idx := examDay[day]
			if left[idx] != 0 {
				fmt.Fprintln(out, -1)
				return
			}
			res[day] = m + 1
			continue
		}
		if pq.Len() > 0 {
			e := heap.Pop(pq).(exam)
			res[day] = e.idx
			left[e.idx]--
			if left[e.idx] > 0 {
				heap.Push(pq, e)
			}
		} else {
			res[day] = 0
		}
	}

	for i := 1; i <= m; i++ {
		if left[i] != 0 {
			fmt.Fprintln(out, -1)
			return
		}
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, res[i])
	}
	out.WriteByte('\n')
}
