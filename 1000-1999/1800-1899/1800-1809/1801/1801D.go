package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to   int
	cost int64
}

type State struct {
	perf  int64
	node  int
	coins int64
}

type PQ []State

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	if pq[i].perf == pq[j].perf {
		return pq[i].coins > pq[j].coins
	}
	return pq[i].perf < pq[j].perf
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(State)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	x := old[len(old)-1]
	*pq = old[:len(old)-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		var p int64
		fmt.Fscan(in, &n, &m, &p)
		w := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &w[i])
		}
		g := make([][]Edge, n+1)
		for i := 0; i < m; i++ {
			var a, b int
			var s int64
			fmt.Fscan(in, &a, &b, &s)
			g[a] = append(g[a], Edge{b, s})
		}

		pq := &PQ{}
		heap.Init(pq)
		heap.Push(pq, State{0, 1, p})

		best := make([]map[int64]int64, n+1)
		for i := 1; i <= n; i++ {
			best[i] = make(map[int64]int64)
		}
		best[1][0] = p
		ans := int64(-1)

		for pq.Len() > 0 {
			st := heap.Pop(pq).(State)
			if st.node == n {
				ans = st.perf
				break
			}
			// check if current state is dominated
			if c, ok := best[st.node][st.perf]; ok && c > st.coins {
				continue
			}
			// move using edges
			for _, e := range g[st.node] {
				if st.coins >= e.cost {
					ncoins := st.coins - e.cost
					if c, ok := best[e.to][st.perf]; !ok || ncoins > c {
						best[e.to][st.perf] = ncoins
						heap.Push(pq, State{st.perf, e.to, ncoins})
					}
				}
			}
			// perform once
			nperf := st.perf + 1
			ncoins := st.coins + w[st.node]
			if c, ok := best[st.node][nperf]; !ok || ncoins > c {
				best[st.node][nperf] = ncoins
				heap.Push(pq, State{nperf, st.node, ncoins})
			}
		}

		fmt.Fprintln(out, ans)
	}
}
