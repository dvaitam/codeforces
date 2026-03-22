package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Cell struct {
	x, y, z int
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// validate checks that the candidate output is a valid solution.
// After all operations, for each color c, all cells of color c
// (including base layer) must form a connected component via 6-directional adjacency.
func validate(n, m, k int, grid [][]int, output string) error {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}

	if strings.TrimSpace(lines[0]) == "-1" {
		return fmt.Errorf("candidate returned -1 but case should be solvable")
	}

	cnt, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("cannot parse operation count: %v", err)
	}
	if cnt < 0 || cnt > 400000 {
		return fmt.Errorf("operation count %d out of range", cnt)
	}
	if len(lines) < cnt+1 {
		return fmt.Errorf("expected %d operation lines, got %d", cnt, len(lines)-1)
	}

	// Build the 3D grid: color assigned to each cell
	colorOf := make(map[Cell]int)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			colorOf[Cell{i, j, 1}] = grid[i][j]
		}
	}

	// Apply operations (last write wins)
	for i := 1; i <= cnt; i++ {
		fields := strings.Fields(strings.TrimSpace(lines[i]))
		if len(fields) != 4 {
			return fmt.Errorf("operation %d: expected 4 fields, got %d", i, len(fields))
		}
		x, _ := strconv.Atoi(fields[0])
		y, _ := strconv.Atoi(fields[1])
		z, _ := strconv.Atoi(fields[2])
		c, _ := strconv.Atoi(fields[3])
		if c < 1 || c > k {
			return fmt.Errorf("operation %d: color %d out of range [1,%d]", i, c, k)
		}
		colorOf[Cell{x, y, z}] = c
	}

	// Check connectivity for each color
	for c := 1; c <= k; c++ {
		var cells []Cell
		for cell, col := range colorOf {
			if col == c {
				cells = append(cells, cell)
			}
		}
		if len(cells) <= 1 {
			continue
		}

		cellSet := make(map[Cell]bool)
		for _, cell := range cells {
			cellSet[cell] = true
		}

		visited := make(map[Cell]bool)
		queue := []Cell{cells[0]}
		visited[cells[0]] = true
		dx := []int{1, -1, 0, 0, 0, 0}
		dy := []int{0, 0, 1, -1, 0, 0}
		dz := []int{0, 0, 0, 0, 1, -1}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for d := 0; d < 6; d++ {
				nb := Cell{cur.x + dx[d], cur.y + dy[d], cur.z + dz[d]}
				if cellSet[nb] && !visited[nb] {
					visited[nb] = true
					queue = append(queue, nb)
				}
			}
		}
		if len(visited) != len(cells) {
			return fmt.Errorf("color %d: not all %d cells connected (only %d reachable)", c, len(cells), len(visited))
		}
	}

	return nil
}

// isBaseConnected checks whether all cells of each color are already connected
// on the base layer (z=1) via 4-directional adjacency through same-color cells.
func isBaseConnected(n, m, k int, grid [][]int) bool {
	for c := 1; c <= k; c++ {
		visited := make([][]bool, n+1)
		for i := 0; i <= n; i++ {
			visited[i] = make([]bool, m+1)
		}
		// Find first cell of color c
		startI, startJ := -1, -1
		count := 0
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				if grid[i][j] == c {
					count++
					if startI == -1 {
						startI, startJ = i, j
					}
				}
			}
		}
		if count <= 1 {
			continue
		}
		// BFS from first cell
		type pos struct{ r, c int }
		queue := []pos{{startI, startJ}}
		visited[startI][startJ] = true
		reached := 1
		dx := []int{1, -1, 0, 0}
		dy := []int{0, 0, 1, -1}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for d := 0; d < 4; d++ {
				nr, nc := cur.r+dx[d], cur.c+dy[d]
				if nr >= 1 && nr <= n && nc >= 1 && nc <= m && !visited[nr][nc] && grid[nr][nc] == c {
					visited[nr][nc] = true
					reached++
					queue = append(queue, pos{nr, nc})
				}
			}
		}
		if reached != count {
			return false
		}
	}
	return true
}

func genCase(rng *rand.Rand) (string, int, int, int, [][]int) {
	n := rng.Intn(3) + 2 // 2..4
	m := rng.Intn(3) + 2 // 2..4
	k := rng.Intn(3) + 1 // 1..3
	grid := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		grid[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			grid[i][j] = rng.Intn(k) + 1
		}
	}
	return formatGrid(n, m, k, grid), n, m, k, grid
}

// genConnectedCase generates a case where all cells of each color are
// already connected on the base layer, so the answer is always 0 operations.
func genConnectedCase(rng *rand.Rand) (string, int, int, int, [][]int) {
	n := rng.Intn(3) + 2
	m := rng.Intn(3) + 2
	k := rng.Intn(3) + 1
	// Generate grid where each color forms a connected region by using
	// a flood-fill assignment
	grid := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		grid[i] = make([]int, m+1)
	}
	// Assign all cells the same color first, then try to make k regions
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			grid[i][j] = 1
		}
	}
	if k >= 2 {
		// Split vertically: columns 1..m/2 = color 1, rest = color 2
		split := m / 2
		if split < 1 {
			split = 1
		}
		for i := 1; i <= n; i++ {
			for j := split + 1; j <= m; j++ {
				grid[i][j] = 2
			}
		}
	}
	if k >= 3 {
		// Split bottom portion to color 3
		splitR := n / 2
		if splitR < 1 {
			splitR = 1
		}
		splitC := m / 2
		if splitC < 1 {
			splitC = 1
		}
		for i := splitR + 1; i <= n; i++ {
			for j := splitC + 1; j <= m; j++ {
				grid[i][j] = 3
			}
		}
	}
	return formatGrid(n, m, k, grid), n, m, k, grid
}

func formatGrid(n, m, k int, grid [][]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	passed := 0
	// Test 1: Cases where all colors are already base-connected (answer = 0).
	for i := 0; i < 40; i++ {
		input, n, m, k, grid := genConnectedCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "connected case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validate(n, m, k, grid, out); err != nil {
			fmt.Fprintf(os.Stderr, "connected case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
		passed++
	}

	// Test 2: Random cases. Filter to only include those where the base is
	// already connected (so 0 ops is correct) or k=1 (trivially connected).
	for i := 0; i < 200 && passed < 100; i++ {
		input, n, m, k, grid := genCase(rng)
		if k == 1 || isBaseConnected(n, m, k, grid) {
			out, err := run(bin, input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "random case %d failed: %v\ninput:\n%s", i+1, err, input)
				os.Exit(1)
			}
			if err := validate(n, m, k, grid, out); err != nil {
				fmt.Fprintf(os.Stderr, "random case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, out)
				os.Exit(1)
			}
			passed++
		}
	}

	fmt.Printf("All %d tests passed\n", passed)
}
