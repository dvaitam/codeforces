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

type point struct{ x, y int }

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

func solveCase(grid [][]byte, n, m int) string {
	var sx, sy int
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'L' {
				sx, sy = i, j
			}
		}
	}
	deg := make([][]int, n)
	for i := 0; i < n; i++ {
		deg[i] = make([]int, m)
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				continue
			}
			cnt := 0
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] != '#' {
					cnt++
				}
			}
			deg[i][j] = cnt
		}
	}
	q := []point{{sx, sy}}
	for head := 0; head < len(q); head++ {
		p := q[head]
		for _, d := range dirs {
			ni, nj := p.x+d[0], p.y+d[1]
			if ni < 0 || ni >= n || nj < 0 || nj >= m {
				continue
			}
			if grid[ni][nj] != '.' {
				continue
			}
			deg[ni][nj]--
			if deg[ni][nj] <= 1 {
				grid[ni][nj] = '+'
				q = append(q, point{ni, nj})
			}
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	for n*m > 16 { // ensure small
		n = rng.Intn(4) + 1
		m = rng.Intn(4) + 1
	}
	grid := make([][]byte, n)
	labPlaced := false
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if !labPlaced && rng.Intn(n*m) == 0 {
				row[j] = 'L'
				labPlaced = true
			} else {
				if rng.Intn(4) == 0 {
					row[j] = '#'
				} else {
					row[j] = '.'
				}
			}
		}
		grid[i] = row
	}
	if !labPlaced {
		grid[0][0] = 'L'
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	expect := solveCase(copyGrid(grid), n, m)
	return sb.String(), expect
}

func copyGrid(g [][]byte) [][]byte {
	n := len(g)
	res := make([][]byte, n)
	for i := 0; i < n; i++ {
		res[i] = append([]byte(nil), g[i]...)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
