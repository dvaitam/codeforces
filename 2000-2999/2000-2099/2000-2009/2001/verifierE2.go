package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "2000-2999/2000-2099/2000-2009/2001/2001E2.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2001E2.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	name     string
	n, k     int
	modPrime int64
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.k, tc.modPrime))
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "sample1", n: 2, k: 1, modPrime: 998244353},
		{name: "sample2", n: 3, k: 2, modPrime: 998244853},
		{name: "sample3", n: 3, k: 3, modPrime: 998244353},
		{name: "sample4", n: 3, k: 4, modPrime: 100000037},
		{name: "sample5", n: 4, k: 2, modPrime: 100000039},
		{name: "sample6", n: 4, k: 3, modPrime: 100000037},
		{name: "edge_zero_k", n: 5, k: 0, modPrime: 999999937},
		{name: "max_n_small_k", n: 10, k: 5, modPrime: 998244353},
	}
}

var primesPool = []int64{
	998244353,
	998244853,
	100000007,
	100000021,
	100000037,
	100000039,
	999999937,
	999999929,
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	totalK := 0
	for len(tests) < 40 {
		nRemaining := 100 - totalN
		kRemaining := 500 - totalK
		if nRemaining <= 0 || kRemaining <= 0 {
			break
		}
		n := rng.Intn(min(10, nRemaining)) + 2 // keep n moderate
		if n > nRemaining {
			n = nRemaining
		}
		k := rng.Intn(min(15, kRemaining)) + 1
		if k > kRemaining {
			k = kRemaining
		}
		mod := primesPool[rng.Intn(len(primesPool))]
		tests = append(tests, testCase{
			name:     fmt.Sprintf("random_%d", len(tests)+1),
			n:        n,
			k:        k,
			modPrime: mod,
		})
		totalN += n
		totalK += k
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2001E2-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2001E2")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseOutputs(out string, expected int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) < expected {
		return nil, fmt.Errorf("expected %d answers, got %d tokens", expected, len(tokens))
	}
	results := make([]int64, expected)
	for i := 0; i < expected; i++ {
		var val int64
		if _, err := fmt.Sscan(tokens[i], &val); err != nil {
			return nil, fmt.Errorf("failed to parse integer %q: %v", tokens[i], err)
		}
		results[i] = val
	}
	if len(tokens) > expected {
		return nil, fmt.Errorf("extra tokens after reading %d answers", expected)
	}
	return results, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatInput(tests)

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseOutputs(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %s (%d) mismatch: expected %d, got %d\ninput case: n=%d k=%d mod=%d\n", tests[i].name, i+1, expected[i], got[i], tests[i].n, tests[i].k, tests[i].modPrime)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
