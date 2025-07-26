package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	output string
}

func solve(n int, s string) string {
	for i := 0; i+1 < n; i++ {
		if s[i] == '1' && s[i+1] == '1' {
			return "No"
		}
	}
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			leftEmpty := i == 0 || s[i-1] == '0'
			rightEmpty := i == n-1 || s[i+1] == '0'
			if leftEmpty && rightEmpty {
				return "No"
			}
		}
	}
	return "Yes"
}

func generateTests() []testCase {
	rand.Seed(1)
	var tests []testCase
	// some fixed edge cases
	fixed := []string{"0", "1", "00", "10", "01", "11", "010", "101", "000", "111"}
	for _, f := range fixed {
		n := len(f)
		tests = append(tests, testCase{
			input:  fmt.Sprintf("%d\n%s\n", n, f),
			output: solve(n, f),
		})
	}
	// random cases
	for len(tests) < 120 { // at least 100
		n := rand.Intn(10) + 1
		b := make([]byte, n)
		for i := range b {
			if rand.Intn(2) == 0 {
				b[i] = '0'
			} else {
				b[i] = '1'
			}
		}
		s := string(b)
		tests = append(tests, testCase{
			input:  fmt.Sprintf("%d\n%s\n", n, s),
			output: solve(n, s),
		})
	}
	return tests
}

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(binary, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.output {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i+1, tc.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
