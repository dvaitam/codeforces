package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

const timeLimit = 1e20

type item struct {
	d   int64
	id  int
	idx int
}

type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].d < pq[j].d
}
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idx, pq[j].idx = i, j
}
func (pq *priorityQueue) Push(x interface{}) {
	it := x.(*item)
	it.idx = len(*pq)
	*pq = append(*pq, it)
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

type grid struct {
	n int
	h [][]int
}

func (g *grid) idx(tx, ty, x, y int) int {
	return (((tx+1)*3 + (ty + 1)) * g.n * g.n) + x*g.n + y
}

// compute minimal cost to translate by (goalTx, goalTy) tiles (goalTx,goalTy in {-1,0,1})
// from any starting cell (x,y) to (x,y) in shifted tile.
func (g *grid) minShiftCost(goalTx, goalTy int) int64 {
	n := g.n
	totalTiles := 3 * 3
	nodeCount := totalTiles * n * n
	inf := int64(1 << 60)
	best := inf

	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for sx := 0; sx < n; sx++ {
		for sy := 0; sy < n; sy++ {
			start := g.idx(0, 0, sx, sy)
			target := g.idx(goalTx, goalTy, sx, sy)

			dist := make([]int64, nodeCount)
			for i := range dist {
				dist[i] = inf
			}
			dist[start] = 0

			pq := priorityQueue{}
			heap.Push(&pq, &item{d: 0, id: start})

			for pq.Len() > 0 {
				it := heap.Pop(&pq).(*item)
				if it.d != dist[it.id] {
					continue
				}
				if it.id == target {
					if it.d < best {
						best = it.d
					}
					break
				}
				// decode
				tmp := it.id
				y := tmp % n
				tmp /= n
				x := tmp % n
				tmp /= n
				ty := tmp%3 - 1
				tx := tmp/3 - 1

				for _, d := range dirs {
					nx, ny := x+d[0], y+d[1]
					ntx, nty := tx, ty
					if nx < 0 {
						nx += n
						ntx--
					}
					if nx >= n {
						nx -= n
						ntx++
					}
					if ny < 0 {
						ny += n
						nty--
					}
					if ny >= n {
						ny -= n
						nty++
					}
					if ntx < -1 || ntx > 1 || nty < -1 || nty > 1 {
						continue
					}
					cost := 1 + absInt(g.h[x][y]-g.h[nx][ny])
					nid := g.idx(ntx, nty, nx, ny)
					nd := it.d + int64(cost)
					if nd < dist[nid] {
						dist[nid] = nd
						heap.Push(&pq, &item{d: nd, id: nid})
					}
				}
			}
		}
	}
	return best
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	h := make([][]int, n)
	for i := 0; i < n; i++ {
		h[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &h[i][j])
		}
	}

	g := grid{n: n, h: h}
	costX := g.minShiftCost(1, 0)
	costY := g.minShiftCost(0, 1)

	// Safety: avoid division by zero
	if costX == 0 || costY == 0 {
		fmt.Println("0")
		return
	}

	ans := 2 * timeLimit * timeLimit / (float64(costX) * float64(costY))
	// Use scientific notation with enough precision
	fmt.Printf("%.10e\n", ans)
	_ = math.NaN() // import math used
}
