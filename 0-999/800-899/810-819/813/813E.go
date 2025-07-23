package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	left, right int
	val         int
}

var seg []node

func clone(idx int) int {
	seg = append(seg, node{})
	newIdx := len(seg) - 1
	if idx != 0 {
		seg[newIdx] = seg[idx]
	}
	return newIdx
}

func update(idx, l, r, pos, delta int) int {
	idx = clone(idx)
	if l == r {
		seg[idx].val += delta
		return idx
	}
	m := (l + r) >> 1
	if pos <= m {
		seg[idx].left = update(seg[idx].left, l, m, pos, delta)
	} else {
		seg[idx].right = update(seg[idx].right, m+1, r, pos, delta)
	}
	seg[idx].val = seg[seg[idx].left].val + seg[seg[idx].right].val
	return idx
}

func query(idx, l, r, ql, qr int) int {
	if idx == 0 || ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return seg[idx].val
	}
	m := (l + r) >> 1
	return query(seg[idx].left, l, m, ql, qr) + query(seg[idx].right, m+1, r, ql, qr)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	a := make([]int, n+1)
	posMap := make(map[int][]int)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
		posMap[a[i]] = append(posMap[a[i]], i)
	}

	removeAt := make([][]int, n+2)
	for _, positions := range posMap {
		for i := 0; i+k < len(positions); i++ {
			rmPos := positions[i]
			at := positions[i+k]
			if at <= n {
				removeAt[at] = append(removeAt[at], rmPos)
			}
		}
	}

	seg = make([]node, 1)
	roots := make([]int, n+1)
	roots[0] = 0
	for r := 1; r <= n; r++ {
		root := update(roots[r-1], 1, n, r, 1)
		for _, idx := range removeAt[r] {
			root = update(root, 1, n, idx, -1)
		}
		roots[r] = root
	}

	var q int
	fmt.Fscan(in, &q)
	last := 0
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		l := (x+last)%n + 1
		r := (y+last)%n + 1
		if l > r {
			l, r = r, l
		}
		ans := query(roots[r], 1, n, l, r)
		fmt.Fprintln(out, ans)
		last = ans
	}
}
