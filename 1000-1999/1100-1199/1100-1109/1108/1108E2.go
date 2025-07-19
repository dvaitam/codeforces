package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func nextInt() int {
	var x int
	var c byte
	var neg bool
	// skip non-digit
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return x
		}
		c = b
		if (c >= '0' && c <= '9') || c == '-' {
			break
		}
	}
	if c == '-' {
		neg = true
	} else {
		x = int(c - '0')
	}
	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		c = b
		if c < '0' || c > '9' {
			break
		}
		x = x*10 + int(c-'0')
	}
	if neg {
		return -x
	}
	return x
}

// Segment tree for range add and range minimum query
type SegTree struct {
	n    int
	it   []int
	lazy []int
}

func NewSegTree(n int, a []int) *SegTree {
	size := 4*n + 5
	st := &SegTree{n: n, it: make([]int, size), lazy: make([]int, size)}
	st.build(1, 1, n, a)
	return st
}

func (st *SegTree) build(id, l, r int, a []int) {
	if l == r {
		st.it[id] = a[l]
		return
	}
	m := (l + r) >> 1
	lc, rc := id<<1, id<<1|1
	st.build(lc, l, m, a)
	st.build(rc, m+1, r, a)
	// use minimum as in original C++ code MAX= min
	if st.it[lc] < st.it[rc] {
		st.it[id] = st.it[lc]
	} else {
		st.it[id] = st.it[rc]
	}
}

func (st *SegTree) apply(id, v int) {
	st.it[id] += v
	st.lazy[id] += v
}

func (st *SegTree) push(id int) {
	if st.lazy[id] != 0 {
		st.apply(id<<1, st.lazy[id])
		st.apply(id<<1|1, st.lazy[id])
		st.lazy[id] = 0
	}
}

// update range [ql, qr] by v
func (st *SegTree) update(id, l, r, ql, qr, v int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(id, v)
		return
	}
	st.push(id)
	m := (l + r) >> 1
	st.update(id<<1, l, m, ql, qr, v)
	st.update(id<<1|1, m+1, r, ql, qr, v)
	lc, rc := id<<1, id<<1|1
	if st.it[lc] < st.it[rc] {
		st.it[id] = st.it[lc]
	} else {
		st.it[id] = st.it[rc]
	}
}

// query full range minimum
func (st *SegTree) queryAll() int {
	return st.it[1]
}

func main() {
	defer writer.Flush()
	n := nextInt()
	m := nextInt()
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = nextInt()
	}
	// prepare structures
	L := make([]int, m)
	R := make([]int, m)
	v := make([][]int, n+2)
	g := make([][]int, n+2)
	// build segment tree
	st := NewSegTree(n, a)
	for i := 0; i < m; i++ {
		l := nextInt()
		r := nextInt()
		L[i], R[i] = l, r
		v[r] = append(v[r], l)
		g[l] = append(g[l], r)
		st.update(1, 1, n, l, r, -1)
	}
	sol := 0
	best := 1
	for i := 1; i <= n; i++ {
		// segments ending at i-1
		for _, l := range v[i-1] {
			st.update(1, 1, n, l, i-1, -1)
		}
		// segments starting at i
		for _, r := range g[i] {
			st.update(1, 1, n, i, r, 1)
		}
		cur := a[i] - st.queryAll()
		if cur > sol {
			sol = cur
			best = i
		}
	}
	fmt.Fprintln(writer, sol)
	if sol == 0 {
		fmt.Fprint(writer, 0)
		return
	}
	var out []int
	for i := 0; i < m; i++ {
		if R[i] < best || L[i] > best {
			out = append(out, i+1)
		}
	}
	fmt.Fprintln(writer, len(out))
	for _, idx := range out {
		fmt.Fprint(writer, idx, " ")
	}
}
