package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./1998E1.go"

type testCase struct {
	n    int
	vals []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := renderInput(tests)

	refOut, err := runWithInput(exec.Command(refBin), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	candCmd := commandFor(candidate)
	candOut, err := runWithInput(candCmd, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	expected, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}
	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\ninput:\n%s", i+1, expected[i], got[i], formatSingleInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1998E1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutputs(output string, expectedCount int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) < expectedCount {
		return nil, fmt.Errorf("expected %d numbers, got %d", expectedCount, len(fields))
	}
	if len(fields) > expectedCount {
		return nil, fmt.Errorf("extra output detected after %d numbers", expectedCount)
	}
	res := make([]int64, expectedCount)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number %d (%q): %v", i+1, f, err)
		}
		res[i] = val
	}
	return res, nil
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.n)) // x = n in easy version
		for i, v := range tc.vals {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func formatSingleInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.n))
	for i, v := range tc.vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildTests() []testCase {
	tests := []testCase{
		{n: 1, vals: []int64{1}},
		{n: 1, vals: []int64{1_000_000_000}},
		{n: 2, vals: []int64{1, 1}},
		{n: 2, vals: []int64{5, 3}},
		{n: 3, vals: []int64{1, 2, 3}},
		{n: 3, vals: []int64{3, 2, 1}},
		{n: 4, vals: []int64{2, 2, 2, 2}},
		{n: 5, vals: []int64{1, 2, 3, 2, 1}},
	}

	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	rng := rand.New(rand.NewSource(1998))
	const maxN = 200_000
	for totalN < maxN {
		remaining := maxN - totalN
		// Mix small random tests with occasional larger ones.
		n := rng.Intn(400) + 1
		if remaining < n {
			n = remaining
		}
		if totalN < 20_000 && rng.Intn(5) == 0 {
			n = minInt(remaining, rng.Intn(3000)+100)
		}
		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			vals[i] = rng.Int63n(1_000_000_000) + 1
		}
		tests = append(tests, testCase{n: n, vals: vals})
		totalN += n
		if len(tests) > 1000 {
			break
		}
	}

	// Ensure at least one large test hits the upper constraint if space remains.
	if totalN < maxN {
		n := maxN - totalN
		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			vals[i] = int64((i % 10) + 1)
		}
		tests = append(tests, testCase{n: n, vals: vals})
	}

	return tests
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
