package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Gap struct {
	L, R int64
	idx  int
}

type GapHeap []Gap

func (h GapHeap) Len() int            { return len(h) }
func (h GapHeap) Less(i, j int) bool  { return h[i].R < h[j].R }
func (h GapHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *GapHeap) Push(x interface{}) { *h = append(*h, x.(Gap)) }
func (h *GapHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	l := make([]int64, n)
	r := make([]int64, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &l[i], &r[i]); err != nil {
			return
		}
	}
	gaps := make([]Gap, n-1)
	for i := 0; i < n-1; i++ {
		gaps[i] = Gap{
			L:   l[i+1] - r[i],
			R:   r[i+1] - l[i],
			idx: i,
		}
	}
	bridges := make([]struct {
		val int64
		idx int
	}, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(reader, &bridges[i].val); err != nil {
			return
		}
		bridges[i].idx = i + 1
	}
	sort.Slice(gaps, func(i, j int) bool { return gaps[i].L < gaps[j].L })
	sort.Slice(bridges, func(i, j int) bool { return bridges[i].val < bridges[j].val })

	res := make([]int, n-1)
	for i := range res {
		res[i] = -1
	}
	h := &GapHeap{}
	heap.Init(h)
	gp := 0
	for _, b := range bridges {
		for gp < len(gaps) && gaps[gp].L <= b.val {
			heap.Push(h, gaps[gp])
			gp++
		}
		for h.Len() > 0 && (*h)[0].R < b.val {
			// smallest R can't fit this or any later bridge
			fmt.Println("No")
			return
		}
		if h.Len() > 0 {
			g := heap.Pop(h).(Gap)
			res[g.idx] = b.idx
		}
	}
	if gp < len(gaps) || h.Len() > 0 {
		fmt.Println("No")
		return
	}

	fmt.Println("Yes")
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i, v := range res {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
}
