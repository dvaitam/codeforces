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

const refSource = "./2065D.go"

type testCase struct {
	n, m int
	a    [][]int64
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2065D-ref-*")
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
		{n: 1, m: 1, a: [][]int64{{1}}},
		{n: 2, m: 1, a: [][]int64{{1}, {2}}},
		{n: 1, m: 3, a: [][]int64{{4, 4, 6}}},
		{n: 2, m: 2, a: [][]int64{{4, 4}, {6, 1}}},
		{n: 3, m: 2, a: [][]int64{{2, 2}, {2, 2}, {3, 2}}},
	}
}

func randomArray(rng *rand.Rand, m int) []int64 {
	arr := make([]int64, m)
	for i := 0; i < m; i++ {
		arr[i] = rng.Int63n(1_000_000) + 1
	}
	return arr
}

func randomTest(rng *rand.Rand, n, m int) testCase {
	a := make([][]int64, n)
	for i := 0; i < n; i++ {
		a[i] = randomArray(rng, m)
	}
	return testCase{n: n, m: m, a: a}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i, row := range tc.a {
			for j, v := range row {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.FormatInt(v, 10))
			}
			if i+1 < len(tc.a) || len(row) > 0 {
				sb.WriteByte('\n')
			}
		}
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	res := make([]int64, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		val, err := strconv.ParseInt(sc.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer on test %d: %v", i+1, err)
		}
		res[i] = val
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected after %d testcases", t)
	}
	return res, nil
}

func appendIfFits(tests []testCase, tc testCase, total int) ([]testCase, int) {
	need := tc.n * tc.m
	if total+need > 200000 {
		return tests, total
	}
	return append(tests, tc), total + need
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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
	total := 0
	for _, tc := range tests {
		total += tc.n * tc.m
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 120; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(50) + 1
		tc := randomTest(rng, n, m)
		tests, total = appendIfFits(tests, tc, total)
	}

	for i := 0; i < 40; i++ {
		n := rng.Intn(400) + 1
		m := rng.Intn(20) + 1
		tc := randomTest(rng, n, m)
		tests, total = appendIfFits(tests, tc, total)
	}

	// Larger stress cases while respecting the overall element bound.
	if total+100000 <= 200000 {
		tests = append(tests, randomTest(rng, 100000, 1))
		total += 100000
	}
	if total+50000 <= 200000 {
		tests = append(tests, randomTest(rng, 1, 50000))
		total += 50000
	}
	if total+30000 <= 200000 {
		tests = append(tests, randomTest(rng, 300, 100))
		total += 30000
	}

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseOutput(wantOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if want[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\nn=%d m=%d\n", i+1, want[i], got[i], tests[i].n, tests[i].m)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
