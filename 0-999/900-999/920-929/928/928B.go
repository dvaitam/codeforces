package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	left, right *node
	sum         int
}

func update(cur *node, l, r, ql, qr int) *node {
	if cur != nil && cur.sum == r-l+1 {
		return cur
	}
	if ql > r || qr < l {
		return cur
	}
	res := &node{}
	if cur != nil {
		*res = *cur
	}
	if ql <= l && r <= qr {
		res.sum = r - l + 1
		res.left = nil
		res.right = nil
		return res
	}
	m := (l + r) >> 1
	res.left = update(res.left, l, m, ql, qr)
	res.right = update(res.right, m+1, r, ql, qr)
	if res.left != nil {
		res.sum = res.left.sum
	}
	if res.right != nil {
		res.sum += res.right.sum
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(reader, &n, &k)

	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	L := make([]int, n+1)
	R := make([]int, n+1)
	for i := 1; i <= n; i++ {
		L[i] = i - k
		if L[i] < 1 {
			L[i] = 1
		}
		R[i] = i + k
		if R[i] > n {
			R[i] = n
		}
	}

	roots := make([]*node, n+1)
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		roots[i] = update(roots[a[i]], 1, n, L[i], R[i])
		ans[i] = roots[i].sum
	}

	writer := bufio.NewWriter(os.Stdout)
	for i := 1; i <= n; i++ {
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[i])
	}
	writer.WriteByte('\n')
	writer.Flush()
}
