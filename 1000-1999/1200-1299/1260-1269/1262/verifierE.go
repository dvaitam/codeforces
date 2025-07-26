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

type cell struct{ x, y int }

var dirs = []cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

func maxT(grid []string) (int, []string) {
	n := len(grid)
	m := len(grid[0])
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, m)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	q := []cell{}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' {
				dist[i][j] = 0
				q = append(q, cell{i, j})
			} else if i == 0 || i == n-1 || j == 0 || j == m-1 {
				dist[i][j] = 1
				q = append(q, cell{i, j})
			}
		}
	}
	for h := 0; h < len(q); h++ {
		c := q[h]
		for _, d := range dirs {
			ni := c.x + d.x
			nj := c.y + d.y
			if ni < 0 || nj < 0 || ni >= n || nj >= m {
				continue
			}
			if dist[ni][nj] == -1 {
				dist[ni][nj] = dist[c.x][c.y] + 1
				q = append(q, cell{ni, nj})
			}
		}
	}
	T := 1 << 30
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'X' {
				if dist[i][j] < T {
					T = dist[i][j]
				}
			}
		}
	}
	if T == 1<<30 {
		T = 0
	} else {
		T--
	}
	seeds := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			if dist[i][j] > T {
				b[j] = 'X'
			} else {
				b[j] = '.'
			}
		}
		seeds[i] = string(b)
	}
	return T, seeds
}

func simulate(seeds []string, T int) []string {
	n := len(seeds)
	m := len(seeds[0])
	burn := make([][]int, n)
	for i := range burn {
		burn[i] = make([]int, m)
	}
	q := []cell{}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if seeds[i][j] == 'X' {
				burn[i][j] = 0
				q = append(q, cell{i, j})
			} else {
				burn[i][j] = -1
			}
		}
	}
	for h := 0; h < len(q); h++ {
		c := q[h]
		if burn[c.x][c.y] == T {
			continue
		}
		for _, d := range dirs {
			ni := c.x + d.x
			nj := c.y + d.y
			if ni < 0 || nj < 0 || ni >= n || nj >= m {
				continue
			}
			if burn[ni][nj] == -1 {
				burn[ni][nj] = burn[c.x][c.y] + 1
				q = append(q, cell{ni, nj})
			}
		}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			if burn[i][j] >= 0 {
				b[j] = 'X'
			} else {
				b[j] = '.'
			}
		}
		res[i] = string(b)
	}
	return res
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
	return out.String(), nil
}

func runCase(bin string, grid []string) error {
	n := len(grid)
	m := len(grid[0])
	input := fmt.Sprintf("%d %d\n", n, m)
	for _, row := range grid {
		input += row + "\n"
	}
	T, _ := maxT(grid)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != n+1 {
		return fmt.Errorf("expected %d lines got %d", n+1, len(lines))
	}
	var gotT int
	if _, err := fmt.Sscan(lines[0], &gotT); err != nil {
		return fmt.Errorf("bad T: %v", err)
	}
	if gotT != T {
		return fmt.Errorf("expected T=%d got %d", T, gotT)
	}
	for i := 0; i < n; i++ {
		if len(lines[i+1]) != m {
			return fmt.Errorf("bad row length")
		}
	}
	sim := simulate(lines[1:], gotT)
	for i := 0; i < n; i++ {
		if sim[i] != grid[i] {
			return fmt.Errorf("simulation mismatch")
		}
	}
	return nil
}

func genCase(rng *rand.Rand) []string {
	n := rng.Intn(3) + 2
	m := rng.Intn(3) + 2
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				b[j] = 'X'
			} else {
				b[j] = '.'
			}
		}
		grid[i] = string(b)
	}
	// ensure at least one X
	xFound := false
	for i := 0; i < n; i++ {
		if strings.Contains(grid[i], "X") {
			xFound = true
			break
		}
	}
	if !xFound {
		grid[0] = "X" + grid[0][1:]
	}
	return grid
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	grid := []string{"XXXXXX", "XXXXXX", "XXXXXX"}
	if err := runCase(bin, grid); err != nil {
		fmt.Fprintf(os.Stderr, "case 1 failed: %v\n", err)
		os.Exit(1)
	}
	for i := 1; i < 100; i++ {
		g := genCase(rng)
		if err := runCase(bin, g); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
