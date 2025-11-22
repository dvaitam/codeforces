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
	"time"
)

const (
	refSource   = "2000-2999/2100-2199/2130-2139/2132/2132C2.go"
	randomCases = 180
	maxT        = 400
	maxValue    = 1_000_000_000
)

type testCase struct {
	n int64
	k int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fail("%v", err)
	}
	defer cleanup()

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference runtime error: %v\n%s", err, refOut)
	}
	expect, err := parseOutput(refOut, len(tests))
	if err != nil {
		fail("reference output invalid: %v\n%s", err, refOut)
	}

	candCmd := commandFor(candidate)
	candOut, err := runProgram(candCmd, input)
	if err != nil {
		fail("candidate runtime error: %v\n%s", err, candOut)
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fail("candidate output invalid: %v\n%s", err, candOut)
	}

	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch on test %d: expected %d, got %d", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2132C2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2132C2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, filepath.Clean(refSource))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func parseOutput(raw string, t int) ([]int64, error) {
	fields := strings.Fields(raw)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok, i+1)
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	add := func(n, k int64) {
		if len(tests) >= maxT {
			return
		}
		if n < 1 {
			n = 1
		}
		if k < 1 {
			k = 1
		}
		tests = append(tests, testCase{n: n, k: k})
	}

	// Statement sample.
	add(1, 1)
	add(3, 3)
	add(8, 3)
	add(2, 4)
	add(10, 10)
	add(20, 14)
	add(3, 2)
	add(9, 1)

	// Deterministic coverage.
	add(1, 5)                   // k much larger than n.
	add(2, 1)                   // minimal k.
	add(3, 1)                   // impossible due to too few deals.
	add(9, 3)                   // exact power of 3 with limited merges.
	add(27, 2)                  // power with very small k.
	add(1000000000, 1)          // largest n with minimal k.
	add(1000000000, 1000000000) // largest n with huge k.
	add(999999937, 30)          // large prime-ish n, moderate k.
	add(81, 4)                  // multiple base3 digits.
	add(50, 50)                 // k >= n => 3n rule.

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for r := 0; r < randomCases && len(tests) < maxT; r++ {
		n := rng.Int63n(maxValue) + 1
		k := rng.Int63n(maxValue) + 1
		add(n, k)
	}

	return tests
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
