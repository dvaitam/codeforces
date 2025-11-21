package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input string
	nums  []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierM.go /path/to/binary")
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
		if err := check(strings.TrimSpace(out), tc.nums); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(out string, nums []int64) error {
	lines := strings.Fields(out)
	if len(lines) != len(nums)*2 {
		return fmt.Errorf("expected %d numbers but got %d", len(nums)*2, len(lines))
	}
	for i, n := range nums {
		a, err := strconv.ParseInt(lines[2*i], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer %q", lines[2*i])
		}
		b, err := strconv.ParseInt(lines[2*i+1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer %q", lines[2*i+1])
		}
		if a <= 0 || b <= 0 {
			return fmt.Errorf("numbers must be positive")
		}
		if a+b != n {
			return fmt.Errorf("pair (%d,%d) doesn't sum to %d", a, b, n)
		}
		if !isOptimal(n, a, b) {
			return fmt.Errorf("pair (%d,%d) is not optimal for %d", a, b, n)
		}
	}
	return nil
}

func isOptimal(n, a, b int64) bool {
	best := int64(1<<63 - 1)
	var besta, bestb int64
	for k := int64(1); k < n; k++ {
		g := gcd(k, n-k)
		l := k / g * (n - k)
		if l < best {
			best = l
			besta, bestb = k, n-k
		}
	}
	lcm := a / gcd(a, b) * b
	return lcm == best && ((a == besta && b == bestb) || (a == bestb && b == besta) || checkWithDivisor(n, a))
}

func checkWithDivisor(n, a int64) bool {
	d := gcd(a, n)
	return d > 1 && a == n/d
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest([]int64{2}),
		makeTest([]int64{3, 5, 10}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Int63n(1000) + 2
		tests = append(tests, makeTest([]int64{n}))
	}
	return tests
}

func makeTest(nums []int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(nums))
	for _, n := range nums {
		fmt.Fprintf(&sb, "%d\n", n)
	}
	return testCase{
		input: sb.String(),
		nums:  nums,
	}
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
