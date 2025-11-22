package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Item struct {
	x, y  int
	cost  int
	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(*pq)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func minimax(a [][]int, sx, sy int) [][]int {
	n := len(a)
	const inf = int(1e18)
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			dist[i][j] = inf
		}
	}
	dist[sx][sy] = a[sx][sy]
	pq := PriorityQueue{}
	heap.Push(&pq, &Item{sx, sy, dist[sx][sy], 0})
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for pq.Len() > 0 {
		cur := heap.Pop(&pq).(*Item)
		if cur.cost != dist[cur.x][cur.y] {
			continue
		}
		for _, d := range dirs {
			nx, ny := cur.x+d[0], cur.y+d[1]
			if nx < 0 || nx >= n || ny < 0 || ny >= n {
				continue
			}
			ncost := cur.cost
			if a[nx][ny] > ncost {
				ncost = a[nx][ny]
			}
			if ncost < dist[nx][ny] {
				dist[nx][ny] = ncost
				heap.Push(&pq, &Item{nx, ny, ncost, 0})
			}
		}
	}
	return dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, r int
		var k int64
		fmt.Fscan(in, &n, &r, &k)
		a := make([][]int, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}
		var tmp string
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &tmp) // colors not used in this simplified strategy
		}

		distM := minimax(a, 0, 0)
		distF := minimax(a, n-1, 0)
		fmt.Fprintln(out, distM[r-1][n-1], distF[r-1][n-1])
	}
}
