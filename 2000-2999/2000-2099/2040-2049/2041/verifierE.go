package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const (
	maxLen = 1000
	maxAbs = 1_000_000
)

type testCase struct {
	a, b int
	desc string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()

	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.a, tc.b)

		if err := validateProgram(oracle, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d (%s): %v\n", idx+1, tc.desc, err)
			os.Exit(1)
		}

		if err := validateProgram(candidate, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.desc, err, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2041E-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", bin, "2041E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return bin, cleanup, nil
}

func validateProgram(path, input string, tc testCase) error {
	out, stderr, err := runBinary(path, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr)
	}
	return validateOutput(out, tc)
}

func runBinary(path, input string) (string, string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func validateOutput(out string, tc testCase) error {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	n64, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid length %q", tokens[0])
	}
	if n64 < 1 || n64 > maxLen {
		return fmt.Errorf("length %d out of bounds", n64)
	}
	n := int(n64)
	if len(tokens) != n+1 {
		return fmt.Errorf("expected %d numbers after length, got %d", n, len(tokens)-1)
	}
	values := make([]int, n)
	sum := int64(0)
	for i := 0; i < n; i++ {
		val64, err := strconv.ParseInt(tokens[i+1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid array element %q", tokens[i+1])
		}
		if val64 < -maxAbs || val64 > maxAbs {
			return fmt.Errorf("array element %d exceeds bounds", val64)
		}
		values[i] = int(val64)
		sum += val64
	}
	if sum != int64(n)*int64(tc.a) {
		return fmt.Errorf("mean mismatch: sum=%d length=%d expected mean=%d", sum, n, tc.a)
	}
	sorted := append([]int(nil), values...)
	sort.Ints(sorted)
	if n%2 == 1 {
		median := sorted[n/2]
		if median != tc.b {
			return fmt.Errorf("median mismatch: got %d expected %d", median, tc.b)
		}
	} else {
		left := sorted[n/2-1]
		right := sorted[n/2]
		if left+right != 2*tc.b {
			return fmt.Errorf("median mismatch: (%d+%d)/2 != %d", left, right, tc.b)
		}
	}
	return nil
}

func generateTests() []testCase {
	tests := []testCase{
		{a: 3, b: 4, desc: "sample-positive"},
		{a: -100, b: -100, desc: "sample-negative"},
		{a: 0, b: 0, desc: "zero-zero"},
		{a: 100, b: -100, desc: "extreme-opposite"},
		{a: -50, b: 50, desc: "extreme-opposite2"},
		{a: 42, b: 7, desc: "mixed"},
		{a: -13, b: -7, desc: "negative"},
	}
	rng := rand.New(rand.NewSource(2041))
	for i := 0; i < 40; i++ {
		a := rng.Intn(201) - 100
		b := rng.Intn(201) - 100
		tests = append(tests, testCase{
			a:    a,
			b:    b,
			desc: fmt.Sprintf("random-%d", i+1),
		})
	}
	return tests
}
