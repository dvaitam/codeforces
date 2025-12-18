package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type cell struct{ r, c int }

func bfs(grid [][]byte, sr, sc int) [][]int {
	n := len(grid)
	m := len(grid[0])
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, m)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	q := make([]cell, 0, n*m)
	q = append(q, cell{sr, sc})
	dist[sr][sc] = 0
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for head := 0; head < len(q); head++ {
		cur := q[head]
		for _, d := range dirs {
			nr := cur.r + d[0]
			nc := cur.c + d[1]
			if nr < 0 || nr >= n || nc < 0 || nc >= m {
				continue
			}
			if grid[nr][nc] == 'X' || dist[nr][nc] != -1 {
				continue
			}
			dist[nr][nc] = dist[cur.r][cur.c] + 1
			q = append(q, cell{nr, nc})
		}
	}
	return dist
}

func expected(grid [][]byte) float64 {
	n := len(grid)
	m := len(grid[0])
	cells := make([]cell, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != 'X' {
				cells = append(cells, cell{i, j})
			}
		}
	}
	S := len(cells)
	if S == 0 {
		return 0.0
	}
	var sum int64
	for _, start := range cells {
		dist := bfs(grid, start.r, start.c)
		for _, end := range cells {
			d := dist[end.r][end.c]
			if d != -1 {
				sum += int64(d)
			}
		}
	}
	// Note: if graph is not connected, d might be -1.
	// But the problem implies connected via "shortest path through unoccupied cells"
	// "All empty cells have the same probability of being selected as the beginning or end of the path."
	// If they are unreachable, what happens?
	// Given the problem constraints (no two static in same row/col, no diagonal),
	// the empty cells are always connected.
	return float64(sum) / float64(S*S)
}

func genGrid(rng *rand.Rand) [][]byte {
	n := rng.Intn(4) + 2
	m := rng.Intn(4) + 2
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	rowsUsed := make(map[int]bool)
	colsUsed := make(map[int]bool)
	var coords []cell
	num := rng.Intn(min(n, m) + 1)
	// Try to place num particles
	for i := 0; i < 100 && len(coords) < num; i++ {
		r := rng.Intn(n)
		c := rng.Intn(m)
		if rowsUsed[r] || colsUsed[c] {
			continue
		}
		ok := true
		for _, p := range coords {
			if abs(p.r-r) == 1 && abs(p.c-c) == 1 {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		coords = append(coords, cell{r, c})
		rowsUsed[r] = true
		colsUsed[c] = true
		grid[r][c] = 'X'
	}
	return grid
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

func genCase(rng *rand.Rand) (string, string) {
	grid := genGrid(rng)
	n := len(grid)
	m := len(grid[0])
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	expectedVal := expected(grid)
	return sb.String(), fmt.Sprintf("%.15f", expectedVal) // Use high precision string for comparison reference, but comparison logic handles float
}

func runCase(bin, input, exp string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	
	expVal, err := strconv.ParseFloat(exp, 64)
	if err != nil {
		return fmt.Errorf("bad expected value: %v", err)
	}
	gotVal, err := strconv.ParseFloat(gotStr, 64)
	if err != nil {
		return fmt.Errorf("bad output value: %v", err)
	}

	// Check absolute or relative error 10^-6
	absErr := math.Abs(expVal - gotVal)
	relErr := 0.0
	if math.Abs(expVal) > 1e-9 {
		relErr = absErr / math.Abs(expVal)
	}

	if absErr > 1e-6 && relErr > 1e-6 {
		return fmt.Errorf("expected %v got %v (abs err %v, rel err %v)", expVal, gotVal, absErr, relErr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}