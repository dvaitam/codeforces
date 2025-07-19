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

const modF int64 = 998244353

func expectedAnswerF(n, x, m int) int64 {
	fib := make([]int64, 31)
	fib[1], fib[2] = 1, 1
	for i := 3; i <= 30; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}
	limit := int(fib[x] * int64(n))
	dp := make([][]int64, n+1)
	for i := range dp {
		dp[i] = make([]int64, limit+1)
	}
	dp[0][0] = 1
	for i := 1; i <= x; i++ {
		fi := int(fib[i])
		for j := 1; j <= n; j++ {
			for l := fi; l <= fi*j; l++ {
				dp[j][l] = (dp[j][l] + dp[j-1][l-fi]) % modF
			}
		}
	}
	ans := int64(0)
	for i := 0; i <= limit; i++ {
		t := i
		c := 0
		for j := 30; j >= 1; j-- {
			fj := int(fib[j])
			c += t / fj
			t %= fj
		}
		if c == m {
			ans = (ans + dp[n][i]) % modF
		}
	}
	return ans
}

func generateCaseF(rng *rand.Rand) (int, int, int) {
	n := rng.Intn(4) + 1
	x := rng.Intn(5) + 1
	m := rng.Intn(5)
	return n, x, m
}

func runCaseF(bin string, n, x, m int) error {
	input := fmt.Sprintf("%d %d %d\n", n, x, m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprint(expectedAnswerF(n, x, m))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, x, m := generateCaseF(rng)
		if err := runCaseF(bin, n, x, m); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
