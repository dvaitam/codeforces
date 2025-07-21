package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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
	var sum int
	for _, start := range cells {
		dist := bfs(grid, start.r, start.c)
		for _, end := range cells {
			d := dist[end.r][end.c]
			sum += d
		}
	}
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
	for len(coords) < num {
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
	out := fmt.Sprintf("%.9f\n", expectedVal)
	return sb.String(), out
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
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s", exp, got)
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
