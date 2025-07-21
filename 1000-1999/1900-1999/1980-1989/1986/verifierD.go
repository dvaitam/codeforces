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

const INF int64 = 1e18

func minResult(arr []int) int64 {
	m := len(arr)
	dp := make([]int64, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = INF
	}
	dp[m] = 0
	for i := m - 1; i >= 0; i-- {
		prod := int64(1)
		for j := i; j < m; j++ {
			prod *= int64(arr[j])
			if prod > INF {
				prod = INF
			}
			if val := prod + dp[j+1]; val < dp[i] {
				dp[i] = val
			}
		}
	}
	return dp[0]
}

func solveCase(n int, s string) int64 {
	digits := make([]int, n)
	for i := 0; i < n; i++ {
		digits[i] = int(s[i] - '0')
	}
	ans := INF
	for merge := 0; merge < n-1; merge++ {
		arr := make([]int, 0, n-1)
		for i := 0; i < n; i++ {
			if i == merge {
				val := digits[i]*10 + digits[i+1]
				arr = append(arr, val)
				i++
				continue
			}
			arr = append(arr, digits[i])
		}
		if val := minResult(arr); val < ans {
			ans = val
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	digits := make([]byte, n)
	for i := 0; i < n; i++ {
		digits[i] = byte('0' + rng.Intn(10))
	}
	s := string(digits)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n%s\n", n, s)
	expected := fmt.Sprintf("%d\n", solveCase(n, s))
	return sb.String(), expected
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
