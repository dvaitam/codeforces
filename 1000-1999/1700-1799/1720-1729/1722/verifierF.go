package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pair struct{ r, c int }

type testCaseF struct {
	n, m int
	grid []string
}

func check(grid [][]byte, n, m int) bool {
	vis := make([][]int, n)
	for i := range vis {
		vis[i] = make([]int, m)
	}
	dirs4 := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	dirs8 := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	id := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '*' && vis[i][j] == 0 {
				id++
				q := []pair{{i, j}}
				vis[i][j] = id
				for k := 0; k < len(q); k++ {
					p := q[k]
					for _, d := range dirs4 {
						nr, nc := p.r+d[0], p.c+d[1]
						if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == '*' && vis[nr][nc] == 0 {
							vis[nr][nc] = id
							q = append(q, pair{nr, nc})
						}
					}
				}
				if len(q) != 3 {
					return false
				}
				minr, maxr := q[0].r, q[0].r
				minc, maxc := q[0].c, q[0].c
				for _, p := range q[1:] {
					if p.r < minr {
						minr = p.r
					}
					if p.r > maxr {
						maxr = p.r
					}
					if p.c < minc {
						minc = p.c
					}
					if p.c > maxc {
						maxc = p.c
					}
				}
				if maxr-minr != 1 || maxc-minc != 1 {
					return false
				}
				count := 0
				for r := minr; r <= maxr; r++ {
					for c := minc; c <= maxc; c++ {
						if grid[r][c] == '*' {
							if vis[r][c] != id {
								return false
							}
							count++
						}
					}
				}
				if count != 3 {
					return false
				}
				for _, p := range q {
					for _, d := range dirs8 {
						nr, nc := p.r+d[0], p.c+d[1]
						if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == '*' && vis[nr][nc] != id {
							return false
						}
					}
				}
			}
		}
	}
	return true
}

func solveF(n, m int, lines []string) string {
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = []byte(lines[i])
	}
	if check(grid, n, m) {
		return "YES"
	}
	return "NO"
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCaseF {
	rand.Seed(47)
	tests := make([]testCaseF, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			b := make([]byte, m)
			for c := 0; c < m; c++ {
				if rand.Intn(4) == 0 {
					b[c] = '*'
				} else {
					b[c] = '.'
				}
			}
			grid[r] = string(b)
		}
		tests[i] = testCaseF{n, m, grid}
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, line := range tc.grid {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
		input := sb.String()
		expected := solveF(tc.n, tc.m, tc.grid)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(output) != expected {
			fmt.Printf("test %d failed:\ninput:%sexpected %s got %s\n", i+1, input, expected, output)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
