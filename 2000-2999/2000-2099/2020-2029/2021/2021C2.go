package main

import (
	"bufio"
	"fmt"
	"os"
)

type minHeap struct {
	data []int
}

func (h *minHeap) push(x int) {
	h.data = append(h.data, x)
	i := len(h.data) - 1
	for i > 0 {
		p := (i - 1) / 2
		if h.data[p] <= h.data[i] {
			break
		}
		h.data[p], h.data[i] = h.data[i], h.data[p]
		i = p
	}
}

func (h *minHeap) pop() int {
	n := len(h.data)
	if n == 0 {
		return -1
	}
	res := h.data[0]
	last := h.data[n-1]
	h.data = h.data[:n-1]
	if n-1 > 0 {
		h.data[0] = last
		i := 0
		for {
			l := 2*i + 1
			r := l + 1
			smallest := i
			if l < len(h.data) && h.data[l] < h.data[smallest] {
				smallest = l
			}
			if r < len(h.data) && h.data[r] < h.data[smallest] {
				smallest = r
			}
			if smallest == i {
				break
			}
			h.data[i], h.data[smallest] = h.data[smallest], h.data[i]
			i = smallest
		}
	}
	return res
}

func (h *minHeap) topValid(val int, b []int, inf int) int {
	for len(h.data) > 0 {
		top := h.data[0]
		if b[top] == val {
			return top
		}
		h.pop()
	}
	return inf
}

func recomputePrefixEdge(i int, present []bool, edge []bool, cnt *int) {
	if i < 1 || i >= len(present)-1 {
		return
	}
	newVal := !present[i] && present[i+1]
	if edge[i] == newVal {
		return
	}
	if newVal {
		*cnt++
	} else {
		*cnt--
	}
	edge[i] = newVal
}

func recomputeOrderEdge(i int, present []bool, firstIdx []int, edge []bool, cnt *int) {
	if i < 1 || i >= len(present)-1 {
		return
	}
	newVal := present[i] && present[i+1] && firstIdx[i] >= firstIdx[i+1]
	if edge[i] == newVal {
		return
	}
	if newVal {
		*cnt++
	} else {
		*cnt--
	}
	edge[i] = newVal
}

func updateFirst(idx, newVal, inf int, firstIdx []int, present []bool, prefixEdge []bool, orderEdge []bool, prefixCnt *int, orderCnt *int) {
	if firstIdx[idx] == newVal {
		return
	}
	oldPresent := present[idx]
	firstIdx[idx] = newVal
	newPresent := newVal != inf
	present[idx] = newPresent
	if oldPresent != newPresent {
		recomputePrefixEdge(idx-1, present, prefixEdge, prefixCnt)
		recomputePrefixEdge(idx, present, prefixEdge, prefixCnt)
	}
	recomputeOrderEdge(idx-1, present, firstIdx, orderEdge, orderCnt)
	recomputeOrderEdge(idx, present, firstIdx, orderEdge, orderCnt)
}

func refreshValue(val int, pos []int, heaps []*minHeap, b []int, inf int, firstIdx []int, present []bool, prefixEdge []bool, orderEdge []bool, prefixCnt *int, orderCnt *int) {
	idx := pos[val]
	newFirst := heaps[val].topValid(val, b, inf)
	updateFirst(idx, newFirst, inf, firstIdx, present, prefixEdge, orderEdge, prefixCnt, orderCnt)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m, q int
		fmt.Fscan(reader, &n, &m, &q)

		a := make([]int, n+1)
		pos := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			pos[a[i]] = i
		}

		b := make([]int, m+1)
		for i := 1; i <= m; i++ {
			fmt.Fscan(reader, &b[i])
		}

		heaps := make([]*minHeap, n+1)
		for i := 1; i <= n; i++ {
			heaps[i] = &minHeap{}
		}
		for i := 1; i <= m; i++ {
			val := b[i]
			heaps[val].push(i)
		}

		inf := m + q + 5
		firstIdx := make([]int, n+2)
		for i := range firstIdx {
			firstIdx[i] = inf
		}
		present := make([]bool, n+2)
		prefixEdge := make([]bool, n+2)
		orderEdge := make([]bool, n+2)
		prefixCnt := 0
		orderCnt := 0

		for val := 1; val <= n; val++ {
			idx := pos[val]
			first := heaps[val].topValid(val, b, inf)
			firstIdx[idx] = first
			if first != inf {
				present[idx] = true
			}
		}

		for i := 1; i < n; i++ {
			recomputePrefixEdge(i, present, prefixEdge, &prefixCnt)
			recomputeOrderEdge(i, present, firstIdx, orderEdge, &orderCnt)
		}

		output := func() {
			if prefixCnt == 0 && orderCnt == 0 {
				fmt.Fprintln(writer, "YA")
			} else {
				fmt.Fprintln(writer, "TIDAK")
			}
		}

		output()

		for ; q > 0; q-- {
			var s, t int
			fmt.Fscan(reader, &s, &t)
			if b[s] != t {
				oldVal := b[s]
				b[s] = t
				heaps[t].push(s)
				refreshValue(oldVal, pos, heaps, b, inf, firstIdx, present, prefixEdge, orderEdge, &prefixCnt, &orderCnt)
				refreshValue(t, pos, heaps, b, inf, firstIdx, present, prefixEdge, orderEdge, &prefixCnt, &orderCnt)
			}
			output()
		}
	}
}
