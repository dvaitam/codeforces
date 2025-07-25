package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	maxN = 31
	maxK = 51
)

var dp [maxN][maxN][maxK]int
var vis [maxN][maxN][maxK]bool

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(n, m, k int) int {
	if k == 0 || k == n*m {
		return 0
	}
	if vis[n][m][k] {
		return dp[n][m][k]
	}
	vis[n][m][k] = true
	ans := 1<<31 - 1
	for i := 1; i <= n/2; i++ {
		maxK1 := min(k, i*m)
		for k1 := 0; k1 <= maxK1; k1++ {
			if k1 > i*m || k-k1 > (n-i)*m {
				continue
			}
			cost := m*m + solve(i, m, k1) + solve(n-i, m, k-k1)
			if cost < ans {
				ans = cost
			}
		}
	}
	for j := 1; j <= m/2; j++ {
		maxK1 := min(k, j*n)
		for k1 := 0; k1 <= maxK1; k1++ {
			if k1 > j*n || k-k1 > n*(m-j) {
				continue
			}
			cost := n*n + solve(n, j, k1) + solve(n, m-j, k-k1)
			if cost < ans {
				ans = cost
			}
		}
	}
	dp[n][m][k] = ans
	return ans
}

func expectedE(cases [][3]int) []int {
	// reset vis each call
	for i := 0; i < maxN; i++ {
		for j := 0; j < maxN; j++ {
			for k := 0; k < maxK; k++ {
				vis[i][j][k] = false
			}
		}
	}
	res := make([]int, len(cases))
	for i, c := range cases {
		res[i] = solve(c[0], c[1], c[2])
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	cases := make([][3]int, t)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t) + "\n")
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		k := rng.Intn(min(n*m, 30)) + 1
		cases[i] = [3]int{n, m, k}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	}
	ans := expectedE(cases)
	var out strings.Builder
	for i, v := range ans {
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(strconv.Itoa(v))
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected\n%s\ngot\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
