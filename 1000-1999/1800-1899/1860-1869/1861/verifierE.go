package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct{ input, expected string }

const mod int64 = 998244353

func solveCase(n, k int) string {
	dp := make([]int64, k)
	sum := make([]int64, k)
	dp[0] = 1
	prefixDp := make([]int64, k+1)
	prefixSum := make([]int64, k+1)
	nextDp := make([]int64, k)
	nextSum := make([]int64, k)
	for step := 0; step < n; step++ {
		prefixDp[k] = 0
		prefixSum[k] = 0
		for i := k - 1; i >= 0; i-- {
			prefixDp[i] = prefixDp[i+1] + dp[i]
			if prefixDp[i] >= mod {
				prefixDp[i] %= mod
			}
			prefixSum[i] = prefixSum[i+1] + sum[i]
			if prefixSum[i] >= mod {
				prefixSum[i] %= mod
			}
		}
		for i := 0; i < k; i++ {
			nextDp[i] = 0
			nextSum[i] = 0
		}
		for r := 1; r < k; r++ {
			nextDp[r] = prefixDp[r]
			nextSum[r] = prefixSum[r]
		}
		for m := 0; m < k; m++ {
			if dp[m] == 0 && sum[m] == 0 {
				continue
			}
			newChoices := int64(k - m)
			if m+1 == k {
				nextDp[0] = (nextDp[0] + dp[m]*newChoices) % mod
				inc := (sum[m] + dp[m]) % mod
				nextSum[0] = (nextSum[0] + inc*newChoices) % mod
			} else {
				nextDp[m+1] = (nextDp[m+1] + dp[m]*newChoices) % mod
				nextSum[m+1] = (nextSum[m+1] + sum[m]*newChoices) % mod
			}
		}
		dp, nextDp = nextDp, dp
		sum, nextSum = nextSum, sum
	}
	var ans int64
	for i := 0; i < k; i++ {
		ans = (ans + sum[i]) % mod
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		k := rng.Intn(5) + 1
		input := fmt.Sprintf("%d %d\n", n, k)
		tests = append(tests, test{input, solveCase(n, k)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
