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

func solveCase(n, m int64) (int64, int64) {
	minIso := n - 2*m
	if minIso < 0 {
		minIso = 0
	}
	low, high := int64(0), n+1
	for low < high {
		mid := (low + high) / 2
		if mid*(mid-1)/2 >= m {
			high = mid
		} else {
			low = mid + 1
		}
	}
	maxIso := n - low
	return minIso, maxIso
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(43))
	tests := make([]testCase, 0, 100)
	tests = append(tests, testCase{input: "1 0\n", expected: "1 1"})
	for len(tests) < 100 {
		n := rng.Int63n(100000) + 1
		maxM := n * (n - 1) / 2
		m := rng.Int63n(maxM + 1)
		in := fmt.Sprintf("%d %d\n", n, m)
		mi, ma := solveCase(n, m)
		exp := fmt.Sprintf("%d %d", mi, ma)
		tests = append(tests, testCase{input: in, expected: exp})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
