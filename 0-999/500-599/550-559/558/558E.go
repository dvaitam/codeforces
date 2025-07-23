package main

import (
	"bufio"
	"fmt"
	"os"
)

const alpha = 26

var n int
var trees [alpha][]int
var lazy [alpha][]int

func push(c int, node, l, r int) {
	if lazy[c][node] != -1 {
		val := lazy[c][node]
		trees[c][node] = (r - l + 1) * val
		if l != r {
			lazy[c][node*2] = val
			lazy[c][node*2+1] = val
		}
		lazy[c][node] = -1
	}
}

func update(c int, node, l, r, ql, qr, val int) {
	push(c, node, l, r)
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		lazy[c][node] = val
		push(c, node, l, r)
		return
	}
	mid := (l + r) / 2
	update(c, node*2, l, mid, ql, qr, val)
	update(c, node*2+1, mid+1, r, ql, qr, val)
	trees[c][node] = trees[c][node*2] + trees[c][node*2+1]
}

func query(c int, node, l, r, ql, qr int) int {
	push(c, node, l, r)
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return trees[c][node]
	}
	mid := (l + r) / 2
	return query(c, node*2, l, mid, ql, qr) + query(c, node*2+1, mid+1, r, ql, qr)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	var q int
	fmt.Fscan(in, &q)
	var s string
	fmt.Fscan(in, &s)

	for i := 0; i < alpha; i++ {
		trees[i] = make([]int, 4*(n+2))
		lazy[i] = make([]int, 4*(n+2))
		for j := range lazy[i] {
			lazy[i][j] = -1
		}
	}

	for i, ch := range s {
		c := int(ch - 'a')
		update(c, 1, 1, n, i+1, i+1, 1)
	}

	for ; q > 0; q-- {
		var l, r, k int
		fmt.Fscan(in, &l, &r, &k)
		if l > r {
			l, r = r, l
		}
		cnt := make([]int, alpha)
		for c := 0; c < alpha; c++ {
			cnt[c] = query(c, 1, 1, n, l, r)
			if cnt[c] > 0 {
				update(c, 1, 1, n, l, r, 0)
			}
		}
		idx := l
		if k == 1 {
			for c := 0; c < alpha; c++ {
				if cnt[c] > 0 {
					update(c, 1, 1, n, idx, idx+cnt[c]-1, 1)
					idx += cnt[c]
				}
			}
		} else {
			for c := alpha - 1; c >= 0; c-- {
				if cnt[c] > 0 {
					update(c, 1, 1, n, idx, idx+cnt[c]-1, 1)
					idx += cnt[c]
				}
			}
		}
	}

	out := make([]byte, n)
	for i := 1; i <= n; i++ {
		for c := 0; c < alpha; c++ {
			if query(c, 1, 1, n, i, i) > 0 {
				out[i-1] = byte('a' + c)
				break
			}
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, string(out))
	writer.Flush()
}
