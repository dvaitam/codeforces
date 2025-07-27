package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Camera struct {
	t int
	s int
}

type job struct {
	release  int
	deadline int
}

type minHeap []int

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	v := old[len(old)-1]
	*h = old[:len(old)-1]
	return v
}

func feasible(cams []Camera, t1, t2 int) bool {
	if t1 < 0 || t2 <= t1 {
		return false
	}
	delta := t2 - t1
	jobs := make([]job, 0, len(cams)*2)
	for _, c := range cams {
		switch c.t {
		case 1:
			r := t1 - c.s
			if r < 0 {
				r = 0
			}
			jobs = append(jobs, job{r, t1 - 1})
		case 2:
			r := t2 - c.s
			if r < 0 {
				r = 0
			}
			jobs = append(jobs, job{r, t2 - 1})
		case 3:
			if c.s > delta {
				r := t2 - c.s
				if r < 0 {
					r = 0
				}
				if r > t1-1 {
					r = t1 - 1
				}
				if r > t1-1 {
					return false
				}
				jobs = append(jobs, job{r, t1 - 1})
			} else {
				r1 := t1 - c.s
				if r1 < 0 {
					r1 = 0
				}
				jobs = append(jobs, job{r1, t1 - 1})
				r2 := t2 - c.s
				if r2 < 0 {
					r2 = 0
				}
				jobs = append(jobs, job{r2, t2 - 1})
			}
		}
	}

	sort.Slice(jobs, func(i, j int) bool { return jobs[i].release < jobs[j].release })

	h := &minHeap{}
	heap.Init(h)
	idx := 0
	for time := 0; time < t2; time++ {
		for idx < len(jobs) && jobs[idx].release <= time {
			heap.Push(h, jobs[idx].deadline)
			idx++
		}
		if time == t1 {
			if h.Len() > 0 && (*h)[0] < time {
				return false
			}
			continue
		}
		if h.Len() == 0 {
			continue
		}
		if (*h)[0] < time {
			return false
		}
		heap.Pop(h)
	}
	return idx == len(jobs) && h.Len() == 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	cams := make([]Camera, n)
	smax := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &cams[i].t, &cams[i].s)
		if cams[i].s > smax {
			smax = cams[i].s
		}
	}

	limit := 2*n + smax + 5
	ans := -1
	for delta := 1; delta <= limit; delta++ {
		low, high := 0, limit-delta
		found := -1
		for low <= high {
			mid := (low + high) / 2
			if feasible(cams, mid, mid+delta) {
				found = mid
				high = mid - 1
			} else {
				low = mid + 1
			}
		}
		if found != -1 {
			t2 := found + delta
			if ans == -1 || t2 < ans {
				ans = t2
			}
		}
	}
	fmt.Fprintln(out, ans)
}
