package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var data []byte
var ptr int

func nextInt() int {
	n := len(data)
	for ptr < n && (data[ptr] < '0' || data[ptr] > '9') {
		ptr++
	}
	v := 0
	for ptr < n && data[ptr] >= '0' && data[ptr] <= '9' {
		v = v*10 + int(data[ptr]-'0')
		ptr++
	}
	return v
}

var hv, htag []int64
var hl, hr, hdist []int
var hcnt int

func newHeap(v int64) int {
	hcnt++
	hv[hcnt] = v
	htag[hcnt] = 0
	hl[hcnt] = 0
	hr[hcnt] = 0
	hdist[hcnt] = 1
	return hcnt
}

func apply(x int, d int64) {
	if x == 0 {
		return
	}
	hv[x] += d
	htag[x] += d
}

func push(x int) {
	if x == 0 || htag[x] == 0 {
		return
	}
	d := htag[x]
	apply(hl[x], d)
	apply(hr[x], d)
	htag[x] = 0
}

func merge(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	if hv[a] < hv[b] {
		a, b = b, a
	}
	push(a)
	hr[a] = merge(hr[a], b)
	if hdist[hl[a]] < hdist[hr[a]] {
		hl[a], hr[a] = hr[a], hl[a]
	}
	hdist[a] = hdist[hr[a]] + 1
	return a
}

func popMax(x int) (int, int) {
	push(x)
	a, b := hl[x], hr[x]
	hl[x], hr[x] = 0, 0
	hdist[x] = 1
	htag[x] = 0
	return merge(a, b), x
}

func main() {
	data, _ = io.ReadAll(os.Stdin)
	n := nextInt()
	k := int64(nextInt())

	head := make([]int, n+1)
	to := make([]int, n)
	nxt := make([]int, n)
	edge := 1
	for i := 2; i <= n; i++ {
		p := nextInt()
		to[edge] = i
		nxt[edge] = head[p]
		head[p] = edge
		edge++
	}

	order := make([]int, 0, n)
	stack := make([]int, 0, n)
	stack = append(stack, 1)
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for e := head[u]; e != 0; e = nxt[e] {
			stack = append(stack, to[e])
		}
	}

	hv = make([]int64, n+5)
	htag = make([]int64, n+5)
	hl = make([]int, n+5)
	hr = make([]int, n+5)
	hdist = make([]int, n+5)
	rootHeap := make([]int, n+1)

	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		h := 0
		for e := head[u]; e != 0; e = nxt[e] {
			h = merge(h, rootHeap[to[e]])
		}
		if u != 1 {
			if h == 0 {
				h = newHeap(1)
			} else {
				var node int
				h, node = popMax(h)
				if h != 0 {
					apply(h, -1)
					if hv[h] <= 0 {
						h = 0
					}
				}
				hv[node]++
				h = merge(h, node)
			}
		}
		rootHeap[u] = h
	}

	h := rootHeap[1]
	saved := int64(0)
	limit := k + 1
	for limit > 0 && h != 0 && hv[h] > 0 {
		var node int
		h, node = popMax(h)
		saved += hv[node]
		limit--
	}

	ans := int64(2*(n-1)) - saved
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	fmt.Fprint(out, ans)
	out.Flush()
}
