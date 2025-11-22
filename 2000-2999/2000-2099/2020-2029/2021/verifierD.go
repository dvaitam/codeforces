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
	refSource2021D = "2000-2999/2000-2099/2020-2029/2021/2021D.go"
	lim            = int64(1_000_000_000)
)

type testCase struct {
	n, m int
	grid [][]int64
}

type testSet struct {
	name  string
	cases []testCase
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource2021D)
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
		refAns, err := parseOutput(refOut, len(ts.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error: failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, ts.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, ts.name, err, input)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut, len(ts.cases))
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
		tmp, err := os.CreateTemp("", "verifier2021D-*")
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
		singleDayTest(),
		extremeValuesTest(),
		randomTest("random_small", rng, 5, 1, 4, 3, 8),
		randomTest("random_medium", rng, 4, 3, 20, 5, 15),
		stripTest(),
	}
}

func sampleTest() testSet {
	return testSet{
		name: "sample",
		cases: []testCase{
			{
				n: 3,
				m: 6,
				grid: [][]int64{
					{79, 20, 49, 5, -1000, 500},
					{-105, 9, 109, 24, -98, -499},
					{14, 47, 12, 39, 23, 50},
				},
			},
		},
	}
}

func singleDayTest() testSet {
	return testSet{
		name: "single_day",
		cases: []testCase{
			{n: 1, m: 3, grid: [][]int64{{-5, 10, -5}}},
			{n: 1, m: 5, grid: [][]int64{{1, 2, 3, 4, 5}}},
		},
	}
}

func extremeValuesTest() testSet {
	return testSet{
		name: "extremes",
		cases: []testCase{
			{n: 2, m: 3, grid: [][]int64{{lim, -lim, lim}, {-lim, lim, -lim}}},
			{n: 2, m: 4, grid: [][]int64{{-lim, -lim, -lim, -lim}, {lim, lim, lim, lim}}},
		},
	}
}

func stripTest() testSet {
	// Dense input near limit but still within n*m <= 2e5.
	n, m := 400, 500 // 200000 entries
	grid := make([][]int64, n)
	for i := 0; i < n; i++ {
		row := make([]int64, m)
		for j := 0; j < m; j++ {
			row[j] = int64((i*37 + j*17) % 2001)
			row[j] -= 1000 // roughly centered
		}
		grid[i] = row
	}
	return testSet{
		name: "dense_limit",
		cases: []testCase{
			{n: n, m: m, grid: grid},
		},
	}
}

func randomTest(name string, rng *rand.Rand, t, nMin, nMax, mMin, mMax int) testSet {
	var cases []testCase
	for len(cases) < t {
		n := rng.Intn(nMax-nMin+1) + nMin
		m := rng.Intn(mMax-mMin+1) + mMin
		if n*m > 200_000 {
			continue
		}
		grid := make([][]int64, n)
		for i := 0; i < n; i++ {
			row := make([]int64, m)
			for j := 0; j < m; j++ {
				row[j] = rng.Int63n(2*lim+1) - lim
			}
			grid[i] = row
		}
		cases = append(cases, testCase{n: n, m: m, grid: grid})
	}
	return testSet{name: name, cases: cases}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for i, row := range tc.grid {
			for j, v := range row {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.FormatInt(v, 10))
			}
			if i+1 != tc.n {
				sb.WriteByte('\n')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = v
	}
	return res, nil
}

func compareAnswers(expected, actual []int64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("expected %d answers, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("test case %d mismatch: expected %d got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}
