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
	modValue          = 1_000_000_007
	randomTests       = 200
	maxN              = 300000
	maxKValue   int64 = 1_000_000_000_000_000_000
)

type caseData struct {
	n int
	k int64
}

type testInput struct {
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("candidate2072G_%d", time.Now().UnixNano()))
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
	src := filepath.Join(dir, "2072G.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2072G_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func runTest(test testInput, candidate, oracle string) error {
	oracleOut, err := runBinary(oracle, test.input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	candOut, err := runBinary(candidate, test.input)
	if err != nil {
		return fmt.Errorf("contestant runtime error: %v", err)
	}

	expect, err := parseOutputs(oracleOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse oracle output: %v", err)
	}
	got, err := parseOutputs(candOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse contestant output: %v", err)
	}

	for i := 0; i < test.t; i++ {
		if expect[i] != got[i] {
			return fmt.Errorf("case %d: expected %d got %d", i+1, expect[i], got[i])
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

func parseOutputs(output string, expected int) ([]int64, error) {
	reader := strings.NewReader(output)
	res := make([]int64, 0, expected)
	for len(res) < expected {
		var val int64
		if _, err := fmt.Fscan(reader, &val); err != nil {
			return nil, fmt.Errorf("need %d integers, got %d (%v)", expected, len(res), err)
		}
		val %= modValue
		if val < 0 {
			val += modValue
		}
		res = append(res, val)
	}
	var extra int64
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("expected %d integers, output has extra data", expected)
	}
	return res, nil
}

func deterministicTests() []testInput {
	tests := make([]testInput, 0, 4)
	tests = append(tests, sampleTest())
	tests = append(tests, buildTestInput([]caseData{
		{n: 1, k: 2},
		{n: 1, k: 10},
		{n: maxN, k: 2},
		{n: maxN, k: maxKValue},
	}))
	tests = append(tests, buildTestInput([]caseData{
		{n: 99991, k: 99991},
		{n: 54321, k: 1_000_000_000},
		{n: 2, k: 2},
		{n: 2, k: maxKValue},
	}))
	tests = append(tests, buildTestInput(fullRangeCases()))
	return tests
}

func sampleTest() testInput {
	cases := []caseData{
		{n: 123, k: 242},
		{n: 42, k: 521},
		{n: 1, k: 10},
		{n: 4, k: 4},
		{n: 16, k: 269},
		{n: 699, k: 319},
		{n: 8, k: 49982},
		{n: 44353, k: 1000000},
		{n: 100000, k: 1000000007},
		{n: 17, k: 1_000_000_000_000_000_000},
	}
	return buildTestInput(cases)
}

func fullRangeCases() []caseData {
	res := make([]caseData, 0, 50)
	steps := []int{1, 2, 3, 5, 7, 11, 13, 17}
	for i := 0; i < len(steps) && len(res) < 40; i++ {
		val := steps[i]
		res = append(res, caseData{n: val, k: int64(val + 1)})
	}
	res = append(res,
		caseData{n: maxN, k: int64(maxN)},
		caseData{n: maxN - 1, k: int64(maxN + 10)},
		caseData{n: 314159, k: 271828},
		caseData{n: 271828, k: 314159},
		caseData{n: 99999, k: maxKValue},
	)
	return res
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(8) + 1
	cases := make([]caseData, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN) + 1
		k := randK(rng)
		cases[i] = caseData{n: n, k: k}
	}
	return buildTestInput(cases)
}

func randK(rng *rand.Rand) int64 {
	switch rng.Intn(4) {
	case 0:
		return int64(rng.Intn(maxN)) + 2
	case 1:
		return int64(rng.Intn(1_000_000_000)) + 2
	case 2:
		return int64(rng.Intn(maxN)+1) * int64(rng.Intn(1000)+1)
	default:
		return rng.Int63n(maxKValue-1) + 2
	}
}

func buildTestInput(cases []caseData) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d %d\n", c.n, c.k)
	}
	return testInput{
		input: sb.String(),
		t:     len(cases),
	}
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
