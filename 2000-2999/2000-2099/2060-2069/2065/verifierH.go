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
	refSource   = "2065H.go"
	mod         = 998244353
	maxTotalN   = 180000
	maxTotalQ   = 180000
	targetTests = 160
	maxSingleN  = 200000
	maxSingleQ  = 200000
)

type testCase struct {
	s   string
	q   int
	ops []int
	n   int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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
	refAns, err := parseAnswers(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d values, got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}
	for i := range refAns {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "mismatch at answer %d: expected %d, got %d\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d test cases).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2065H-ref-*")
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

func parseAnswers(out string, tests []testCase) ([]int64, error) {
	total := 0
	for _, tc := range tests {
		total += tc.q
	}
	tokens := strings.Fields(out)
	if len(tokens) != total {
		return nil, fmt.Errorf("expected %d integers, got %d", total, len(tokens))
	}
	ans := make([]int64, total)
	for i, t := range tokens {
		val, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse integer %q at position %d: %v", t, i+1, err)
		}
		val %= mod
		if val < 0 {
			val += mod
		}
		ans[i] = val
	}
	return ans, nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%s\n", tc.s)
		fmt.Fprintf(&b, "%d\n", tc.q)
		for i, v := range tc.ops {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0
	totalQ := 0

	add := func(tc testCase) {
		if tc.n <= 0 || tc.q <= 0 {
			return
		}
		if totalN+tc.n > maxTotalN || totalQ+tc.q > maxTotalQ {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
		totalQ += tc.q
	}

	// Manual small cases.
	add(manualCase("0", []int{1}))
	add(manualCase("1", []int{1, 1, 1}))
	add(manualCase("101", []int{1, 2}))
	add(manualCase("111000", []int{1, 3, 6}))

	// Random cases.
	for len(tests) < targetTests && totalN < maxTotalN && totalQ < maxTotalQ {
		remainingN := maxTotalN - totalN
		remainingQ := maxTotalQ - totalQ
		if remainingN <= 0 || remainingQ <= 0 {
			break
		}
		n := rng.Intn(min(4000, remainingN)) + 1
		if n > maxSingleN {
			n = maxSingleN
		}
		maxQ := min(min(3*n+5, 5000), remainingQ)
		if maxQ <= 0 {
			break
		}
		q := rng.Intn(maxQ) + 1
		if q > maxSingleQ {
			q = maxSingleQ
		}
		ops := make([]int, q)
		for i := 0; i < q; i++ {
			ops[i] = rng.Intn(n) + 1
		}
		sb := make([]byte, n)
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb[i] = '0'
			} else {
				sb[i] = '1'
			}
		}
		add(testCase{
			s:   string(sb),
			q:   q,
			ops: ops,
			n:   n,
		})
	}

	if len(tests) == 0 {
		add(manualCase("0", []int{1}))
	}
	return tests
}

func manualCase(s string, ops []int) testCase {
	return testCase{
		s:   s,
		q:   len(ops),
		ops: ops,
		n:   len(s),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
