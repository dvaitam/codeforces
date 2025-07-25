package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Item struct {
	val int64
	id  int
	ver int
}

type MinHeap struct {
	items []Item
	ver   []int
}

func (h MinHeap) Len() int            { return len(h.items) }
func (h MinHeap) Less(i, j int) bool  { return h.items[i].val < h.items[j].val }
func (h MinHeap) Swap(i, j int)       { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *MinHeap) Push(x interface{}) { h.items = append(h.items, x.(Item)) }
func (h *MinHeap) Pop() interface{} {
	n := len(h.items)
	x := h.items[n-1]
	h.items = h.items[:n-1]
	return x
}
func (h *MinHeap) top() *Item {
	for len(h.items) > 0 && h.ver[h.items[0].id] != h.items[0].ver {
		heap.Pop(h)
	}
	if len(h.items) == 0 {
		return nil
	}
	return &h.items[0]
}
func (h *MinHeap) popValid() *Item {
	for len(h.items) > 0 {
		it := heap.Pop(h).(Item)
		if h.ver[it.id] == it.ver {
			return &it
		}
	}
	return nil
}

type MaxHeap struct {
	items []Item
	ver   []int
}

func (h MaxHeap) Len() int            { return len(h.items) }
func (h MaxHeap) Less(i, j int) bool  { return h.items[i].val > h.items[j].val }
func (h MaxHeap) Swap(i, j int)       { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *MaxHeap) Push(x interface{}) { h.items = append(h.items, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	n := len(h.items)
	x := h.items[n-1]
	h.items = h.items[:n-1]
	return x
}
func (h *MaxHeap) top() *Item {
	for len(h.items) > 0 && h.ver[h.items[0].id] != h.items[0].ver {
		heap.Pop(h)
	}
	if len(h.items) == 0 {
		return nil
	}
	return &h.items[0]
}
func (h *MaxHeap) popValid() *Item {
	for len(h.items) > 0 {
		it := heap.Pop(h).(Item)
		if h.ver[it.id] == it.ver {
			return &it
		}
	}
	return nil
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solve(a []int64, b int64, k1, k2 int) int64 {
	n := len(a)
	val := append([]int64(nil), a...)
	ver := make([]int, n)
	halved := make([]bool, n)
	inBig := make([]bool, n)
	t := make([]int64, n)
	for i := 0; i < n; i++ {
		t[i] = minInt64(b, val[i])
	}
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return t[idx[i]] > t[idx[j]] })
	big := &MinHeap{ver: ver}
	small := &MaxHeap{ver: ver}
	var bigSum int64
	for rank, id := range idx {
		if rank < k2 {
			heap.Push(big, Item{t[id], id, ver[id]})
			inBig[id] = true
			bigSum += t[id]
		} else {
			heap.Push(small, Item{t[id], id, ver[id]})
		}
	}
	sumVal := int64(0)
	for _, v := range val {
		sumVal += v
	}
	for ; k1 > 0; k1-- {
		st := small.top()
		smallTop := int64(-1)
		if st != nil {
			smallTop = st.val
		}
		bestIdx := -1
		var bestDelta int64
		for i := 0; i < n; i++ {
			if halved[i] {
				continue
			}
			old := val[i]
			newv := (old + 1) / 2
			deltaVal := old - newv
			tmp := bigSum
			if inBig[i] {
				oldT := minInt64(b, old)
				newT := minInt64(b, newv)
				tmp = bigSum - oldT
				if smallTop > newT {
					tmp += smallTop
				} else {
					tmp += newT
				}
			}
			delta := deltaVal + (tmp - bigSum)
			if delta > bestDelta {
				bestDelta = delta
				bestIdx = i
			}
		}
		if bestIdx == -1 || bestDelta <= 0 {
			break
		}
		old := val[bestIdx]
		newv := (old + 1) / 2
		val[bestIdx] = newv
		halved[bestIdx] = true
		sumVal -= old - newv
		ver[bestIdx]++
		oldT := minInt64(b, old)
		newT := minInt64(b, newv)
		if inBig[bestIdx] {
			bigSum -= oldT
			st := small.top()
			smallTop := int64(-1)
			if st != nil {
				smallTop = st.val
			}
			if smallTop > newT {
				item := small.popValid()
				if item != nil {
					heap.Push(big, Item{item.val, item.id, item.ver})
					bigSum += item.val
					inBig[item.id] = true
				}
				heap.Push(small, Item{newT, bestIdx, ver[bestIdx]})
				inBig[bestIdx] = false
			} else {
				heap.Push(big, Item{newT, bestIdx, ver[bestIdx]})
				bigSum += newT
			}
		} else {
			heap.Push(small, Item{newT, bestIdx, ver[bestIdx]})
		}
	}
	return sumVal - bigSum
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var b int64
		var k1, k2 int
		fmt.Fscan(in, &n, &b, &k1, &k2)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		ans := solve(a, b, k1, k2)
		fmt.Fprintln(out, ans)
	}
}
