package main

import (
	"bufio"
	"container/heap"
	"os"
)

type Edge struct {
	to  int
	idx int
}

type Item struct {
	node int
	time int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].time < pq[j].time }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
	x := 0
	sign := 1
	b, err := reader.ReadByte()
	for err == nil && (b < '0' || b > '9') && b != '-' {
		b, err = reader.ReadByte()
	}
	if err != nil {
		return 0
	}
	if b == '-' {
		sign = -1
		b, _ = reader.ReadByte()
	}
	for ; err == nil && b >= '0' && b <= '9'; b, err = reader.ReadByte() {
		x = x*10 + int(b-'0')
	}
	return x * sign
}

func main() {
	defer writer.Flush()
	n := readInt()
	m := readInt()
	q := readInt()

	adj := make([][]Edge, n+1)
	for i := 1; i <= m; i++ {
		v := readInt()
		u := readInt()
		adj[v] = append(adj[v], Edge{to: u, idx: i})
		adj[u] = append(adj[u], Edge{to: v, idx: i})
	}

	const INF = int(1<<31 - 1)
	dist := make([]int, n+1)

	for ; q > 0; q-- {
		l := readInt()
		r := readInt()
		s := readInt()
		t := readInt()
		for i := 1; i <= n; i++ {
			dist[i] = INF
		}
		pq := &PriorityQueue{}
		heap.Init(pq)
		dist[s] = l
		heap.Push(pq, Item{node: s, time: l})
		for pq.Len() > 0 {
			it := heap.Pop(pq).(Item)
			if it.time != dist[it.node] {
				continue
			}
			if it.time > r {
				continue
			}
			if it.node == t {
				break
			}
			for _, e := range adj[it.node] {
				if e.idx < l || e.idx > r {
					continue
				}
				if e.idx < it.time {
					continue
				}
				nt := e.idx
				if nt < dist[e.to] {
					dist[e.to] = nt
					heap.Push(pq, Item{node: e.to, time: nt})
				}
			}
		}
		if dist[t] <= r {
			writer.WriteString("Yes\n")
		} else {
			writer.WriteString("No\n")
		}
	}
}
