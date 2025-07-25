package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const inf int64 = 1<<63 - 1

type query struct {
	l, r, idx int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	x := make([]int64, n+1)
	w := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &x[i], &w[i])
	}
	qs := make([]query, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &qs[i].l, &qs[i].r)
		qs[i].idx = i
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].r < qs[j].r })

	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]int64, 2*size)
	for i := range seg {
		seg[i] = inf
	}
	update := func(pos int, val int64) {
		p := pos + size - 1
		if val < seg[p] {
			seg[p] = val
			p >>= 1
			for p > 0 {
				newVal := seg[p<<1]
				if seg[p<<1|1] < newVal {
					newVal = seg[p<<1|1]
				}
				if seg[p] == newVal {
					break
				}
				seg[p] = newVal
				p >>= 1
			}
		}
	}
	querySeg := func(l, r int) int64 {
		if l > r {
			return inf
		}
		l += size - 1
		r += size - 1
		res := inf
		for l <= r {
			if l&1 == 1 {
				if seg[l] < res {
					res = seg[l]
				}
				l++
			}
			if r&1 == 0 {
				if seg[r] < res {
					res = seg[r]
				}
				r--
			}
			l >>= 1
			r >>= 1
		}
		return res
	}

	ans := make([]int64, q)
	stack := make([]int, 0)
	qi := 0
	for i := 1; i <= n; i++ {
		for len(stack) > 0 {
			j := stack[len(stack)-1]
			val := (x[i] - x[j]) * (w[i] + w[j])
			update(j, val)
			if w[j] <= w[i] {
				break
			}
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
		for qi < q && qs[qi].r == i {
			l := qs[qi].l
			ans[qs[qi].idx] = querySeg(l, i-1)
			qi++
		}
	}
	for i := 0; i < q; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
