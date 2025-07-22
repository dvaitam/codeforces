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

type TestCase struct {
	n, m int
	grid []string
	ans  int
}

func compute(grid []string) int {
	n := len(grid)
	m := len(grid[0])
	g := make([][]byte, n)
	for i := range g {
		g[i] = []byte(grid[i])
	}
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if g[i][j] == 'W' {
				for _, d := range dirs {
					ni, nj := i+d[0], j+d[1]
					if ni >= 0 && ni < n && nj >= 0 && nj < m && g[ni][nj] == 'P' {
						cnt++
						g[ni][nj] = '.'
						break
					}
				}
			}
		}
	}
	return cnt
}

func genCase() TestCase {
	n := rand.Intn(10) + 1
	m := rand.Intn(10) + 1
	g := make([][]byte, n)
	for i := range g {
		g[i] = make([]byte, m)
		for j := range g[i] {
			if rand.Intn(3) == 0 {
				g[i][j] = 'W'
			} else {
				g[i][j] = '.'
			}
		}
	}
	// place pigs respecting the constraint
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if g[i][j] != '.' {
				continue
			}
			wcnt := 0
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && g[ni][nj] == 'W' {
					wcnt++
				}
			}
			if wcnt <= 1 && rand.Intn(3) == 0 {
				g[i][j] = 'P'
			}
		}
	}
	grid := make([]string, n)
	for i := range g {
		grid[i] = string(g[i])
	}
	return TestCase{n, m, grid, compute(grid)}
}

func genCases(n int) []TestCase {
	rand.Seed(time.Now().UnixNano())
	cs := make([]TestCase, n)
	for i := 0; i < n; i++ {
		cs[i] = genCase()
	}
	return cs
}

func buildInput(tc TestCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for _, row := range tc.grid {
		fmt.Fprintln(&sb, row)
	}
	return sb.String()
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases(100)
	for i, tc := range cases {
		input := buildInput(tc)
		expected := fmt.Sprint(tc.ans)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
