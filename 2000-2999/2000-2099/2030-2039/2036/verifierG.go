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
)

type testCase struct {
	n, a, b, c int64
	desc       string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
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
	input := buildInput(tests)

	if err := ensureOracleWorks(oracle, input, len(tests)); err != nil {
		fmt.Fprintf(os.Stderr, "oracle validation failed: %v\n", err)
		os.Exit(1)
	}

	out, stderr, err := runBinary(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\nstderr:\n%s\n", err, stderr)
		os.Exit(1)
	}
	answers, err := parseCandidateOutput(out, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s\n", err, out)
		os.Exit(1)
	}

	for i, triple := range answers {
		if err := verifyTriple(tests[i], triple); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\noutput triple: %v\n", i+1, tests[i].desc, err, singleInput(tests[i]), triple)
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
	tmpDir, err := os.MkdirTemp("", "oracle-2036G-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleG")
	cmd := exec.Command("go", "build", "-o", path, "2036G.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return path, cleanup, nil
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

func ensureOracleWorks(oracle, input string, t int) error {
	out, stderr, err := runBinary(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v\nstderr:\n%s", err, stderr)
	}
	if _, err := parseOracleOutput(out, t); err != nil {
		return fmt.Errorf("oracle output invalid: %v\noutput:\n%s", err, out)
	}
	return nil
}

func parseOracleOutput(out string, t int) ([][]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t*3 {
		return nil, fmt.Errorf("expected %d integers, got %d", t*3, len(tokens))
	}
	result := make([][]int64, t)
	pos := 0
	for i := 0; i < t; i++ {
		triple := make([]int64, 3)
		for j := 0; j < 3; j++ {
			val, err := strconv.ParseInt(tokens[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("token %d is not an integer (%v)", pos+1, err)
			}
			triple[j] = val
			pos++
		}
		result[i] = triple
	}
	return result, nil
}

func parseCandidateOutput(out string, t int) ([][]int64, error) {
	tokens := strings.Fields(out)
	result := make([][]int64, t)
	pos := 0
	for i := 0; i < t; i++ {
		triple := make([]int64, 0, 3)
		for len(triple) < 3 {
			if pos >= len(tokens) {
				return nil, fmt.Errorf("not enough numbers for test %d", i+1)
			}
			tok := tokens[pos]
			pos++
			if strings.EqualFold(tok, "ans") {
				continue
			}
			if strings.EqualFold(tok, "xor") {
				return nil, fmt.Errorf("unexpected token %q (queries are not allowed)", tok)
			}
			val, err := strconv.ParseInt(tok, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("token %d is not an integer (%v)", pos, err)
			}
			triple = append(triple, val)
		}
		result[i] = triple
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("extra tokens after reading all answers")
	}
	return result, nil
}

func verifyTriple(tc testCase, triple []int64) error {
	if len(triple) != 3 {
		return fmt.Errorf("expected 3 numbers, got %d", len(triple))
	}
	seen := make(map[int64]bool)
	for _, v := range triple {
		if v < 1 || v > tc.n {
			return fmt.Errorf("value %d outside [1, %d]", v, tc.n)
		}
		if seen[v] {
			return fmt.Errorf("value %d repeated", v)
		}
		seen[v] = true
	}
	expected := map[int64]int{
		tc.a: 1,
		tc.b: 1,
		tc.c: 1,
	}
	for _, v := range triple {
		if expected[v] == 0 {
			return fmt.Errorf("value %d is not one of the missing tomes", v)
		}
		expected[v]--
	}
	return nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d %d\n", tc.n, tc.a, tc.b, tc.c)
	}
	return sb.String()
}

func singleInput(tc testCase) string {
	return fmt.Sprintf("%d %d %d %d\n", tc.n, tc.a, tc.b, tc.c)
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 6, a: 2, b: 3, c: 5, desc: "sample-like"},
		{n: 3, a: 1, b: 2, c: 3, desc: "small-all-missing"},
		{n: 10, a: 1, b: 5, c: 9, desc: "small-mixed"},
		{n: 100, a: 10, b: 20, c: 30, desc: "mid-range"},
		{n: 10_000, a: 1, b: 9999, c: 5000, desc: "spread"},
		{n: 1_000_000_000_000, a: 123456789, b: 234567890, c: 345678901, desc: "large"},
		{n: 999_999_999_999_999_999, a: 111111111111111111, b: 222222222222222222, c: 333333333333333333, desc: "near-limit"},
	}

	rng := rand.New(rand.NewSource(2036))
	for i := 0; i < 20; i++ {
		n := randRange(rng, 3, 1_000_000_000_000_000_000) // up to 1e18
		a, b, c := distinctTriple(rng, n)
		tests = append(tests, testCase{
			n:    n,
			a:    a,
			b:    b,
			c:    c,
			desc: fmt.Sprintf("random-%d", i+1),
		})
	}
	return tests
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}

func distinctTriple(rng *rand.Rand, n int64) (int64, int64, int64) {
	if n < 3 {
		return 1, 1, 1
	}
	choose := func() int64 { return randRange(rng, 1, n) }
	a := choose()
	var b int64
	for {
		b = choose()
		if b != a {
			break
		}
	}
	var c int64
	for {
		c = choose()
		if c != a && c != b {
			break
		}
	}
	return a, b, c
}
