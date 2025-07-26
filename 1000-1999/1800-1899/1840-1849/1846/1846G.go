package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

type Item struct {
	state int
	dist  int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Item))
}
func (pq *PriorityQueue) Pop() interface{} {
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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var s string
		fmt.Fscan(reader, &s)
		initState := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				initState |= 1 << i
			}
		}
		type Med struct {
			d   int
			rem int
			add int
		}
		meds := make([]Med, m)
		for i := 0; i < m; i++ {
			var d int
			fmt.Fscan(reader, &d)
			var remStr, addStr string
			fmt.Fscan(reader, &remStr)
			fmt.Fscan(reader, &addStr)
			remMask := 0
			addMask := 0
			for j := 0; j < n; j++ {
				if remStr[j] == '1' {
					remMask |= 1 << j
				}
				if addStr[j] == '1' {
					addMask |= 1 << j
				}
			}
			meds[i] = Med{d: d, rem: remMask, add: addMask}
		}

		const INF = math.MaxInt32
		dist := make([]int, 1<<uint(n))
		for i := range dist {
			dist[i] = INF
		}
		pq := &PriorityQueue{}
		heap.Init(pq)
		dist[initState] = 0
		heap.Push(pq, Item{state: initState, dist: 0})
		ans := -1
		for pq.Len() > 0 {
			it := heap.Pop(pq).(Item)
			if it.dist != dist[it.state] {
				continue
			}
			if it.state == 0 {
				ans = it.dist
				break
			}
			for _, med := range meds {
				newState := (it.state &^ med.rem) | med.add
				nd := it.dist + med.d
				if nd < dist[newState] {
					dist[newState] = nd
					heap.Push(pq, Item{state: newState, dist: nd})
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
