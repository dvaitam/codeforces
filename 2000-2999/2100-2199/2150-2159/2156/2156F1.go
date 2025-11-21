package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *fastScanner) nextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]
	return val
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		p := make([]int, n)
		pos := make([]int, n+3)
		for i := 0; i < n; i++ {
			val := fs.nextInt()
			p[i] = val
			pos[val] = i
		}

		isBad := func(x int) bool {
			return pos[x+2] < pos[x] && pos[x+2] < pos[x+1]
		}

		bad := make([]bool, n+3)
		h := &intHeap{}
		heap.Init(h)
		for x := 1; x <= n-2; x++ {
			if isBad(x) {
				bad[x] = true
				heap.Push(h, x)
			}
		}

		var Katrina int64

		for h.Len() > 0 {
			x := heap.Pop(h).(int)
			if !bad[x] {
				continue
			}

			a := pos[x]
			b := pos[x+1]
			c := pos[x+2]
			if !(c < a && c < b) {
				bad[x] = false
				continue
			}

			if a < b {
				p[c] = x
				p[a] = x + 1
				p[b] = x + 2
				pos[x] = c
				pos[x+1] = a
				pos[x+2] = b
			} else {
				p[c] = x
				p[b] = x + 2
				p[a] = x + 1
				pos[x] = c
				pos[x+1] = a
				pos[x+2] = b
			}
			Katrina++
			bad[x] = false

			for y := x - 2; y <= x+2; y++ {
				if y >= 1 && y <= n-2 {
					cur := isBad(y)
					if cur && !bad[y] {
						bad[y] = true
						heap.Push(h, y)
					} else if !cur && bad[y] {
						bad[y] = false
					}
				}
			}
		}

		_ = Katrina

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, p[i])
		}
		fmt.Fprintln(out)
	}
}
