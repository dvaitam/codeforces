package main

import (
	"bufio"
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
	refSource   = "./2114A.go"
	targetTests = 200
	maxTests    = 10000
)

type testCase struct {
	s string
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
	refResults, err := parseReference(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candTokens := strings.Fields(candOut)
	if len(candTokens) == 0 {
		fmt.Fprintln(os.Stderr, "candidate produced empty output")
		os.Exit(1)
	}

	// Validate per test using candidate lines.
	candScanner := bufio.NewScanner(strings.NewReader(candOut))
	candScanner.Split(bufio.ScanLines)
	var idx int
	for idx = 0; idx < len(tests) && candScanner.Scan(); idx++ {
		line := strings.TrimSpace(candScanner.Text())
		if line == "" {
			// skip empty lines
			idx--
			continue
		}
		val := yearValue(tests[idx].s)
		expectSquare := refResults[idx]
		if err := validateLine(line, val, expectSquare); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid output: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	if scanErr := candScanner.Err(); scanErr != nil {
		fmt.Fprintf(os.Stderr, "error reading candidate output: %v\n", scanErr)
		os.Exit(1)
	}
	if idx != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d output lines, got %d\n", len(tests), idx)
		os.Exit(1)
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2114A-ref-*")
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

func parseReference(out string, t int) ([]bool, error) {
	lines := strings.Fields(out)
	if len(lines) < t {
		return nil, fmt.Errorf("expected at least %d outputs, got %d tokens", t, len(lines))
	}
	res := make([]bool, 0, t)
	reader := bufio.NewScanner(strings.NewReader(out))
	reader.Split(bufio.ScanLines)
	for reader.Scan() && len(res) < t {
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}
		if line == "-1" {
			res = append(res, false)
		} else {
			res = append(res, true)
		}
	}
	if err := reader.Err(); err != nil {
		return nil, err
	}
	if len(res) != t {
		return nil, fmt.Errorf("expected %d lines, parsed %d", t, len(res))
	}
	return res, nil
}

func validateLine(line string, val int, expectSquare bool) error {
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return fmt.Errorf("empty line")
	}
	if tokens[0] == "-1" {
		if len(tokens) != 1 {
			return fmt.Errorf("extra tokens after -1")
		}
		if expectSquare {
			return fmt.Errorf("expected square representation, got -1")
		}
		return nil
	}
	if len(tokens) != 2 {
		return fmt.Errorf("expected two integers or -1, got %d tokens", len(tokens))
	}
	a, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse first integer: %v", err)
	}
	b, err := strconv.ParseInt(tokens[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse second integer: %v", err)
	}
	if a < 0 || b < 0 {
		return fmt.Errorf("a and b must be non-negative")
	}
	sum := a + b
	if sum*sum != int64(val) {
		return fmt.Errorf("(a+b)^2 != %d (got a=%d b=%d)", val, a, b)
	}
	if !expectSquare {
		return fmt.Errorf("value is not a perfect square but pair provided")
	}
	return nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%s\n", tc.s)
	}
	return b.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	add := func(tc testCase) {
		if len(tests) >= maxTests {
			return
		}
		tests = append(tests, tc)
	}

	// Sample tests.
	add(testCase{s: "0001"})
	add(testCase{s: "1001"})
	add(testCase{s: "1000"})
	add(testCase{s: "4900"})
	add(testCase{s: "2025"})

	// Cover 0 and max.
	add(testCase{s: "0000"})
	add(testCase{s: "9999"})

	for len(tests) < targetTests && len(tests) < maxTests {
		val := rng.Intn(10000)
		s := fmt.Sprintf("%04d", val)
		add(testCase{s: s})
	}

	if len(tests) == 0 {
		add(testCase{s: "0000"})
	}
	return tests
}

func yearValue(s string) int {
	val := 0
	for i := 0; i < len(s); i++ {
		val = val*10 + int(s[i]-'0')
	}
	return val
}
