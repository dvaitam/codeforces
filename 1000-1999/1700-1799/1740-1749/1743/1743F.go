package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

// Matrix represents a 2x2 matrix.
type Matrix struct {
	a, b, c, d int64
}

func mul(x, y Matrix) Matrix {
	return Matrix{
		a: (x.a*y.a + x.b*y.c) % MOD,
		b: (x.a*y.b + x.b*y.d) % MOD,
		c: (x.c*y.a + x.d*y.c) % MOD,
		d: (x.c*y.b + x.d*y.d) % MOD,
	}
}

var (
	M0 = Matrix{3, 1, 0, 2}
	M1 = Matrix{1, 1, 2, 2}
	I  = Matrix{1, 0, 0, 1}
)

type SegTree struct {
	n    int
	size int
	tree []Matrix
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]Matrix, 2*size)
	st := &SegTree{n: n, size: size, tree: tree}
	for i := 0; i < n; i++ {
		st.tree[st.size+i] = M0
	}
	for i := st.size - 1; i >= 1; i-- {
		st.tree[i] = mul(st.tree[i<<1|1], st.tree[i<<1])
	}
	return st
}

func (st *SegTree) Set(pos int, val Matrix) {
	p := pos + st.size
	st.tree[p] = val
	for p > 1 {
		p >>= 1
		st.tree[p] = mul(st.tree[p<<1|1], st.tree[p<<1])
	}
}

func (st *SegTree) Root() Matrix {
	if st.n == 0 {
		return I
	}
	return st.tree[1]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	l := make([]int, n)
	r := make([]int, n)
	maxR := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &l[i], &r[i])
		if r[i] > maxR {
			maxR = r[i]
		}
	}

	start := make([][]int, maxR+2)
	end := make([][]int, maxR+2)
	for i := 0; i < n; i++ {
		start[l[i]] = append(start[l[i]], i)
		if r[i]+1 <= maxR+1 {
			end[r[i]+1] = append(end[r[i]+1], i)
		}
	}

	st := NewSegTree(n - 1)
	var dp0, dp1 int64 = 1, 0 // values after processing first set

	ans := int64(0)
	for x := 0; x <= maxR; x++ {
		for _, idx := range end[x] {
			if idx == 0 {
				dp0, dp1 = 1, 0
			} else {
				st.Set(idx-1, M0)
			}
		}
		for _, idx := range start[x] {
			if idx == 0 {
				dp0, dp1 = 0, 1
			} else {
				st.Set(idx-1, M1)
			}
		}
		root := st.Root()
		res1 := (root.c*dp0 + root.d*dp1) % MOD
		ans = (ans + res1) % MOD
	}

	fmt.Fprintln(writer, ans%MOD)
}
