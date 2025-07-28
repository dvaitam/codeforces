package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	k int
	s string
}

func solve(k int, s string) string {
	base := "FBFFBFFB"
	pattern := strings.Repeat(base, 20)
	if strings.Contains(pattern, s) {
		return "YES"
	}
	return "NO"
}

func genTests() []TestCase {
	tests := make([]TestCase, 0, 100)
	for k := 1; k <= 10 && len(tests) < 100; k++ {
		for mask := 0; mask < (1<<k) && len(tests) < 100; mask++ {
			b := make([]byte, k)
			for i := 0; i < k; i++ {
				if mask&(1<<i) == 0 {
					b[i] = 'F'
				} else {
					b[i] = 'B'
				}
			}
			tests = append(tests, TestCase{k, string(b)})
		}
	}
	return tests
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input strings.Builder
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %s\n", tc.k, tc.s)
	}

	expected := make([]string, len(tests))
	for i, tc := range tests {
		expected[i] = solve(tc.k, tc.s)
	}

	out, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		fmt.Print(out)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(expected) {
		fmt.Printf("wrong number of lines: got %d want %d\n", len(lines), len(expected))
		os.Exit(1)
	}
	for i, got := range lines {
		got = strings.TrimSpace(strings.ToUpper(got))
		if got != expected[i] {
			fmt.Printf("test %d failed: input (%d %s) expected %s got %s\n", i+1, tests[i].k, tests[i].s, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
