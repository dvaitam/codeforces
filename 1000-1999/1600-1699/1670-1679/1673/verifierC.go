package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MOD int = 1e9 + 7
const MAXN int = 40000

func generatePalindromes(limit int) []int {
	res := make([]int, 0)
	for i := 1; i <= limit; i++ {
		if isPalindrome(i) {
			res = append(res, i)
		}
	}
	return res
}

func isPalindrome(x int) bool {
	orig := x
	rev := 0
	for x > 0 {
		rev = rev*10 + x%10
		x /= 10
	}
	return orig == rev
}

func precompute() []int {
	pals := generatePalindromes(MAXN)
	dp := make([]int, MAXN+1)
	dp[0] = 1
	for _, v := range pals {
		for i := v; i <= MAXN; i++ {
			dp[i] += dp[i-v]
			if dp[i] >= MOD {
				dp[i] -= MOD
			}
		}
	}
	return dp
}

var dp = precompute()

type testCase struct {
	n int
}

func generateTests() []testCase {
	r := rand.New(rand.NewSource(3))
	tests := []testCase{{1}, {2}, {3}, {4}, {5}}
	for len(tests) < 100 {
		tests = append(tests, testCase{n: r.Intn(MAXN) + 1})
	}
	return tests
}

func expected(n int) string {
	return fmt.Sprintf("%d", dp[n])
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%d\n", t.n)
		want := expected(t.n)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
