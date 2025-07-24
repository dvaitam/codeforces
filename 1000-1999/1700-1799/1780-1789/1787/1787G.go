package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	u, v int
	w    int
}

type ColorInfo struct {
	length  int64
	nodes   []int
	blocked int
	isPath  bool
	active  bool
}

type Item struct {
	length int64
	color  int
}

type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].length > h[j].length }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	edgesByColor := make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var u, v, w, c int
		fmt.Fscan(reader, &u, &v, &w, &c)
		edgesByColor[c] = append(edgesByColor[c], Edge{u, v, w})
	}

	colors := make([]ColorInfo, n+1)
	colorsPerNode := make([][]int, n+1)

	pq := &MaxHeap{}
	heap.Init(pq)

	// preprocess each color
	for c := 1; c <= n; c++ {
		edges := edgesByColor[c]
		if len(edges) == 0 {
			continue
		}
		deg := make(map[int]int)
		nodesMap := make(map[int]struct{})
		var length int64
		for _, e := range edges {
			nodesMap[e.u] = struct{}{}
			nodesMap[e.v] = struct{}{}
			deg[e.u]++
			deg[e.v]++
			length += int64(e.w)
		}
		if len(nodesMap) != len(edges)+1 {
			continue
		}
		endpoints := 0
		ok := true
		for _, d := range deg {
			if d > 2 {
				ok = false
				break
			}
			if d == 1 {
				endpoints++
			}
		}
		if !ok || endpoints != 2 {
			continue
		}
		nodes := make([]int, 0, len(nodesMap))
		for v := range nodesMap {
			nodes = append(nodes, v)
		}
		colors[c] = ColorInfo{length: length, nodes: nodes, isPath: true, active: true}
		for _, v := range nodes {
			colorsPerNode[v] = append(colorsPerNode[v], c)
		}
		heap.Push(pq, Item{length: length, color: c})
	}

	blockedNode := make([]bool, n+1)

	getMax := func() int64 {
		for pq.Len() > 0 {
			top := (*pq)[0]
			if !colors[top.color].active || !colors[top.color].isPath {
				heap.Pop(pq)
				continue
			}
			return top.length
		}
		return 0
	}

	for i := 0; i < q; i++ {
		var p, x int
		fmt.Fscan(reader, &p, &x)
		if p == 0 {
			// block
			if !blockedNode[x] {
				blockedNode[x] = true
				for _, c := range colorsPerNode[x] {
					if !colors[c].isPath {
						continue
					}
					colors[c].blocked++
					if colors[c].blocked == 1 {
						colors[c].active = false
					}
				}
			}
		} else {
			// unblock
			if blockedNode[x] {
				blockedNode[x] = false
				for _, c := range colorsPerNode[x] {
					if !colors[c].isPath {
						continue
					}
					colors[c].blocked--
					if colors[c].blocked == 0 {
						colors[c].active = true
						heap.Push(pq, Item{length: colors[c].length, color: c})
					}
				}
			}
		}
		ans := getMax()
		fmt.Fprintln(writer, ans)
	}
}
