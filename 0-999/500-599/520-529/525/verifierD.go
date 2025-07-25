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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(n, m int, grid [][]byte) string {
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			if grid[r][c] == '.' && !visited[r][c] {
				stack := [][2]int{{r, c}}
				visited[r][c] = true
				minr, maxr := r, r
				minc, maxc := c, c
				for len(stack) > 0 {
					x := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					cr, cc := x[0], x[1]
					if cr < minr {
						minr = cr
					}
					if cr > maxr {
						maxr = cr
					}
					if cc < minc {
						minc = cc
					}
					if cc > maxc {
						maxc = cc
					}
					for _, d := range dirs {
						nr, nc := cr+d[0], cc+d[1]
						if nr >= 0 && nr < n && nc >= 0 && nc < m && !visited[nr][nc] && grid[nr][nc] == '.' {
							visited[nr][nc] = true
							stack = append(stack, [2]int{nr, nc})
						}
					}
				}
				for i := minr; i <= maxr; i++ {
					for j := minc; j <= maxc; j++ {
						grid[i][j] = '.'
					}
				}
			}
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.Write(grid[i])
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	m := rng.Intn(6) + 1
	grid := make([][]byte, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '.'
			} else {
				row[j] = '*'
			}
		}
		grid[i] = append([]byte(nil), row...)
		sb.Write(row)
		sb.WriteByte('\n')
	}
	expect := solveCase(n, m, grid)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
