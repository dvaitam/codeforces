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

type testCase struct {
	x int
	y int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2067A-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", outPath, "2067A.go")
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
	sb.Grow(len(tests) * 32)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.x))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.y))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(lines))
	}
	res := make([]string, expected)
	for i, s := range lines {
		ls := strings.ToLower(s)
		switch ls {
		case "yes", "no":
			res[i] = ls
		default:
			return nil, fmt.Errorf("invalid answer %q at position %d", s, i+1)
		}
	}
	return res, nil
}

func compareAnswers(expected, actual []string) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("answer count mismatch")
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("answer %d mismatch: expected %s, got %s", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	tests := []testCase{
		{x: 1, y: 2},   // sample yes
		{x: 7, y: 7},   // sample no
		{x: 9, y: 9},   // same sums
		{x: 999, y: 1}, // wrap from many digits
		{x: 18, y: 1},  // carry overflow
		{x: 27, y: 27},
		{x: 1, y: 1},
		{x: 25, y: 26},
	}
	// edge of constraints
	tests = append(tests, testCase{x: 1000, y: 1000}, testCase{x: 1000, y: 1})
	return tests
}

func randomTests(total int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, total)
	for i := range tests {
		tests[i] = testCase{
			x: rng.Intn(1000) + 1,
			y: rng.Intn(1000) + 1,
		}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	tests = append(tests, randomTests(400)...)

	input := buildInput(tests)

	expOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	expectedAns, err := parseOutput(expOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expOut)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actOut)
		os.Exit(1)
	}

	if err := compareAnswers(expectedAns, actualAns); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}
