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

func computeS(n, m, x int, skipX bool) int64 {
	dp := make([][]int64, n+1)
	for i := range dp {
		dp[i] = make([]int64, m+1)
	}
	dp[0][0] = 1
	for l := 1; l <= m; l++ {
		dp2 := make([][]int64, n+1)
		for i := 0; i <= n; i++ {
			row := make([]int64, m+1)
			copy(row, dp[i])
			dp2[i] = row
		}
		if !(skipX && l == x) {
			for i := 0; i < n; i++ {
				pre := make([]int64, m+1)
				pre[0] = dp[i][0]
				for r := 1; r <= m; r++ {
					pre[r] = pre[r-1] + dp[i][r]
					if pre[r] >= mod {
						pre[r] -= mod
					}
				}
				for r := l; r <= m; r++ {
					dp2[i+1][r] = (dp2[i+1][r] + pre[r-1]) % mod
				}
			}
		}
		dp = dp2
	}
	var total int64
	for r := 1; r <= m; r++ {
		total = (total + dp[n][r]) % mod
	}
	return total
}

func solveE(n, m, x int) int64 {
	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	total := computeS(n, m, x, false)
	bad := computeS(n, m, x, true)
	ways := (total - bad + mod) % mod
	return ways * fact[n] % mod
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(6) + 1
	for n*m > 20 {
		n = rng.Intn(4) + 1
		m = rng.Intn(6) + 1
	}
	x := rng.Intn(m) + 1
	input := fmt.Sprintf("%d %d %d\n", n, m, x)
	ans := solveE(n, m, x)
	expected := fmt.Sprintf("%d", ans)
	return input, expected
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
