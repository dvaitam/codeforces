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

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// genValidGrid generates an n x m grid where X cells have no common points
// (no edge or corner adjacency).
func genValidGrid(n, m int, rng *rand.Rand) []string {
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	// Try to place some X cells randomly
	for attempts := 0; attempts < n*m*2; attempts++ {
		r := rng.Intn(n)
		c := rng.Intn(m)
		if grid[r][c] == 'X' {
			continue
		}
		ok := true
		for d := 0; d < 8; d++ {
			nr, nc := r+dx[d], c+dy[d]
			if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == 'X' {
				ok = false
				break
			}
		}
		if ok {
			grid[r][c] = 'X'
		}
	}
	rows := make([]string, n)
	for i := range grid {
		rows[i] = string(grid[i])
	}
	return rows
}

// validate checks candidate output for a given input grid.
// Returns "" if valid, or an error description.
func validate(n, m int, inputGrid, outputGrid []string) string {
	if len(outputGrid) != n {
		return fmt.Sprintf("expected %d rows, got %d", n, len(outputGrid))
	}
	for i := 0; i < n; i++ {
		if len(outputGrid[i]) != m {
			return fmt.Sprintf("row %d: expected length %d, got %d", i+1, m, len(outputGrid[i]))
		}
	}
	// Check no planting: if input has 'X', output must have 'X'
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if inputGrid[i][j] == 'X' && outputGrid[i][j] != 'X' {
				return fmt.Sprintf("cell (%d,%d) was empty but got sunflower in output", i+1, j+1)
			}
			if outputGrid[i][j] != 'X' && outputGrid[i][j] != '.' {
				return fmt.Sprintf("cell (%d,%d) has invalid char %c", i+1, j+1, outputGrid[i][j])
			}
		}
	}
	// Count empty cells and check connectivity + tree property
	emptyCells := 0
	id := make([][]int, n)
	for i := range id {
		id[i] = make([]int, m)
		for j := range id[i] {
			id[i][j] = -1
		}
	}
	idx := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if outputGrid[i][j] == 'X' {
				id[i][j] = idx
				idx++
				emptyCells++
			}
		}
	}
	if emptyCells == 0 {
		return "no empty cells"
	}
	// BFS to check connectivity and count edges
	dx := []int{-1, 0, 1, 0}
	dy := []int{0, 1, 0, -1}
	visited := make([]bool, emptyCells)
	// Find first empty cell
	startR, startC := -1, -1
	for i := 0; i < n && startR == -1; i++ {
		for j := 0; j < m; j++ {
			if outputGrid[i][j] == 'X' {
				startR, startC = i, j
				break
			}
		}
	}
	queue := [][2]int{{startR, startC}}
	visited[id[startR][startC]] = true
	visitCount := 1
	edgeCount := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for d := 0; d < 4; d++ {
			nr, nc := cur[0]+dx[d], cur[1]+dy[d]
			if nr >= 0 && nr < n && nc >= 0 && nc < m && outputGrid[nr][nc] == 'X' {
				edgeCount++
				if !visited[id[nr][nc]] {
					visited[id[nr][nc]] = true
					visitCount++
					queue = append(queue, [2]int{nr, nc})
				}
			}
		}
	}
	if visitCount != emptyCells {
		return fmt.Sprintf("empty cells not connected: reachable %d of %d", visitCount, emptyCells)
	}
	// Each edge counted twice in BFS traversal
	edgeCount /= 2
	if edgeCount != emptyCells-1 {
		return fmt.Sprintf("not a tree: %d nodes, %d edges (expected %d edges)", emptyCells, edgeCount, emptyCells-1)
	}
	return ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type tc struct {
		n, m int
		grid []string
	}

	var cases []tc

	// Fixed cases
	cases = append(cases, tc{3, 3, []string{"X.X", "...", "X.X"}})
	cases = append(cases, tc{1, 1, []string{"X"}})
	cases = append(cases, tc{1, 1, []string{"."}})
	cases = append(cases, tc{1, 5, []string{"X.X.X"}})
	cases = append(cases, tc{3, 3, []string{"...", "...", "..."}})
	cases = append(cases, tc{3, 3, []string{"...", ".X.", "..."}})
	cases = append(cases, tc{5, 5, []string{".....", ".X.X.", ".....", ".X.X.", "....."}})
	cases = append(cases, tc{1, 3, []string{"..."}})
	cases = append(cases, tc{3, 1, []string{".", ".", "."}})
	cases = append(cases, tc{2, 2, []string{"..", ".."}})

	// Random cases
	for i := 0; i < 90; i++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(6) + 1
		grid := genValidGrid(n, m, rng)
		cases = append(cases, tc{n, m, grid})
	}

	for i, c := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", c.n, c.m))
		for _, row := range c.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		input := sb.String()

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		outLines := strings.Split(got, "\n")
		// Filter empty lines
		var filtered []string
		for _, l := range outLines {
			if l != "" {
				filtered = append(filtered, l)
			}
		}

		errMsg := validate(c.n, c.m, c.grid, filtered)
		if errMsg != "" {
			fmt.Fprintf(os.Stderr, "case %d failed: %s\ninput:\n%soutput:\n%s\n", i+1, errMsg, input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
