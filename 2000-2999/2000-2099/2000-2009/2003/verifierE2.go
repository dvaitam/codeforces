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

type interval struct {
	l int
	r int
}

type singleCase struct {
	n   int
	seg []interval
}

type testCase struct {
	name    string
	input   string
	answers int
}

func main() {
	candidate, err := candidatePathFromArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	refSrc := referencePath()
	refBin, cleanup, err := buildReferenceBinary(refSrc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("reference failed on case %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		expected, err := parseAnswers(refOut, tc.answers)
		if err != nil {
			fmt.Printf("reference output invalid on case %d (%s): %v\nraw output:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Printf("case %d (%s): runtime error: %v\n", idx+1, tc.name, err)
			fmt.Println(previewInput(tc.input))
			os.Exit(1)
		}
		got, err := parseAnswers(candOut, tc.answers)
		if err != nil {
			fmt.Printf("case %d (%s): invalid output: %v\nraw output:\n%s\n", idx+1, tc.name, err, candOut)
			fmt.Println(previewInput(tc.input))
			os.Exit(1)
		}
		for i := 0; i < tc.answers; i++ {
			if got[i] != expected[i] {
				fmt.Printf("case %d (%s) failed at test %d: expected %d got %d\n", idx+1, tc.name, i+1, expected[i], got[i])
				fmt.Println(previewInput(tc.input))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func referencePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "2000-2999/2000-2099/2000-2009/2003/2003E2.go"
	}
	return filepath.Join(filepath.Dir(file), "2003E2.go")
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2003E2")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2003e2")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference solution: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseAnswers(out string, expected int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = v
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{}
	tests = append(tests, sampleTest())
	tests = append(tests, customTest("single-case-no-intervals", []singleCase{{n: 5}}))
	tests = append(tests, customTest("identical-intervals", []singleCase{{n: 6, seg: []interval{{l: 1, r: 6}, {l: 1, r: 6}, {l: 2, r: 5}}}}))
	tests = append(tests, customTest("dense-small", []singleCase{{n: 8, seg: fullIntervals(8)}}))

	rng := rand.New(rand.NewSource(0x6f3c2b91))
	tests = append(tests, randomBatch("random-small", rng, 200, 200))
	tests = append(tests, randomBatch("random-medium", rng, 1200, 1200))
	tests = append(tests, randomBatch("random-large", rng, 5000, 5000))

	// Add a few targeted stress cases
	tests = append(tests, customTest("chain-intervals", []singleCase{{n: 10, seg: chainIntervals(10)}}))
	tests = append(tests, customTest("zero-m", []singleCase{{n: 2}, {n: 9, seg: []interval{}}}))

	return tests
}

func sampleTest() testCase {
	cases := []singleCase{
		{n: 2, seg: []interval{}},
		{n: 2, seg: []interval{{l: 1, r: 2}}},
		{n: 5, seg: []interval{{l: 1, r: 2}}},
		{n: 8, seg: []interval{{l: 1, r: 4}, {l: 2, r: 5}, {l: 7, r: 8}}},
		{n: 7, seg: []interval{{l: 1, r: 4}, {l: 4, r: 7}}},
		{n: 7, seg: []interval{{l: 1, r: 2}, {l: 1, r: 7}, {l: 3, r: 7}}},
	}
	return customTest("sample-like", cases)
}

func customTest(name string, cases []singleCase) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		m := len(cs.seg)
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, m))
		for _, seg := range cs.seg {
			sb.WriteString(fmt.Sprintf("%d %d\n", seg.l, seg.r))
		}
	}
	return testCase{name: name, input: sb.String(), answers: len(cases)}
}

func fullIntervals(n int) []interval {
	segs := make([]interval, 0)
	for l := 1; l < n; l++ {
		segs = append(segs, interval{l: l, r: n})
	}
	return segs
}

func chainIntervals(n int) []interval {
	segs := make([]interval, 0, n-1)
	for i := 1; i < n; i++ {
		segs = append(segs, interval{l: i, r: i + 1})
	}
	return segs
}

func randomBatch(name string, rng *rand.Rand, totalN, totalM int) testCase {
	cases := make([]singleCase, 0)
	usedN, usedM := 0, 0
	for usedN < totalN {
		remainingN := totalN - usedN
		if remainingN < 2 {
			break
		}
		maxStep := min(600, remainingN)
		if maxStep < 2 {
			maxStep = remainingN
		}
		n := 2
		if maxStep > 2 {
			n += rng.Intn(maxStep - 1)
		}
		if n > remainingN {
			n = remainingN
		}
		maxPossibleM := n * (n - 1) / 2
		remainingM := totalM - usedM
		if remainingM < 0 {
			remainingM = 0
		}
		maxM := min(remainingM, maxPossibleM)
		m := 0
		if maxM > 0 {
			m = rng.Intn(maxM + 1)
		}
		intervals := make([]interval, m)
		for i := 0; i < m; i++ {
			l := 1 + rng.Intn(n-1)
			r := l + 1 + rng.Intn(n-l)
			intervals[i] = interval{l: l, r: r}
		}
		cases = append(cases, singleCase{n: n, seg: intervals})
		usedN += n
		usedM += m
	}
	if len(cases) == 0 {
		cases = append(cases, singleCase{n: max(2, totalN), seg: []interval{}})
	}
	return customTest(name, cases)
}

func previewInput(in string) string {
	const limit = 400
	if len(in) <= limit {
		return "Input:\n" + in
	}
	return fmt.Sprintf("Input (first %d chars):\n%s...\n", limit, in[:limit])
}

func candidatePathFromArgs() (string, error) {
	if len(os.Args) == 2 {
		return os.Args[1], nil
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], nil
	}
	return "", fmt.Errorf("usage: go run verifierE2.go /path/to/solution")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
