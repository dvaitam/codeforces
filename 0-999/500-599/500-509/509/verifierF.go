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

const MOD = 1000000007

func solveTree(b []int) string {
	n := len(b) - 1
	if n <= 1 {
		return "1"
	}
	dp := make([][]int, n+2)
	for i := range dp {
		dp[i] = make([]int, n+2)
	}
	for i := 1; i <= n; i++ {
		dp[i][i] = 1
	}
	for length := 2; length <= n; length++ {
		for l := 1; l+length-1 <= n; l++ {
			r := l + length - 1
			for k := l + 1; k <= r; k++ {
				if k == r || b[k+1] > b[l+1] {
					dp[l][r] = (dp[l][r] + dp[l+1][k]*dp[k][r]) % MOD
				}
			}
		}
	}
	return strconv.Itoa(dp[1][n])
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(6) + 1
	b := make([]int, n+1)
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		b[i+1] = perm[i] + 1
	}
	b[0] = 0
	if b[1] != 1 {
		idx := 1
		for i := 1; i <= n; i++ {
			if b[i] == 1 {
				idx = i
				break
			}
		}
		b[1], b[idx] = b[idx], b[1]
	}
	return b
}

func runCase(bin string, b []int) error {
	var sb strings.Builder
	n := len(b) - 1
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(b[i]))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solveTree(b)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		b := generateCase(rng)
		if err := runCase(bin, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
