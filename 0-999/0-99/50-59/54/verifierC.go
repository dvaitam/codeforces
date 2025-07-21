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

func getn(a int64) int64 {
	if a <= 0 {
		return 0
	}
	var c int64 = 1
	for i := 0; i < 20; i++ {
		if a/c == 0 {
			c /= 10
			break
		}
		c *= 10
	}
	var ans int64
	if a/c != 1 {
		ans += c
	} else {
		ans += a%c + 1
	}
	c /= 10
	for c > 0 {
		ans += c
		c /= 10
	}
	return ans
}

func probability(l, r int64) float64 {
	total := r - l + 1
	cnt := getn(r) - getn(l-1)
	return float64(cnt) / float64(total)
}

func solve(ranges [][2]int64, kPercent int) float64 {
	n := len(ranges)
	pro := make([]float64, n)
	for i, lr := range ranges {
		pro[i] = probability(lr[0], lr[1])
	}
	k := (n*kPercent + 99) / 100
	dp := make([][]float64, n+1)
	for i := range dp {
		dp[i] = make([]float64, k+1)
	}
	dp[0][0] = 1
	for i := 1; i <= n; i++ {
		for j := 0; j <= i && j <= k; j++ {
			if j == 0 {
				dp[i][j] = dp[i-1][j] * (1 - pro[i-1])
			} else {
				dp[i][j] = dp[i-1][j-1]*pro[i-1] + dp[i-1][j]*(1-pro[i-1])
			}
		}
	}
	sum := 0.0
	for j := 0; j < k; j++ {
		sum += dp[n][j]
	}
	ans := 1.0 - sum
	if ans < 0 {
		ans = 0
	}
	if ans > 1 {
		ans = 1
	}
	return ans
}

func runCase(exe string, input string, expected float64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if math.Abs(got-expected) > 1e-6*math.Max(1, math.Abs(expected)) {
		return fmt.Errorf("expected %f got %f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		ranges := make([][2]int64, n)
		for j := 0; j < n; j++ {
			l := rng.Int63n(1_000_000) + 1
			r := l + rng.Int63n(1000)
			ranges[j] = [2]int64{l, r}
			sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
		}
		k := rng.Intn(101)
		sb.WriteString(fmt.Sprintf("%d\n", k))
		input := sb.String()
		exp := solve(ranges, k)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
