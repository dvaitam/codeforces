package main

import (
	"io"
	"os"
	"strconv"
)

type FastScanner struct {
	data []byte
	idx  int
}

func (fs *FastScanner) NextInt() int {
	n := len(fs.data)
	for fs.idx < n && fs.data[fs.idx] <= ' ' {
		fs.idx++
	}
	sign := 1
	if fs.data[fs.idx] == '-' {
		sign = -1
		fs.idx++
	}
	val := 0
	for fs.idx < n {
		c := fs.data[fs.idx]
		if c < '0' || c > '9' {
			break
		}
		val = val*10 + int(c-'0')
		fs.idx++
	}
	return sign * val
}

type Item struct {
	val  int64
	list int
	idx  int
}

type MaxHeap struct {
	a []Item
}

func (h *MaxHeap) Len() int {
	return len(h.a)
}

func (h *MaxHeap) Push(x Item) {
	h.a = append(h.a, x)
	i := len(h.a) - 1
	for i > 0 {
		p := (i - 1) >> 1
		if h.a[p].val >= h.a[i].val {
			break
		}
		h.a[p], h.a[i] = h.a[i], h.a[p]
		i = p
	}
}

func (h *MaxHeap) Pop() Item {
	root := h.a[0]
	last := h.a[len(h.a)-1]
	h.a = h.a[:len(h.a)-1]
	if len(h.a) > 0 {
		h.a[0] = last
		i := 0
		for {
			l := i*2 + 1
			if l >= len(h.a) {
				break
			}
			r := l + 1
			j := l
			if r < len(h.a) && h.a[r].val > h.a[l].val {
				j = r
			}
			if h.a[i].val >= h.a[j].val {
				break
			}
			h.a[i], h.a[j] = h.a[j], h.a[i]
			i = j
		}
	}
	return root
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	fs := FastScanner{data: data}
	t := fs.NextInt()
	out := make([]byte, 0, 1<<20)

	for ; t > 0; t-- {
		n := fs.NextInt()
		k := fs.NextInt()

		a := make([][]int64, n+2)
		for i := 0; i <= n+1; i++ {
			a[i] = make([]int64, n+2)
		}
		for i := 1; i <= n; i++ {
			for j := i; j <= n; j++ {
				a[i][j] = int64(fs.NextInt())
			}
		}

		dp := make([][]int64, n+1)
		dp[0] = []int64{0}

		for i := 1; i <= n; i++ {
			m := i + 1
			listSource := make([]int, m)
			listAdd := make([]int64, m)

			listSource[0] = i - 1
			listAdd[0] = 0

			listSource[1] = -1
			listAdd[1] = a[1][i]

			p := 2
			for l := 2; l <= i; l++ {
				listSource[p] = l - 2
				listAdd[p] = a[l][i]
				p++
			}

			h := MaxHeap{a: make([]Item, 0, m)}
			for id := 0; id < m; id++ {
				src := listSource[id]
				if src == -1 {
					h.Push(Item{val: listAdd[id], list: id, idx: 0})
				} else if len(dp[src]) > 0 {
					h.Push(Item{val: dp[src][0] + listAdd[id], list: id, idx: 0})
				}
			}

			res := make([]int64, 0, k)
			for h.Len() > 0 && len(res) < k {
				it := h.Pop()
				res = append(res, it.val)
				src := listSource[it.list]
				if src != -1 {
					ni := it.idx + 1
					if ni < len(dp[src]) {
						h.Push(Item{val: dp[src][ni] + listAdd[it.list], list: it.list, idx: ni})
					}
				}
			}
			dp[i] = res
		}

		ans := dp[n]
		cnt := k
		if cnt > len(ans) {
			cnt = len(ans)
		}
		for i := 0; i < cnt; i++ {
			if i > 0 {
				out = append(out, ' ')
			}
			out = strconv.AppendInt(out, ans[i], 10)
		}
		out = append(out, '\n')
	}

	os.Stdout.Write(out)
}
