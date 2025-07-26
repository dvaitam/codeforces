package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(n, k int, a []int) int64 {
	dp := make([][]int64, k+1)
	for i := range dp {
		dp[i] = make([]int64, n+1)
		for j := range dp[i] {
			dp[i][j] = math.MaxInt64 / 4
		}
	}
	dp[0][0] = 0
	for seg := 1; seg <= k; seg++ {
		for i := seg; i <= n; i++ {
			maxVal := 0
			for j := i; j >= seg; j-- {
				if a[j-1] > maxVal {
					maxVal = a[j-1]
				}
				val := dp[seg-1][j-1]
				if val != math.MaxInt64/4 {
					cost := val + int64(maxVal*(i-j+1))
					if cost < dp[seg][i] {
						dp[seg][i] = cost
					}
				}
			}
		}
	}
	return dp[k][n]
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	k := rng.Intn(n) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(20) + 1
	}
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, k)
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", v)
	}
	input.WriteByte('\n')
	out := fmt.Sprintf("%d\n", solve(n, k, arr))
	return input.String(), out
}

func runCase(bin string, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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
