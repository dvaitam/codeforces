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
	refSource     = "2000-2999/2100-2199/2110-2119/2119/2119E.go"
	maxTotalN     = 4000
	randomTests   = 120
	maxNPerRandom = 60
)

type testCase struct {
	n int
	a []int64
	b []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
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
	expected, err := parseOutput(refOut, len(tests))
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

	for i := range expected {
		if expected[i] != got[i] {
			fail("mismatch on test %d: expected %d, got %d", i+1, expected[i], got[i])
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2119E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2119E.bin")
	cmd := exec.Command("go", "build", "-o", binPath, filepath.Clean(refSource))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
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
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(fields))
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
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		if totalN+tc.n > maxTotalN {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Deterministic cases.
	add(testCase{
		n: 2,
		a: []int64{0},
		b: []int64{0, 0},
	})
	add(testCase{
		n: 2,
		a: []int64{7},
		b: []int64{5, 2},
	})
	// Clearly impossible: need a_1 bit 0 but both neighbors already have it set.
	add(testCase{
		n: 3,
		a: []int64{0, 1},
		b: []int64{1, 1, 1},
	})
	// High bits scenario near the limit.
	limit := int64(1<<29) - 1
	add(testCase{
		n: 4,
		a: []int64{limit, limit, limit},
		b: []int64{0, 0, limit, 0},
	})
	// Mixed zeros/ones to test forced propagation.
	add(testCase{
		n: 5,
		a: []int64{3, 0, 5, 1},
		b: []int64{1, 2, 3, 4, 5},
	})

	for len(tests) < randomTests && totalN < maxTotalN {
		maxN := maxNPerRandom
		if rem := maxTotalN - totalN; maxN > rem {
			maxN = rem
		}
		if maxN < 2 {
			break
		}
		n := rng.Intn(maxN-1) + 2
		a := make([]int64, n-1)
		b := make([]int64, n)
		for i := 0; i < n-1; i++ {
			a[i] = rng.Int63n(1 << 29)
		}
		for i := 0; i < n; i++ {
			b[i] = rng.Int63n(1 << 29)
		}
		add(testCase{n: n, a: a, b: b})
	}

	return tests
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
