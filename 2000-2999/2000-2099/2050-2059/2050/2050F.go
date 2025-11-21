package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

type segTree struct {
	n  int
	tr []int
}

func newSeg(arr []int) *segTree {
	n := len(arr)
	tr := make([]int, 4*n)
	st := &segTree{n: n, tr: tr}
	if n > 0 {
		st.build(1, 0, n-1, arr)
	}
	return st
}

func (st *segTree) build(node, l, r int, arr []int) {
	if l == r {
		st.tr[node] = arr[l]
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid, arr)
	st.build(node<<1|1, mid+1, r, arr)
	st.tr[node] = gcd(st.tr[node<<1], st.tr[node<<1|1])
}

func (st *segTree) query(l, r int) int {
	if l > r || st.n == 0 {
		return 0
	}
	return st.q(1, 0, st.n-1, l, r)
}

func (st *segTree) q(node, nl, nr, l, r int) int {
	if r < nl || nr < l {
		return 0
	}
	if l <= nl && nr <= r {
		return st.tr[node]
	}
	mid := (nl + nr) >> 1
	left := st.q(node<<1, nl, mid, l, r)
	right := st.q(node<<1|1, mid+1, nr, l, r)
	return gcd(left, right)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		diff := make([]int, 0)
		if n > 1 {
			diff = make([]int, n-1)
			for i := 0; i < n-1; i++ {
				d := a[i+1] - a[i]
				if d < 0 {
					d = -d
				}
				diff[i] = d
			}
		}
		st := newSeg(diff)
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			r--
			if l == r {
				fmt.Fprint(out, 0, " ")
				continue
			}
			g := st.query(l, r-1)
			fmt.Fprint(out, g, " ")
		}
		fmt.Fprintln(out)
	}
}
