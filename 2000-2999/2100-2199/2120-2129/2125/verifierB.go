package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource2125B = "./2125B.go"

type testCase struct {
	a int64
	b int64
	k int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	input := formatInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d, got %d\nInput:\n%s", i+1, expected[i], got[i], stringifyCase(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2125B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2125B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2125B)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.Write(errBuf.Bytes())
	}
	return out.String(), err
}

func parseOutput(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int, t)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.a, tc.b, tc.k)
	}
	return []byte(sb.String())
}

func stringifyCase(tc testCase) string {
	return fmt.Sprintf("%d %d %d\n", tc.a, tc.b, tc.k)
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 120)

	// Deterministic coverage.
	tests = append(tests,
		testCase{a: 1, b: 1, k: 1},                       // immediate reach with cost 1
		testCase{a: 1, b: 2, k: 1},                       // needs two patterns
		testCase{a: 5, b: 5, k: 10},                      // both within k after gcd
		testCase{a: 5, b: 5, k: 4},                       // require two distinct moves
		testCase{a: 10, b: 15, k: 5},                     // gcd 5, both <= k
		testCase{a: 10, b: 14, k: 3},                     // gcd 2, 5>k -> cost 2
		testCase{a: 1_000_000_000_000000000, b: 1, k: 1}, // huge a, needs 2
		testCase{a: 12, b: 18, k: 6},                     // gcd 6, within k -> 1
		testCase{a: 12, b: 18, k: 5},                     // gcd 6, 3<=5,2<=5? wait a/g=2 b/g=3 <=5 ->1
		testCase{a: 9, b: 7, k: 5},                       // gcd1, both<=5 ->1? 9>5 no ->2
	)

	// Randomized tests.
	rng := rand.New(rand.NewSource(2125_2024))
	for len(tests) < 100 {
		a := rng.Int63n(1_000_000_000_000000000) + 1 // up to 1e18
		b := rng.Int63n(1_000_000_000_000000000) + 1
		k := rng.Int63n(1_000_000_000_000000000) + 1
		tests = append(tests, testCase{a: a, b: b, k: k})
	}

	return tests
}
