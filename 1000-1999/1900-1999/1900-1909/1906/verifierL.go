package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input string
	n, k  int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(strings.TrimSpace(out), tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(out string, tc testCase) error {
	if out == "-1" {
		if hasSolution(tc.n, tc.k) {
			return fmt.Errorf("solution exists but got -1")
		}
		return nil
	}
	if !hasSolution(tc.n, tc.k) {
		return fmt.Errorf("no solution should exist but got %q", out)
	}
	if len(out) != tc.n {
		return fmt.Errorf("output length %d != n %d", len(out), tc.n)
	}
	for _, ch := range out {
		if ch != '(' && ch != ')' {
			return fmt.Errorf("invalid character %q", ch)
		}
	}
	if !isBalanced(out) {
		return fmt.Errorf("sequence not balanced")
	}
	if lps(out) != tc.k {
		return fmt.Errorf("expected LPS %d but got %d", tc.k, lps(out))
	}
	return nil
}

func hasSolution(n, k int) bool {
	if k == n {
		return false
	}
	return 2*k > n
}

func isBalanced(s string) bool {
	balance := 0
	for _, ch := range s {
		if ch == '(' {
			balance++
		} else {
			balance--
		}
		if balance < 0 {
			return false
		}
	}
	return balance == 0
}

func lps(s string) int {
	n := len(s)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
		dp[i][i] = 1
	}
	for length := 2; length <= n; length++ {
		for i := 0; i+length <= n; i++ {
			j := i + length - 1
			if s[i] == s[j] {
				if length == 2 {
					dp[i][j] = 2
				} else {
					dp[i][j] = dp[i+1][j-1] + 2
				}
			} else {
				if dp[i+1][j] > dp[i][j-1] {
					dp[i][j] = dp[i+1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}
	return dp[0][n-1]
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		{input: "4 2\n", n: 4, k: 2},
		{input: "6 3\n", n: 6, k: 3},
		{input: "6 4\n", n: 6, k: 4},
		{input: "2 1\n", n: 2, k: 1},
		{input: "2 2\n", n: 2, k: 2},
	}
	for i := 0; i < 100; i++ {
		n := (rand.Intn(10) + 1) * 2
		k := rand.Intn(n) + 1
		tests = append(tests, testCase{
			input: fmt.Sprintf("%d %d\n", n, k),
			n:     n,
			k:     k,
		})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
