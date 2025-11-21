package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
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
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	exp, _ := strconv.ParseInt(expect, 10, 64)
	val, err := strconv.ParseInt(actual, 10, 64)
	if err != nil {
		return fmt.Errorf("output is not integer: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase(5, 2, []int64{3, 1, 2, 1, 1}),
		makeCase(5, 3, []int64{2, 3, 1, 2, 4}),
		makeCase(3, 1, []int64{5, 4, 3}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 2
		k := rand.Intn(n-1) + 1
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = rand.Int63n(100) + 1
		}
		tests = append(tests, makeCase(n, k, a))
	}
	return tests
}

func makeCase(n, k int, a []int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", solveReference(n, k, a)),
	}
}

func solveReference(n, k int, a []int64) int64 {
	ans := int64(1<<63 - 1)
	for s := 0; s < n; s++ {
		cost := solveWithStart(n, k, s, a)
		if cost < ans {
			ans = cost
		}
	}
	return ans
}

func solveWithStart(n, k, start int, a []int64) int64 {
	b := append([]int64(nil), a[start:]...)
	b = append(b, a[:start]...)
	dp := make([]int64, n+1)
	for i := range dp {
		dp[i] = 1<<63 - 1
	}
	dp[0] = 0
	deque := []int{0}
	for i := 1; i <= n; i++ {
		for len(deque) > 0 && deque[0] < i-k {
			deque = deque[1:]
		}
		best := deque[0]
		dp[i] = dp[best] + b[i-1]
		for len(deque) > 0 && dp[deque[len(deque)-1]] >= dp[i] {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i)
	}
	res := dp[n]
	last := b[n-1]
	for i := max(0, n-k); i < n; i++ {
		if dp[i]+last < res {
			res = dp[i] + last
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
