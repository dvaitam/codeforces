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

func countComponents(grid [][]int) int {
	n := len(grid)
	m := len(grid[0])
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	comps := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if !visited[i][j] {
				comps++
				stack := [][2]int{{i, j}}
				visited[i][j] = true
				for len(stack) > 0 {
					cur := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					x, y := cur[0], cur[1]
					for _, d := range dirs {
						nx, ny := x+d[0], y+d[1]
						if nx >= 0 && nx < n && ny >= 0 && ny < m && !visited[nx][ny] && grid[nx][ny] == grid[x][y] {
							visited[nx][ny] = true
							stack = append(stack, [2]int{nx, ny})
						}
					}
				}
			}
		}
	}
	return comps
}

func solveCase(n, m, q int, queries [][3]int) string {
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, m)
	}
	var sb strings.Builder
	for _, qu := range queries {
		x, y, c := qu[0], qu[1], qu[2]
		grid[x][y] = c
		comps := countComponents(grid)
		fmt.Fprintf(&sb, "%d\n", comps)
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	q := rng.Intn(5) + 1
	queries := make([][3]int, q)
	lastC := 0
	for i := 0; i < q; i++ {
		x := rng.Intn(n)
		y := rng.Intn(m)
		c := lastC + rng.Intn(3) + 1
		lastC = c
		queries[i] = [3]int{x, y, c}
	}
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
	for _, qu := range queries {
		fmt.Fprintf(&sb, "%d %d %d\n", qu[0]+1, qu[1]+1, qu[2])
	}
	out := solveCase(n, m, q, queries)
	return sb.String(), out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
