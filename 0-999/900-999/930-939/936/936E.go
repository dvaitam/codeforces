package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int
}

type Query struct {
	t  int
	id int
}

type Item struct {
	dist int
	v    int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

const (
	inf       = int(1e9)
	blockSize = 1500
	keyMul    = int64(1_000_000)
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	xs := make([]int, n)
	ys := make([]int, n)
	coordToID := make(map[int64]int, n*2)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &ys[i])
		key := int64(xs[i])*keyMul + int64(ys[i])
		coordToID[key] = i
	}

	adj := make([][]int, n)
	dir := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < n; i++ {
		x, y := xs[i], ys[i]
		for _, d := range dir {
			nx, ny := x+d[0], y+d[1]
			key := int64(nx)*keyMul + int64(ny)
			if id, ok := coordToID[key]; ok {
				adj[i] = append(adj[i], id)
			}
		}
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([]Query, q)
	type2cnt := 0
	for i := 0; i < q; i++ {
		var t, x, y int
		fmt.Fscan(in, &t, &x, &y)
		id := coordToID[int64(x)*keyMul+int64(y)]
		queries[i] = Query{t: t, id: id}
		if t == 2 {
			type2cnt++
		}
	}

	distBase := make([]int, n)
	root := make([]int, n)
	tmpDist := make([]int, n)
	queue := make([]int, n)

	shopActive := make([]bool, n)
	globalSources := make([]int, 0, n)
	answers := make([]int, 0, type2cnt)

	for blockStart := 0; blockStart < q; blockStart += blockSize {
		blockEnd := blockStart + blockSize
		if blockEnd > q {
			blockEnd = q
		}

		// Collect all vertices mentioned in this block.
		idxMap := make(map[int]int, blockEnd-blockStart)
		terminals := make([]int, 0, blockEnd-blockStart)
		for i := blockStart; i < blockEnd; i++ {
			id := queries[i].id
			if _, ok := idxMap[id]; !ok {
				idxMap[id] = len(terminals)
				terminals = append(terminals, id)
			}
		}
		k := len(terminals)

		// Multi-source BFS from shops opened before this block.
		for i := range distBase {
			distBase[i] = inf
		}
		head, tail := 0, 0
		for _, s := range globalSources {
			if distBase[s] == 0 {
				continue
			}
			distBase[s] = 0
			queue[tail] = s
			tail++
		}
		for head < tail {
			u := queue[head]
			head++
			du := distBase[u]
			for _, v := range adj[u] {
				if distBase[v] == inf {
					distBase[v] = du + 1
					queue[tail] = v
					tail++
				}
			}
		}

		// Build compressed graph on terminals using multi-source BFS.
		for i := range root {
			root[i] = 0
			tmpDist[i] = -1
		}
		head, tail = 0, 0
		for idx, v := range terminals {
			root[v] = idx + 1 // +1 to keep zero as unvisited
			tmpDist[v] = 0
			queue[tail] = v
			tail++
		}

		edgeMap := make(map[uint64]int)
		for head < tail {
			u := queue[head]
			head++
			ru := root[u]
			du := tmpDist[u]
			for _, v := range adj[u] {
				if root[v] == 0 {
					root[v] = ru
					tmpDist[v] = du + 1
					queue[tail] = v
					tail++
				} else if root[v] != ru {
					a := ru - 1
					b := root[v] - 1
					if a > b {
						a, b = b, a
					}
					key := (uint64(a) << 32) | uint64(b)
					w := du + tmpDist[v] + 1
					if cur, ok := edgeMap[key]; !ok || w < cur {
						edgeMap[key] = w
					}
				}
			}
		}

		compressed := make([][]Edge, k)
		for key, w := range edgeMap {
			a := int(key >> 32)
			b := int(key & 0xffffffff)
			compressed[a] = append(compressed[a], Edge{to: b, w: w})
			compressed[b] = append(compressed[b], Edge{to: a, w: w})
		}

		distSmall := make([]int, k)
		for i := range distSmall {
			distSmall[i] = inf
		}
		pq := &PriorityQueue{}
		heap.Init(pq)

		for i := blockStart; i < blockEnd; i++ {
			qr := queries[i]
			idx := idxMap[qr.id]
			if qr.t == 1 {
				if distSmall[idx] > 0 {
					distSmall[idx] = 0
					heap.Push(pq, Item{dist: 0, v: idx})
				}
				for pq.Len() > 0 {
					it := heap.Pop(pq).(Item)
					if it.dist != distSmall[it.v] {
						continue
					}
					for _, e := range compressed[it.v] {
						nd := it.dist + e.w
						if nd < distSmall[e.to] {
							distSmall[e.to] = nd
							heap.Push(pq, Item{dist: nd, v: e.to})
						}
					}
				}
				if !shopActive[qr.id] {
					shopActive[qr.id] = true
					globalSources = append(globalSources, qr.id)
				}
			} else {
				best := distBase[qr.id]
				if distSmall[idx] < best {
					best = distSmall[idx]
				}
				if best >= inf {
					answers = append(answers, -1)
				} else {
					answers = append(answers, best)
				}
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for _, v := range answers {
		fmt.Fprintln(out, v)
	}
	out.Flush()
}
