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

func expected(n int, p float64, t int) string {
	dp := make([]float64, n+1)
	dp[0] = 1.0
	for timeStep := 1; timeStep <= t; timeStep++ {
		dp[n] += p * dp[n-1]
		for i := n - 1; i > 0; i-- {
			dp[i] = p*dp[i-1] + (1-p)*dp[i]
		}
		dp[0] *= (1 - p)
	}
	ans := 0.0
	for i := 0; i <= n; i++ {
		ans += float64(i) * dp[i]
	}
	return fmt.Sprintf("%.10f", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	t := rng.Intn(20) + 1
	p := rng.Float64()
	p = float64(int(p*100)) / 100.0
	input := fmt.Sprintf("%d %.2f %d\n", n, p, t)
	exp := expected(n, p, t)
	return input, exp
}

func runCase(bin, input, exp string) error {
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
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
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
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
