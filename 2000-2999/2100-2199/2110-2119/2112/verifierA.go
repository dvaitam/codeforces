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
	refSource   = "./2112A.go"
	maxTests    = 500
	randomCases = 400
)

type testCase struct {
	a int
	x int
	y int
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
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (a=%d x=%d y=%d): expected %s, got %s\n",
				i+1, tc.a, tc.x, tc.y, boolToStr(refAns[i]), boolToStr(candAns[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2112A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
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
		fmt.Fprintf(&b, "%d %d %d\n", tc.a, tc.x, tc.y)
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
			return nil, fmt.Errorf("token %d: %v", i+1, err)
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
		return "YES"
	}
	return "NO"
}

func generateTests() []testCase {
	// Exhaustive for small values ensures coverage.
	var tests []testCase
	for a := 1; a <= 8; a++ {
		for x := 1; x <= 8; x++ {
			for y := 1; y <= 8; y++ {
				if a != x && a != y && x != y {
					tests = append(tests, testCase{a: a, x: x, y: y})
					if len(tests) >= maxTests {
						return tests
					}
				}
			}
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < maxTests && len(tests) < randomCases {
		a := rng.Intn(100) + 1
		x := rng.Intn(100) + 1
		y := rng.Intn(100) + 1
		if a == x || a == y || x == y {
			continue
		}
		tests = append(tests, testCase{a: a, x: x, y: y})
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{a: 1, x: 2, y: 3})
	}
	return tests
}
