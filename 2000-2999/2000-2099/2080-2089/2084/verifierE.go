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

const mod int64 = 1_000_000_007

type testCase struct {
	n   int
	arr []int
}

// mex returns the minimum excludant of a slice.
func mex(a []int) int {
	seen := make(map[int]bool, len(a))
	for _, v := range a {
		seen[v] = true
	}
	for i := 0; ; i++ {
		if !seen[i] {
			return i
		}
	}
}

// value computes sum of MEX over all non-empty subsegments of a.
func value(a []int) int64 {
	n := len(a)
	var total int64
	for l := 0; l < n; l++ {
		for r := l; r < n; r++ {
			total += int64(mex(a[l : r+1]))
		}
	}
	return total % mod
}

// bruteForce enumerates all completions of arr (positions with -1 filled by
// missing values) and sums their values mod MOD.
func bruteForce(tc testCase) int64 {
	n := tc.n
	arr := tc.arr

	missingSet := make([]bool, n)
	for i := range missingSet {
		missingSet[i] = true
	}
	for _, v := range arr {
		if v != -1 {
			missingSet[v] = false
		}
	}

	missing := []int{}
	for v := 0; v < n; v++ {
		if missingSet[v] {
			missing = append(missing, v)
		}
	}
	freePos := []int{}
	for i, v := range arr {
		if v == -1 {
			freePos = append(freePos, i)
		}
	}

	filled := make([]int, n)
	copy(filled, arr)

	var total int64
	var permute func(k int)
	permute = func(k int) {
		if k == len(missing) {
			total = (total + value(filled)) % mod
			return
		}
		for i := k; i < len(missing); i++ {
			missing[k], missing[i] = missing[i], missing[k]
			filled[freePos[k]] = missing[k]
			permute(k + 1)
			missing[k], missing[i] = missing[i], missing[k]
		}
	}
	permute(0)
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	tests := generateTests()
	input := buildInput(tests)

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s", err, candOut)
		os.Exit(1)
	}
	candVals, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		expected := bruteForce(tc)
		if expected != candVals[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (n=%d arr=%v): expected %d got %d\n",
				i+1, tc.n, tc.arr, expected, candVals[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
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

func parseOutputs(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not integer", tok)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Sample-based deterministic cases from the problem statement.
	add := func(tc testCase) { tests = append(tests, tc) }
	add(testCase{n: 2, arr: []int{0, -1}})
	add(testCase{n: 2, arr: []int{-1, -1}})
	add(testCase{n: 3, arr: []int{2, 0, 1}})
	add(testCase{n: 3, arr: []int{-1, 2, -1}})
	add(testCase{n: 5, arr: []int{-1, 0, -1, 2, -1}})

	// Random small cases — keep n ≤ 7 so brute force (n! permutations) is fast.
	for len(tests) < 200 {
		n := rng.Intn(7) + 1
		add(randomCase(rng, n))
	}

	return tests
}

func randomCase(rng *rand.Rand, n int) testCase {
	arr := make([]int, n)
	used := make(map[int]bool)
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 {
			arr[i] = -1
			continue
		}
		val := rng.Intn(n)
		for used[val] {
			val = rng.Intn(n)
		}
		used[val] = true
		arr[i] = val
	}
	return testCase{n: n, arr: arr}
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
