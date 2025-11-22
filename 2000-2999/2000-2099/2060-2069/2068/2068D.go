package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type node struct {
	w      int
	left   int
	right  int
	leafID int
}

type item struct {
	w   int
	idx int
}

type minHeap []item

func (h minHeap) Len() int { return len(h) }
func (h minHeap) Less(i, j int) bool {
	return h[i].w < h[j].w || (h[i].w == h[j].w && h[i].idx < h[j].idx)
}
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

func parseFreq(s string) int {
	if strings.IndexByte(s, '.') == -1 {
		v, _ := strconv.Atoi(s)
		return v * 10000
	}
	parts := strings.SplitN(s, ".", 2)
	whole, _ := strconv.Atoi(parts[0])
	frac := parts[1]
	for len(frac) < 4 {
		frac += "0"
	}
	if len(frac) > 4 {
		frac = frac[:4]
	}
	f, _ := strconv.Atoi(frac)
	return whole*10000 + f
}

func buildCodes(weights []int) []string {
	n := len(weights)
	nodes := make([]node, 0, 2*n)
	pq := minHeap{}
	for i, w := range weights {
		nodes = append(nodes, node{w: w, left: -1, right: -1, leafID: i})
		heap.Push(&pq, item{w: w, idx: i})
	}

	nextIdx := n
	for pq.Len() > 1 {
		a := heap.Pop(&pq).(item)
		b := heap.Pop(&pq).(item)
		// ensure heavier goes to dot (cost 1), lighter to dash (cost 2)
		dotIdx, dashIdx := b.idx, a.idx
		if nodes[a.idx].w > nodes[b.idx].w {
			dotIdx, dashIdx = a.idx, b.idx
		}
		newNode := node{w: nodes[a.idx].w + nodes[b.idx].w, left: dotIdx, right: dashIdx, leafID: -1}
		nodes = append(nodes, newNode)
		heap.Push(&pq, item{w: newNode.w, idx: nextIdx})
		nextIdx++
	}

	root := pq[0].idx
	codes := make([]string, n)
	var dfs func(idx int, pref string)
	dfs = func(idx int, pref string) {
		nd := nodes[idx]
		if nd.leafID != -1 {
			codes[nd.leafID] = pref
			return
		}
		if nd.left != -1 {
			dfs(nd.left, pref+".")
		}
		if nd.right != -1 {
			dfs(nd.right, pref+"-")
		}
	}
	dfs(root, "")
	return codes
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		weights[i] = parseFreq(s)
	}

	codes := buildCodes(weights)
	for i := 0; i < n; i++ {
		if codes[i] == "" {
			// single symbol edge case not present due to n>=2, but keep safe.
			codes[i] = "."
		}
		fmt.Fprintln(out, codes[i])
	}
}
