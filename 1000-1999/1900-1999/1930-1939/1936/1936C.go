package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type item struct {
	node int
	cost int64
}
type pq []item

func (h pq) Len() int            { return len(h) }
func (h pq) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h pq) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *pq) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *pq) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
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
		var n, m int
		fmt.Fscan(in, &n, &m)
		c := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}

		if n == 1 {
			fmt.Fprintln(out, 0)
			continue
		}

		const INF = int64(1<<62 - 1)
		dist := make([]int64, n)
		for i := range dist {
			dist[i] = INF
		}
		dist[0] = 0

		h := &pq{{0, 0}}
		for h.Len() > 0 {
			cur := heap.Pop(h).(item)
			u := cur.node
			if cur.cost > dist[u] {
				continue
			}
			if u == n-1 {
				break
			}
			for v := 0; v < n; v++ {
				if v == u {
					continue
				}
				// cost to hire pokemon v to beat pokemon u: c[v] + min_j max(0, a[u][j] - a[v][j])
				best := INF
				for j := 0; j < m; j++ {
					d := a[u][j] - a[v][j]
					if d < 0 {
						d = 0
					}
					if d < best {
						best = d
					}
				}
				w := c[v] + best
				if dist[u]+w < dist[v] {
					dist[v] = dist[u] + w
					heap.Push(h, item{v, dist[v]})
				}
			}
		}
		fmt.Fprintln(out, dist[n-1])
	}
}
