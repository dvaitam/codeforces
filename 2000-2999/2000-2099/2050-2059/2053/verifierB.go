package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	refSource2053B = "2053B.go"
	refBinary2053B = "ref2053B.bin"
	maxTests       = 160
	maxTotalN      = 200000
)

type testCase struct {
	intervals [][2]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on test %d:\nexpected: %s\ngot: %s\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2053B, refSource2053B)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2053B), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d lines, got %d", t, len(lines))
	}
	return lines, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", len(tc.intervals))
		for _, interval := range tc.intervals {
			fmt.Fprintf(&sb, "%d %d\n", interval[0], interval[1])
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2053))
	var tests []testCase
	totalN := 0

	add := func(intervals [][2]int) {
		tests = append(tests, testCase{intervals: intervals})
		totalN += len(intervals)
	}

	add([][2]int{{1, 1}})
	add([][2]int{{1, 2}, {2, 4}})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		maxN := 4000
		if remain < maxN {
			maxN = remain
		}
		n := rnd.Intn(maxN) + 1
		intervals := make([][2]int, n)
		for i := 0; i < n; i++ {
			l := rnd.Intn(2*n) + 1
			r := l + rnd.Intn(2*n-l+1)
			intervals[i] = [2]int{l, r}
		}
		add(intervals)
	}
	return tests
}
