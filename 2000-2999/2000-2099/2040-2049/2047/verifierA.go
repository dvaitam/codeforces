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

const (
	maxDays      = 100
	maxPerDay    = 100
	maxTestCases = 500
)

type testCase struct {
	n    int
	days []int
	desc string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	if len(tests) > maxTestCases {
		tests = tests[:maxTestCases]
	}

	input := buildInput(tests)

	expOut, stderr, err := runBinary(oraclePath, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\nstderr:\n%s\n", err, stderr)
		os.Exit(1)
	}
	expected, err := parseOutputs(expOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\noutput:\n%s\n", err, expOut)
		os.Exit(1)
	}

	gotOut, gotErr, err := runBinary(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\nstderr:\n%s\n", err, gotErr)
		os.Exit(1)
	}
	actual, err := parseOutputs(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\noutput:\n%s\n", err, gotOut)
		os.Exit(1)
	}

	for i := range tests {
		if actual[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d, got %d\ninput:\n%s\n", i+1, tests[i].desc, expected[i], actual[i], singleInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2047A-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", bin, "2047A.go")
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

func parseOutputs(out string, t int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(tokens))
	}
	res := make([]int, t)
	for i, tok := range tokens {
		val64, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d is not an integer (%v)", i+1, err)
		}
		res[i] = int(val64)
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, v := range tc.days {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func singleInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", tc.n)
	for i, v := range tc.days {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 1, days: []int{1}, desc: "single-day"},
		{n: 2, days: []int{1, 8}, desc: "two-days"},
		{n: 3, days: []int{1, 4, 4}, desc: "three-days"},
		{n: 5, days: []int{1, 3, 5, 7, 9}, desc: "increasing"},
	}
	rng := rand.New(rand.NewSource(2047))
	for i := 0; i < 80 && len(tests) < maxTestCases; i++ {
		if tc, ok := randomCase(rng, i+1); ok {
			tests = append(tests, tc)
		}
	}
	return tests
}

func randomCase(rng *rand.Rand, idx int) (testCase, bool) {
	for attempts := 0; attempts < 200; attempts++ {
		k := rng.Intn(50)*2 + 1 // odd number from 1 to 99
		total := k * k
		if seq, ok := buildDays(total, rng); ok && len(seq) <= maxDays {
			return testCase{
				n:    len(seq),
				days: seq,
				desc: fmt.Sprintf("random-%d", idx),
			}, true
		}
	}
	return testCase{}, false
}

func buildDays(total int, rng *rand.Rand) ([]int, bool) {
	if total <= 0 {
		return nil, false
	}
	days := []int{1}
	remaining := total - 1
	for remaining > 0 {
		slotsLeft := maxDays - len(days)
		if slotsLeft <= 0 {
			return nil, false
		}
		maxTake := maxPerDay
		if remaining < maxTake {
			maxTake = remaining
		}
		minTake := 1
		limit := remaining - (slotsLeft-1)*maxPerDay
		if limit > minTake {
			minTake = limit
		}
		if minTake < 1 {
			minTake = 1
		}
		if minTake > maxTake {
			return nil, false
		}
		take := minTake
		if maxTake > minTake {
			take += rng.Intn(maxTake - minTake + 1)
		}
		remaining -= take
		days = append(days, take)
	}
	return days, true
}
