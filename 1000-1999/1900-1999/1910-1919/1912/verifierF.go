package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	u int
	v int
}

type testCase struct {
	n     int
	edges []edge
	s     int
}

func buildOracle() (string, func(), error) {
	dir := filepath.Dir(callerFile())
	tmpDir, err := os.MkdirTemp("", "oracle-1912F-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", outPath, "1912F.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func callerFile() string {
	_, file, _, _ := runtime.Caller(0)
	return file
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(tc.n*16 + len(tc.edges)*16)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(strconv.Itoa(e.u))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e.v))
		sb.WriteByte('\n')
	}
	sb.WriteString(strconv.Itoa(tc.s))
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) (float64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q: %v", fields[0], err)
	}
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, fmt.Errorf("non-finite output %v", val)
	}
	return val, nil
}

func compareAnswers(expected, actual float64) error {
	diff := math.Abs(expected - actual)
	den := math.Max(1, math.Abs(expected))
	if diff/den > 1e-6+1e-12 {
		return fmt.Errorf("expected %.10f, got %.10f (diff %.3e)", expected, actual, diff)
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2,
			edges: []edge{
				{1, 2},
			},
			s: 1,
		},
		{
			n: 3,
			edges: []edge{
				{1, 2},
				{1, 3},
			},
			s: 2,
		},
		{
			n: 4,
			edges: []edge{
				{1, 2},
				{2, 3},
				{2, 4},
			},
			s: 3,
		},
		{
			n: 6,
			edges: []edge{
				{1, 2},
				{1, 3},
				{3, 4},
				{3, 5},
				{5, 6},
			},
			s: 5,
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		n := rng.Intn(98) + 2 // [2,99]
		if rng.Intn(5) == 0 {
			n = rng.Intn(90) + 10
		}
		edges := make([]edge, 0, n-1)
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges = append(edges, edge{p, i})
		}
		s := rng.Intn(n) + 1
		tests = append(tests, testCase{n: n, edges: edges, s: s})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := buildInput(tc)
		expectedOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		actualOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expectedAns, err := parseOutput(expectedOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\noutput:\n%s", idx+1, err, expectedOut)
			os.Exit(1)
		}
		actualAns, err := parseOutput(actualOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s", idx+1, err, actualOut)
			os.Exit(1)
		}
		if err := compareAnswers(expectedAns, actualAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %d mismatch: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
