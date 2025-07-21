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

func solveCase(n, h int) uint64 {
	dp := make([][]uint64, n+1)
	for i := range dp {
		dp[i] = make([]uint64, n+1)
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = 1
	}
	for j := 1; j <= n; j++ {
		for i := 1; i <= n; i++ {
			var cnt uint64
			for left := 0; left < i; left++ {
				right := i - 1 - left
				cnt += dp[left][j-1] * dp[right][j-1]
			}
			dp[i][j] = cnt
		}
	}
	total := dp[n][n]
	less := uint64(0)
	if h > 0 {
		less = dp[n][h-1]
	}
	return total - less
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	h := rng.Intn(n) + 1
	input := fmt.Sprintf("%d %d\n", n, h)
	expect := fmt.Sprintf("%d", solveCase(n, h))
	return input, expect
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
