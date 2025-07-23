package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Point represents a coordinate in the grid.
type Point struct {
	x, y int
}

// Item is an element in the priority queue for Dijkstra.
type Item struct {
	d   int
	idx int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].d < pq[j].d }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, p int
	if _, err := fmt.Fscan(reader, &n, &m, &p); err != nil {
		return
	}

	grid := make([][]int, n)
	pos := make([][]Point, p+1)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &grid[i][j])
			val := grid[i][j]
			if val >= 1 && val <= p {
				pos[val] = append(pos[val], Point{i, j})
			}
		}
	}

	const inf = int(1e9)
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, m)
	}

	// initialize distances for type 1 cells
	for i := range dist {
		for j := range dist[i] {
			dist[i][j] = inf
		}
	}
	for _, pt := range pos[1] {
		dist[pt.x][pt.y] = abs(pt.x) + abs(pt.y)
	}

	// process types 2..p
	for t := 2; t <= p; t++ {
		prev := pos[t-1]
		cur := pos[t]
		if len(prev)*len(cur) <= n*m {
			// direct pair computation
			for _, pt := range cur {
				best := inf
				for _, pr := range prev {
					cand := dist[pr.x][pr.y] + abs(pt.x-pr.x) + abs(pt.y-pr.y)
					if cand < best {
						best = cand
					}
				}
				dist[pt.x][pt.y] = best
			}
		} else {
			// Dijkstra from all previous positions
			for i := range dist {
				for j := range dist[i] {
					dist[i][j] = inf
				}
			}
			pq := &PriorityQueue{}
			heap.Init(pq)
			for _, pr := range prev {
				d := dist[pr.x][pr.y]
				dist[pr.x][pr.y] = min(dist[pr.x][pr.y], d)
				heap.Push(pq, Item{d: d, idx: pr.x*m + pr.y})
			}
			dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
			for pq.Len() > 0 {
				item := heap.Pop(pq).(Item)
				d := item.d
				x := item.idx / m
				y := item.idx % m
				if d != dist[x][y] {
					continue
				}
				for _, dir := range dirs {
					nx := x + dir[0]
					ny := y + dir[1]
					if nx < 0 || nx >= n || ny < 0 || ny >= m {
						continue
					}
					nd := d + 1
					if nd < dist[nx][ny] {
						dist[nx][ny] = nd
						heap.Push(pq, Item{d: nd, idx: nx*m + ny})
					}
				}
			}
		}
	}

	finalPos := pos[p][0]
	fmt.Fprintln(writer, dist[finalPos.x][finalPos.y])
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
