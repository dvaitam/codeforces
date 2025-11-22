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
	n int
	m int
	k int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	ref, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := serializeInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	refSeqs, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	candSeqs, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		refVal := computeF(refSeqs[i], tc.m, tc.k)
		candVal := computeF(candSeqs[i], tc.m, tc.k)
		if candVal != refVal {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected f=%d got %d\nn=%d m=%d k=%d\n", i+1, refVal, candVal, tc.n, tc.m, tc.k)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine current path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2084D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref2084D")
	cmd := exec.Command("go", "build", "-o", outPath, "2084D.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	tokens := strings.Fields(out)
	res := make([][]int64, len(tests))
	pos := 0
	for i, tc := range tests {
		need := tc.n
		if pos+need > len(tokens) {
			return nil, fmt.Errorf("test %d: not enough numbers", i+1)
		}
		res[i] = make([]int64, need)
		for j := 0; j < need; j++ {
			val, err := strconv.ParseInt(tokens[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer %q", i+1, tokens[pos])
			}
			if val < 0 || val > 1_000_000_000 {
				return nil, fmt.Errorf("test %d: value out of bounds %d", i+1, val)
			}
			res[i][j] = val
			pos++
		}
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("extra output tokens: expected %d got %d", pos, len(tokens))
	}
	return res, nil
}

// computeF returns minimal possible mex after at most m deletions of length k.
// Brute-force BFS over states; only used on small tests.
func computeF(seq []int64, m, k int) int {
	type state struct {
		arr []int64
		ops int
	}
	encode := func(a []int64) string {
		var sb strings.Builder
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		return sb.String()
	}
	mex := func(a []int64) int {
		seen := make(map[int64]struct{}, len(a))
		for _, v := range a {
			seen[v] = struct{}{}
		}
		for x := 0; ; x++ {
			if _, ok := seen[int64(x)]; !ok {
				return x
			}
		}
	}

	q := []state{{arr: seq, ops: 0}}
	vis := make(map[string]struct{})
	vis[encode(seq)] = struct{}{}
	best := mex(seq)

	for head := 0; head < len(q); head++ {
		cur := q[head]
		if cur.ops == m {
			continue
		}
		l := len(cur.arr)
		for i := 0; i+k <= l; i++ {
			nextArr := make([]int64, 0, l-k)
			nextArr = append(nextArr, cur.arr[:i]...)
			nextArr = append(nextArr, cur.arr[i+k:]...)
			key := encode(nextArr)
			if _, ok := vis[key]; ok {
				continue
			}
			vis[key] = struct{}{}
			val := mex(nextArr)
			if val < best {
				best = val
			}
			q = append(q, state{arr: nextArr, ops: cur.ops + 1})
		}
	}
	return best
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, m: 1, k: 1},
		{n: 5, m: 2, k: 2},
		{n: 6, m: 1, k: 3},
		{n: 4, m: 1, k: 2},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		n := rng.Intn(8) + 2
		k := rng.Intn(n-1) + 1
		m := rng.Intn(n-1) + 1
		if m*k >= n {
			continue
		}
		tests = append(tests, testCase{n: n, m: m, k: k})
	}
	return tests
}
