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
	"time"
)

type point struct {
	x, y int64
}

type testCase struct {
	points []point
	start  point
	target point
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2002C-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, "2002C.go")
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

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.points)))
		sb.WriteByte('\n')
		for _, p := range tc.points {
			sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
		}
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.start.x, tc.start.y, tc.target.x, tc.target.y))
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			points: []point{{1, 1}},
			start:  point{0, 0},
			target: point{10, 0},
		},
		{
			points: []point{{5, 5}, {7, 7}},
			start:  point{0, 0},
			target: point{3, 4},
		},
		{
			points: []point{{10, 10}, {20, 20}, {30, 40}},
			start:  point{1, 1},
			target: point{1, 2},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 60)
	total := 0
	for len(tests) < 50 && total < 100000 {
		n := rng.Intn(50) + 1
		if total+n > 100000 {
			n = 100000 - total
		}
		points := make([]point, n)
		for i := 0; i < n; i++ {
			points[i] = point{
				x: rng.Int63n(1_000_000_000) + 1,
				y: rng.Int63n(1_000_000_000) + 1,
			}
		}
		start := point{x: rng.Int63n(1_000_000_000) + 1, y: rng.Int63n(1_000_000_000) + 1}
		target := point{x: rng.Int63n(1_000_000_000) + 1, y: rng.Int63n(1_000_000_000) + 1}
		tests = append(tests, testCase{points: points, start: start, target: target})
		total += n
	}
	return tests
}

func stressTests() []testCase {
	n := 100000
	points := make([]point, n)
	for i := 0; i < n; i++ {
		points[i] = point{x: int64(i + 1), y: int64(i + 2)}
	}
	return []testCase{
		{
			points: points,
			start:  point{x: 1, y: 1},
			target: point{x: 1, y: 3},
		},
	}
}

func compareOutputs(expected, actual string) error {
	exp := strings.Fields(strings.TrimSpace(expected))
	act := strings.Fields(strings.TrimSpace(actual))
	if len(exp) != len(act) {
		return fmt.Errorf("expected %d answers, got %d", len(exp), len(act))
	}
	for i := range exp {
		if !strings.EqualFold(exp[i], act[i]) {
			return fmt.Errorf("mismatch at test %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	tests = append(tests, stressTests()...)
	input := buildInput(tests)

	expected, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	actual, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := compareOutputs(expected, actual); err != nil {
		fmt.Fprintf(os.Stderr, "%v\nInput:\n%s\nExpected:\n%s\nActual:\n%s\n", err, input, expected, actual)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
