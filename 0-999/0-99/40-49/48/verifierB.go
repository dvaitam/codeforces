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

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func rectSum(ps [][]int, r, c, h, w int) int {
	r0, c0 := r, c
	r1, c1 := r+h, c+w
	return ps[r1][c1] - ps[r0][c1] - ps[r1][c0] + ps[r0][c0]
}

func expected(n, m int, grid [][]int, a, b int) string {
	ps := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		ps[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			ps[i][j] = grid[i-1][j-1] + ps[i-1][j] + ps[i][j-1] - ps[i-1][j-1]
		}
	}
	inf := n*m + 1
	best := inf
	if a <= n && b <= m {
		for i := 0; i+a <= n; i++ {
			for j := 0; j+b <= m; j++ {
				s := rectSum(ps, i, j, a, b)
				if s < best {
					best = s
				}
			}
		}
	}
	if b <= n && a <= m {
		for i := 0; i+b <= n; i++ {
			for j := 0; j+a <= m; j++ {
				s := rectSum(ps, i, j, b, a)
				if s < best {
					best = s
				}
			}
		}
	}
	if best == inf {
		best = 0
	}
	return fmt.Sprintf("%d", best)
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(5) + 1
	m := r.Intn(5) + 1
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, m)
		for j := range grid[i] {
			grid[i][j] = r.Intn(2)
		}
	}
	a := r.Intn(n) + 1
	b := r.Intn(m) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	return sb.String(), expected(n, m, grid, a, b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
