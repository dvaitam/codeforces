package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Op struct {
	mask uint16
	cost int
}

type Item struct {
	mask uint16
	cost int
}

type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].cost < pq[j].cost }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{}   { old := *pq; n := len(old); x := old[n-1]; *pq = old[:n-1]; return x }

var (
	costs  [5]int
	opsByK [5][]Op
	cache  [5]map[uint16]map[uint16]int
)

func buildOps() {
	for k := 1; k <= 4; k++ {
		cache[k] = make(map[uint16]map[uint16]int)
	}
	for sz := 1; sz <= 4; sz++ {
		for r := 0; r <= 4-sz; r++ {
			var m uint16
			for rr := r; rr < r+sz; rr++ {
				for c := 0; c < sz; c++ {
					m |= 1 << uint(4*c+rr)
				}
			}
			op := Op{mask: m, cost: costs[sz]}
			for k := sz; k <= 4; k++ {
				opsByK[k] = append(opsByK[k], op)
			}
		}
	}
}

func transitions(k int, state uint16) map[uint16]int {
	if t, ok := cache[k][state]; ok {
		return t
	}
	res := make(map[uint16]int)
	dist := map[uint16]int{state: 0}
	pq := &PQ{{mask: state, cost: 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		if cur.cost != dist[cur.mask] {
			continue
		}
		if cur.mask&0xF == 0 {
			key := cur.mask >> 4
			if v, ok := res[key]; !ok || cur.cost < v {
				res[key] = cur.cost
			}
		}
		for _, op := range opsByK[k] {
			next := cur.mask &^ op.mask
			nc := cur.cost + op.cost
			if d, ok := dist[next]; !ok || nc < d {
				dist[next] = nc
				heap.Push(pq, Item{mask: next, cost: nc})
			}
		}
	}
	cache[k][state] = res
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	for i := 1; i <= 4; i++ {
		fmt.Fscan(in, &costs[i])
	}
	grid := make([]string, 4)
	for i := 0; i < 4; i++ {
		fmt.Fscan(in, &grid[i])
	}
	buildOps()
	cols := make([]uint16, n+4)
	for c := 0; c < n; c++ {
		var m uint16
		for r := 0; r < 4; r++ {
			if grid[r][c] == '*' {
				m |= 1 << uint(r)
			}
		}
		cols[c] = m
	}
	start := cols[0] | cols[1]<<4 | cols[2]<<8 | cols[3]<<12
	dp := map[uint16]int{start: 0}
	for i := 0; i < n; i++ {
		k := 4
		if n-i < k {
			k = n - i
		}
		next := make(map[uint16]int)
		for st, val := range dp {
			trans := transitions(k, st)
			for p, c := range trans {
				ns := p | (cols[i+4] << 12)
				nc := val + c
				if v, ok := next[ns]; !ok || nc < v {
					next[ns] = nc
				}
			}
		}
		dp = next
	}
	fmt.Println(dp[0])
}
