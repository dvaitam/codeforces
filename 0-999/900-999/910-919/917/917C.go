package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type State struct {
	cost  int64
	pos   int
	mask  int
	index int // heap index
}

type PQ []*State

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].cost < pq[j].cost }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i]; pq[i].index = i; pq[j].index = j }
func (pq *PQ) Push(x interface{}) { item := x.(*State); item.index = len(*pq); *pq = append(*pq, item) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var x, k, n, q int
	if _, err := fmt.Fscan(in, &x, &k, &n, &q); err != nil {
		return
	}
	c := make([]int64, k+1)
	for i := 1; i <= k; i++ {
		fmt.Fscan(in, &c[i])
	}
	special := make(map[int]int64)
	for i := 0; i < q; i++ {
		var p int
		var w int64
		fmt.Fscan(in, &p, &w)
		special[p] = w
	}

	startMask := 0
	for i := 0; i < x; i++ {
		startMask |= 1 << i
	}
	targetPos := n - x + 1
	targetMask := startMask

	dist := make(map[[2]int]int64)
	pq := &PQ{}
	heap.Push(pq, &State{cost: 0, pos: 1, mask: startMask})
	dist[[2]int{1, startMask}] = 0

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*State)
		key := [2]int{cur.pos, cur.mask}
		if cur.cost != dist[key] {
			continue
		}
		if cur.pos == targetPos && cur.mask == targetMask {
			fmt.Println(cur.cost)
			return
		}
		for j := 1; j <= k; j++ {
			if j < k && (cur.mask>>j)&1 == 1 {
				continue
			}
			// find next occupied >0
			nextPos := k + 1
			for i := 1; i < k; i++ {
				if (cur.mask>>i)&1 == 1 {
					nextPos = i
					break
				}
			}
			delta := j
			if nextPos < j {
				delta = nextPos
			}
			newPos := cur.pos + delta
			if newPos > targetPos {
				continue
			}
			newMask := ((cur.mask &^ 1) | (1 << j)) >> delta
			newMask &= (1 << k) - 1
			addCost := c[j]
			if w, ok := special[cur.pos+j]; ok {
				addCost += w
			}
			newCost := cur.cost + addCost
			nk := [2]int{newPos, newMask}
			if d, ok := dist[nk]; !ok || newCost < d {
				dist[nk] = newCost
				heap.Push(pq, &State{cost: newCost, pos: newPos, mask: newMask})
			}
		}
	}
}
