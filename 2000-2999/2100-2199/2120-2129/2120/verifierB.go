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
	refSource   = "./2120B.go"
	maxTests    = 300
	totalNLimit = 800
)

type ball struct {
	dx int64
	dy int64
	x  int64
	y  int64
}

type testCase struct {
	n     int
	s     int64
	balls []ball
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d, got %d\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2120B-ref-*")
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
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.s)
		for _, bl := range tc.balls {
			fmt.Fprintf(&b, "%d %d %d %d\n", bl.dx, bl.dy, bl.x, bl.y)
		}
	}
	return b.String()
}

func parseOutput(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(fields))
	}
	res := make([]int, t)
	for i := 0; i < t; i++ {
		val, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("failed to parse integer %q at position %d: %v", fields[i], i+1, err)
		}
		res[i] = val
		if res[i] < 0 || res[i] > 1000 {
			return nil, fmt.Errorf("parsed value out of range at test %d: %d", i+1, res[i])
		}
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sumN := 0

	add := func(tc testCase) {
		if len(tests) >= maxTests || sumN+tc.n > totalNLimit {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	// Deterministic scenarios.
	add(testCase{n: 1, s: 2, balls: []ball{{dx: 1, dy: 1, x: 1, y: 1}}})                              // direct to corner
	add(testCase{n: 1, s: 5, balls: []ball{{dx: 1, dy: -1, x: 2, y: 3}}})                             // not potted
	add(testCase{n: 2, s: 4, balls: []ball{{dx: 1, dy: 1, x: 1, y: 1}, {dx: -1, dy: 1, x: 3, y: 1}}}) // mixed

	for len(tests) < maxTests && sumN < totalNLimit {
		n := rng.Intn(20) + 1
		if sumN+n > totalNLimit {
			n = totalNLimit - sumN
		}
		s := int64(rng.Intn(1_000_000_000-2) + 2)
		balls := make([]ball, n)
		for i := 0; i < n; i++ {
			dx := int64(1)
			if rng.Intn(2) == 0 {
				dx = -1
			}
			dy := int64(1)
			if rng.Intn(2) == 0 {
				dy = -1
			}
			x := int64(rng.Intn(int(s-1)) + 1)
			y := int64(rng.Intn(int(s-1)) + 1)
			balls[i] = ball{dx: dx, dy: dy, x: x, y: y}
		}
		add(testCase{n: n, s: s, balls: balls})
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{n: 1, s: 2, balls: []ball{{dx: 1, dy: 1, x: 1, y: 1}}})
	}
	return tests
}
