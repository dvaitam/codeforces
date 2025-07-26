package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf int64 = -1 << 60

type Node struct {
	v [2][2]int64
}

func makeLeaf(w int64) Node {
	var n Node
	n.v[0][0] = 0
	n.v[0][1] = w
	n.v[1][0] = 0
	n.v[1][1] = negInf
	return n
}

func merge(a, b Node) Node {
	var c Node
	for s := 0; s < 2; s++ {
		for t := 0; t < 2; t++ {
			best := negInf
			for m := 0; m < 2; m++ {
				val := a.v[s][m] + b.v[m][t]
				if val > best {
					best = val
				}
			}
			c.v[s][t] = best
		}
	}
	return c
}

type SegTree struct {
	n int
	t []Node
}

func NewSegTree(arr []int64) *SegTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	t := make([]Node, 2*n)
	st := &SegTree{n: n, t: t}
	for i := 0; i < len(arr); i++ {
		t[n+i] = makeLeaf(arr[i])
	}
	for i := n - 1; i > 0; i-- {
		t[i] = merge(t[i<<1], t[i<<1|1])
	}
	return st
}

func (st *SegTree) Update(pos int, val int64) {
	i := st.n + pos
	st.t[i] = makeLeaf(val)
	for i >>= 1; i > 0; i >>= 1 {
		st.t[i] = merge(st.t[i<<1], st.t[i<<1|1])
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &a[i])
	}

	m := n - 3
	var st *SegTree
	if m > 0 {
		arr := make([]int64, m)
		for i := 0; i < m; i++ {
			arr[i] = a[i+1]
		}
		st = NewSegTree(arr)
	} else {
		st = NewSegTree([]int64{})
	}

	var sum int64
	for _, v := range a {
		sum += v
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var k int
		var x int64
		fmt.Fscan(in, &k, &x)
		sum += x - a[k-1]
		if k >= 2 && k <= n-2 {
			st.Update(k-2, x)
		}
		a[k-1] = x
		var F int64
		if m > 0 {
			res := st.t[1]
			if res.v[0][0] > res.v[0][1] {
				F = res.v[0][0]
			} else {
				F = res.v[0][1]
			}
		}
		ans := 2*sum - 2*F
		fmt.Fprintln(out, ans)
	}
}
