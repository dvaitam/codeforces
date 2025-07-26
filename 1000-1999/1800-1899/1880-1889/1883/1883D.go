package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// MaxHeap implements a max-heap for ints using container/heap.
type MaxHeap struct{ sort.IntSlice }

func (h MaxHeap) Less(i, j int) bool  { return h.IntSlice[i] > h.IntSlice[j] }
func (h *MaxHeap) Push(x interface{}) { h.IntSlice = append(h.IntSlice, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := h.IntSlice
	x := old[len(old)-1]
	h.IntSlice = old[:len(old)-1]
	return x
}
func (h *MaxHeap) Peek() int { return h.IntSlice[0] }

// MinHeap implements a min-heap for ints using container/heap.
type MinHeap struct{ sort.IntSlice }

func (h *MinHeap) Push(x interface{}) { h.IntSlice = append(h.IntSlice, x.(int)) }
func (h *MinHeap) Pop() interface{} {
	old := h.IntSlice
	x := old[len(old)-1]
	h.IntSlice = old[:len(old)-1]
	return x
}
func (h *MinHeap) Peek() int { return h.IntSlice[0] }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	lcnt := make(map[int]int)
	rcnt := make(map[int]int)
	var lMax MaxHeap
	var rMin MinHeap
	heap.Init(&lMax)
	heap.Init(&rMin)
	segments := 0

	for ; q > 0; q-- {
		var op string
		var l, r int
		fmt.Fscan(reader, &op, &l, &r)
		if op == "+" {
			heap.Push(&lMax, l)
			heap.Push(&rMin, r)
			lcnt[l]++
			rcnt[r]++
			segments++
		} else { // op == "-"
			lcnt[l]--
			if lcnt[l] == 0 {
				delete(lcnt, l)
			}
			rcnt[r]--
			if rcnt[r] == 0 {
				delete(rcnt, r)
			}
			segments--
		}

		for lMax.Len() > 0 {
			top := lMax.Peek()
			if lcnt[top] == 0 {
				heap.Pop(&lMax)
			} else {
				break
			}
		}
		for rMin.Len() > 0 {
			top := rMin.Peek()
			if rcnt[top] == 0 {
				heap.Pop(&rMin)
			} else {
				break
			}
		}

		if segments >= 2 && lMax.Peek() > rMin.Peek() {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
