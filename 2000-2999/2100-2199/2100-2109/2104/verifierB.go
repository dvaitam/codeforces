package main

import (
	"bufio"
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

const refSource = "./2104B.go"

type testCase struct {
	n int
	a []int64
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2104B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
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
	return stdout.String(), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, a: []int64{42}},
		{n: 2, a: []int64{1, 2}},
		{n: 3, a: []int64{3, 5, 1}},
		{n: 4, a: []int64{10, 10, 10, 10}},
		{n: 5, a: []int64{5, 4, 3, 2, 1}},
		{n: 6, a: []int64{13, 5, 10, 14, 8, 13}}, // from sample prefix
		{n: 7, a: []int64{1, 2, 3, 4, 5, 6, 7}},
		{n: 8, a: []int64{1000000000, 1, 1000000000, 1, 1000000000, 1, 1000000000, 1}},
	}
}

func randomTest(rng *rand.Rand, n int, maxVal int64) testCase {
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(maxVal) + 1
	}
	return testCase{n: n, a: a}
}

func monotoneTest(n int, inc bool) testCase {
	a := make([]int64, n)
	if inc {
		for i := 0; i < n; i++ {
			a[i] = int64(i + 1)
		}
	} else {
		for i := 0; i < n; i++ {
			a[i] = int64(n - i)
		}
	}
	return testCase{n: n, a: a}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
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

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	result := make([][]int64, len(tests))
	for i, tc := range tests {
		row := make([]int64, tc.n)
		for j := 0; j < tc.n; j++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("test %d: expected %d numbers, got %d", i+1, tc.n, j)
			}
			val, err := strconv.ParseInt(sc.Text(), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer %q", i+1, sc.Text())
			}
			row[j] = val
		}
		result[i] = row
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected after %d testcases", len(tests))
	}
	return result, nil
}

func totalN(tests []testCase) int {
	s := 0
	for _, tc := range tests {
		s += tc.n
	}
	return s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 30; i++ {
		n := rng.Intn(50) + 1
		tests = append(tests, randomTest(rng, n, 1_000_000_000))
	}
	for i := 0; i < 20; i++ {
		n := rng.Intn(300) + 50
		tests = append(tests, randomTest(rng, n, 1_000_000_000))
	}
	tests = append(tests, monotoneTest(5000, true))
	tests = append(tests, monotoneTest(5000, false))
	tests = append(tests, randomTest(rng, 10000, 1_000_000_000))

	// Ensure we respect total n <= 2e5 constraint
	if totalN(tests) > 180000 {
		tests = tests[:len(tests)-1]
	}

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseOutput(wantOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(gotOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if len(want[i]) != len(got[i]) {
			fmt.Fprintf(os.Stderr, "test %d failed: length mismatch expected %d got %d\n", i+1, len(want[i]), len(got[i]))
			os.Exit(1)
		}
		for j := range want[i] {
			if want[i][j] != got[i][j] {
				fmt.Fprintf(os.Stderr, "test %d failed at position %d: expected %d got %d\nn=%d a_sample=%v\n",
					i+1, j+1, want[i][j], got[i][j], tests[i].n, sampleArray(tests[i].a))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func sampleArray(a []int64) []int64 {
	if len(a) <= 10 {
		return a
	}
	return append(append(a[:5:5], int64(-1)), a[len(a)-4:]...)
}
