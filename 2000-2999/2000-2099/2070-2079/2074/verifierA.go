package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	maxTests     = 500
	domainLimit  = 10
	randomTrials = 400
)

type testCase struct {
	l int
	r int
	d int
	u int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	input := buildInput(tests)

	refAns := make([]bool, len(tests))
	for i, tc := range tests {
		refAns[i] = isSquare(tc)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
		os.Exit(1)
	}

	for i := range tests {
		if candAns[i] != refAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %s, got %s (l=%d r=%d d=%d u=%d)\n",
				i+1, boolToStr(refAns[i]), boolToStr(candAns[i]), tests[i].l, tests[i].r, tests[i].d, tests[i].u)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func isSquare(tc testCase) bool {
	return tc.l == tc.r && tc.r == tc.d && tc.d == tc.u
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
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

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d %d %d\n", tc.l, tc.r, tc.d, tc.u)
	}
	return b.String()
}

func parseOutput(out string, t int) ([]bool, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(fields))
	}
	res := make([]bool, t)
	for i, f := range fields {
		val, err := parseBoolToken(f)
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func parseBoolToken(s string) (bool, error) {
	low := strings.ToLower(strings.TrimSpace(s))
	switch low {
	case "yes", "y", "true", "1":
		return true, nil
	case "no", "n", "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean token %q", s)
	}
}

func boolToStr(v bool) string {
	if v {
		return "Yes"
	}
	return "No"
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Basic deterministic cases.
	tests = append(tests,
		testCase{l: 1, r: 1, d: 1, u: 1},
		testCase{l: 2, r: 2, d: 3, u: 3},
		testCase{l: 2, r: 3, d: 2, u: 3},
		testCase{l: 5, r: 5, d: 1, u: 2},
		testCase{l: 1, r: 2, d: 1, u: 2},
	)

	for len(tests) < maxTests && len(tests) < randomTrials {
		l := rng.Intn(domainLimit) + 1
		r := rng.Intn(domainLimit) + 1
		d := rng.Intn(domainLimit) + 1
		u := rng.Intn(domainLimit) + 1
		tests = append(tests, testCase{l: l, r: r, d: d, u: u})
	}

	if len(tests) > maxTests {
		tests = tests[:maxTests]
	}

	return tests
}
