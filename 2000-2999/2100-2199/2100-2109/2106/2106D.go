package main

import (
	"bufio"
	"fmt"
	"os"
)

const infVal int64 = 1<<62 - 1

type segTree struct {
	size int
	data []int64
}

func newSegTree(vals []int) *segTree {
	size := 1
	for size < len(vals) {
		size <<= 1
	}
	data := make([]int64, 2*size)
	for i := range data {
		data[i] = infVal
	}
	for i, v := range vals {
		data[size+i] = int64(v)
	}
	for i := size - 1; i >= 1; i-- {
		if data[i<<1] < data[i<<1|1] {
			data[i] = data[i<<1]
		} else {
			data[i] = data[i<<1|1]
		}
	}
	return &segTree{size: size, data: data}
}

func (st *segTree) query(l, r int) int64 {
	if l > r {
		return infVal
	}
	l += st.size
	r += st.size
	res := infVal
	for l <= r {
		if l&1 == 1 {
			if st.data[l] < res {
				res = st.data[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.data[r] < res {
				res = st.data[r]
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}

		pref := make([]int, n+1)
		ptr := 0
		for i := 0; i < n; i++ {
			if ptr < m && a[i] >= b[ptr] {
				ptr++
			}
			pref[i+1] = ptr
		}

		if pref[n] == m {
			fmt.Fprintln(out, 0)
			continue
		}

		suff := make([]int, n+1)
		ptr = m - 1
		matched := 0
		for i := n - 1; i >= 0; i-- {
			if ptr >= 0 && a[i] >= b[ptr] {
				ptr--
				matched++
			}
			suff[i] = matched
		}

		st := newSegTree(b)
		ans := infVal
		for i := 0; i <= n; i++ {
			L := m - suff[i]
			if L < 1 {
				L = 1
			}
			R := pref[i] + 1
			if R > m {
				R = m
			}
			if L > R {
				continue
			}
			val := st.query(L-1, R-1)
			if val < ans {
				ans = val
			}
		}

		if ans >= infVal {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}
