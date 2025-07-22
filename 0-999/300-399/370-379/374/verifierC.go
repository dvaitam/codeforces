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

type pair struct{ x, y int }

var dirs = [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
var nextCh = map[byte]byte{'D': 'I', 'I': 'M', 'M': 'A', 'A': 'D'}

func expected(grid [][]byte) string {
	n := len(grid)
	m := len(grid[0])
	dp := make([][]int, n)
	vis := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, m)
		vis[i] = make([]int, m)
	}
	infinite := false
	var dfs func(int, int) int
	dfs = func(i, j int) int {
		if infinite {
			return 0
		}
		if vis[i][j] == 1 {
			infinite = true
			return 0
		}
		if vis[i][j] == 2 {
			return dp[i][j]
		}
		vis[i][j] = 1
		best := 1
		cur := grid[i][j]
		nxt, ok := nextCh[cur]
		if ok {
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == nxt {
					v := dfs(ni, nj)
					if infinite {
						return 0
					}
					if v+1 > best {
						best = v + 1
					}
				}
			}
		}
		vis[i][j] = 2
		dp[i][j] = best
		return best
	}
	maxLen := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'D' {
				v := dfs(i, j)
				if infinite {
					return "Poor Inna!"
				}
				if v > maxLen {
					maxLen = v
				}
			}
		}
	}
	if maxLen < 4 {
		return "Poor Dima!"
	}
	return fmt.Sprintf("%d", maxLen/4)
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []byte{'D', 'I', 'M', 'A'}
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(6) + 1
		grid := make([][]byte, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				row[j] = letters[rng.Intn(4)]
			}
			grid[i] = row
			sb.WriteString(string(row))
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp := expected(grid)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
