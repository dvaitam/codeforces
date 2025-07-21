package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXN = 200005

var a [MAXN]int64

type STree struct {
	n  int
	ty int
	t  []int64
}

func oper(i, j int64, ty int) int64 {
	if ty == 1 {
		return i + j
	}
	if ty == 0 {
		if i == -1 {
			return j
		}
		if j == -1 {
			return i
		}
		if a[i] >= a[j] {
			return i
		}
		return j
	}
	if i > j {
		return i
	}
	return j
}

func newSTree(n int, ty int) *STree {
	neut := int64(0)
	if ty != 1 {
		neut = -1
	}
	t := make([]int64, 2*n+5)
	for i := range t {
		t[i] = neut
	}
	return &STree{n: n, ty: ty, t: t}
}

func (st *STree) upd(p int, v int64) {
	p += st.n
	st.t[p] = v
	for p > 1 {
		p >>= 1
		st.t[p] = oper(st.t[p<<1], st.t[p<<1|1], st.ty)
	}
}

func (st *STree) query(l, r int) int64 {
	var res int64
	if st.ty == 1 {
		res = 0
	} else {
		res = -1
	}
	l += st.n
	r += st.n
	for l < r {
		if l&1 == 1 {
			res = oper(res, st.t[l], st.ty)
			l++
		}
		if r&1 == 1 {
			r--
			res = oper(res, st.t[r], st.ty)
		}
		l >>= 1
		r >>= 1
	}
	return res
}

type pair struct {
	first  int64
	second int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		stm := newSTree(n, 0)
		sts := newSTree(n, 1)
		tim := newSTree(n, 2)
		for i := 0; i < n; i++ {
			sts.upd(i, a[i])
			stm.upd(i, int64(i))
		}
		cnt := int64(0)
		dp := make([]map[int]pair, n)
		for i := 0; i < n; i++ {
			dp[i] = make(map[int]pair)
		}
		var f func(int, int) int64
		f = func(l, r int) int64 {
			if r-l <= 0 {
				return -1
			}
			res, ok := dp[l][r]
			if ok && tim.query(l, r) <= res.second {
				return res.first
			}
			res.second = cnt
			p := int(stm.query(l, r))
			s := sts.query(l, r)
			if s-a[p] > a[p] {
				res.first = int64(r - l)
			} else {
				left := f(l, p)
				right := f(p+1, r)
				if left > right {
					res.first = left
				} else {
					res.first = right
				}
			}
			dp[l][r] = res
			return res.first
		}
		for ; q > 0; q-- {
			var ty int
			var l, r int
			fmt.Fscan(reader, &ty, &l, &r)
			l--
			if ty == 1 {
				ans := f(l, r)
				fmt.Fprintln(writer, ans)
			} else {
				a[l] = int64(r)
				sts.upd(l, int64(r))
				stm.upd(l, int64(l))
				tim.upd(l, cnt)
			}
			cnt++
		}
		for i := 0; i < n; i++ {
			dp[i] = nil
		}
	}
}
