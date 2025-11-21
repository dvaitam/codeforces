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

type testCase struct {
	n    int
	m    int
	rows []string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1773G-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleG")
	cmd := exec.Command("go", "build", "-o", outPath, "1773G.go")
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
	sb.Grow(tc.n*(tc.m+2) + 32)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(tc.m))
	sb.WriteByte('\n')
	for _, row := range tc.rows {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string) (float64, error) {
	lines := strings.Fields(out)
	if len(lines) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q: %v", lines[0], err)
	}
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, fmt.Errorf("non-finite value: %v", val)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1, m: 5,
			rows: []string{
				"11010",
			},
		},
		{
			n: 3, m: 3,
			rows: []string{
				"011",
				"101",
				"110",
			},
		},
		{
			n: 6, m: 4,
			rows: []string{
				"1011",
				"0110",
				"1111",
				"0110",
				"0000",
				"1101",
			},
		},
		{
			n: 2, m: 2,
			rows: []string{
				"00",
				"11",
			},
		},
		{
			n: 4, m: 2,
			rows: []string{
				"01",
				"01",
				"01",
				"10",
			},
		},
		{
			n: 5, m: 3,
			rows: []string{
				"000",
				"000",
				"000",
				"000",
				"000",
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		m := rng.Intn(5) + 2
		if rng.Intn(5) == 0 {
			m = rng.Intn(16) + 2
		}
		n := rng.Intn(60) + 1
		if rng.Intn(4) == 0 {
			n = rng.Intn(200) + 1
		}
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			var sb strings.Builder
			for j := 0; j < m; j++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			rows[i] = sb.String()
		}
		tests = append(tests, testCase{n: n, m: m, rows: rows})
	}
	return tests
}

func compareFloats(exp, act float64) error {
	diff := math.Abs(exp - act)
	tol := 1e-9 * math.Max(1, math.Abs(exp))
	if diff > tol {
		return fmt.Errorf("expected %.15f, got %.15f (diff %.3e > tol %.3e)", exp, act, diff, tol)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
		expectedVal, err := parseOutput(expectedOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\noutput:\n%s", idx+1, err, expectedOut)
			os.Exit(1)
		}
		actualVal, err := parseOutput(actualOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s", idx+1, err, actualOut)
			os.Exit(1)
		}
		if err := compareFloats(expectedVal, actualVal); err != nil {
			fmt.Fprintf(os.Stderr, "test %d mismatch: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
