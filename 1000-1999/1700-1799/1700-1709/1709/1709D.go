package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type SegTree struct {
	n    int
	data []int
}

func NewSegTree(arr []int) *SegTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	data := make([]int, 2*n)
	for i := 0; i < len(arr); i++ {
		data[n+i] = arr[i]
	}
	for i := n - 1; i > 0; i-- {
		if data[2*i] > data[2*i+1] {
			data[i] = data[2*i]
		} else {
			data[i] = data[2*i+1]
		}
	}
	return &SegTree{n: n, data: data}
}

func (t *SegTree) Query(l, r int) int {
	l += t.n
	r += t.n
	res := 0
	for l <= r {
		if l%2 == 1 {
			if t.data[l] > res {
				res = t.data[l]
			}
			l++
		}
		if r%2 == 0 {
			if t.data[r] > res {
				res = t.data[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &a[i])
	}
	st := NewSegTree(a)

	var q int
	fmt.Fscan(in, &q)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; q > 0; q-- {
		var xs, ys, xf, yf, k int
		fmt.Fscan(in, &xs, &ys, &xf, &yf, &k)
		ys--
		yf--
		if (xs-xf)%k != 0 || (ys-yf)%k != 0 {
			fmt.Fprintln(out, "NO")
			continue
		}
		high := xs + (n-xs)/k*k
		l := ys
		r := yf
		if l > r {
			l, r = r, l
		}
		maxA := st.Query(l, r)
		if maxA < high {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
