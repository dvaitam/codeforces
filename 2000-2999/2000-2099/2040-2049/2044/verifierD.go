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
	refSource    = "./2044D.go"
	totalNLimit  = 200000
	maxTestCases = 1000
)

type testCase struct {
	n int
	a []int
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

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refTokens := parseOutputs(refOut, len(tests))
	candTokens := parseOutputs(candOut, len(tests))

	if len(refTokens) != len(candTokens) {
		fmt.Fprintf(os.Stderr, "token count mismatch: expected %d got %d\n", len(refTokens), len(candTokens))
		os.Exit(1)
	}

	for idx := range tests {
		refAns := refTokens[idx]
		candAns := candTokens[idx]
		if err := validateAnswer(tests[idx], candAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		if len(refAns) != len(candAns) {
			fmt.Fprintf(os.Stderr, "test %d length mismatch\n", idx+1)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2044D-ref-*")
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

	add := func(a []int) bool {
		n := len(a)
		if n == 0 {
			return false
		}
		if totalN+n > totalNLimit {
			return false
		}
		tests = append(tests, testCase{n: n, a: a})
		totalN += n
		return true
	}

	// Deterministic tests
	add([]int{1})
	add([]int{1, 2})
	add([]int{1, 1, 1, 1})
	add([]int{4, 5, 5, 5, 1, 1, 2, 2})
	add([]int{1, 1, 2, 2, 1, 1, 3, 3, 1, 1})

	for totalN < totalNLimit && len(tests) < maxTestCases {
		remaining := totalNLimit - totalN
		if remaining == 0 {
			break
		}
		maxN := remaining
		if maxN > 5000 {
			maxN = 5000
		}
		n := rng.Intn(maxN) + 1
		a := make([]int, n)
		heavy := rng.Intn(n) + 1
		for i := 0; i < n; i++ {
			if rng.Intn(5) == 0 {
				a[i] = heavy
			} else {
				a[i] = rng.Intn(n) + 1
			}
		}
		add(a)
	}
	return tests
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		for i, x := range tc.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", x)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func parseOutputs(out string, t int) [][]int {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	res := make([][]int, 0, t)
	idx := 0
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
				fmt.Fprintf(os.Stderr, "failed to parse output value %q: %v\n", f, err)
				os.Exit(1)
			}
			row[i] = val
		}
		res = append(res, row)
		idx++
		if idx == t {
			break
		}
	}
	if len(res) != t {
		fmt.Fprintf(os.Stderr, "expected %d test outputs, got %d\n", t, len(res))
		os.Exit(1)
	}
	return res
}

func validateAnswer(tc testCase, ans []int) error {
	if len(ans) != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(ans))
	}
	freq := make([]int, tc.n+2)
	maxFreq := 0
	for i, val := range ans {
		if val < 1 || val > tc.n {
			return fmt.Errorf("b[%d]=%d out of range", i+1, val)
		}
		freq[val]++
		if freq[val] > maxFreq {
			maxFreq = freq[val]
		}
		if freq[tc.a[i]] != maxFreq {
			return fmt.Errorf("a[%d]=%d is not a mode (max freq %d, freq=%d)", i+1, tc.a[i], maxFreq, freq[tc.a[i]])
		}
	}
	return nil
}
