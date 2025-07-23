package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// prepNode represents a preparation day cost.
type prepNode struct {
	cost int
	idx  int
}

type prepHeap []prepNode

func (h prepHeap) Len() int            { return len(h) }
func (h prepHeap) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h prepHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *prepHeap) Push(x interface{}) { *h = append(*h, x.(prepNode)) }
func (h *prepHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// pairNode represents a candidate pair (preparation, print).
type pairNode struct {
	cost int
	ai   int
	j    int
}

type pairHeap []pairNode

func (h pairHeap) Len() int            { return len(h) }
func (h pairHeap) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h pairHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *pairHeap) Push(x interface{}) { *h = append(*h, x.(pairNode)) }
func (h *pairHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	usedPrep := make([]bool, n)
	usedPrint := make([]bool, n)

	ph := &prepHeap{}
	prh := &pairHeap{}
	heap.Init(ph)
	heap.Init(prh)

	ans := int64(0)
	cnt := 0
	for j := 0; j < n; j++ {
		heap.Push(ph, prepNode{a[j], j})
		for ph.Len() > 0 && usedPrep[(*ph)[0].idx] {
			heap.Pop(ph)
		}
		if ph.Len() > 0 {
			p := (*ph)[0]
			heap.Push(prh, pairNode{p.cost + b[j], p.idx, j})
		}
		for cnt < k && prh.Len() > 0 {
			node := heap.Pop(prh).(pairNode)
			if usedPrep[node.ai] || usedPrint[node.j] {
				continue
			}
			usedPrep[node.ai] = true
			usedPrint[node.j] = true
			ans += int64(node.cost)
			cnt++
			for ph.Len() > 0 && usedPrep[(*ph)[0].idx] {
				heap.Pop(ph)
			}
		}
	}
	for cnt < k && prh.Len() > 0 {
		node := heap.Pop(prh).(pairNode)
		if usedPrep[node.ai] || usedPrint[node.j] {
			continue
		}
		usedPrep[node.ai] = true
		usedPrint[node.j] = true
		ans += int64(node.cost)
		cnt++
	}

	if cnt < k {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}
