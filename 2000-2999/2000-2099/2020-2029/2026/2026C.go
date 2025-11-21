package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func (fs *FastScanner) NextBytes() []byte {
	c, _ := fs.r.ReadByte()
	for c == ' ' || c == '\n' || c == '\r' || c == '\t' {
		c, _ = fs.r.ReadByte()
	}
	buf := []byte{c}
	for {
		c, err := fs.r.ReadByte()
		if err != nil || c == ' ' || c == '\n' || c == '\r' || c == '\t' {
			break
		}
		buf = append(buf, c)
	}
	return buf
}

type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]
	return val
}

type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]
	return val
}

func popValidMin(h *MinHeap, used []bool) int {
	for h.Len() > 0 {
		val := heap.Pop(h).(int)
		if !used[val] {
			return val
		}
	}
	return -1
}

func popValidMax(h *MaxHeap, used []bool) int {
	for h.Len() > 0 {
		val := heap.Pop(h).(int)
		if !used[val] {
			return val
		}
	}
	return -1
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt()
	for ; t > 0; t-- {
		n := in.NextInt()
		s := in.NextBytes()

		candidates := make([]int, 0)
		for i := 1; i <= n; i++ {
			if s[i-1] == '1' && i >= 2 {
				candidates = append(candidates, i)
			}
		}

		cnt := 0
		for _, day := range candidates {
			if day >= 2*(cnt+1) {
				cnt++
			}
		}

		selected := make([]int, cnt)
		need := cnt
		for i := len(candidates) - 1; i >= 0 && need > 0; i-- {
			day := candidates[i]
			if day >= 2*need {
				selected[need-1] = day
				need--
			}
		}

		total := int64(n) * int64(n+1) / 2
		used := make([]bool, n+1)
		minH := &MinHeap{}
		maxH := &MaxHeap{}
		heap.Init(minH)
		heap.Init(maxH)

		available := 0
		ptr := 0
		for day := 1; day <= n; day++ {
			heap.Push(minH, day)
			heap.Push(maxH, day)
			available++
			if ptr < len(selected) && selected[ptr] == day {
				if available < 2 {
					continue
				}
				free := popValidMax(maxH, used)
				if free == -1 {
					continue
				}
				used[free] = true
				available--
				total -= int64(free)

				pay := popValidMin(minH, used)
				if pay == -1 {
					continue
				}
				used[pay] = true
				available--
				ptr++
			}
		}

		fmt.Fprintln(out, total)
	}
}
