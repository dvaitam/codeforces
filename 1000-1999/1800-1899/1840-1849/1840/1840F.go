package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

const inf int64 = 1 << 60

type item struct {
	t    int64
	i, j int
}

type priorityQueue []item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].t < pq[j].t }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func nextAfter(arr []int64, t int64) int64 {
	idx := sort.Search(len(arr), func(i int) bool { return arr[i] > t })
	if idx == len(arr) {
		return inf
	}
	return arr[idx]
}

func isShot(arr []int64, t int64) bool {
	idx := sort.Search(len(arr), func(i int) bool { return arr[i] >= t })
	return idx < len(arr) && arr[idx] == t
}

func earliest(t int64, i, j, ni, nj int, rowShots, colShots [][]int64) (int64, bool) {
	deadline := nextAfter(rowShots[i], t)
	tmp := nextAfter(colShots[j], t)
	if tmp < deadline {
		deadline = tmp
	}
	cand := t + 1
	for cand <= deadline {
		if !isShot(rowShots[ni], cand) && !isShot(colShots[nj], cand) {
			return cand, true
		}
		cand++
	}
	return 0, false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		rowShots := make([][]int64, n+1)
		colShots := make([][]int64, m+1)
		var r int
		fmt.Fscan(reader, &r)
		for k := 0; k < r; k++ {
			var tt int64
			var d int
			var coord int
			fmt.Fscan(reader, &tt, &d, &coord)
			if d == 1 {
				rowShots[coord] = append(rowShots[coord], tt)
			} else {
				colShots[coord] = append(colShots[coord], tt)
			}
		}
		for i := 0; i <= n; i++ {
			sort.Slice(rowShots[i], func(a, b int) bool { return rowShots[i][a] < rowShots[i][b] })
		}
		for j := 0; j <= m; j++ {
			sort.Slice(colShots[j], func(a, b int) bool { return colShots[j][a] < colShots[j][b] })
		}

		dist := make([][]int64, n+1)
		for i := 0; i <= n; i++ {
			dist[i] = make([]int64, m+1)
			for j := 0; j <= m; j++ {
				dist[i][j] = inf
			}
		}
		dist[0][0] = 0
		pq := &priorityQueue{}
		heap.Push(pq, item{0, 0, 0})

		ans := int64(-1)
		for pq.Len() > 0 {
			it := heap.Pop(pq).(item)
			t, i, j := it.t, it.i, it.j
			if t != dist[i][j] {
				continue
			}
			if i == n && j == m {
				ans = t
				break
			}
			if i < n {
				if nt, ok := earliest(t, i, j, i+1, j, rowShots, colShots); ok && nt < dist[i+1][j] {
					dist[i+1][j] = nt
					heap.Push(pq, item{nt, i + 1, j})
				}
			}
			if j < m {
				if nt, ok := earliest(t, i, j, i, j+1, rowShots, colShots); ok && nt < dist[i][j+1] {
					dist[i][j+1] = nt
					heap.Push(pq, item{nt, i, j + 1})
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
