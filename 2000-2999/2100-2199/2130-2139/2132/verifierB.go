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
	refSource   = "2000-2999/2100-2199/2130-2139/2132/2132B.go"
	targetTests = 160
	maxTests    = 300
	maxN        = int64(1_000_000_000_000_000_000)
	minN        = int64(11)
)

type testCase struct {
	n int64
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
	refAns, err := parseReference(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseReference(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d test outputs, got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}
	for i := range refAns {
		if !equalSlices(refAns[i], candAns[i]) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %v, got %v\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2132B-ref-*")
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

func parseReference(out string, t int) ([][]int64, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	res := make([][]int64, 0, t)
	i := 0
	for i < len(lines) && len(res) < t {
		line := strings.TrimSpace(lines[i])
		i++
		if line == "" {
			continue
		}
		cnt, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse count on test %d: %v", len(res)+1, err)
		}
		if cnt == 0 {
			res = append(res, []int64{})
			continue
		}
		if i >= len(lines) {
			return nil, fmt.Errorf("missing line with values for test %d", len(res)+1)
		}
		valLine := strings.TrimSpace(lines[i])
		i++
		tokens := strings.Fields(valLine)
		if len(tokens) != cnt {
			return nil, fmt.Errorf("expected %d values on test %d, got %d", cnt, len(res)+1, len(tokens))
		}
		arr := make([]int64, cnt)
		for j, tok := range tokens {
			v, err := strconv.ParseInt(tok, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse value %q on test %d: %v", tok, len(res)+1, err)
			}
			arr[j] = v
		}
		res = append(res, arr)
	}
	if len(res) != t {
		return nil, fmt.Errorf("expected %d test outputs, parsed %d", t, len(res))
	}
	return res, nil
}

func equalSlices(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
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

	// Sample values.
	add(testCase{n: 1111})
	add(testCase{n: 12})
	add(testCase{n: 999999999999999999})
	add(testCase{n: 1000000000000000000})

	for len(tests) < targetTests && len(tests) < maxTests {
		// random in range [minN, maxN]
		n := rng.Int63n(maxN-minN+1) + minN
		add(testCase{n: n})
	}

	if len(tests) == 0 {
		add(testCase{n: 1111})
	}
	return tests
}
