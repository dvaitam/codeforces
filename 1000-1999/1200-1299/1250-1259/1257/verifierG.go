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

const mod int64 = 998244353

func expected(counts []int) int64 {
	n := 0
	for _, c := range counts {
		n += c
	}
	limit := n / 2
	dp := make([]int64, limit+1)
	dp[0] = 1
	for _, c := range counts {
		if c > limit {
			c = limit
		}
		next := make([]int64, limit+1)
		for i := 0; i <= limit; i++ {
			if dp[i] == 0 {
				continue
			}
			for j := 0; j <= c && i+j <= limit; j++ {
				next[i+j] = (next[i+j] + dp[i]) % mod
			}
		}
		dp = next
	}
	return dp[limit] % mod
}

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(5) + 1
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		counts[i] = rng.Intn(4) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", counts[i]))
	}
	sb.WriteByte('\n')
	exp := []string{fmt.Sprintf("%d", expected(counts))}
	return sb.String(), exp
}

func runCase(bin, input string, exp []string) error {
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, line := range lines {
		if strings.TrimSpace(line) != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], strings.TrimSpace(line))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
