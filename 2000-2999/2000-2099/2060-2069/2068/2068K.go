package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type event struct {
	val  int64
	ride int
	k    int
}

type minHeap []event

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].val < h[j].val }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(event)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

var (
	timeAdj   [][]int
	matchTime []int
	matchRide []int
	timeVis   []int
	rideVis   []int
	visitTok  int
)

func dfs(tid int) bool {
	if timeVis[tid] == visitTok {
		return false
	}
	timeVis[tid] = visitTok
	for _, r := range timeAdj[tid] {
		if rideVis[r] == visitTok {
			continue
		}
		rideVis[r] = visitTok
		if matchRide[r] == -1 || dfs(matchRide[r]) {
			matchRide[r] = tid
			matchTime[tid] = r
			return true
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		h := &minHeap{}
		heap.Init(h)
		for i := 0; i < n; i++ {
			heap.Push(h, event{val: a[i], ride: i, k: 1})
		}

		timeAdj = timeAdj[:0]
		matchTime = matchTime[:0]
		timeVis = timeVis[:0]
		matchRide = make([]int, n)
		rideVis = make([]int, n)
		for i := range matchRide {
			matchRide[i] = -1
		}
		visitTok = 0

		matched := 0
		var ans int64
		for matched < n {
			curVal := (*h)[0].val
			rides := make([]int, 0)
			for h.Len() > 0 && (*h)[0].val == curVal {
				ev := heap.Pop(h).(event)
				rides = append(rides, ev.ride)
				if ev.k < n {
					nextVal := ev.val + a[ev.ride]
					heap.Push(h, event{val: nextVal, ride: ev.ride, k: ev.k + 1})
				}
			}
			timeAdj = append(timeAdj, rides)
			matchTime = append(matchTime, -1)
			timeVis = append(timeVis, 0)

			visitTok++
			if dfs(len(timeAdj) - 1) {
				matched++
				if matched == n {
					ans = curVal
					break
				}
			}
		}

		fmt.Fprintln(out, ans)
	}
}
