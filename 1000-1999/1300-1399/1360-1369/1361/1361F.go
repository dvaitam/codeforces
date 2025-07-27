package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	l, r int
}

func buildTree(n int, w []int) ([]Node, int) {
	pos := make([]int, n) // positions of weight value
	for i := 1; i < n; i++ {
		pos[w[i-1]] = i // w[i-1] value (1..n-1)
	}
	// DSU arrays
	parent := make([]int, n+1)
	rootId := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		rootId[i] = i - 1 // leaf indices 0..n-1
	}
	find := func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}

	nodes := make([]Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = Node{-1, -1}
	}
	next := n
	for wv := 1; wv <= n-1; wv++ {
		idx := pos[wv]
		a := find(idx)
		b := find(idx + 1)
		if a == b {
			continue
		}
		nodes = append(nodes, Node{rootId[a], rootId[b]})
		newIdx := next
		next++
		parent[b] = a
		rootId[a] = newIdx
	}
	root := rootId[find(1)]
	return nodes, root
}

func mergeCount(left, right []int) ([]int, int64) {
	i, j := 0, 0
	m, n := len(left), len(right)
	merged := make([]int, 0, m+n)
	var inv int64
	for i < m && j < n {
		if left[i] <= right[j] {
			merged = append(merged, left[i])
			i++
		} else {
			merged = append(merged, right[j])
			// all remaining left elements are greater than right[j]
			inv += int64(m - i)
			j++
		}
	}
	for i < m {
		merged = append(merged, left[i])
		i++
	}
	for j < n {
		merged = append(merged, right[j])
		j++
	}
	return merged, inv
}

func solve(node int, nodes []Node, p []int) ([]int, int64) {
	if nodes[node].l == -1 && nodes[node].r == -1 {
		return []int{p[node]}, 0
	}
	leftArr, leftInv := solve(nodes[node].l, nodes, p)
	rightArr, rightInv := solve(nodes[node].r, nodes, p)
	merged, cross0 := mergeCount(leftArr, rightArr)
	totalPairs := int64(len(leftArr) * len(rightArr))
	cross1 := totalPairs - cross0
	if cross1 < cross0 {
		cross0 = cross1
	}
	return merged, leftInv + rightInv + cross0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	w := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &w[i])
	}
	nodes, root := buildTree(n, w)

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		x--
		y--
		p[x], p[y] = p[y], p[x]
		_, ans := solve(root, nodes, p)
		fmt.Fprintln(writer, ans)
	}
}
