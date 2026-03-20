package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ===== Embedded reference solver for 1025E =====

type refPoint struct{ x, y int }

type refMove struct {
	x1, y1, x2, y2 int
}

type refItem struct {
	p refPoint
	d int
	i int
}

type refPQ []*refItem

func (pq refPQ) Len() int           { return len(pq) }
func (pq refPQ) Less(i, j int) bool { return pq[i].d < pq[j].d }
func (pq refPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].i = i
	pq[j].i = j
}
func (pq *refPQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*refItem)
	item.i = n
	*pq = append(*pq, item)
}
func (pq *refPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.i = -1
	*pq = old[0 : n-1]
	return item
}

func refGetPathP(start, target refPoint, fixed [][]bool, grid [][]int, n int) []refPoint {
	dist := make([][]int, n+1)
	parent := make([][]refPoint, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = make([]int, n+1)
		parent[i] = make([]refPoint, n+1)
		for j := 1; j <= n; j++ {
			dist[i][j] = 1e9
		}
	}

	pq := make(refPQ, 0)
	heap.Init(&pq)

	dist[start.x][start.y] = 0
	heap.Push(&pq, &refItem{p: start, d: 0})

	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(*refItem)
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
					heap.Push(&pq, &refItem{p: refPoint{nx, ny}, d: dist[nx][ny]})
				}
			}
		}
	}

	var path []refPoint
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

func refGetPathQ(start refPoint, fixed [][]bool, grid [][]int, pos_c refPoint, n int) []refPoint {
	visited := make([][]bool, n+1)
	for i := 1; i <= n; i++ {
		visited[i] = make([]bool, n+1)
	}
	parent := make([][]refPoint, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = make([]refPoint, n+1)
	}

	queue := []refPoint{start}
	visited[start.x][start.y] = true

	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	var emptyCell refPoint

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
					queue = append(queue, refPoint{nx, ny})
				}
			}
		}
	}

	var path []refPoint
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

func refSolveInternal(initialPos []refPoint, n int, m int) []refMove {
	grid := make([][]int, n+1)
	fixed := make([][]bool, n+1)
	for i := 1; i <= n; i++ {
		grid[i] = make([]int, n+1)
		fixed[i] = make([]bool, n+1)
	}

	pos := make([]refPoint, m+1)
	for i := 1; i <= m; i++ {
		p := initialPos[i]
		grid[p.x][p.y] = i
		pos[i] = p
	}

	moves := []refMove{}

	for c := 1; c <= m; c++ {
		target := refPoint{1, c}
		if pos[c] == target {
			fixed[target.x][target.y] = true
			continue
		}

		P := refGetPathP(pos[c], target, fixed, grid, n)
		for _, p_next := range P {
			if grid[p_next.x][p_next.y] == 0 {
				grid[pos[c].x][pos[c].y] = 0
				grid[p_next.x][p_next.y] = c
				moves = append(moves, refMove{pos[c].x, pos[c].y, p_next.x, p_next.y})
				pos[c] = p_next
			} else {
				Q := refGetPathQ(p_next, fixed, grid, pos[c], n)
				for j := len(Q) - 2; j >= 0; j-- {
					q_curr := Q[j]
					q_next := Q[j+1]
					cube_id := grid[q_curr.x][q_curr.y]
					grid[q_curr.x][q_curr.y] = 0
					grid[q_next.x][q_next.y] = cube_id
					pos[cube_id] = q_next
					moves = append(moves, refMove{q_curr.x, q_curr.y, q_next.x, q_next.y})
				}
				grid[pos[c].x][pos[c].y] = 0
				grid[p_next.x][p_next.y] = c
				moves = append(moves, refMove{pos[c].x, pos[c].y, p_next.x, p_next.y})
				pos[c] = p_next
			}
		}
		fixed[target.x][target.y] = true
	}
	return moves
}

func refSolve(input string) string {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	nextInt := func() int {
		sc.Scan()
		v, _ := strconv.Atoi(sc.Text())
		return v
	}

	n := nextInt()
	m := nextInt()

	initial := make([]refPoint, m+1)
	for i := 1; i <= m; i++ {
		initial[i] = refPoint{nextInt(), nextInt()}
	}

	final := make([]refPoint, m+1)
	for i := 1; i <= m; i++ {
		final[i] = refPoint{nextInt(), nextInt()}
	}

	moves1 := refSolveInternal(initial, n, m)
	moves2 := refSolveInternal(final, n, m)

	for i := 0; i < len(moves2)/2; i++ {
		j := len(moves2) - 1 - i
		moves2[i], moves2[j] = moves2[j], moves2[i]
	}
	for i := range moves2 {
		moves2[i] = refMove{moves2[i].x2, moves2[i].y2, moves2[i].x1, moves2[i].y1}
	}

	totalMoves := append(moves1, moves2...)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(totalMoves))
	for _, mv := range totalMoves {
		fmt.Fprintf(&sb, "%d %d %d %d\n", mv.x1, mv.y1, mv.x2, mv.y2)
	}
	return strings.TrimSpace(sb.String())
}

