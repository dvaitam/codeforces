package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type event struct {
	s, e int
	kind int // 0 lecture, 1 seminar
	idx  int
}

type item struct {
	end int
	idx int
}

type pq []item

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].end < p[j].end }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(item)) }
func (p *pq) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, x, y int
		if _, err := fmt.Fscan(in, &n, &m, &x, &y); err != nil {
			return
		}
		lectures := make([][2]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &lectures[i][0], &lectures[i][1])
		}
		seminars := make([][2]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &seminars[i][0], &seminars[i][1])
		}
		events := make([]event, 0, n+m)
		for i := 0; i < n; i++ {
			events = append(events, event{s: lectures[i][0], e: lectures[i][1], kind: 0, idx: i})
		}
		for i := 0; i < m; i++ {
			events = append(events, event{s: seminars[i][0], e: seminars[i][1], kind: 1, idx: i})
		}
		sort.Slice(events, func(i, j int) bool {
			if events[i].s == events[j].s {
				return events[i].kind < events[j].kind
			}
			return events[i].s < events[j].s
		})

		availableHD := make([]int, x)
		for i := 0; i < x; i++ {
			availableHD[i] = i + 1
		}
		availableOrd := make([]int, y)
		for i := 0; i < y; i++ {
			availableOrd[i] = x + i + 1
		}
		var hdpq, ordpq pq
		assignmentsL := make([]int, n)
		assignmentsS := make([]int, m)
		possible := true

		for _, ev := range events {
			for len(hdpq) > 0 && hdpq[0].end <= ev.s {
				it := heap.Pop(&hdpq).(item)
				availableHD = append(availableHD, it.idx)
			}
			for len(ordpq) > 0 && ordpq[0].end <= ev.s {
				it := heap.Pop(&ordpq).(item)
				availableOrd = append(availableOrd, it.idx)
			}

			if ev.kind == 0 {
				if len(availableHD) == 0 {
					possible = false
					break
				}
				idx := availableHD[len(availableHD)-1]
				availableHD = availableHD[:len(availableHD)-1]
				assignmentsL[ev.idx] = idx
				heap.Push(&hdpq, item{end: ev.e, idx: idx})
			} else {
				if len(availableOrd) > 0 {
					idx := availableOrd[len(availableOrd)-1]
					availableOrd = availableOrd[:len(availableOrd)-1]
					assignmentsS[ev.idx] = idx
					heap.Push(&ordpq, item{end: ev.e, idx: idx})
				} else if len(availableHD) > 0 {
					idx := availableHD[len(availableHD)-1]
					availableHD = availableHD[:len(availableHD)-1]
					assignmentsS[ev.idx] = idx
					heap.Push(&hdpq, item{end: ev.e, idx: idx})
				} else {
					possible = false
					break
				}
			}
		}

		if possible {
			fmt.Fprintln(out, "YES")
			for i := 0; i < n; i++ {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, assignmentsL[i])
			}
			for j := 0; j < m; j++ {
				fmt.Fprint(out, " ", assignmentsS[j])
			}
			fmt.Fprintln(out)
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
