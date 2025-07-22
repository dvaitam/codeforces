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

const mod = 1000000007

func runCandidate(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func computeDP(n, m int, grid [][]byte, sx, sy int) (int32, int32) {
	t1x, t1y := n-2, m-1
	t2x, t2y := n-1, m-2
	dp := make([][]int32, n)
	for i := range dp {
		dp[i] = make([]int32, m)
	}
	if grid[sx][sy] == '#' {
		return 0, 0
	}
	dp[sx][sy] = 1
	for i := sx; i < n; i++ {
		for j := sy; j < m; j++ {
			if i == sx && j == sy {
				continue
			}
			if grid[i][j] == '#' {
				dp[i][j] = 0
				continue
			}
			var v int32
			if i > sx {
				v += dp[i-1][j]
			}
			if j > sy {
				v += dp[i][j-1]
			}
			if v >= mod {
				v %= mod
			}
			dp[i][j] = v
		}
	}
	var t1, t2 int32
	if t1x >= sx && t1y >= sy {
		t1 = dp[t1x][t1y]
	}
	if t2x >= sx && t2y >= sy {
		t2 = dp[t2x][t2y]
	}
	return t1, t2
}

func solveCase(n, m int, grid [][]byte) int64 {
	a11, a12 := computeDP(n, m, grid, 0, 1)
	a21, a22 := computeDP(n, m, grid, 1, 0)
	res := (int64(a11)*int64(a22) - int64(a12)*int64(a21)) % mod
	if res < 0 {
		res += mod
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	m := rng.Intn(5) + 2
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			if (i == 0 && j == 0) || (i == n-1 && j == m-1) {
				grid[i][j] = '.'
			} else {
				if rng.Intn(4) == 0 {
					grid[i][j] = '#'
				} else {
					grid[i][j] = '.'
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			sb.WriteByte(grid[i][j])
		}
		sb.WriteByte('\n')
	}
	expect := fmt.Sprintf("%d", solveCase(n, m, grid))
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
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