// ===== Verifier infrastructure =====

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	var errBuf strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func checkSolution(input, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	nextInt := func() int {
		scanner.Scan()
		v, _ := strconv.Atoi(scanner.Text())
		return v
	}
	n := nextInt()
	m := nextInt()
	type Point struct{ x, y int }
	start := make([]Point, m)
	target := make([]Point, m)
	for i := 0; i < m; i++ {
		start[i] = Point{nextInt(), nextInt()}
	}
	for i := 0; i < m; i++ {
		target[i] = Point{nextInt(), nextInt()}
	}

	outScanner := bufio.NewScanner(strings.NewReader(output))
	outScanner.Split(bufio.ScanWords)
	if !outScanner.Scan() {
		return fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(outScanner.Text())
	if err != nil {
		return fmt.Errorf("invalid number of moves: %v", err)
	}
	if k > 10800 {
		return fmt.Errorf("too many moves: %d > 10800", k)
	}

	grid := make([][]int, n+1)
	for i := range grid {
		grid[i] = make([]int, n+1)
	}
	for i, p := range start {
		grid[p.x][p.y] = i + 1
	}

	for i := 0; i < k; i++ {
		if !outScanner.Scan() {
			return fmt.Errorf("expected more moves")
		}
		x1, _ := strconv.Atoi(outScanner.Text())
		if !outScanner.Scan() {
			return fmt.Errorf("incomplete move")
		}
		y1, _ := strconv.Atoi(outScanner.Text())
		if !outScanner.Scan() {
			return fmt.Errorf("incomplete move")
		}
		x2, _ := strconv.Atoi(outScanner.Text())
		if !outScanner.Scan() {
			return fmt.Errorf("incomplete move")
		}
		y2, _ := strconv.Atoi(outScanner.Text())

		if x1 < 1 || x1 > n || y1 < 1 || y1 > n {
			return fmt.Errorf("move %d: start (%d,%d) out of bounds", i+1, x1, y1)
		}
		if x2 < 1 || x2 > n || y2 < 1 || y2 > n {
			return fmt.Errorf("move %d: end (%d,%d) out of bounds", i+1, x2, y2)
		}
		if abs(x1-x2)+abs(y1-y2) != 1 {
			return fmt.Errorf("move %d: not adjacent (%d,%d)->(%d,%d)", i+1, x1, y1, x2, y2)
		}
		if grid[x1][y1] == 0 {
			return fmt.Errorf("move %d: no cube at (%d,%d)", i+1, x1, y1)
		}
		if grid[x2][y2] != 0 {
			return fmt.Errorf("move %d: destination (%d,%d) occupied", i+1, x2, y2)
		}
		grid[x2][y2] = grid[x1][y1]
		grid[x1][y1] = 0
	}

	for i, p := range target {
		if grid[p.x][p.y] == 0 {
			return fmt.Errorf("target (%d,%d) empty", p.x, p.y)
		}
		if grid[p.x][p.y] != i+1 {
			return fmt.Errorf("target (%d,%d) has cube %d, expected %d", p.x, p.y, grid[p.x][p.y], i+1)
		}
	}
	return nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 4
	maxM := n - 1 // ensure enough free cells
	if maxM < 1 {
		maxM = 1
	}
	m := r.Intn(maxM) + 1

	type cell struct{ x, y int }
	cells := make([]cell, 0, n*n)
	for x := 1; x <= n; x++ {
		for y := 1; y <= n; y++ {
			cells = append(cells, cell{x, y})
		}
	}
	r.Shuffle(len(cells), func(i, j int) { cells[i], cells[j] = cells[j], cells[i] })

	start := make([]cell, m)
	target := make([]cell, m)
	for i := 0; i < m; i++ {
		start[i] = cells[i]
	}
	for i := 0; i < m; i++ {
		target[i] = cells[m+i]
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, c := range start {
		fmt.Fprintf(&sb, "%d %d\n", c.x, c.y)
	}
	for _, c := range target {
		fmt.Fprintf(&sb, "%d %d\n", c.x, c.y)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]

	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkSolution(input, got); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%sgot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
