package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const INF = 1e9
const MOD = 998244353

// Item for priority queue: state (node u, parity p) with cost (flips f, moves m)
type Item struct {
	u, p, f, m int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].f != pq[j].f {
		return pq[i].f < pq[j].f
	}
	return pq[i].m < pq[j].m
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Item))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func powmod(a, e int) int {
	res := 1
	base := a % MOD
	for e > 0 {
		if e&1 != 0 {
			res = int((int64(res) * int64(base)) % MOD)
		}
		base = int((int64(base) * int64(base)) % MOD)
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	// graph: original and reversed adjacency lists
	g0 := make([][]int, n+1)
	g1 := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g0[u] = append(g0[u], v)
		g1[v] = append(g1[v], u)
	}
	// dist flips and moves
	flips := make([][2]int, n+1)
	moves := make([][2]int, n+1)
	for i := 1; i <= n; i++ {
		flips[i][0], flips[i][1] = INF, INF
		moves[i][0], moves[i][1] = INF, INF
	}
	// initial state
	pq := &PriorityQueue{}
	heap.Init(pq)
	flips[1][0] = 0
	moves[1][0] = 0
	heap.Push(pq, &Item{u: 1, p: 0, f: 0, m: 0})

	for pq.Len() > 0 {
		it := heap.Pop(pq).(*Item)
		u, p, f, mm := it.u, it.p, it.f, it.m
		// skip if outdated
		if f != flips[u][p] || mm != moves[u][p] {
			continue
		}
		// try flip
		np := 1 - p
		nf := f + 1
		nm := mm
		if nf < flips[u][np] || (nf == flips[u][np] && nm < moves[u][np]) {
			flips[u][np] = nf
			moves[u][np] = nm
			heap.Push(pq, &Item{u: u, p: np, f: nf, m: nm})
		}
		// try movements
		var adj [][]int
		if p == 0 {
			adj = g0
		} else {
			adj = g1
		}
		for _, v := range adj[u] {
			nf = f
			nm = mm + 1
			if nf < flips[v][p] || (nf == flips[v][p] && nm < moves[v][p]) {
				flips[v][p] = nf
				moves[v][p] = nm
				heap.Push(pq, &Item{u: v, p: p, f: nf, m: nm})
			}
		}
	}
	// choose best at n
	f0, m0 := flips[n][0], moves[n][0]
	f1, m1 := flips[n][1], moves[n][1]
	bf, bm := f0, m0
	if f1 < bf || (f1 == bf && m1 < bm) {
		bf, bm = f1, m1
	}
	// ans = (2^bf -1 + bm) mod MOD
	t := powmod(2, bf)
	ans := (int((t-1+MOD)%MOD) + bm) % MOD
	fmt.Println(ans)
}
