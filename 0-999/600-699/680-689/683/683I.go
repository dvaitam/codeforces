package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Item struct {
	p, s   int
	bx, by int
	lx, ly int
	idx    int
}

type PQ []Item

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	if pq[i].p != pq[j].p {
		return pq[i].p < pq[j].p
	}
	return pq[i].s < pq[j].s
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	v := old[len(old)-1]
	*pq = old[:len(old)-1]
	return v
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	var sx, sy, bx, by, tx, ty int
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(in, &line)
		grid[i] = []byte(line)
		for j := 0; j < m; j++ {
			switch grid[i][j] {
			case 'Y':
				sx, sy = i, j
				grid[i][j] = '.'
			case 'B':
				bx, by = i, j
				grid[i][j] = '.'
			case 'T':
				tx, ty = i, j
				grid[i][j] = '.'
			}
		}
	}
	size := n * m
	idx := func(bx, by, lx, ly int) int {
		return (((bx*m)+by)*n+lx)*m + ly
	}
	const INF = int(1e9)
	distP := make([]int, size*size)
	distS := make([]int, size*size)
	parent := make([]int, size*size)
	op := make([]byte, size*size)
	for i := range distP {
		distP[i] = INF
		distS[i] = INF
		parent[i] = -1
	}
	start := idx(bx, by, sx, sy)
	distP[start], distS[start] = 0, 0
	pq := &PQ{}
	heap.Push(pq, Item{0, 0, bx, by, sx, sy, start})

	dirs := []struct {
		dx, dy     int
		move, push byte
	}{
		{-1, 0, 'n', 'N'},
		{1, 0, 's', 'S'},
		{0, -1, 'w', 'W'},
		{0, 1, 'e', 'E'},
	}

	targetIdx := -1
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		if cur.p != distP[cur.idx] || cur.s != distS[cur.idx] {
			continue
		}
		if cur.bx == tx && cur.by == ty {
			targetIdx = cur.idx
			break
		}
		for _, d := range dirs {
			nx, ny := cur.lx+d.dx, cur.ly+d.dy
			if nx >= 0 && nx < n && ny >= 0 && ny < m && grid[nx][ny] != 'X' && !(nx == cur.bx && ny == cur.by) {
				ni := idx(cur.bx, cur.by, nx, ny)
				np, ns := cur.p, cur.s+1
				if np < distP[ni] || (np == distP[ni] && ns < distS[ni]) {
					distP[ni] = np
					distS[ni] = ns
					parent[ni] = cur.idx
					op[ni] = d.move
					heap.Push(pq, Item{np, ns, cur.bx, cur.by, nx, ny, ni})
				}
			}
			if cur.lx == cur.bx-d.dx && cur.ly == cur.by-d.dy {
				nbx, nby := cur.bx+d.dx, cur.by+d.dy
				if nbx >= 0 && nbx < n && nby >= 0 && nby < m && grid[nbx][nby] != 'X' {
					ni := idx(nbx, nby, cur.bx, cur.by)
					np, ns := cur.p+1, cur.s+1
					if np < distP[ni] || (np == distP[ni] && ns < distS[ni]) {
						distP[ni] = np
						distS[ni] = ns
						parent[ni] = cur.idx
						op[ni] = d.push
						heap.Push(pq, Item{np, ns, nbx, nby, cur.bx, cur.by, ni})
					}
				}
			}
		}
	}

	if targetIdx == -1 {
		fmt.Println("NO")
		return
	}
	path := make([]byte, 0)
	for i := targetIdx; parent[i] != -1; i = parent[i] {
		path = append(path, op[i])
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	fmt.Println("YES")
	fmt.Println(string(path))
}
