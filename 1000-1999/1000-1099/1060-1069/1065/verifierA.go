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
	input    string
	expected string
}

func solveCase(s, a, b, c int64) int64 {
	x := s / c
	return x + (x/a)*b
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(42))
	tests := make([]testCase, 0, 100)
	tests = append(tests, testCase{input: "1\n1 1 1 1\n", expected: fmt.Sprint(solveCase(1, 1, 1, 1))})
	tests = append(tests, testCase{input: fmt.Sprintf("1\n%d %d %d %d\n", int64(1e9), int64(1e9), int64(1e9), int64(1)), expected: fmt.Sprint(solveCase(1e9, 1e9, 1e9, 1))})
	for len(tests) < 100 {
		s := rng.Int63n(1e9) + 1
		a := rng.Int63n(1e9) + 1
		b := rng.Int63n(1e9) + 1
		c := rng.Int63n(1e9) + 1
		input := fmt.Sprintf("1\n%d %d %d %d\n", s, a, b, c)
		exp := fmt.Sprint(solveCase(s, a, b, c))
		tests = append(tests, testCase{input: input, expected: exp})
	}
	return tests
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
