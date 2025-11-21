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

const refSource = "2103E.go"

type singleCase struct {
	n int
	k int64
	a []int64
}

type operation struct {
	i int
	j int
	x int64
}

type testCase struct {
	name  string
	input string
	cases []singleCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate_binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refPossible, err := parseReference(tc.cases, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := verifyCandidate(tc.cases, refPossible, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n", idx+1, tc.name, err, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "oracle-2103E-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleE")

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
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

func runProgram(bin, input string) (string, error) {
	cmd := commandFor(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseReference(cases []singleCase, output string) ([]bool, error) {
	tokens := strings.Fields(output)
	pos := 0
	res := make([]bool, 0, len(cases))
	for idx, cs := range cases {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("missing output for case %d", idx+1)
		}
		if tokens[pos] == "-1" {
			res = append(res, false)
			pos++
			continue
		}
		m, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return nil, fmt.Errorf("invalid op count for case %d: %v", idx+1, err)
		}
		pos++
		need := 3 * m
		if pos+need > len(tokens) {
			return nil, fmt.Errorf("not enough operation tokens for case %d", idx+1)
		}
		pos += need
		res = append(res, true)
		if m > 3*cs.n {
			return nil, fmt.Errorf("reference used too many operations (%d) for case %d", m, idx+1)
		}
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("extra tokens in reference output")
	}
	return res, nil
}

func verifyCandidate(cases []singleCase, refPossible []bool, output string) error {
	tokens := strings.Fields(output)
	pos := 0
	for idx, cs := range cases {
		if pos >= len(tokens) {
			return fmt.Errorf("missing output for case %d", idx+1)
		}
		if strings.EqualFold(tokens[pos], "-1") {
			if refPossible[idx] {
				return fmt.Errorf("case %d declared impossible but reference has a solution", idx+1)
			}
			pos++
			continue
		}
		m, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return fmt.Errorf("invalid operation count for case %d: %v", idx+1, err)
		}
		pos++
		if m < 0 || m > 3*cs.n {
			return fmt.Errorf("operation count %d out of bounds for case %d", m, idx+1)
		}
		need := 3 * m
		if pos+need > len(tokens) {
			return fmt.Errorf("not enough tokens for operations in case %d", idx+1)
		}
		ops := make([]operation, m)
		for i := 0; i < m; i++ {
			ii, err1 := strconv.Atoi(tokens[pos])
			jj, err2 := strconv.Atoi(tokens[pos+1])
			xv, err3 := strconv.ParseInt(tokens[pos+2], 10, 64)
			if err1 != nil || err2 != nil || err3 != nil {
				return fmt.Errorf("invalid operation #%d in case %d", i+1, idx+1)
			}
			pos += 3
			ops[i] = operation{i: ii, j: jj, x: xv}
		}
		if err := simulate(cs, ops); err != nil {
			return fmt.Errorf("case %d invalid operations: %v", idx+1, err)
		}
	}
	if pos != len(tokens) {
		return fmt.Errorf("extra tokens in candidate output")
	}
	return nil
}

func simulate(cs singleCase, ops []operation) error {
	a := append([]int64(nil), cs.a...)
	n := cs.n
	k := cs.k
	for idx, op := range ops {
		i := op.i - 1
		j := op.j - 1
		if i < 0 || i >= n || j < 0 || j >= n {
			return fmt.Errorf("operation %d uses out-of-range indices", idx+1)
		}
		if i == j {
			return fmt.Errorf("operation %d uses equal indices", idx+1)
		}
		if a[i]+a[j] != k {
			return fmt.Errorf("operation %d violated sum condition: a[%d]+a[%d]=%d, k=%d", idx+1, i+1, j+1, a[i]+a[j], k)
		}
		x := op.x
		if x < -a[j] || x > a[i] {
			return fmt.Errorf("operation %d uses invalid x=%d for a[%d]=%d a[%d]=%d", idx+1, x, i+1, a[i], j+1, a[j])
		}
		a[i] -= x
		a[j] += x
		if a[i] < 0 || a[i] > k || a[j] < 0 || a[j] > k {
			return fmt.Errorf("operation %d leaves array out of bounds", idx+1)
		}
	}
	for i := 0; i+1 < n; i++ {
		if a[i] > a[i+1] {
			return fmt.Errorf("array not non-decreasing after operations")
		}
	}
	return nil
}

func buildTests() []testCase {
	tests := []testCase{
		sampleTests(),
		impossibleSmall(),
		alreadySorted(),
		edgeValues(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTests() testCase {
	cases := []singleCase{
		{n: 5, k: 10, a: []int64{1, 2, 3, 4, 5}},
		{n: 5, k: 6, a: []int64{1, 2, 3, 5, 4}},
		{n: 5, k: 7, a: []int64{7, 1, 5, 3, 1}},
		{n: 10, k: 10, a: []int64{2, 5, 3, 2, 7, 3, 1, 10, 4, 0}},
	}
	return buildTest("sample", cases)
}

func impossibleSmall() testCase {
	// k is large enough that no pair sums to k; array is not sorted.
	cases := []singleCase{
		{n: 4, k: 3, a: []int64{1, 0, 1, 0}},
		{n: 6, k: 5, a: []int64{4, 0, 0, 0, 0, 0}},
	}
	return buildTest("impossible", cases)
}

func alreadySorted() testCase {
	cases := []singleCase{
		{n: 6, k: 8, a: []int64{0, 1, 2, 3, 4, 8}},
		{n: 5, k: 4, a: []int64{0, 0, 1, 3, 4}},
	}
	return buildTest("sorted", cases)
}

func edgeValues() testCase {
	cases := []singleCase{
		{n: 6, k: 1, a: []int64{0, 1, 0, 1, 0, 1}},
		{n: 6, k: 1000000000, a: []int64{0, 1000000000, 0, 1000000000, 500000000, 500000000}},
	}
	return buildTest("edge_values", cases)
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCnt := rng.Intn(5) + 1
	cases := make([]singleCase, caseCnt)
	for i := 0; i < caseCnt; i++ {
		n := rng.Intn(30) + 4
		k := rng.Int63n(50) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Int63n(k + 1)
		}
		// Occasionally make array sorted to allow zero ops as valid.
		if rng.Intn(4) == 0 {
			for j := 1; j < n; j++ {
				if arr[j] < arr[j-1] {
					arr[j] = arr[j-1]
				}
			}
		}
		cases[i] = singleCase{n: n, k: k, a: arr}
	}
	return buildTest(fmt.Sprintf("random_%d", idx), cases)
}

func buildTest(name string, cases []singleCase) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&sb, "%d %d\n", cs.n, cs.k)
		for i, v := range cs.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:  name,
		input: sb.String(),
		cases: cases,
	}
}
