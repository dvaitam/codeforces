package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type testCase struct {
	n, m int
	grid [][]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n%s", err, refOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s", err, candOut)
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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-1993E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle1993E")
	cmd := exec.Command("go", "build", "-o", outPath, "1993E.go")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
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
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("unexpected stderr output")
	}
	return out.String(), nil
}

func parseOutputs(output string, t int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[t])
	}
	res := make([]int64, t)
	for i := 0; i < t; i++ {
		val, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not integer", tokens[i])
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(1993))
	budget := 500 // sum of (n^2 + m^2) across tests must not exceed 500
	var tests []testCase

	add := func(tc testCase) {
		cost := tc.n*tc.n + tc.m*tc.m
		if cost > budget {
			return
		}
		budget -= cost
		tests = append(tests, tc)
	}

	add(testCase{n: 1, m: 1, grid: [][]int{{0}}})
	add(testCase{n: 1, m: 2, grid: [][]int{{1, 3}}})
	add(testCase{n: 2, m: 2, grid: [][]int{{5, 5}, {7, 9}}})
	add(testCase{n: 2, m: 3, grid: [][]int{{0, 1, 0}, {4, 7, 2}}})
	add(testCase{n: 3, m: 3, grid: [][]int{{6, 6, 6}, {2, 9, 5}, {3, 0, 12}}})

	for attempts := 0; budget >= 2 && attempts < 2000; attempts++ {
		n := rng.Intn(15) + 1
		m := rng.Intn(15) + 1
		cost := n*n + m*m
		if cost > budget {
			continue
		}
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			row := make([]int, m)
			for j := 0; j < m; j++ {
				row[j] = rng.Intn(1 << 20)
			}
			grid[i] = row
		}
		add(testCase{n: n, m: m, grid: grid})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.m; j++ {
				if j > 0 {
					b.WriteByte(' ')
				}
				fmt.Fprintf(&b, "%d", tc.grid[i][j])
			}
			b.WriteByte('\n')
		}
	}
	return b.String()
}
