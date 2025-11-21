package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	w int64
	b int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refPath, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: w=%d b=%d expected=%d got=%d\nInput used:\n%s\n", i+1, tests[i].w, tests[i].b, expected[i], got[i], input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2041B-ref-*")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", path, "2041B.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return path, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(w, b int64) {
		tests = append(tests, testCase{w: w, b: b})
	}
	add(0, 0)
	add(1, 0)
	add(0, 1)
	add(6, 4)
	add(12, 0)
	add(1_000_000_000, 1_000_000_000)
	add(1_000_000_000, 0)
	add(0, 1_000_000_000)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 100 {
		w := rng.Int63n(1_000_000_001)
		b := rng.Int63n(1_000_000_001)
		add(w, b)
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.w, tc.b))
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok, i+1)
		}
		if val < 0 {
			return nil, fmt.Errorf("negative side length %d at position %d", val, i+1)
		}
		res[i] = val
	}
	return res, nil
}
