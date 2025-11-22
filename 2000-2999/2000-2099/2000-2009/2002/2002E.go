package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type fastScanner struct {
	r   *bufio.Reader
	buf []byte
}

func newScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *fastScanner) nextInt64() int64 {
	sign, val := int64(1), int64(0)
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

type event struct {
	t  int64
	id int
}

type eventHeap []event

func (h eventHeap) Len() int            { return len(h) }
func (h eventHeap) Less(i, j int) bool  { return h[i].t < h[j].t }
func (h eventHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *eventHeap) Push(x interface{}) { *h = append(*h, x.(event)) }
func (h *eventHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := newScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	t := int(in.nextInt64())
	for ; t > 0; t-- {
		n := int(in.nextInt64())

		val := make([]int, n+2)
		head := make([]int, n+2)
		tail := make([]int, n+2)
		prev := make([]int, n+2)
		next := make([]int, n+2)
		active := make([]bool, n+2)

		compLen := make([]int64, n+2) // remaining length for component id
		compNext := make([]int, n+2)
		dead := make([]int64, n+2)

		rem := make([]int64, n+2)  // remaining total length of current run
		last := make([]int64, n+2) // last time the run was updated

		var h eventHeap
		h = make([]event, 0, n)

		for i := 1; i <= n; i++ {
			a := in.nextInt64()
			b := int(in.nextInt64())
			val[i] = b
			head[i] = i
			tail[i] = i
			compLen[i] = a
			rem[i] = a
			active[i] = true
			if i > 1 {
				prev[i] = i - 1
			}
			if i < n {
				next[i] = i + 1
			}
			heap.Push(&h, event{t: a, id: i})
		}

		decay := func(id int, to int64) {
			if !active[id] {
				return
			}
			dt := to - last[id]
			if dt <= 0 {
				return
			}
			curr := head[id]
			currTime := last[id]
			for dt > 0 && curr != 0 {
				l := compLen[curr]
				if l > dt {
					compLen[curr] = l - dt
					rem[id] -= dt
					currTime += dt
					dt = 0
				} else {
					rem[id] -= l
					currTime += l
					dt -= l
					dead[curr] = currTime
					curr = compNext[curr]
				}
			}
			head[id] = curr
			if curr == 0 {
				tail[id] = 0
				rem[id] = 0
			}
			last[id] = to
		}

		for len(h) > 0 {
			ev := heap.Pop(&h).(event)
			id := ev.id
			if !active[id] {
				continue
			}
			if ev.t != last[id]+rem[id] {
				continue // outdated event
			}

			decay(id, ev.t)
			// run id should now be empty
			active[id] = false
			l := prev[id]
			r := next[id]
			if l != 0 {
				next[l] = r
			}
			if r != 0 {
				prev[r] = l
			}

			if l != 0 && r != 0 && active[l] && active[r] && val[l] == val[r] {
				decay(l, ev.t)
				decay(r, ev.t)
				if !active[l] || !active[r] || head[l] == 0 || head[r] == 0 {
					continue
				}
				compNext[tail[l]] = head[r]
				tail[l] = tail[r]
				rem[l] += rem[r]
				last[l] = ev.t
				next[l] = next[r]
				if next[r] != 0 {
					prev[next[r]] = l
				}
				active[r] = false
				heap.Push(&h, event{t: last[l] + rem[l], id: l})
			}
		}

		ans := make([]int64, n+1)
		var mx int64
		for i := 1; i <= n; i++ {
			if dead[i] > mx {
				mx = dead[i]
			}
			ans[i] = mx
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
