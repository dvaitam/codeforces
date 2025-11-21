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
	n int
	s string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2013C-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, "2013C.go")
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
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, s: "0"},
		{n: 1, s: "1"},
		{n: 2, s: "00"},
		{n: 2, s: "11"},
		{n: 3, s: "101"},
		{n: 4, s: "0101"},
		{n: 5, s: "00111"},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	totalN := 0
	for len(tests) < 40 && totalN < 5000 {
		n := rng.Intn(10) + 1
		if totalN+n > 5000 {
			n = 5000 - totalN
		}
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		tests = append(tests, testCase{n: n, s: sb.String()})
		totalN += n
	}
	for len(tests) < 60 {
		n := rng.Intn(100) + 1
		if totalN+n > 10000 {
			n = 10000 - totalN
		}
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		tests = append(tests, testCase{n: n, s: sb.String()})
		totalN += n
		if totalN >= 10000 {
			break
		}
	}
	return tests
}

func stressTests() []testCase {
	return []testCase{
		{n: 100, s: strings.Repeat("0", 100)},
		{n: 100, s: strings.Repeat("1", 100)},
		{n: 100, s: strings.Repeat("01", 50)},
	}
}

func compareOutputs(expected, actual string, count int) error {
	exp := strings.Fields(expected)
	act := strings.Fields(actual)
	if len(exp) != count {
		return fmt.Errorf("oracle produced %d outputs, expected %d", len(exp), count)
	}
	if len(act) != count {
		return fmt.Errorf("expected %d outputs, got %d", count, len(act))
	}
	for i := 0; i < count; i++ {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at case %d: expected %s got %s", i+1, exp[i], act[i])
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
	if err := compareOutputs(expected, actual, len(tests)); err != nil {
		fmt.Fprintf(os.Stderr, "%v\nInput:\n%s\nExpected:\n%s\nActual:\n%s\n", err, input, expected, actual)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
