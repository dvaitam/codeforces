package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	limit           = 1_000_000_000
	randomTests     = 200
	maxRandA        = 900000
	maxInputBValue  = 1_000_000
	largeCaseLength = 10_000
)

type caseData struct {
	n       int
	a       []int
	b       []int
	countsB map[int]int
	maxB    int
}

type testInput struct {
	input string
	cases []caseData
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("failed to prepare contestant binary:", err)
		return
	}
	if candCleanup != nil {
		defer candCleanup()
	}

	oracle, oracleCleanup, err := prepareOracle()
	if err != nil {
		fmt.Println("failed to prepare reference solution:", err)
		return
	}
	defer oracleCleanup()

	tests := deterministicTests()
	total := 0
	for idx, test := range tests {
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("deterministic test %d failed: %v\ninput:\n%s", idx+1, err, test.input)
			return
		}
		total++
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTests; i++ {
		test := randomTest(rng)
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("random test %d failed: %v\ninput:\n%s", i+1, err, test.input)
			return
		}
		total++
	}

	fmt.Printf("All %d tests passed.\n", total)
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("candidate2085E_%d", time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, abs)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, nil, nil
}

func prepareOracle() (string, func(), error) {
	dir := sourceDir()
	src := filepath.Join(dir, "2085E.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2085E_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func runTest(test testInput, candidate, oracle string) error {
	candOut, err := runBinary(candidate, test.input)
	if err != nil {
		return fmt.Errorf("contestant runtime error: %v", err)
	}
	oracleOut, err := runBinary(oracle, test.input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}

	t := len(test.cases)
	candVals, err := parseOutputs(candOut, t)
	if err != nil {
		return fmt.Errorf("failed to parse contestant output: %v", err)
	}
	oracleVals, err := parseOutputs(oracleOut, t)
	if err != nil {
		return fmt.Errorf("failed to parse oracle output: %v", err)
	}

	for idx, c := range test.cases {
		kCand := candVals[idx]
		kOracle := oracleVals[idx]
		if kCand == -1 {
			if kOracle != -1 {
				return fmt.Errorf("case %d: contestant answered -1 but oracle found %d", idx+1, kOracle)
			}
			continue
		}
		if !checkCandidateCase(c, kCand) {
			return fmt.Errorf("case %d: answer %d is invalid", idx+1, kCand)
		}
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutputs(output string, expected int) ([]int, error) {
	reader := strings.NewReader(output)
	res := make([]int, 0, expected)
	for len(res) < expected {
		var val int
		if _, err := fmt.Fscan(reader, &val); err != nil {
			return nil, fmt.Errorf("need %d integers, got %d (%v)", expected, len(res), err)
		}
		res = append(res, val)
	}
	var extra int
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("expected %d integers, output has extra data", expected)
	}
	return res, nil
}

func checkCandidateCase(c caseData, k int) bool {
	if k < 1 || k > limit {
		return false
	}
	if k <= c.maxB {
		return false
	}
	temp := make(map[int]int, len(c.countsB))
	for _, val := range c.a {
		r := val % k
		need, ok := c.countsB[r]
		if !ok {
			return false
		}
		temp[r]++
		if temp[r] > need {
			return false
		}
	}
	for key, need := range c.countsB {
		if temp[key] != need {
			return false
		}
	}
	return true
}

func deterministicTests() []testInput {
	return []testInput{
		sampleStyleTest(),
		allZeroTest(),
		largeTest(),
	}
}

func sampleStyleTest() testInput {
	cases := []caseData{
		newCaseData(
			[]int{3, 5, 2, 7},
			[]int{1, 1, 0, 1},
		),
		newCaseData(
			[]int{4, 1, 3, 2, 5, 6},
			[]int{6, 5, 4, 3, 2, 1},
		),
		newCaseData(
			[]int{0, 0, 0},
			[]int{1, 0, 0},
		),
	}
	return buildTestInput(cases)
}

func allZeroTest() testInput {
	cases := []caseData{
		newCaseData(
			[]int{0, 2, 4, 6},
			[]int{0, 0, 0, 0},
		),
		newCaseData(
			[]int{5, 5, 5},
			[]int{5, 5, 6},
		),
	}
	return buildTestInput(cases)
}

func largeTest() testInput {
	n := largeCaseLength
	k := 721
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		r := i % k
		q := i % 900
		a[i] = q*k + r
		if a[i] > maxRandA {
			a[i] = maxRandA - (i % 500)
		}
		b[i] = r
	}
	rng := rand.New(rand.NewSource(42))
	rng.Shuffle(n, func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})
	return buildTestInput([]caseData{newCaseData(a, b)})
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(4) + 1
	cases := make([]caseData, t)
	for i := 0; i < t; i++ {
		if rng.Intn(3) == 0 {
			cases[i] = randomInvalidCase(rng)
		} else {
			cases[i] = randomValidCase(rng)
		}
	}
	return buildTestInput(cases)
}

func randomValidCase(rng *rand.Rand) caseData {
	a, b := randomValidArrays(rng)
	return newCaseData(a, b)
}

func randomInvalidCase(rng *rand.Rand) caseData {
	a, b := randomValidArrays(rng)
	maxA := 0
	for _, val := range a {
		if val > maxA {
			maxA = val
		}
	}
	invalidVal := maxA + 1
	if invalidVal > maxInputBValue {
		invalidVal = maxInputBValue
	}
	b[0] = invalidVal
	return newCaseData(a, b)
}

func randomValidArrays(rng *rand.Rand) ([]int, []int) {
	n := rng.Intn(8) + 1
	k := rng.Intn(maxRandA) + 1
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		var r int
		if k == 1 {
			r = 0
		} else {
			r = rng.Intn(k)
		}
		maxQ := 0
		if k > 0 {
			maxQ = (maxRandA - r) / k
		}
		q := 0
		if maxQ > 0 {
			q = rng.Intn(maxQ + 1)
		}
		a[i] = q*k + r
		b[i] = r
	}
	rng.Shuffle(n, func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})
	return a, b
}

func newCaseData(a, b []int) caseData {
	if len(a) != len(b) {
		panic("a and b must have same length")
	}
	aCopy := append([]int(nil), a...)
	bCopy := append([]int(nil), b...)
	counts := make(map[int]int, len(bCopy))
	maxB := 0
	for _, val := range bCopy {
		counts[val]++
		if val > maxB {
			maxB = val
		}
	}
	return caseData{
		n:       len(aCopy),
		a:       aCopy,
		b:       bCopy,
		countsB: counts,
		maxB:    maxB,
	}
}

func buildTestInput(cases []caseData) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d\n", c.n)
		writeArray(&sb, c.a)
		writeArray(&sb, c.b)
	}
	return testInput{input: sb.String(), cases: cases}
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, val := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(sb, "%d", val)
	}
	sb.WriteByte('\n')
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
