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

const (
	refSource1998E2 = "1000-1999/1900-1999/1990-1999/1998/1998E2.go"
	maxValue        = int64(1_000_000_000)
)

type testCase struct {
	n int
	a []int64
}

type testSet struct {
	name  string
	cases []testCase
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource1998E2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, ts := range tests {
		input := buildInput(ts.cases)

		refOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d (%s): %v\ninput:\n%s", idx+1, ts.name, err, input)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut, ts.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error: failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, ts.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, ts.name, err, input)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut, ts.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: cannot parse output: %v\ninput:\n%soutput:\n%s", idx+1, ts.name, err, input, candOut)
			os.Exit(1)
		}

		if err := compareAnswers(refAns, candAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%sreference:\n%soutput:\n%s", idx+1, ts.name, err, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier1998E2-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
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

func buildTests() []testSet {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return []testSet{
		sampleTest(),
		edgeCasesTest(),
		randomTest("random_small", rng, 6, 1, 15, 50),
		randomTest("random_medium", rng, 5, 40, 120, maxValue),
		wavePatternTest(),
		largeRandomTest(rng),
	}
}

func sampleTest() testSet {
	return testSet{
		name: "sample",
		cases: []testCase{
			{n: 5, a: []int64{1, 2, 3, 2, 1}},
			{n: 7, a: []int64{4, 5, 1, 2, 1, 4, 5}},
			{n: 11, a: []int64{1, 2, 3, 1, 1, 9, 3, 2, 4, 1, 3}},
		},
	}
}

func edgeCasesTest() testSet {
	return testSet{
		name: "edges",
		cases: []testCase{
			{n: 1, a: []int64{1}},
			{n: 2, a: []int64{maxValue, maxValue}},
			{n: 3, a: []int64{1, maxValue, 1}},
		},
	}
}

func wavePatternTest() testSet {
	return testSet{
		name: "wave",
		cases: []testCase{
			{n: 10, a: []int64{5, 1, 5, 1, 5, 1, 5, 1, 5, 1}},
			{n: 12, a: []int64{3, 6, 9, 6, 3, 6, 9, 6, 3, 6, 9, 6}},
		},
	}
}

func randomTest(name string, rng *rand.Rand, t, nMin, nMax int, maxA int64) testSet {
	var cases []testCase
	for i := 0; i < t; i++ {
		n := rng.Intn(nMax-nMin+1) + nMin
		cases = append(cases, testCase{n: n, a: randomArray(rng, n, maxA)})
	}
	return testSet{name: name, cases: cases}
}

func largeRandomTest(rng *rand.Rand) testSet {
	n := 180000
	return testSet{
		name: "large_random",
		cases: []testCase{
			{n: n, a: randomArray(rng, n, maxValue)},
		},
	}
}

func randomArray(rng *rand.Rand, n int, maxA int64) []int64 {
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		// Ensure values stay within [1, maxA].
		arr[i] = rng.Int63n(maxA) + 1
	}
	return arr
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d 1\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, cases []testCase) ([][]int64, error) {
	fields := strings.Fields(out)
	expectedCount := 0
	for _, tc := range cases {
		expectedCount += tc.n
	}
	if len(fields) != expectedCount {
		return nil, fmt.Errorf("expected %d numbers, got %d", expectedCount, len(fields))
	}
	res := make([][]int64, len(cases))
	idx := 0
	for i, tc := range cases {
		res[i] = make([]int64, tc.n)
		for j := 0; j < tc.n; j++ {
			val, err := strconv.ParseInt(fields[idx], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[idx])
			}
			res[i][j] = val
			idx++
		}
	}
	return res, nil
}

func compareAnswers(expected, actual [][]int64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("expected %d testcases, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if len(expected[i]) != len(actual[i]) {
			return fmt.Errorf("test %d: expected %d values, got %d", i+1, len(expected[i]), len(actual[i]))
		}
		for j := range expected[i] {
			if expected[i][j] != actual[i][j] {
				return fmt.Errorf("test %d: position %d mismatch: expected %d got %d", i+1, j+1, expected[i][j], actual[i][j])
			}
		}
	}
	return nil
}
