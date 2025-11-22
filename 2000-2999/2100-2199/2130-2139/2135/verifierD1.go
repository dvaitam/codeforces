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

const (
	refSource = "2135D1.go"
	maxW      = 100000
)

type testCase struct {
	w int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s\n", err, refOut)
		os.Exit(1)
	}
	expect, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s\n", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expect[i] != got[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d (W=%d)\n", i+1, expect[i], got[i], tests[i].w)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2135D1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Join(dir, refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutputs(out string, expected int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) < expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	if len(tokens) > expected {
		return nil, fmt.Errorf("extra output starting at token %q", tokens[expected])
	}
	res := make([]int, expected)
	for i := 0; i < expected; i++ {
		val, err := strconv.Atoi(tokens[i])
		if err != nil {
			return nil, fmt.Errorf("token %q is not an integer", tokens[i])
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte(' ')
	sb.WriteString("manual")
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.w))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0)
	add := func(w int) {
		if w < 1 {
			w = 1
		}
		if w > maxW {
			w = maxW
		}
		tests = append(tests, testCase{w: w})
	}

	// Deterministic edge cases.
	add(1)
	add(maxW)
	add(50000)
	add(99999)

	// Randomized cases.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 10 {
		add(1 + rng.Intn(maxW))
	}
	return tests
}
