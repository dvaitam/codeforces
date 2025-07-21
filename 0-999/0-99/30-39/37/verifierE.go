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

type caseE struct {
	n, m int
	grid [][]byte
}

func generateCase(rng *rand.Rand) caseE {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = 'B'
			} else {
				row[j] = 'W'
			}
		}
		grid[i] = row
	}
	return caseE{n, m, grid}
}

func countComp(grid [][]byte, color byte) int {
	n := len(grid)
	m := len(grid[0])
	vis := make([][]bool, n)
	for i := range vis {
		vis[i] = make([]bool, m)
	}
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	cnt := 0
	var q [][2]int
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if !vis[i][j] && grid[i][j] == color {
				cnt++
				vis[i][j] = true
				q = q[:0]
				q = append(q, [2]int{i, j})
				for qi := 0; qi < len(q); qi++ {
					x, y := q[qi][0], q[qi][1]
					for _, d := range dirs {
						nx, ny := x+d[0], y+d[1]
						if nx >= 0 && nx < n && ny >= 0 && ny < m && !vis[nx][ny] && grid[nx][ny] == color {
							vis[nx][ny] = true
							q = append(q, [2]int{nx, ny})
						}
					}
				}
			}
		}
	}
	return cnt
}

func solveCase(tc caseE) int {
	b := countComp(tc.grid, 'B')
	w := countComp(tc.grid, 'W')
	ans := b
	if w+1 < ans {
		ans = w + 1
	}
	return ans
}

func runCase(bin string, tc caseE) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			input.WriteByte(tc.grid[i][j])
		}
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := fmt.Sprintf("%d", solveCase(tc))
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
