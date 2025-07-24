package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Node struct {
	left, right int
	val         int
}

var nodes []Node
var roots []int
var values []int
var pos map[int]int

func update(prev, l, r, idx int) int {
	cur := len(nodes)
	nodes = append(nodes, nodes[prev])
	if l == r {
		nodes[cur].val ^= 1
		return cur
	}
	mid := (l + r) >> 1
	if idx <= mid {
		nodes[cur].left = update(nodes[prev].left, l, mid, idx)
	} else {
		nodes[cur].right = update(nodes[prev].right, mid+1, r, idx)
	}
	nodes[cur].val = nodes[nodes[cur].left].val ^ nodes[nodes[cur].right].val
	return cur
}

func query(u, v, l, r int) int {
	if nodes[u].val^nodes[v].val == 0 {
		return -1
	}
	if l == r {
		return l
	}
	mid := (l + r) >> 1
	res := query(nodes[u].left, nodes[v].left, l, mid)
	if res != -1 {
		return res
	}
	return query(nodes[u].right, nodes[v].right, mid+1, r)
}

func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	sort.Ints(a)
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n+1)
	vals := make([]int, n)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
		vals[i-1] = arr[i]
	}
	values = uniqueInts(vals)
	pos = make(map[int]int, len(values))
	for i, v := range values {
		pos[v] = i + 1
	}

	nodes = make([]Node, 1)
	roots = make([]int, n+1)
	for i := 1; i <= n; i++ {
		roots[i] = update(roots[i-1], 1, len(values), pos[arr[i]])
	}

	var q int
	fmt.Fscan(reader, &q)
	ans := 0
	for ; q > 0; q-- {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		l := a ^ ans
		r := b ^ ans
		if l > r {
			l, r = r, l
		}
		idx := query(roots[r], roots[l-1], 1, len(values))
		if idx == -1 {
			ans = 0
		} else {
			ans = values[idx-1]
		}
		fmt.Fprintln(writer, ans)
	}
}
