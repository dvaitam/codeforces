package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Point struct{ x, y int }

type Move struct {
	x1, y1, x2, y2 int
}

type Item struct {
	p Point
	d int
	i int
}

type PQ []*Item

func (pq PQ) Len() int           { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].d < pq[j].d }
func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].i = i
	pq[j].i = j
}
func (pq *PQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.i = n
	*pq = append(*pq, item)
}
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.i = -1
	*pq = old[0 : n-1]
	return item
}

func getPathP(start, target Point, fixed [][]bool, grid [][]int, n int) []Point {
	dist := make([][]int, n+1)
	parent := make([][]Point, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = make([]int, n+1)
		parent[i] = make([]Point, n+1)
		for j := 1; j <= n; j++ {
			dist[i][j] = 1e9
		}
	}

	pq := make(PQ, 0)
	heap.Init(&pq)

	dist[start.x][start.y] = 0
	heap.Push(&pq, &Item{p: start, d: 0})

	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(*Item)
		u := curr.p
		if curr.d > dist[u.x][u.y] {
			continue
		}
		if u == target {
			break
		}

		for k := 0; k < 4; k++ {
			nx, ny := u.x+dx[k], u.y+dy[k]
			if nx >= 1 && nx <= n && ny >= 1 && ny <= n {
				if fixed[nx][ny] {
					continue
				}
				cost := 1
				if grid[nx][ny] != 0 {
					cost = 100
				}
				if dist[u.x][u.y]+cost < dist[nx][ny] {
					dist[nx][ny] = dist[u.x][u.y] + cost
					parent[nx][ny] = u
					heap.Push(&pq, &Item{p: Point{nx, ny}, d: dist[nx][ny]})
				}
			}
		}
	}

	var path []Point
	curr := target
	for curr != start {
		path = append(path, curr)
		curr = parent[curr.x][curr.y]
	}

	for i := 0; i < len(path)/2; i++ {
		j := len(path) - 1 - i
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func getPathQ(start Point, fixed [][]bool, grid [][]int, pos_c Point, n int) []Point {
	visited := make([][]bool, n+1)
	for i := 1; i <= n; i++ {
		visited[i] = make([]bool, n+1)
	}
	parent := make([][]Point, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = make([]Point, n+1)
	}

	queue := []Point{start}
	visited[start.x][start.y] = true

	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	var emptyCell Point

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		if grid[u.x][u.y] == 0 {
			emptyCell = u
			break
		}

		for k := 0; k < 4; k++ {
			nx, ny := u.x+dx[k], u.y+dy[k]
			if nx >= 1 && nx <= n && ny >= 1 && ny <= n {
				if fixed[nx][ny] {
					continue
				}
				if nx == pos_c.x && ny == pos_c.y {
					continue
				}
				if !visited[nx][ny] {
					visited[nx][ny] = true
					parent[nx][ny] = u
					queue = append(queue, Point{nx, ny})
				}
			}
		}
	}

	var path []Point
	curr := emptyCell
	for {
		path = append(path, curr)
		if curr == start {
			break
		}
		curr = parent[curr.x][curr.y]
	}
	for i := 0; i < len(path)/2; i++ {
		j := len(path) - 1 - i
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func solve(initialPos []Point, n int, m int) []Move {
	grid := make([][]int, n+1)
	fixed := make([][]bool, n+1)
	for i := 1; i <= n; i++ {
		grid[i] = make([]int, n+1)
		fixed[i] = make([]bool, n+1)
	}

	pos := make([]Point, m+1)
	for i := 1; i <= m; i++ {
		p := initialPos[i]
		grid[p.x][p.y] = i
		pos[i] = p
	}

	moves := []Move{}

	for c := 1; c <= m; c++ {
		target := Point{1, c}
		if pos[c] == target {
			fixed[target.x][target.y] = true
			continue
		}

		P := getPathP(pos[c], target, fixed, grid, n)
		for _, p_next := range P {
			if grid[p_next.x][p_next.y] == 0 {
				grid[pos[c].x][pos[c].y] = 0
				grid[p_next.x][p_next.y] = c
				moves = append(moves, Move{pos[c].x, pos[c].y, p_next.x, p_next.y})
				pos[c] = p_next
			} else {
				Q := getPathQ(p_next, fixed, grid, pos[c], n)
				for j := len(Q) - 2; j >= 0; j-- {
					q_curr := Q[j]
					q_next := Q[j+1]
					cube_id := grid[q_curr.x][q_curr.y]
					grid[q_curr.x][q_curr.y] = 0
					grid[q_next.x][q_next.y] = cube_id
					pos[cube_id] = q_next
					moves = append(moves, Move{q_curr.x, q_curr.y, q_next.x, q_next.y})
				}
				grid[pos[c].x][pos[c].y] = 0
				grid[p_next.x][p_next.y] = c
				moves = append(moves, Move{pos[c].x, pos[c].y, p_next.x, p_next.y})
				pos[c] = p_next
			}
		}
		fixed[target.x][target.y] = true
	}
	return moves
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	initial := make([]Point, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(reader, &initial[i].x, &initial[i].y)
	}

	final := make([]Point, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(reader, &final[i].x, &final[i].y)
	}

	moves1 := solve(initial, n, m)
	moves2 := solve(final, n, m)

	for i := 0; i < len(moves2)/2; i++ {
		j := len(moves2) - 1 - i
		moves2[i], moves2[j] = moves2[j], moves2[i]
	}
	for i := range moves2 {
		moves2[i] = Move{moves2[i].x2, moves2[i].y2, moves2[i].x1, moves2[i].y1}
	}

	totalMoves := append(moves1, moves2...)

	fmt.Println(len(totalMoves))
	writer := bufio.NewWriter(os.Stdout)
	for _, mv := range totalMoves {
		fmt.Fprintf(writer, "%d %d %d %d\n", mv.x1, mv.y1, mv.x2, mv.y2)
	}
	writer.Flush()
}