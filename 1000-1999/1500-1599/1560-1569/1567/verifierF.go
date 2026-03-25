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

func genGrid(rng *rand.Rand) []string {
	n := rng.Intn(5) + 2
	m := rng.Intn(5) + 2
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			grid[i][j] = '.'
		}
	}
	for i := 1; i < n-1; i++ {
		for j := 1; j < m-1; j++ {
			if rng.Intn(3) == 0 {
				grid[i][j] = 'X'
			}
		}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = string(grid[i])
	}
	return res
}

func buildInput(grid []string) string {
	n := len(grid)
	m := len(grid[0])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i] + "\n")
	}
	return sb.String()
}

// hasSolution checks if a valid solution exists using bipartite check.
func hasSolution(grid []string) bool {
	n := len(grid)
	m := len(grid[0])
	dirs := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != 'X' {
				continue
			}
			cnt := 0
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
					cnt++
				}
			}
			if cnt%2 == 1 {
				return false
			}
		}
	}

	// Build adjacency for bipartite check
	type Edge struct{ u, v int }
	id := func(r, c int) int { return r*m + c }
	adj := make([][]int, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != 'X' {
				continue
			}
			var ne []int
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
					ne = append(ne, id(ni, nj))
				}
			}
			if len(ne) == 2 {
				adj[ne[0]] = append(adj[ne[0]], ne[1])
				adj[ne[1]] = append(adj[ne[1]], ne[0])
			} else if len(ne) == 4 {
				// pair vertical and horizontal
				adj[ne[0]] = append(adj[ne[0]], ne[2])
				adj[ne[2]] = append(adj[ne[2]], ne[0])
				adj[ne[1]] = append(adj[ne[1]], ne[3])
				adj[ne[3]] = append(adj[ne[3]], ne[1])
			}
		}
	}

	color := make([]int, n*m)
	for i := range color {
		color[i] = -1
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != '.' || color[id(i, j)] != -1 {
				continue
			}
			queue := []int{id(i, j)}
			color[id(i, j)] = 0
			for len(queue) > 0 {
				u := queue[0]
				queue = queue[1:]
				for _, v := range adj[u] {
					if color[v] == -1 {
						color[v] = 1 - color[u]
						queue = append(queue, v)
					} else if color[v] == color[u] {
						return false
					}
				}
			}
		}
	}
	return true
}

// validateOutput checks if the candidate's output is valid for the given grid.
func validateOutput(grid []string, output string) error {
	n := len(grid)
	m := len(grid[0])
	dirs := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}

	verdict := strings.TrimSpace(lines[0])
	canSolve := hasSolution(grid)

	if verdict == "NO" {
		if canSolve {
			return fmt.Errorf("candidate said NO but solution exists")
		}
		return nil
	}
	if verdict != "YES" {
		return fmt.Errorf("expected YES or NO, got %q", verdict)
	}
	if !canSolve {
		return fmt.Errorf("candidate said YES but no solution exists")
	}

	if len(lines) < n+1 {
		return fmt.Errorf("expected %d grid rows, got %d lines", n, len(lines)-1)
	}

	vals := make([][]int, n)
	for i := 0; i < n; i++ {
		parts := strings.Fields(strings.TrimSpace(lines[i+1]))
		if len(parts) != m {
			return fmt.Errorf("row %d: expected %d values, got %d", i, m, len(parts))
		}
		vals[i] = make([]int, m)
		for j := 0; j < m; j++ {
			var v int
			if _, err := fmt.Sscan(parts[j], &v); err != nil {
				return fmt.Errorf("row %d col %d: bad int %q", i, j, parts[j])
			}
			vals[i][j] = v
		}
	}

	// Check unmarked cells have value 1 or 4
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' {
				if vals[i][j] != 1 && vals[i][j] != 4 {
					return fmt.Errorf("cell (%d,%d) is unmarked but has value %d (need 1 or 4)", i, j, vals[i][j])
				}
			}
		}
	}

	// Check marked cells: value = sum of adjacent unmarked, and divisible by 5
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != 'X' {
				continue
			}
			sum := 0
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
					sum += vals[ni][nj]
				}
			}
			if vals[i][j] != sum {
				return fmt.Errorf("cell (%d,%d) is marked: expected sum %d, got %d", i, j, sum, vals[i][j])
			}
			if sum%5 != 0 {
				return fmt.Errorf("cell (%d,%d) is marked: sum %d not divisible by 5", i, j, sum)
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		grid := genGrid(rng)
		input := buildInput(grid)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validateOutput(grid, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s\noutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
