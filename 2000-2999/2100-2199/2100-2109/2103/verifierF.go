package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	k int
	a []int
}

type testRun struct {
	input string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	targetPath := os.Args[1]

	refBin, refCleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := prepareCandidate(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	seed := time.Now().UnixNano()
	tests := generateTests(seed)

	for i, tr := range tests {
		expRaw, err := runBinary(refBin, tr.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tr.input)
			os.Exit(1)
		}
		expAns, err := parseOutputs(expRaw, tr.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, expRaw)
			os.Exit(1)
		}

		actRaw, err := runBinary(candBin, tr.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i+1, err, tr.input)
			os.Exit(1)
		}
		actAns, err := parseOutputs(actRaw, tr.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tr.input, actRaw)
			os.Exit(1)
		}

		if len(expAns) != len(actAns) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d test cases, got %d\ninput:\n%sreference:\n%v\ncandidate:\n%v\n",
				i+1, len(expAns), len(actAns), tr.input, expAns, actAns)
			os.Exit(1)
		}
		for tcIdx := range expAns {
			if len(expAns[tcIdx]) != len(actAns[tcIdx]) {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected %d numbers, got %d\ninput:\n%sreference:\n%v\ncandidate:\n%v\n",
					i+1, tcIdx+1, len(expAns[tcIdx]), len(actAns[tcIdx]), tr.input, expAns, actAns)
				os.Exit(1)
			}
			for pos := range expAns[tcIdx] {
				if expAns[tcIdx][pos] != actAns[tcIdx][pos] {
					fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d position %d: expected %d, got %d\ninput:\n%sreference:\n%v\ncandidate:\n%v\n",
						i+1, tcIdx+1, pos+1, expAns[tcIdx][pos], actAns[tcIdx][pos], tr.input, expAns, actAns)
					os.Exit(1)
				}
			}
		}
	}

	fmt.Printf("All %d tests passed (seed %d).\n", len(tests), seed)
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("unable to determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "verifier-2103F-ref-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref2103F")
	cmd := exec.Command("go", "build", "-o", outPath, "2103F.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, out)
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	dir := filepath.Dir(abs)
	tmpDir, err := os.MkdirTemp("", "verifier-2103F-cand-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "candidate2103F")
	cmd := exec.Command("go", "build", "-o", outPath, abs)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build candidate: %v\n%s", err, out)
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, tests []testCase) ([][]int, error) {
	res := make([][]int, len(tests))
	r := bufio.NewReader(strings.NewReader(out))
	for idx, tc := range tests {
		ans := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(r, &ans[i]); err != nil {
				return nil, fmt.Errorf("test case %d: failed to read value %d/%d: %v", idx+1, i+1, tc.n, err)
			}
		}
		res[idx] = ans
	}
	var extra string
	if _, err := fmt.Fscan(r, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected: %q", extra)
	} else if err != nil && err != io.EOF {
		return nil, err
	}
	return res, nil
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
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

func deterministicRuns() []testRun {
	var runs []testRun

	runs = append(runs, makeRun([]testCase{
		{n: 2, k: 2, a: []int{1, 3}},
		{n: 5, k: 3, a: []int{1, 7, 4, 6, 2}},
	}))

	runs = append(runs, makeRun([]testCase{
		{n: 1, k: 1, a: []int{0}},
		{n: 1, k: 5, a: []int{31}},
		{n: 4, k: 4, a: []int{0, 0, 0, 0}},
		{n: 4, k: 4, a: []int{15, 15, 15, 15}},
	}))

	runs = append(runs, makeRun([]testCase{
		{n: 3, k: 3, a: []int{1, 2, 4}},
		{n: 6, k: 5, a: []int{1, 31, 0, 31, 0, 31}},
	}))

	runs = append(runs, makeRun([]testCase{
		{n: 8, k: 6, a: []int{1, 2, 3, 4, 5, 6, 7, 8}},
	}))

	return runs
}

func randomCase(rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN) + 1
	k := rng.Intn(17) + 1
	limit := 1 << k
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(limit)
	}
	return testCase{n: n, k: k, a: a}
}

func randomRuns(rng *rand.Rand) []testRun {
	var runs []testRun

	for i := 0; i < 30; i++ {
		t := rng.Intn(4) + 1
		maxN := []int{5, 20, 60}[rng.Intn(3)]
		cases := make([]testCase, t)
		for j := 0; j < t; j++ {
			cases[j] = randomCase(rng, maxN)
		}
		runs = append(runs, makeRun(cases))
	}

	for i := 0; i < 5; i++ {
		t := rng.Intn(3) + 1
		cases := make([]testCase, t)
		for j := 0; j < t; j++ {
			cases[j] = randomCase(rng, 500)
		}
		runs = append(runs, makeRun(cases))
	}

	for i := 0; i < 3; i++ {
		n := rng.Intn(40000) + 60000
		k := rng.Intn(17) + 1
		limit := 1 << k
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(limit)
		}
		runs = append(runs, makeRun([]testCase{{n: n, k: k, a: a}}))
	}

	return runs
}

func makeRun(cases []testCase) testRun {
	return testRun{
		input: buildInput(cases),
		cases: cases,
	}
}

func generateTests(seed int64) []testRun {
	rng := rand.New(rand.NewSource(seed))
	tests := deterministicRuns()
	tests = append(tests, randomRuns(rng)...)
	return tests
}
