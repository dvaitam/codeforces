package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, exitErr.Stderr)
		}
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func countSubsequences(s, p string) int64 {
	n := len(s)
	m := len(p)
	if m == 0 {
		return 0
	}
	if n == 0 {
		return 0
	}
	
	// dp[j] will store the number of ways to form p[:j] using a prefix of s
	dp := make([]int64, m+1)
	dp[0] = 1

	for i := 0; i < n; i++ {
		// Traverse backwards to avoid using the same character of s multiple times for the same position in p
		for j := m; j >= 1; j-- {
			if s[i] == p[j-1] {
				dp[j] += dp[j-1]
			}
		}
	}
	return dp[m]
}

func verify(n int, output string) error {
	parts := strings.Split(output, " ")
	if len(parts) != 2 {
		return fmt.Errorf("expected 2 strings separated by space, got %d parts", len(parts))
	}
	s := parts[0]
	p := parts[1]

	if len(s) > 200 {
		return fmt.Errorf("string s too long: %d > 200", len(s))
	}
	if len(p) > 200 {
		return fmt.Errorf("string p too long: %d > 200", len(p))
	}
	if len(s) == 0 || len(p) == 0 {
		return fmt.Errorf("empty strings not allowed")
	}

	count := countSubsequences(s, p)
	if count != int64(n) {
		return fmt.Errorf("subsequence count mismatch: expected %d, got %d", n, count)
	}
	return nil
}

func genTests() []int {
	rand.Seed(time.Now().UnixNano())
	tests := []int{1, 2, 3, 39, 100, 1000, 1000000}
	for i := 0; i < 93; i++ {
		tests = append(tests, rand.Intn(1000000)+1)
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := genTests()
	for i, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d (n=%d) failed: %v\n", i+1, n, err)
			os.Exit(1)
		}
		if err := verify(n, out); err != nil {
			fmt.Printf("test %d (n=%d) failed: %v\nOutput: %s\n", i+1, n, err, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
