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
	V1 int
	V2 int
	V3 int
	Vm int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-907A-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", outPath, "907A.go")
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
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.V1, tc.V2, tc.V3, tc.Vm))
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{V1: 50, V2: 30, V3: 10, Vm: 10},
		{V1: 100, V2: 50, V3: 10, Vm: 20},
		{V1: 10, V2: 9, V3: 8, Vm: 7},
		{V1: 9, V2: 6, V3: 5, Vm: 4},
		{V1: 7, V2: 4, V3: 2, Vm: 1},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	for i := 0; i < 100; i++ {
		v1 := rng.Intn(100) + 3
		v2 := rng.Intn(v1-2) + 2
		v3 := rng.Intn(v2-1) + 1
		vm := rng.Intn(100) + 1
		tests = append(tests, testCase{V1: v1, V2: v2, V3: v3, Vm: vm})
	}
	return tests
}

func parseOutput(output string) (int, int, int, error) {
	output = strings.TrimSpace(output)
	if output == "-1" {
		return -1, -1, -1, nil
	}
	parts := strings.Fields(output)
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("expected 3 values, got %v", output)
	}
	l, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, err
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, err
	}
	s, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, err
	}
	return l, m, s, nil
}

func validate(tc testCase, l, m, s int) error {
	if tc.V1 <= tc.V2 || tc.V2 <= tc.V3 {
		return fmt.Errorf("invalid input V1>%d V2>%d V3>%d", tc.V1, tc.V2, tc.V3)
	}
	// check existence of solution
	minS := tc.Vm
	if minS < tc.V3 {
		minS = tc.V3
	}
	maxS := min(2*tc.V3, 2*tc.Vm)
	if minS > maxS {
		if l != -1 {
			return fmt.Errorf("expected -1 but got %d %d %d", l, m, s)
		}
		return nil
	}
	if l == -1 {
		return fmt.Errorf("solution exists but got -1")
	}
	if !(l > m && m > s) {
		return fmt.Errorf("l>%d m>%d s>%d not strictly decreasing", l, m, s)
	}
	if !(tc.V1 <= l && l <= 2*tc.V1) {
		return fmt.Errorf("l not suitable for father: %d", l)
	}
	if !(tc.V2 <= m && m <= 2*tc.V2) {
		return fmt.Errorf("m not suitable for mother: %d", m)
	}
	if !(tc.V3 <= s && s <= 2*tc.V3) {
		return fmt.Errorf("s not suitable for son: %d", s)
	}
	if !(tc.Vm <= s && s <= 2*tc.Vm) {
		return fmt.Errorf("s not suitable for Masha: %d", s)
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	expected, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}

	outputs := strings.Split(strings.TrimSpace(expected), "\n")
	if len(outputs) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d lines, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}

	for i, line := range outputs {
		l, m, s, err := parseOutput(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validate(tests[i], l, m, s); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput: %d %d %d %d\noutput: %s\n", i+1, err, tests[i].V1, tests[i].V2, tests[i].V3, tests[i].Vm, line)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
