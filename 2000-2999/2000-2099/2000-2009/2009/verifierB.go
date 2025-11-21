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
	rows []string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2009B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2009B.go")
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
		sb.WriteString(strconv.Itoa(len(tc.rows)))
		sb.WriteByte('\n')
		for _, row := range tc.rows {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{rows: []string{"#...", "...#", "..#.", ".#.."}},
		{rows: []string{"...#", "..#.", ".#..", "#..."}},
		{rows: []string{".#.."}},
		{rows: []string{"#...", ".#..", "..#.", "...#"}},
		{rows: []string{"...#", "...#", "...#"}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 80)
	totalRows := 0
	for len(tests) < 80 && totalRows < 500 {
		n := rng.Intn(12) + 1
		if totalRows+n > 500 {
			n = 500 - totalRows
		}
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			cols := []byte("....")
			cols[rng.Intn(4)] = '#'
			rows[i] = string(cols)
		}
		tests = append(tests, testCase{rows: rows})
		totalRows += n
	}
	return tests
}

func compareOutputs(expected, actual string) error {
	exp := strings.TrimSpace(expected)
	act := strings.TrimSpace(actual)
	if exp != act {
		return fmt.Errorf("expected %q got %q", exp, act)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
