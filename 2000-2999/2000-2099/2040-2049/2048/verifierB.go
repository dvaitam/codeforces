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
	refSource   = "./2048B.go"
	totalNLimit = 100000
	maxTests    = 1000
)

type testCase struct {
	n int
	k int
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
	refAns := parseOutputs(refOut, len(tests))

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns := parseOutputs(candOut, len(tests))

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "test count mismatch: expected %d got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}

	for i, tc := range tests {
		refPerm := refAns[i]
		candPerm := candAns[i]
		if err := validatePermutation(tc, candPerm); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid: %v\n", i+1, err)
			os.Exit(1)
		}
		if score(tc, candPerm) != score(tc, refPerm) {
			fmt.Fprintf(os.Stderr, "test %d wrong answer: expected score %d, got %d\n", i+1, score(tc, refPerm), score(tc, candPerm))
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2048B-ref-*")
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

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0

	add := func(n, k int) {
		if totalN+n > totalNLimit {
			return
		}
		tests = append(tests, testCase{n: n, k: k})
		totalN += n
	}

	add(1, 1)
	add(4, 2)
	add(6, 1)
	add(8, 3)
	add(10, 5)

	for len(tests) < maxTests && totalN < totalNLimit {
		maxN := totalNLimit - totalN
		if maxN == 0 {
			break
		}
		if maxN > 1000 {
			maxN = 1000
		}
		n := rng.Intn(maxN) + 1
		k := rng.Intn(n) + 1
		add(n, k)
	}
	return tests
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.k)
	}
	return b.String()
}

func parseOutputs(out string, t int) [][]int {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	res := make([][]int, 0, t)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		row := make([]int, len(fields))
		for i, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to parse value %q: %v\n", f, err)
				os.Exit(1)
			}
			row[i] = val
		}
		res = append(res, row)
		if len(res) == t {
			break
		}
	}
	if len(res) != t {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", t, len(res))
		os.Exit(1)
	}
	return res
}

func validatePermutation(tc testCase, p []int) error {
	if len(p) != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(p))
	}
	seen := make([]bool, tc.n+1)
	for i, v := range p {
		if v < 1 || v > tc.n {
			return fmt.Errorf("p[%d]=%d out of range", i+1, v)
		}
		if seen[v] {
			return fmt.Errorf("p[%d]=%d duplicated", i+1, v)
		}
		seen[v] = true
	}
	return nil
}

func score(tc testCase, p []int) int {
	n, k := tc.n, tc.k
	if k == 0 {
		return 0
	}
	sum := 0
	for i := 0; i <= n-k; i++ {
		minVal := p[i]
		for j := i + 1; j < i+k; j++ {
			if p[j] < minVal {
				minVal = p[j]
			}
		}
		sum += minVal
	}
	return sum
}
