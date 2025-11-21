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
	str   string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
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
		if err := checkOutput(strings.TrimSpace(out), tc.str); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest("abba"),
		makeTest("bbbb"),
		makeTest("aba"),
	}
	for i := 0; i < 200; i++ {
		length := rand.Intn(20) + 3
		var sb strings.Builder
		for j := 0; j < length; j++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('a')
			} else {
				sb.WriteByte('b')
			}
		}
		tests = append(tests, makeTest(sb.String()))
	}
	return tests
}

func makeTest(s string) testCase {
	return testCase{
		input: fmt.Sprintf("1\n%s\n", s),
		str:   s,
	}
}

func checkOutput(out, s string) error {
	if out == ":(" {
		if hasSolution(s) {
			return fmt.Errorf("solution exists but got ':('")
		}
		return nil
	}
	parts := strings.Fields(out)
	if len(parts) != 3 {
		return fmt.Errorf("expected three substrings or ':(', got %d parts", len(parts))
	}
	a, b, c := parts[0], parts[1], parts[2]
	if len(a)+len(b)+len(c) != len(s) || a+b+c != s {
		return fmt.Errorf("concatenation mismatch")
	}
	if !propertyHolds(a, b, c) {
		return fmt.Errorf("lexicographic property violated")
	}
	return nil
}

func hasSolution(s string) bool {
	n := len(s)
	for i := 1; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			a := s[:i]
			b := s[i:j]
			c := s[j:]
			if propertyHolds(a, b, c) {
				return true
			}
		}
	}
	return false
}

func propertyHolds(a, b, c string) bool {
	if (cmp(b, a) >= 0 && cmp(b, c) >= 0) || (cmp(b, a) <= 0 && cmp(b, c) <= 0) {
		return true
	}
	return false
}

func cmp(x, y string) int {
	if x == y {
		return 0
	}
	if x < y {
		return -1
	}
	return 1
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
