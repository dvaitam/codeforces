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
)

type testCase struct {
	n int
	k int
	a []int
}

type triple struct {
	v int
	l int
	r int
}

type node struct {
	left  int
	right int
	sum   int
}

const (
	randSeed   = 2128
	maxTests   = 120
	usageGuide = "usage: go run verifierE2.go /path/to/candidate"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usageGuide)
		return
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Printf("failed to build oracle: %v\n", err)
		return
	}
	defer cleanup()

	tests := buildTests()
	input := buildInput(tests)

	oracleOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Printf("oracle runtime error: %v\ninput:\n%s", err, input)
		return
	}
	oracleAns, err := parseOutputs(oracleOut, tests)
	if err != nil {
		fmt.Printf("oracle output parse error: %v\noutput:\n%s", err, oracleOut)
		return
	}

	candOut, err := runBinary(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\ninput:\n%s", err, input)
		return
	}
	candAns, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for idx, tc := range tests {
		expectedSet := make(map[int]struct{})
		for _, t := range oracleAns[idx] {
			expectedSet[t.v] = struct{}{}
		}
		if len(expectedSet) != len(oracleAns[idx]) {
			fmt.Printf("oracle reported duplicate submedian in test %d\n", idx+1)
			return
		}

		gotSet := make(map[int]triple)
		for _, t := range candAns[idx] {
			if _, ok := gotSet[t.v]; ok {
				fmt.Printf("test %d: duplicate value %d in candidate output\n", idx+1, t.v)
				return
			}
			gotSet[t.v] = t
		}

		if len(gotSet) != len(expectedSet) {
			fmt.Printf("test %d: expected %d submedians, got %d\ninput:\n%s", idx+1, len(expectedSet), len(gotSet), input)
			return
		}
		for v := range expectedSet {
			if _, ok := gotSet[v]; !ok {
				fmt.Printf("test %d: missing submedian %d\ninput:\n%s", idx+1, v, input)
				return
			}
		}

		roots, nodes := buildPersistent(tc.a, tc.n)
		for _, t := range candAns[idx] {
			if t.v < 1 || t.v > tc.n {
				fmt.Printf("test %d: value %d out of range [1,%d]\n", idx+1, t.v, tc.n)
				return
			}
			if t.l < 1 || t.r < t.l || t.r > tc.n {
				fmt.Printf("test %d: invalid interval (%d,%d)\n", idx+1, t.l, t.r)
				return
			}
			if t.r-t.l+1 < tc.k {
				fmt.Printf("test %d: interval too short for k, value %d\n", idx+1, t.v)
				return
			}
			length := t.r - t.l + 1
			lowIdx := (length + 1) / 2
			highIdx := length/2 + 1
			lowVal := kth(nodes, roots[t.l-1], roots[t.r], 1, tc.n, lowIdx)
			highVal := kth(nodes, roots[t.l-1], roots[t.r], 1, tc.n, highIdx)
			if t.v < lowVal || t.v > highVal {
				fmt.Printf("test %d: value %d is not a median of subarray [%d,%d] (median range [%d,%d])\n", idx+1, t.v, t.l, t.r, lowVal, highVal)
				return
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2128E2-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2128E2")
	cmd := exec.Command("go", "build", "-o", outPath, "2128E2.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, string(out))
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return outPath, cleanup, nil
}

func runBinary(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, tests []testCase) ([][]triple, error) {
	fields := strings.Fields(out)
	pos := 0
	ans := make([][]triple, len(tests))
	for idx, tc := range tests {
		if pos >= len(fields) {
			return nil, fmt.Errorf("test %d: missing count", idx+1)
		}
		c, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid count: %v", idx+1, err)
		}
		pos++
		if c < 0 || c > tc.n {
			return nil, fmt.Errorf("test %d: invalid count %d", idx+1, c)
		}
		if pos+3*c > len(fields) {
			return nil, fmt.Errorf("test %d: not enough tokens for triples", idx+1)
		}
		cur := make([]triple, c)
		for i := 0; i < c; i++ {
			v, err1 := strconv.Atoi(fields[pos])
			l, err2 := strconv.Atoi(fields[pos+1])
			r, err3 := strconv.Atoi(fields[pos+2])
			if err := firstErr(err1, err2, err3); err != nil {
				return nil, fmt.Errorf("test %d, triple %d: %v", idx+1, i+1, err)
			}
			cur[i] = triple{v: v, l: l, r: r}
			pos += 3
		}
		ans[idx] = cur
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extraneous tokens in output")
	}
	return ans, nil
}

func firstErr(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(256 * len(tests))
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.k))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTests() []testCase {
	var tests []testCase
	rng := rand.New(rand.NewSource(randSeed))
	totalN := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Deterministic small cases
	add(testCase{n: 1, k: 1, a: []int{1}})
	add(testCase{n: 2, k: 1, a: []int{2, 1}})
	add(testCase{n: 3, k: 2, a: []int{1, 2, 3}})
	add(testCase{n: 4, k: 3, a: []int{4, 1, 3, 2}})
	add(testCase{n: 5, k: 2, a: []int{2, 2, 2, 2, 2}})

	// Random cases
	for len(tests) < maxTests && totalN < 300000 {
		remain := 300000 - totalN
		maxN := 30000
		if remain < maxN {
			maxN = remain
		}
		n := rng.Intn(maxN) + 1
		k := rng.Intn(n) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(n) + 1
		}
		add(testCase{n: n, k: k, a: a})
	}
	return tests
}

func buildPersistent(a []int, maxVal int) ([]int, []node) {
	nodes := make([]node, 1, (len(a)+1)*20)
	roots := make([]int, len(a)+1)
	for i, v := range a {
		roots[i+1] = update(&nodes, roots[i], 1, maxVal, v, 1)
	}
	return roots, nodes
}

func update(nodes *[]node, root, l, r, pos, delta int) int {
	newRoot := len(*nodes)
	*nodes = append(*nodes, node{})
	(*nodes)[newRoot] = (*nodes)[root]
	(*nodes)[newRoot].sum += delta
	if l != r {
		mid := (l + r) >> 1
		if pos <= mid {
			child := update(nodes, (*nodes)[newRoot].left, l, mid, pos, delta)
			(*nodes)[newRoot].left = child
		} else {
			child := update(nodes, (*nodes)[newRoot].right, mid+1, r, pos, delta)
			(*nodes)[newRoot].right = child
		}
	}
	return newRoot
}

func kth(nodes []node, leftRoot, rightRoot, l, r, k int) int {
	if l == r {
		return l
	}
	mid := (l + r) >> 1
	leftCount := nodes[nodes[rightRoot].left].sum - nodes[nodes[leftRoot].left].sum
	if k <= leftCount {
		return kth(nodes, nodes[leftRoot].left, nodes[rightRoot].left, l, mid, k)
	}
	return kth(nodes, nodes[leftRoot].right, nodes[rightRoot].right, mid+1, r, k-leftCount)
}
