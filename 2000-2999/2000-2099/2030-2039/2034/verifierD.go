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
	n   int
	arr []int
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
	if err := validateOutput(refOut, tests); err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid solution: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s", err, candOut)
		os.Exit(1)
	}
	if err := validateOutput(candOut, tests); err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2034D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2034D")
	cmd := exec.Command("go", "build", "-o", outPath, "2034D.go")
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

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const maxTotal = 200000
	var tests []testCase
	total := 0

	add := func(arr []int) {
		if total+len(arr) > maxTotal {
			return
		}
		tests = append(tests, testCase{n: len(arr), arr: append([]int(nil), arr...)})
		total += len(arr)
	}

	add([]int{1})
	add([]int{0, 1})
	add([]int{2, 1})
	add([]int{0, 2, 1})
	add([]int{1, 0, 2})
	add([]int{2, 2, 1, 0, 0})

	for len(tests) < 60 && total < maxTotal {
		n := rng.Intn(5000) + 1
		if total+n > maxTotal {
			n = maxTotal - total
		}
		arr := make([]int, n)
		hasOne := false
		for i := 0; i < n; i++ {
			val := rng.Intn(3)
			arr[i] = val
			if val == 1 {
				hasOne = true
			}
		}
		if !hasOne {
			arr[rng.Intn(n)] = 1
		}
		add(arr)
	}
	return tests
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func validateOutput(output string, tests []testCase) error {
	tokens := strings.Fields(output)
	pos := 0
	for ti, tc := range tests {
		if pos >= len(tokens) {
			return fmt.Errorf("test %d: missing k", ti+1)
		}
		k, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return fmt.Errorf("test %d: k is not integer (%q)", ti+1, tokens[pos])
		}
		pos++
		if k < 0 || k > tc.n {
			return fmt.Errorf("test %d: k=%d out of bounds [0,%d]", ti+1, k, tc.n)
		}
		need := 2 * k
		if pos+need > len(tokens) {
			return fmt.Errorf("test %d: expected %d tokens for moves, got %d", ti+1, need, len(tokens)-pos)
		}

		arr := append([]int(nil), tc.arr...)
		cntStart := countVals(arr)

		for i := 0; i < k; i++ {
			u, err1 := strconv.Atoi(tokens[pos])
			v, err2 := strconv.Atoi(tokens[pos+1])
			if err1 != nil || err2 != nil {
				return fmt.Errorf("test %d move %d: indices must be integers", ti+1, i+1)
			}
			pos += 2
			if u < 1 || u > tc.n || v < 1 || v > tc.n {
				return fmt.Errorf("test %d move %d: index out of range (%d,%d)", ti+1, i+1, u, v)
			}
			u--
			v--
			diff := arr[u] - arr[v]
			if diff < 0 {
				diff = -diff
			}
			if diff != 1 {
				return fmt.Errorf("test %d move %d: |a[u]-a[v]|=%d, expected 1", ti+1, i+1, diff)
			}
			if arr[u] > arr[v] {
				arr[u]--
				arr[v]++
			} else {
				arr[u]++
				arr[v]--
			}
		}

		if !nonDecreasing(arr) {
			return fmt.Errorf("test %d: final array not non-decreasing: %v", ti+1, arr)
		}
		if !equalCounts(cntStart, countVals(arr)) {
			return fmt.Errorf("test %d: value counts changed", ti+1)
		}
	}
	if pos != len(tokens) {
		return fmt.Errorf("extra output detected starting at token %q", tokens[pos])
	}
	return nil
}

func nonDecreasing(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i] < a[i-1] {
			return false
		}
	}
	return true
}

func countVals(a []int) [3]int {
	var c [3]int
	for _, v := range a {
		if v < 0 || v > 2 {
			// should not happen, but make counts obviously different
			return [3]int{-1, -1, -1}
		}
		c[v]++
	}
	return c
}

func equalCounts(a, b [3]int) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}
