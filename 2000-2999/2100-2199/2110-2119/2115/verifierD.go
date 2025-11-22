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
	refSource         = "2115D.go"
	totalNLimit       = 100000
	maxVal            = uint64(1) << 60
	deterministicSeed = 2115
)

type testCase struct {
	a []uint64
	b []uint64
	c string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\n", i+1, expect[i], got[i])
			fmt.Fprintf(os.Stderr, "n=%d c=%s\n", len(tc.a), summarizeString(tc.c))
			fmt.Fprintf(os.Stderr, "a=%s\n", summarizeSlice(tc.a))
			fmt.Fprintf(os.Stderr, "b=%s\n", summarizeSlice(tc.b))
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
	tmp, err := os.CreateTemp("", "2115D-ref-*")
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

func parseOutputs(out string, expected int) ([]uint64, error) {
	tokens := strings.Fields(out)
	if len(tokens) < expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	if len(tokens) > expected {
		return nil, fmt.Errorf("extra output starting at token %q", tokens[expected])
	}
	res := make([]uint64, expected)
	for i := 0; i < expected; i++ {
		val, err := strconv.ParseUint(tokens[i], 10, 64)
		if err != nil {
			// allow signed formatting if non-negative
			signed, err2 := strconv.ParseInt(tokens[i], 10, 64)
			if err2 != nil || signed < 0 {
				return nil, fmt.Errorf("token %q is not a non-negative integer", tokens[i])
			}
			val = uint64(signed)
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		n := len(tc.a)
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		writeSlice(&sb, tc.a)
		writeSlice(&sb, tc.b)
		sb.WriteString(tc.c)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func writeSlice(sb *strings.Builder, arr []uint64) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatUint(v, 10))
	}
	sb.WriteByte('\n')
}

func buildTests() []testCase {
	tests := make([]testCase, 0)
	totalN := 0
	add := func(tc testCase) {
		if totalN+len(tc.a) > totalNLimit {
			return
		}
		tests = append(tests, tc)
		totalN += len(tc.a)
	}

	for _, tc := range deterministicTests() {
		add(tc)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Include a large stress near the limit.
	remaining := totalNLimit - totalN
	if remaining > 50000 {
		n := 50000
		add(randomCase(rand.New(rand.NewSource(deterministicSeed)), n, allOnes))
		remaining = totalNLimit - totalN
	}

	for remaining := totalNLimit - totalN; remaining > 0; remaining = totalNLimit - totalN {
		n := 1 + rng.Intn(500)
		if n > remaining {
			n = remaining
		}
		mode := rng.Intn(4)
		add(randomCase(rng, n, mode))
	}

	return tests
}

const (
	randomMixed = iota
	allZeros
	allOnes
	alt01
)

func deterministicTests() []testCase {
	return []testCase{
		// Single round, each player active.
		{a: []uint64{5}, b: []uint64{7}, c: "0"},
		{a: []uint64{10}, b: []uint64{3}, c: "1"},

		// Equal options.
		{a: []uint64{42, 13}, b: []uint64{42, 99}, c: "01"},

		// Only one player acts.
		{a: []uint64{1, 2, 3, 4}, b: []uint64{8, 9, 10, 11}, c: "0000"},
		{a: []uint64{1, 2, 3, 4}, b: []uint64{8, 9, 10, 11}, c: "1111"},

		// Alternating control.
		{a: []uint64{0, 1, 2, 3, 4, 5}, b: []uint64{6, 7, 8, 9, 10, 11}, c: "010101"},

		// High bits near 2^60.
		{a: []uint64{maxVal - 1, maxVal - 2, 1}, b: []uint64{maxVal - 3, 5, maxVal - 4}, c: "101"},
	}
}

func randomCase(rng *rand.Rand, n int, mode int) testCase {
	a := make([]uint64, n)
	b := make([]uint64, n)
	cBytes := make([]byte, n)
	for i := 0; i < n; i++ {
		a[i] = randValue(rng)
		b[i] = randValue(rng)
		switch mode {
		case allZeros:
			cBytes[i] = '0'
		case allOnes:
			cBytes[i] = '1'
		case alt01:
			if i%2 == 0 {
				cBytes[i] = '0'
			} else {
				cBytes[i] = '1'
			}
		default:
			if rng.Intn(2) == 0 {
				cBytes[i] = '0'
			} else {
				cBytes[i] = '1'
			}
		}
	}
	return testCase{a: a, b: b, c: string(cBytes)}
}

func randValue(rng *rand.Rand) uint64 {
	switch rng.Intn(3) {
	case 0:
		return uint64(rng.Int63n(256)) // small numbers
	case 1:
		return uint64(rng.Int63()) % maxVal
	default:
		// emphasize higher bits
		return maxVal - uint64(rng.Intn(1<<16))
	}
}

func summarizeString(s string) string {
	if len(s) <= 80 {
		return s
	}
	return s[:80] + "... (truncated)"
}

func summarizeSlice(arr []uint64) string {
	const limit = 10
	if len(arr) <= limit {
		return fmt.Sprint(arr)
	}
	head := fmt.Sprint(arr[:limit])
	return head[:len(head)-1] + fmt.Sprintf(" ... total=%d]", len(arr))
}
