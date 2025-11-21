package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type edge struct {
	to int
}

type state struct {
	time int64
	wait int64
	id   int
}

type priorityQueue []state

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].time == pq[j].time {
		return pq[i].wait < pq[j].wait
	}
	return pq[i].time < pq[j].time
}
func (pq priorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(state))
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const inf int64 = 1<<62 - 1

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		adj := make([][]edge, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], edge{to: v})
			adj[v] = append(adj[v], edge{to: u})
		}

		// Build state mapping.
		base := make([]int, n)
		deg := make([]int, n)
		totalStates := 0
		for i := 0; i < n; i++ {
			deg[i] = len(adj[i])
			base[i] = totalStates
			totalStates += deg[i]
		}

		nodeOf := make([]int, totalStates)
		resOf := make([]int, totalStates)
		for u := 0; u < n; u++ {
			for r := 0; r < deg[u]; r++ {
				id := base[u] + r
				nodeOf[id] = u
				resOf[id] = r
			}
		}

		if deg[0] == 0 || deg[n-1] == 0 {
			fmt.Fprintln(out, -1, -1)
			continue
		}

		dist := make([]int64, totalStates)
		wait := make([]int64, totalStates)
		for i := 0; i < totalStates; i++ {
			dist[i] = inf
			wait[i] = inf
		}

		startID := base[0] + 0
		dist[startID] = 0
		wait[startID] = 0

		pq := priorityQueue{{time: 0, wait: 0, id: startID}}
		heap.Init(&pq)

		for pq.Len() > 0 {
			cur := heap.Pop(&pq).(state)
			if cur.time != dist[cur.id] || cur.wait != wait[cur.id] {
				continue
			}

			u := nodeOf[cur.id]
			r := resOf[cur.id]

			d := deg[u]
			if d == 0 {
				continue
			}

			// Wait transition.
			nextResid := (r + 1) % d
			waitID := base[u] + nextResid
			if cur.time+1 < dist[waitID] || (cur.time+1 == dist[waitID] && cur.wait+1 < wait[waitID]) {
				dist[waitID] = cur.time + 1
				wait[waitID] = cur.wait + 1
				heap.Push(&pq, state{time: cur.time + 1, wait: cur.wait + 1, id: waitID})
			}

			// Move transition if there is corresponding edge.
			if r < len(adj[u]) {
				to := adj[u][r].to
				if deg[to] == 0 {
					continue
				}
				newTime := cur.time + 1
				newWait := cur.wait
				nextResidTo := int(newTime % int64(deg[to]))
				nextID := base[to] + nextResidTo
				if newTime < dist[nextID] || (newTime == dist[nextID] && newWait < wait[nextID]) {
					dist[nextID] = newTime
					wait[nextID] = newWait
					heap.Push(&pq, state{time: newTime, wait: newWait, id: nextID})
				}
			}
		}

		bestTime := inf
		bestWait := inf
		for r := 0; r < deg[n-1]; r++ {
			id := base[n-1] + r
			if dist[id] < bestTime || (dist[id] == bestTime && wait[id] < bestWait) {
				bestTime = dist[id]
				bestWait = wait[id]
			}
		}

		fmt.Fprintln(out, bestTime, bestWait)
	}
}
