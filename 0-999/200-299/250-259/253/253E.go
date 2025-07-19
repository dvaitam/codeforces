package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// Task holds start time t, required duration s, and priority p
type Task struct {
	t, s, p int
}

// Info holds a key x and an identifier y
type Info struct {
	x, y int
}

// MaxHeap implements a max-heap of Info based on x
type MaxHeap []Info

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].x > h[j].x }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Info)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

var (
	n, w  int
	T     int64
	tasks []Task
	q     []Info
	used  []int64
	last  []int64
)

// calc returns overlap length of intervals [s1,t1] and [s2,t2]
func calc(s1, t1, s2, t2 int64) int64 {
	low := s1
	if s2 > low {
		low = s2
	}
	high := t1
	if t2 < high {
		high = t2
	}
	res := high - low + 1
	if res > 0 {
		return res
	}
	return 0
}

// solve simulates task processing; if updateUsed is true, accumulates used times
func solve(updateUsed bool) {
	var cur int64
	h := &MaxHeap{}
	heap.Init(h)
	need := make([]int, n)
	for i := 0; i < n; i++ {
		heap.Push(h, Info{tasks[q[i].y].p, q[i].y})
		need[q[i].y] = tasks[q[i].y].s
		if cur < int64(q[i].x) {
			cur = int64(q[i].x)
		}
		for (i == n-1 || cur < int64(q[i+1].x)) && h.Len() > 0 {
			top := (*h)[0]
			x := top.y
			var d int
			if i == n-1 {
				d = need[x]
			} else {
				avail := int(int64(q[i+1].x) - cur)
				if need[x] < avail {
					d = need[x]
				} else {
					d = avail
				}
			}
			added := calc(int64(tasks[w].t), T-1, cur, cur+int64(d)-1)
			if updateUsed {
				used[x] += added
			}
			cur += int64(d)
			need[x] -= d
			if need[x] == 0 {
				last[x] = cur
				heap.Pop(h)
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fscan(reader, &n)
	tasks = make([]Task, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &tasks[i].t, &tasks[i].s, &tasks[i].p)
		if tasks[i].p == -1 {
			w = i
		}
	}
	fmt.Fscan(reader, &T)
	q = make([]Info, n)
	for i := 0; i < n; i++ {
		q[i] = Info{tasks[i].t, i}
	}
	sort.Slice(q, func(i, j int) bool { return q[i].x < q[j].x })
	used = make([]int64, n)
	last = make([]int64, n)
	solve(true)
	g := make([]Info, n)
	for i := 0; i < n; i++ {
		g[i] = Info{tasks[i].p, i}
	}
	sort.Slice(g, func(i, j int) bool { return g[i].x < g[j].x })
	var sum int64
	j := 0
	for sum < int64(tasks[w].s) {
		sum += used[g[j].y]
		j++
	}
	j--
	newP := g[j].x + 1
	if newP < 1 {
		newP = 1
	}
	for j+1 < n && g[j+1].x == newP {
		newP++
		j++
	}
	tasks[w].p = newP
	fmt.Fprintln(writer, newP)
	solve(false)
	for i := 0; i < n; i++ {
		fmt.Fprintf(writer, "%d ", last[i])
	}
	fmt.Fprintln(writer)
}
