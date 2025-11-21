package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	refSource        = "761D.go"
	tempOraclePrefix = "oracle-761D-"
	randomTestsCount = 120
	maxRandomN       = 2000
)

type testCase struct {
	name string
	n    int
	l    int64
	r    int64
	a    []int64
	p    []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(randomTestsCount, rng)...)
	tests = append(tests, largeTests()...)

	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		refResult := strings.TrimSpace(refOut)

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		candResult := strings.TrimSpace(candOut)

		if refResult == "-1" {
			if candResult != "-1" {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected -1, got %q\n", idx+1, tc.name, candResult)
				fmt.Println("Input:")
				fmt.Print(input)
				os.Exit(1)
			}
			continue
		}

		if candResult == "-1" {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: candidate reported -1 but solution exists\n", idx+1, tc.name)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		bVals, err := parseSequence(candResult, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if err := validateSequence(tc, bVals); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Candidate output:")
			fmt.Print(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseSequence(out string, n int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	res := make([]int64, n)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func validateSequence(tc testCase, b []int64) error {
	if len(b) != tc.n {
		return fmt.Errorf("expected %d elements, got %d", tc.n, len(b))
	}
	minVal, maxVal := b[0], b[0]
	for _, v := range b {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	if minVal < tc.l || maxVal > tc.r {
		return fmt.Errorf("values out of range [%d,%d]", tc.l, tc.r)
	}

	c := make([]int64, tc.n)
	for i := 0; i < tc.n; i++ {
		c[i] = b[i] - tc.a[i]
	}

	if hasDuplicates(c) {
		return fmt.Errorf("sequence c contains duplicate values")
	}

	compressed := compressRanks(c)
	for i := 0; i < tc.n; i++ {
		if compressed[i] != tc.p[i] {
			return fmt.Errorf("compressed value mismatch at position %d: expected %d got %d", i+1, tc.p[i], compressed[i])
		}
	}
	return nil
}

func hasDuplicates(arr []int64) bool {
	set := make(map[int64]struct{}, len(arr))
	for _, v := range arr {
		if _, ok := set[v]; ok {
			return true
		}
		set[v] = struct{}{}
	}
	return false
}

func compressRanks(arr []int64) []int {
	type pair struct {
		val int64
		idx int
	}
	n := len(arr)
	ps := make([]pair, n)
	for i, v := range arr {
		ps[i] = pair{val: v, idx: i}
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].val == ps[j].val {
			return ps[i].idx < ps[j].idx
		}
		return ps[i].val < ps[j].val
	})
	res := make([]int, n)
	for rank, p := range ps {
		res[p.idx] = rank + 1
	}
	return res
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.l, tc.r)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "simple_valid", n: 3, l: 1, r: 10, a: []int64{2, 4, 6}, p: []int{1, 2, 3}},
		{name: "simple_shift", n: 4, l: 1, r: 10, a: []int64{1, 2, 3, 4}, p: []int{4, 3, 2, 1}},
		{name: "already_sorted", n: 5, l: 5, r: 15, a: []int64{5, 6, 7, 8, 9}, p: []int{1, 2, 3, 4, 5}},
		{name: "reverse", n: 5, l: 5, r: 20, a: []int64{10, 9, 8, 7, 6}, p: []int{5, 4, 3, 2, 1}},
		{name: "edge_l_r", n: 3, l: 100, r: 105, a: []int64{100, 102, 104}, p: []int{3, 1, 2}},
	}
}

func randomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxRandomN-1) + 1
		l := rng.Int63n(1000) + 1
		r := l + rng.Int63n(1000)
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = l + rng.Int63n(r-l+1)
		}
		p := randomPermutation(n, rng)
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			n:    n,
			l:    l,
			r:    r,
			a:    a,
			p:    p,
		})
	}
	return tests
}

func largeTests() []testCase {
	n := 100000
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(i + 1)
	}
	pInc := make([]int, n)
	pDec := make([]int, n)
	for i := 0; i < n; i++ {
		pInc[i] = i + 1
		pDec[i] = n - i
	}
	return []testCase{
		{name: "large_inc", n: n, l: 1, r: int64(2 * n), a: a, p: pInc},
		{name: "large_dec", n: n, l: 1, r: int64(2 * n), a: a, p: pDec},
	}
}

func randomPermutation(n int, rng *rand.Rand) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	return perm
}
