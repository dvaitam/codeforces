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
	n int
	h [][]int
}

func callerFile() (string, bool) {
	_, file, _, ok := runtime.Caller(0)
	return file, ok
}

func buildOracle() (string, func(), error) {
	_, file, ok := callerFile()
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2068G-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleG")
	cmd := exec.Command("go", "build", "-o", outPath, "2068G.go")
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
	sb.Grow(tc.n*tc.n*6 + 32)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.h[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseAnswer(out string) (float64, error) {
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
		return fmt.Errorf("expected %.10e, got %.10e (diff %.3e)", expected, actual, diff)
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2,
			h: [][]int{{3, 3}, {3, 3}},
		},
		{
			n: 3,
			h: [][]int{
				{0, 0, 0},
				{0, 1545, 0},
				{0, 0, 0},
			},
		},
		{
			n: 4,
			h: [][]int{
				{0, 1, 2, 3},
				{5, 6, 7, 4},
				{10, 11, 8, 9},
				{15, 12, 13, 14},
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		n := rng.Intn(19) + 2 // [2,20]
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, n)
			for j := 0; j < n; j++ {
				val := rng.Intn(1546)
				if rng.Intn(10) == 0 {
					val = rng.Intn(10)
				}
				grid[i][j] = val
			}
		}
		tests = append(tests, testCase{n: n, h: grid})
	}
	return tests
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
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expAns, err := parseAnswer(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s", idx+1, err, expOut)
			os.Exit(1)
		}
		gotAns, err := parseAnswer(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\n%s", idx+1, err, gotOut)
			os.Exit(1)
		}
		if err := compareAnswers(expAns, gotAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %d mismatch: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
