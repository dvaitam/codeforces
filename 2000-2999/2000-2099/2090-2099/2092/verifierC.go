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
	"time"
)

const (
	refSource    = "2000-2999/2000-2099/2090-2099/2092/2092C.go"
	maxTests     = 400
	totalNLimit  = 30000
	randomTrials = 500
)

type testCase struct {
	n int
	a []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
		os.Exit(1)
	}

	for i := range tests {
		if candAns[i] != refAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d, got %d\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2092C-ref-*")
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse integer at position %d (%q): %v", i+1, f, err)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sumN := 0

	add := func(tc testCase) {
		if len(tests) >= maxTests || sumN+tc.n > totalNLimit {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	// Deterministic edge cases.
	add(testCase{n: 1, a: []int64{1}})
	add(testCase{n: 1, a: []int64{10}})
	add(testCase{n: 2, a: []int64{5, 7}})
	add(testCase{n: 2, a: []int64{4, 8}})
	add(testCase{n: 3, a: []int64{1, 2, 3}})
	add(testCase{n: 4, a: []int64{1, 2, 2, 1}})

	for attempts := 0; attempts < randomTrials && len(tests) < maxTests && sumN < totalNLimit; attempts++ {
		n := rng.Intn(400) + 1
		if sumN+n > totalNLimit {
			n = totalNLimit - sumN
		}

		a := make([]int64, n)
		for i := 0; i < n; i++ {
			switch rng.Intn(6) {
			case 0:
				a[i] = 1
			case 1:
				a[i] = int64(rng.Intn(10)+1) * int64(rng.Intn(10)+1)
			case 2:
				a[i] = int64(rng.Intn(1_000_000_000) + 1)
			case 3:
				a[i] = int64(rng.Intn(1000)+1) * int64(rng.Intn(1000)+1)
			case 4:
				a[i] = int64(rng.Intn(1_000_000_000)+1) | 1 // force odd
			default:
				a[i] = int64(rng.Intn(1_000_000_000)+1) & ^int64(1) // force even
			}
			if a[i] <= 0 {
				a[i] = 1
			}
		}

		add(testCase{n: n, a: a})
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{n: 1, a: []int64{1}})
	}
	return tests
}
