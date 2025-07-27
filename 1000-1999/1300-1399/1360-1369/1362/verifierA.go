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
	a, b     int64
	expected int
}

func computeExpected(a, b int64) int {
	if a == b {
		return 0
	}
	var bigger, smaller int64
	if a > b {
		bigger = a
		smaller = b
	} else {
		bigger = b
		smaller = a
	}
	if bigger%smaller != 0 {
		return -1
	}
	ratio := bigger / smaller
	k := 0
	for ratio > 1 && ratio%2 == 0 {
		ratio /= 2
		k++
	}
	if ratio != 1 {
		return -1
	}
	return (k + 2) / 3
}

func generateTests() []testCase {
	const numTests = 100
	rand.Seed(1)
	tests := make([]testCase, 0, numTests+5)
	for i := 0; i < numTests; i++ {
		a := rand.Int63n(1e18) + 1
		b := rand.Int63n(1e18) + 1
		tests = append(tests, testCase{a: a, b: b, expected: computeExpected(a, b)})
	}
	edge := []testCase{
		{1, 1, computeExpected(1, 1)},
		{1, 2, computeExpected(1, 2)},
		{2, 1, computeExpected(2, 1)},
		{3, 6, computeExpected(3, 6)},
		{96, 3, computeExpected(96, 3)},
	}
	tests = append(tests, edge...)
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	var input strings.Builder
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.a, tc.b)
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
		got, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d (%d %d)\n", i+1, tc.expected, got, tc.a, tc.b)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
