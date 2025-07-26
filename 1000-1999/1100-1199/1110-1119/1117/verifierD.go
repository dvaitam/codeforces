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

func expectedD(N, M int) int64 {
	dp := make([]int64, N+1)
	dp[0] = 1
	for i := 1; i <= N; i++ {
		dp[i] = dp[i-1]
		if i >= M {
			dp[i] = (dp[i] + dp[i-M]) % MOD
		}
	}
	return dp[N]
}

func generateCase(rng *rand.Rand) (string, int64) {
	N := rng.Intn(40) + 1
	M := rng.Intn(8) + 2
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", N, M)
	return sb.String(), expectedD(N, M)
}

func runCase(bin, input string, exp int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
