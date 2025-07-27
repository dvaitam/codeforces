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

type testCaseC struct {
	n        int64
	expected int64
}

func computeExpectedC(n int64) int64 {
	var ans int64
	for n > 0 {
		ans += n
		n >>= 1
	}
	return ans
}

func generateTestsC() []testCaseC {
	const numTests = 100
	rand.Seed(3)
	tests := make([]testCaseC, 0, numTests+5)
	for i := 0; i < numTests; i++ {
		n := rand.Int63n(1e18) + 1
		tests = append(tests, testCaseC{n: n, expected: computeExpectedC(n)})
	}
	edge := []testCaseC{
		{1, computeExpectedC(1)},
		{5, computeExpectedC(5)},
		{10, computeExpectedC(10)},
	}
	tests = append(tests, edge...)
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTestsC()
	var input strings.Builder
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d\n", tc.n)
	}

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running binary: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(&out)
	for i, tc := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		got, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, tc.expected, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
