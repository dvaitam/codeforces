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
	refSource    = "2000-2999/2000-2099/2060-2069/2061/2061E.go"
	totalNLimit  = 4000
	totalTCLimit = 200
)

type testCase struct {
	n int
	m int
	k int
	a []int
	b []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2061E-ref-*")
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
		fmt.Fprintf(&b, "%d %d %d\n", tc.n, tc.m, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		for i, v := range tc.b {
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
	for i := 0; i < t; i++ {
		val, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse integer %q at position %d: %v", fields[i], i+1, err)
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
		if sumN+tc.n > totalNLimit || len(tests) >= totalTCLimit {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	// Small deterministic cases.
	add(testCase{n: 1, m: 1, k: 0, a: []int{5}, b: []int{7}})
	add(testCase{n: 1, m: 3, k: 2, a: []int{7}, b: []int{7, 3, 1}})
	add(testCase{n: 2, m: 2, k: 2, a: []int{10, 5}, b: []int{8, 7}})

	for attempts := 0; attempts < 500 && sumN < totalNLimit && len(tests) < totalTCLimit; attempts++ {
		n := rng.Intn(120) + 1
		if sumN+n > totalNLimit {
			n = totalNLimit - sumN
		}
		m := rng.Intn(10) + 1
		k := rng.Intn(n*m + 1)

		a := make([]int, n)
		for i := 0; i < n; i++ {
			switch rng.Intn(5) {
			case 0:
				a[i] = 0
			case 1:
				a[i] = (1 << 30) - 1
			default:
				a[i] = rng.Intn(1 << 30)
			}
		}

		b := make([]int, m)
		for i := 0; i < m; i++ {
			switch rng.Intn(6) {
			case 0:
				b[i] = 0
			case 1:
				b[i] = (1 << 30) - 1
			default:
				b[i] = rng.Intn(1 << 30)
			}
		}

		// Occasionally force aggressive k to allow full minimization paths.
		if rng.Intn(5) == 0 {
			k = n * m
		}

		add(testCase{n: n, m: m, k: k, a: a, b: b})
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{n: 1, m: 1, k: 0, a: []int{0}, b: []int{0}})
	}
	return tests
}
