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

func maxK(arr []int, k int) int {
	n := len(arr)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, k+1)
		for j := range dp[i] {
			dp[i][j] = math.MinInt32
		}
	}
	dp[0][0] = 0
	prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + arr[i]
	}
	for i := 1; i <= n; i++ {
		for j := 0; j <= k; j++ {
			if dp[i-1][j] > dp[i][j] {
				dp[i][j] = dp[i-1][j]
			}
		}
		for j := 1; j <= k; j++ {
			for t := 0; t < i; t++ {
				if dp[t][j-1] == math.MinInt32 {
					continue
				}
				sum := prefix[i] - prefix[t]
				if val := dp[t][j-1] + sum; val > dp[i][j] {
					dp[i][j] = val
				}
			}
		}
	}
	ans := 0
	for j := 1; j <= k; j++ {
		if dp[n][j] > ans {
			ans = dp[n][j]
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(21) - 10
	}
	m := rng.Intn(15) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", m))
	expOutputs := []string{}
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			pos := rng.Intn(n) + 1
			val := rng.Intn(21) - 10
			arr[pos-1] = val
			sb.WriteString(fmt.Sprintf("0 %d %d\n", pos, val))
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			k := rng.Intn(3) + 1
			ans := maxK(arr[l-1:r], k)
			sb.WriteString(fmt.Sprintf("1 %d %d %d\n", l, r, k))
			expOutputs = append(expOutputs, fmt.Sprintf("%d", ans))
		}
	}
	sbInput := sb.String()
	exp := strings.Join(expOutputs, "\n")
	if len(expOutputs) > 0 {
		exp += "\n"
	}
	return sbInput, exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
