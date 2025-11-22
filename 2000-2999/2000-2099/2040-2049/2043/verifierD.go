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
	"time"
)

type testCase struct {
	l, r, g int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n%s", err, refOut)
		os.Exit(1)
	}
	refPairs, err := parsePairs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s", err, candOut)
		os.Exit(1)
	}
	candPairs, err := parsePairs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		refA, refB := refPairs[2*i], refPairs[2*i+1]
		cA, cB := candPairs[2*i], candPairs[2*i+1]

		if refA == -1 && refB == -1 {
			if !(cA == -1 && cB == -1) {
				fmt.Fprintf(os.Stderr, "test %d: expected -1 -1, got %d %d\n", i+1, cA, cB)
				os.Exit(1)
			}
			continue
		}

		if cA < tc.l || cA > tc.r || cB < tc.l || cB > tc.r || cA > cB {
			fmt.Fprintf(os.Stderr, "test %d: pair out of range or unordered: %d %d\n", i+1, cA, cB)
			os.Exit(1)
		}
		if gcd(cA, cB) != tc.g {
			fmt.Fprintf(os.Stderr, "test %d: gcd(%d,%d) != %d\n", i+1, cA, cB, tc.g)
			os.Exit(1)
		}

		refDiff := refB - refA
		cDiff := cB - cA
		if cDiff != refDiff || cA != refA {
			fmt.Fprintf(os.Stderr, "test %d: candidate not optimal. expected diff %d and A %d, got diff %d and A %d\n", i+1, refDiff, refA, cDiff, cA)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2043D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2043D")
	cmd := exec.Command("go", "build", "-o", outPath, "2043D.go")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return outPath, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
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

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("unexpected stderr output")
	}
	return out.String(), nil
}

func parsePairs(output string, t int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) != t*2 {
		return nil, fmt.Errorf("expected %d tokens, got %d", t*2, len(tokens))
	}
	res := make([]int64, len(tokens))
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not integer", tok)
		}
		res[i] = val
	}
	return res, nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	dets := []testCase{
		{l: 4, r: 8, g: 2},
		{l: 4, r: 8, g: 3},
		{l: 5, r: 7, g: 6},
		{l: 1, r: 1, g: 1},
		{l: 10, r: 10, g: 2},
		{l: 1, r: 20, g: 1},
	}
	tests = append(tests, dets...)

	for len(tests) < 200 {
		l := rng.Int63n(1_000_000_000_000_000_000) + 1
		width := rng.Int63n(1_000_000_000) // keep ranges reasonable and within limits
		r := l + width
		if r > 1_000_000_000_000_000_000 {
			r = 1_000_000_000_000_000_000
		}
		g := rng.Int63n(1_000_000_000_000_000_000) + 1
		tests = append(tests, testCase{l: l, r: r, g: g})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d %d\n", tc.l, tc.r, tc.g)
	}
	return b.String()
}
