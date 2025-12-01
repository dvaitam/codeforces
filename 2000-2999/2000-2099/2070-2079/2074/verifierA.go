package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	refSource    = "2074A.go"
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

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
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

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2074A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
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
