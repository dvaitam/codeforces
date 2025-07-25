package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	sum     int64
	minPref int64
}

var tree []Node
var diff []int64
var n int

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func merge(left, right Node) Node {
	res := Node{}
	res.sum = left.sum + right.sum
	res.minPref = left.minPref
	if left.sum+right.minPref < res.minPref {
		res.minPref = left.sum + right.minPref
	}
	return res
}

func build(idx, l, r int) {
	if l == r {
		v := diff[l]
		tree[idx] = Node{v, min64(0, v)}
		return
	}
	m := (l + r) / 2
	build(idx*2, l, m)
	build(idx*2+1, m+1, r)
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
}

func update(idx, l, r, pos int, val int64) {
	if l == r {
		tree[idx] = Node{val, min64(0, val)}
		return
	}
	m := (l + r) / 2
	if pos <= m {
		update(idx*2, l, m, pos, val)
	} else {
		update(idx*2+1, m+1, r, pos, val)
	}
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
}

func query(idx, l, r, L, R int) Node {
	if L <= l && r <= R {
		return tree[idx]
	}
	m := (l + r) / 2
	if R <= m {
		return query(idx*2, l, m, L, R)
	}
	if L > m {
		return query(idx*2+1, m+1, r, L, R)
	}
	left := query(idx*2, l, m, L, R)
	right := query(idx*2+1, m+1, r, L, R)
	return merge(left, right)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	a := make([]int64, n+1)
	b := make([]int64, n+1)
	c := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &b[i])
	}
	for i := 1; i < n; i++ {
		fmt.Fscan(in, &c[i])
	}

	diff = make([]int64, n+1)
	tree = make([]Node, 4*(n+1))

	var sumA, sumB int64
	for i := 1; i <= n; i++ {
		sumA += a[i]
		sumB += b[i]
		diff[i] = a[i] - b[i]
	}

	build(1, 1, n)

	for ; q > 0; q-- {
		var p int
		var x, y, z int64
		fmt.Fscan(in, &p, &x, &y, &z)

		sumA += x - a[p]
		sumB += y - b[p]
		a[p] = x
		b[p] = y
		diff[p] = x - y
		update(1, 1, n, p, diff[p])
		if p < n {
			c[p] = z
		}

		totalD := sumA - sumB
		minPref := int64(0)
		if n > 1 {
			node := query(1, 1, n, 1, n-1)
			minPref = node.minPref
			if minPref > 0 {
				minPref = 0
			}
		}
		leftover := totalD - minPref
		if leftover < 0 {
			leftover = 0
		}
		ans := sumA - leftover
		fmt.Fprintln(out, ans)
	}
}
