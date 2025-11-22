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
	refSource       = "2000-2999/2000-2099/2080-2089/2089/2089C2.go"
	mod             = 1000000007
	randomTests     = 120
	maxTotalLocks   = 5000
	maxParticipants = 100
)

type testCase struct {
	n int
	l int
	k int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference runtime failure: %v\n%s", err, refOut)
	}
	refAns, err := parseOutput(refOut, tests)
	if err != nil {
		fail("could not parse reference output: %v\n%s", err, refOut)
	}

	candOut, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate runtime failure: %v\n%s", err, candOut)
	}
	candAns, err := parseOutput(candOut, tests)
	if err != nil {
		fail("could not parse candidate output: %v\n%s", err, candOut)
	}

	for idx := range tests {
		r := refAns[idx]
		c := candAns[idx]
		if len(r) != len(c) {
			tc := tests[idx]
			fail("test %d length mismatch (expected %d numbers, got %d) for n=%d l=%d k=%d", idx+1, len(r), len(c), tc.n, tc.l, tc.k)
		}
		for i := range r {
			if r[i] != c[i]%mod {
				tc := tests[idx]
				fail("test %d mismatch at position %d: expected %d, got %d (n=%d l=%d k=%d)\nreference: %v\ncandidate: %v", idx+1, i+1, r[i], c[i]%mod, tc.n, tc.l, tc.k, r, c)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2089C2-ref-*")
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

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+10)

	totalL := 0
	add := func(n, l, k int) {
		if l <= 0 || totalL+l > maxTotalLocks {
			return
		}
		if n < 1 {
			n = 1
		}
		if n > maxParticipants {
			n = maxParticipants
		}
		if k < 0 {
			k = 0
		}
		if k > 25 {
			k = 25
		}
		tests = append(tests, testCase{n: n, l: l, k: k})
		totalL += l
	}

	// Sample-inspired and edge coverage.
	add(3, 1, 4)
	add(3, 2, 0)
	add(25, 2, 5)
	add(4, 10, 2)
	add(5, 1, 25)
	add(1, 1, 0)
	add(100, 50, 0)
	add(100, 75, 25)

	for len(tests) < randomTests && totalL < maxTotalLocks {
		maxL := maxTotalLocks - totalL
		l := rng.Intn(maxL) + 1
		n := rng.Intn(maxParticipants) + 1
		k := rng.Intn(26)
		add(n, l, k)
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.l))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.k))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("stderr not empty")
	}
	return out.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(tests) {
		return nil, fmt.Errorf("expected %d lines, got %d", len(tests), len(lines))
	}
	res := make([][]int64, len(tests))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != tests[idx].n {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, tests[idx].n, len(fields))
		}
		row := make([]int64, len(fields))
		for i, f := range fields {
			val, err := strconv.ParseInt(f, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q on line %d", f, idx+1)
			}
			val %= mod
			if val < 0 {
				val += mod
			}
			row[i] = val
		}
		res[idx] = row
	}
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
