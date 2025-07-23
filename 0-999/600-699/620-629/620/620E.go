package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type Node struct {
	mask uint64
	lazy uint64
}

var (
	n, m   int
	colors []int
	adj    [][]int
	tin    []int
	tout   []int
	order  []int
	st     []Node
)

func apply(pos int, mask uint64) {
	st[pos].mask = mask
	st[pos].lazy = mask
}

func push(pos int) {
	if st[pos].lazy != 0 {
		apply(pos<<1, st[pos].lazy)
		apply(pos<<1|1, st[pos].lazy)
		st[pos].lazy = 0
	}
}

func build(pos, l, r int, arr []uint64) {
	if r-l == 1 {
		st[pos].mask = arr[l]
		return
	}
	mid := (l + r) >> 1
	build(pos<<1, l, mid, arr)
	build(pos<<1|1, mid, r, arr)
	st[pos].mask = st[pos<<1].mask | st[pos<<1|1].mask
}

func update(pos, l, r, ql, qr int, mask uint64) {
	if ql <= l && r <= qr {
		apply(pos, mask)
		return
	}
	push(pos)
	mid := (l + r) >> 1
	if ql < mid {
		update(pos<<1, l, mid, ql, qr, mask)
	}
	if qr > mid {
		update(pos<<1|1, mid, r, ql, qr, mask)
	}
	st[pos].mask = st[pos<<1].mask | st[pos<<1|1].mask
}

func query(pos, l, r, ql, qr int) uint64 {
	if ql <= l && r <= qr {
		return st[pos].mask
	}
	push(pos)
	mid := (l + r) >> 1
	var res uint64
	if ql < mid {
		res |= query(pos<<1, l, mid, ql, qr)
	}
	if qr > mid {
		res |= query(pos<<1|1, mid, r, ql, qr)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &m)
	colors = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &colors[i])
	}
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	tin = make([]int, n+1)
	tout = make([]int, n+1)
	order = make([]int, 0, n)
	parent := make([]int, n+1)
	it := make([]int, n+1)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		if it[v] == 0 {
			tin[v] = len(order)
			order = append(order, v)
		}
		if it[v] < len(adj[v]) {
			to := adj[v][it[v]]
			it[v]++
			if to == parent[v] {
				continue
			}
			parent[to] = v
			stack = append(stack, to)
		} else {
			tout[v] = len(order)
			stack = stack[:len(stack)-1]
		}
	}

	arr := make([]uint64, n)
	for i, v := range order {
		arr[i] = 1 << uint(colors[v]-1)
	}
	st = make([]Node, 4*n)
	build(1, 0, n, arr)

	for i := 0; i < m; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var v, c int
			fmt.Fscan(reader, &v, &c)
			l, r := tin[v], tout[v]
			mask := uint64(1) << uint(c-1)
			update(1, 0, n, l, r, mask)
		} else {
			var v int
			fmt.Fscan(reader, &v)
			l, r := tin[v], tout[v]
			mask := query(1, 0, n, l, r)
			fmt.Fprintln(writer, bits.OnesCount64(mask))
		}
	}
}
