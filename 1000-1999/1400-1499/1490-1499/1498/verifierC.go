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

const MOD int64 = 1000000007

func precompute(maxN, maxK int) [][]int64 {
	dp := make([][]int64, maxN+1)
	for i := range dp {
		dp[i] = make([]int64, maxK+1)
	}
	for n := 0; n <= maxN; n++ {
		dp[n][1] = 1
	}
	for k := 1; k <= maxK; k++ {
		dp[0][k] = 1
	}
	for k := 2; k <= maxK; k++ {
		prefix := dp[0][k-1] % MOD
		for n := 1; n <= maxN; n++ {
			dp[n][k] = (1 + prefix) % MOD
			prefix = (prefix + dp[n][k-1]) % MOD
		}
	}
	return dp
}

func generateCase(rng *rand.Rand, dp [][]int64) (string, string) {
	n := rng.Intn(len(dp))
	k := rng.Intn(len(dp[0])-1) + 1
	in := fmt.Sprintf("1\n%d %d\n", n, k)
	out := fmt.Sprintf("%d\n", dp[n][k]%MOD)
	return in, out
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	maxN, maxK := 1000, 1000
	dp := precompute(maxN, maxK)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng, dp)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
