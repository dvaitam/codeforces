package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to      int
	l, r    int64
	incWith bool // true if edge touches special node k
}

type Item struct {
	v    int
	dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

const inf int64 = 1 << 60

// mode: 0 -> all l; 1 -> r on edges incident to k else l; 2 -> l on incident, r otherwise; 3 -> all r
func dijkstra(n int, g [][]Edge, src int, skipK int, mode int) []int64 {
	dist := make([]int64, n)
	for i := range dist {
		dist[i] = inf
	}
	if src == skipK {
		return dist
	}
	dist[src] = 0
	pq := PriorityQueue{{v: src, dist: 0}}
	for len(pq) > 0 {
		cur := heap.Pop(&pq).(Item)
		if cur.dist != dist[cur.v] {
			continue
		}
		for _, e := range g[cur.v] {
			if e.to == skipK {
				continue
			}
			w := e.l
			switch mode {
			case 1:
				if e.incWith {
					w = e.r
				}
			case 2:
				if e.incWith {
					w = e.l
				} else {
					w = e.r
				}
			case 3:
				w = e.r
			}
			if dist[e.to] > cur.dist+w {
				dist[e.to] = cur.dist + w
				heap.Push(&pq, Item{v: e.to, dist: dist[e.to]})
			}
		}
	}
	return dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		k-- // zero-based

		g := make([][]Edge, n)
		for i := 0; i < m; i++ {
			var u, v int
			var l, r int64
			fmt.Fscan(in, &u, &v, &l, &r)
			u--
			v--
			inc := u == k || v == k
			g[u] = append(g[u], Edge{to: v, l: l, r: r, incWith: inc})
			g[v] = append(g[v], Edge{to: u, l: l, r: r, incWith: inc})
		}

		// 1) all weights = l, direct check
		d1AllL := dijkstra(n, g, 0, -1, 0)
		dkAllL := dijkstra(n, g, k, -1, 0)
		if d1AllL[n-1] != d1AllL[k]+dkAllL[n-1] {
			fmt.Fprintln(out, "YES")
			continue
		}

		// 2) make path avoiding k as small as possible, via k as large as possible
		avoidMin := dijkstra(n, g, 0, k, 0)[n-1]
		d1ViaMax := dijkstra(n, g, 0, -1, 1)[k]
		dkViaMax := dijkstra(n, g, k, -1, 1)[n-1]
		if avoidMin < d1ViaMax+dkViaMax {
			fmt.Fprintln(out, "YES")
			continue
		}

		// 3) make path avoiding k as large as possible, via k as small as possible (with incident edges cheap)
		avoidMax := dijkstra(n, g, 0, k, 3)[n-1]
		d1ViaMin := dijkstra(n, g, 0, -1, 2)[k]
		dkViaMin := dijkstra(n, g, k, -1, 2)[n-1]
		if d1ViaMin+dkViaMin < avoidMax {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
