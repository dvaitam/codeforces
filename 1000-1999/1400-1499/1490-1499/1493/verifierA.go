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

type Test struct {
	n int
	k int
}

func generateTests() []Test {
	rand.Seed(1)
	tests := make([]Test, 0, 100)
	// include some edge cases
	edge := []Test{
		{1, 1}, {2, 1}, {2, 2}, {5, 3}, {1000, 1}, {1000, 1000}, {10, 5},
	}
	tests = append(tests, edge...)
	for len(tests) < 100 {
		n := rand.Intn(1000) + 1
		k := rand.Intn(n) + 1
		tests = append(tests, Test{n, k})
	}
	return tests
}

func solve(n, k int) string {
	nums := []int{}
	for i := k + 1; i <= n; i++ {
		nums = append(nums, i)
	}
	for i := k/2 + 1; i < k; i++ {
		nums = append(nums, i)
	}
	var b strings.Builder
	fmt.Fprintln(&b, len(nums))
	for i, v := range nums {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	return b.String()
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	// prepare single run input with T=len(tests)
	var in strings.Builder
	fmt.Fprintln(&in, len(tests))
	for _, t := range tests {
		fmt.Fprintf(&in, "%d %d\n", t.n, t.k)
	}
	expectedParts := make([]string, len(tests))
	for i, t := range tests {
		expectedParts[i] = solve(t.n, t.k)
	}
	expect := strings.Join(expectedParts, "")

	got, err := run(binary, in.String())
	if err != nil {
		fmt.Printf("runtime error: %v\noutput:\n%s", err, got)
		os.Exit(1)
	}
	// normalize newlines
	got = strings.ReplaceAll(strings.TrimSpace(got), "\r\n", "\n")
	expect = strings.ReplaceAll(strings.TrimSpace(expect), "\r\n", "\n")
	if got != expect {
		fmt.Println("wrong answer")
		fmt.Println("input:")
		fmt.Print(in.String())
		fmt.Println("expected:")
		fmt.Print(expect)
		fmt.Println("got:")
		fmt.Print(got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
	time.Sleep(0) // avoid exit handling
}
