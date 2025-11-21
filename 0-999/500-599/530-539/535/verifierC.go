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

const refSource = "0-999/500-599/530-539/535/535C.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := compareOutputs(refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, err, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-535C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref535C.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func compareOutputs(refOut, candOut string) error {
	expTokens := strings.Fields(refOut)
	gotTokens := strings.Fields(candOut)
	if len(expTokens) != len(gotTokens) {
		return fmt.Errorf("expected %d tokens got %d", len(expTokens), len(gotTokens))
	}
	for i := range expTokens {
		expVal, err := strconv.ParseInt(expTokens[i], 10, 64)
		if err != nil {
			return fmt.Errorf("reference token %q invalid: %v", expTokens[i], err)
		}
		gotVal, err := strconv.ParseInt(gotTokens[i], 10, 64)
		if err != nil {
			return fmt.Errorf("candidate token %q invalid: %v", gotTokens[i], err)
		}
		if expVal != gotVal {
			return fmt.Errorf("mismatch at position %d: expected %d got %d", i+1, expVal, gotVal)
		}
	}
	return nil
}

func buildTests() []testCase {
	tests := append([]testCase{}, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	tests = append(tests, stressCase())
	return tests
}

func manualTests() []testCase {
	return []testCase{
		{
			name:  "single_easy",
			input: formatInput(1, 1, []query{{1, 1, 1}}),
		},
		{
			name:  "impossible_first",
			input: formatInput(5, 3, []query{{1, 4, 2}}),
		},
		{
			name: "increasing_sequence",
			input: formatInput(2, 2, []query{
				{1, 10, 3},
				{2, 6, 2},
				{3, 4, 1},
			}),
		},
		{
			name: "multiple_queries",
			input: formatInput(3, 5, []query{
				{1, 20, 4},
				{5, 1000000, 1000000},
				{10, 15, 2},
				{7, 12, 1},
			}),
		},
	}
}

type query struct {
	l int64
	t int64
	m int64
}

func formatInput(A, B int64, qs []query) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", A, B, len(qs))
	for _, q := range qs {
		fmt.Fprintf(&sb, "%d %d %d\n", q.l, q.t, q.m)
	}
	return sb.String()
}

func randomCase(rng *rand.Rand, idx int) testCase {
	A := randInRange(rng, 1, 1_000_000)
	B := randInRange(rng, 1, 1_000_000)
	n := rng.Intn(50) + 1
	var qs []query
	for i := 0; i < n; i++ {
		l := randInRange(rng, 1, 1_000_000)
		t := randInRange(rng, 1, 1_000_000)
		m := randInRange(rng, 1, 1_000_000)
		qs = append(qs, query{l: l, t: t, m: m})
	}
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: formatInput(A, B, qs),
	}
}

func stressCase() testCase {
	const A = 1_000_000
	const B = 999_999
	const n = 100000
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", A, B, n)
	for i := 1; i <= n; i++ {
		l := int64(i)
		t := int64(1_000_000)
		m := int64((i % 1000) + 1)
		fmt.Fprintf(&sb, "%d %d %d\n", l, t, m)
	}
	return testCase{
		name:  "stress_large",
		input: sb.String(),
	}
}

func randInRange(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}
