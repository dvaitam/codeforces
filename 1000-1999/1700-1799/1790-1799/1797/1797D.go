package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Child struct {
	size  int
	index int
}

type ChildHeap []Child

func (h ChildHeap) Len() int { return len(h) }
func (h ChildHeap) Less(i, j int) bool {
	if h[i].size == h[j].size {
		return h[i].index < h[j].index
	}
	return h[i].size > h[j].size
}
func (h ChildHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *ChildHeap) Push(x interface{}) { *h = append(*h, x.(Child)) }
func (h *ChildHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

type Node struct {
	parent    int
	size      int
	sum       int64
	childSize map[int]int
	heap      ChildHeap
}

var nodes []Node

func getHeavy(v int) int {
	h := &nodes[v].heap
	m := nodes[v].childSize
	for h.Len() > 0 {
		top := (*h)[0]
		if sz, ok := m[top.index]; ok && sz == top.size {
			return top.index
		}
		heap.Pop(h)
	}
	return 0
}

func rotate(x int) {
	hv := getHeavy(x)
	if hv == 0 {
		return
	}
	p := nodes[x].parent
	oldXSize := nodes[x].size
	oldXSum := nodes[x].sum
	oldHVSize := nodes[hv].size
	oldHVSum := nodes[hv].sum

	// parent p: replace child x with hv
	if p != 0 {
		delete(nodes[p].childSize, x)
		nodes[p].childSize[hv] = oldXSize
		heap.Push(&nodes[p].heap, Child{size: oldXSize, index: hv})
	}

	// update parents
	nodes[hv].parent = p
	nodes[x].parent = hv

	// remove hv from children of x
	delete(nodes[x].childSize, hv)

	// update sizes and sums
	nodes[x].size = oldXSize - oldHVSize
	nodes[x].sum = oldXSum - oldHVSum
	nodes[hv].size = oldXSize
	nodes[hv].sum = oldXSum

	// add x as child of hv with new size
	nodes[hv].childSize[x] = nodes[x].size
	heap.Push(&nodes[hv].heap, Child{size: nodes[x].size, index: x})
}

func dfs(v, p int, adj [][]int, vals []int64) {
	nodes[v].parent = p
	nodes[v].size = 1
	nodes[v].sum = vals[v]
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		dfs(to, v, adj, vals)
		nodes[v].size += nodes[to].size
		nodes[v].sum += nodes[to].sum
		nodes[v].childSize[to] = nodes[to].size
		heap.Push(&nodes[v].heap, Child{size: nodes[to].size, index: to})
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	vals := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &vals[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	nodes = make([]Node, n+1)
	for i := range nodes {
		nodes[i].childSize = make(map[int]int)
	}
	dfs(1, 0, adj, vals)

	for i := 0; i < m; i++ {
		var t, x int
		fmt.Fscan(reader, &t, &x)
		if t == 1 {
			fmt.Fprintln(writer, nodes[x].sum)
		} else {
			rotate(x)
		}
	}
}
