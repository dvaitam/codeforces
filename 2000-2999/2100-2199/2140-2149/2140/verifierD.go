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
	refSource   = "2000-2999/2100-2199/2140-2149/2140/2140D.go"
	targetTests = 120
	maxTotalN   = 200000
	maxCoord    = int64(1_000_000_000)
)

type segment struct {
	l int64
	r int64
}

type testCase struct {
	segs []segment
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
	refAns, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d, got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}
	for i := range refAns {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d, got %d\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2140D-ref-*")
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

func parseAnswers(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i := 0; i < t; i++ {
		val, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse answer %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", len(tc.segs))
		for _, sg := range tc.segs {
			fmt.Fprintf(&b, "%d %d\n", sg.l, sg.r)
		}
	}
	return b.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		if totalN+len(tc.segs) > maxTotalN {
			return
		}
		tests = append(tests, tc)
		totalN += len(tc.segs)
	}

	// Sample cases from statement.
	add(testCase{segs: []segment{{1, 1_000_000_000}, {1, 1_000_000_000}}})
	add(testCase{segs: []segment{{1, 10}, {2, 15}, {3, 9}}})
	add(testCase{segs: []segment{{1, 1}, {2, 7}, {1, 3}, {11, 15}, {1, 1}}})
	add(testCase{segs: []segment{{1, 1}}})

	// Small random cases.
	add(randomCase(1, rng))
	add(randomCase(2, rng))
	add(randomCase(3, rng))

	// Random cases up to limits.
	for len(tests) < targetTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		if remain < 1 {
			break
		}
		n := rng.Intn(min(5000, remain)) + 1
		add(randomCase(n, rng))
	}

	if len(tests) == 0 {
		add(randomCase(1, rng))
	}
	return tests
}

func randomCase(n int, rng *rand.Rand) testCase {
	segs := make([]segment, n)
	for i := 0; i < n; i++ {
		l := rng.Int63n(maxCoord) + 1
		r := l + rng.Int63n(maxCoord-l+1)
		segs[i] = segment{l: l, r: r}
	}
	return testCase{segs: segs}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
