package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const INF = int(1e9)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveB(n, m int, grid [][]byte) int {
	dp1H := make([][]int, n)
	dp1V := make([][]int, n)
	for i := 0; i < n; i++ {
		dp1H[i] = make([]int, m)
		dp1V[i] = make([]int, m)
		for j := 0; j < m; j++ {
			dp1H[i][j] = INF
			dp1V[i][j] = INF
		}
	}
	dp1H[0][0] = 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if i == 0 && j == 0 {
				continue
			}
			if j > 0 {
				dp1H[i][j] = min(dp1H[i][j], dp1H[i][j-1])
				if grid[i][j-1] == '#' && dp1V[i][j-1] < INF {
					dp1H[i][j] = min(dp1H[i][j], dp1V[i][j-1]+1)
				}
			}
			if i > 0 {
				dp1V[i][j] = min(dp1V[i][j], dp1V[i-1][j])
				if grid[i-1][j] == '#' && dp1H[i-1][j] < INF {
					dp1V[i][j] = min(dp1V[i][j], dp1H[i-1][j]+1)
				}
			}
		}
	}
	ds := make([][]int, n)
	for i := 0; i < n; i++ {
		ds[i] = make([]int, m)
		for j := 0; j < m; j++ {
			ds[i][j] = min(dp1H[i][j], dp1V[i][j])
		}
	}

	dp2H := make([][]int, n)
	dp2V := make([][]int, n)
	for i := 0; i < n; i++ {
		dp2H[i] = make([]int, m)
		dp2V[i] = make([]int, m)
		for j := 0; j < m; j++ {
			dp2H[i][j] = INF
			dp2V[i][j] = INF
		}
	}
	dp2H[n-1][m-1] = 0
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if i == n-1 && j == m-1 {
				continue
			}
			if j < m-1 {
				dp2H[i][j] = min(dp2H[i][j], dp2H[i][j+1])
				if grid[i][j+1] == '#' && dp2V[i][j+1] < INF {
					dp2H[i][j] = min(dp2H[i][j], dp2V[i][j+1]+1)
				}
			}
			if i < n-1 {
				dp2V[i][j] = min(dp2V[i][j], dp2V[i+1][j])
				if grid[i+1][j] == '#' && dp2H[i+1][j] < INF {
					dp2V[i][j] = min(dp2V[i][j], dp2H[i+1][j]+1)
				}
			}
		}
	}
	db := make([][]int, n)
	for i := 0; i < n; i++ {
		db[i] = make([]int, m)
		for j := 0; j < m; j++ {
			db[i][j] = min(dp2H[i][j], dp2V[i][j])
		}
	}
	ans := INF
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' && ds[i][j] < INF && db[i][j] < INF {
				cost := ds[i][j] + db[i][j]
				if cost < ans {
					ans = cost
				}
			}
		}
	}
	if ans >= INF {
		return -1
	}
	return ans
}

func runCase(bin string, n, m int, grid [][]byte) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		input.Write(grid[i])
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(&out, &got); err != nil {
		return fmt.Errorf("parse error: %v", err)
	}
	want := solveB(n, m, grid)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const tests = 100
	for t := 0; t < tests; t++ {
		n := rng.Intn(10) + 2
		m := rng.Intn(10) + 2
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				if rng.Intn(3) == 0 {
					row[j] = '#'
				} else {
					row[j] = '.'
				}
			}
			grid[i] = row
		}
		if err := runCase(bin, n, m, grid); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
