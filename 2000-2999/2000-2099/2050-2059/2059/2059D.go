package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type edge struct {
	u, v int
}

type item struct {
	idx  int
	dist int64
}

type priorityQueue []item

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(item))
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, s1, s2 int
		fmt.Fscan(in, &n, &s1, &s2)

		g1 := make([][]int, n+1)
		g2 := make([][]int, n+1)
		g2Mat := make([][]bool, n+1)
		for i := range g2Mat {
			g2Mat[i] = make([]bool, n+1)
		}

		var m1 int
		fmt.Fscan(in, &m1)
		edges1 := make([]edge, 0, m1)
		for i := 0; i < m1; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			g1[a] = append(g1[a], b)
			g1[b] = append(g1[b], a)
			edges1 = append(edges1, edge{a, b})
		}

		var m2 int
		fmt.Fscan(in, &m2)
		for i := 0; i < m2; i++ {
			var c, d int
			fmt.Fscan(in, &c, &d)
			g2[c] = append(g2[c], d)
			g2[d] = append(g2[d], c)
			g2Mat[c][d] = true
			g2Mat[d][c] = true
		}

		commonDeg := make([]int, n+1)
		for _, e := range edges1 {
			if g2Mat[e.u][e.v] {
				commonDeg[e.u]++
				commonDeg[e.v]++
			}
		}

		hasCommon := false
		for v := 1; v <= n; v++ {
			if commonDeg[v] > 0 {
				hasCommon = true
				break
			}
		}

		if !hasCommon {
			fmt.Fprintln(out, -1)
			continue
		}

		totalStates := n * n
		dist := make([]int64, totalStates)
		for i := range dist {
			dist[i] = inf
		}

		startIdx := (s1-1)*n + (s2 - 1)
		dist[startIdx] = 0
		pq := priorityQueue{{idx: startIdx, dist: 0}}
		heap.Init(&pq)

		answer := inf
		for pq.Len() > 0 {
			cur := heap.Pop(&pq).(item)
			if cur.dist != dist[cur.idx] {
				continue
			}
			u := cur.idx/n + 1
			v := cur.idx%n + 1

			if u == v && commonDeg[u] > 0 {
				answer = cur.dist
				break
			}

			for _, nu := range g1[u] {
				for _, nv := range g2[v] {
					nextIdx := (nu-1)*n + (nv - 1)
					nd := cur.dist + int64(absInt(nu-nv))
					if nd < dist[nextIdx] {
						dist[nextIdx] = nd
						heap.Push(&pq, item{idx: nextIdx, dist: nd})
					}
				}
			}
		}

		if answer == inf {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, answer)
		}
	}
}
