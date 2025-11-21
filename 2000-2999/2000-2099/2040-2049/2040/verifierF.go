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

const mod = 998244353

type testCase struct {
	a, b, c int
	d       []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2040F-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", outPath, "2040F.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
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
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 256)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		k := len(tc.d)
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.a, tc.b, tc.c, k))
		for i, val := range tc.d {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != len(tests) {
		return nil, fmt.Errorf("expected %d answers, got %d", len(tests), len(fields))
	}
	res := make([]int, len(tests))
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at test %d", f, i+1)
		}
		val %= mod
		if val < 0 {
			val += mod
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			a: 1, b: 1, c: 1,
			d: []int{1},
		},
		{
			a: 1, b: 2, c: 3,
			d: []int{1, 1, 1, 1, 1, 1},
		},
		{
			a: 2, b: 2, c: 1,
			d: []int{1, 1, 1, 1},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 150)
	totalVolume := 0
	for len(tests) < cap(tests) && totalVolume < 2_800_000 {
		a := rng.Intn(30) + 1
		b := rng.Intn(30) + 1
		c := rng.Intn(10) + 1
		if a*b*c > 50_000 {
			continue
		}
		vol := a * b * c
		if totalVolume+vol > 3_000_000 {
			break
		}
		k := rng.Intn(10) + 1
		if k > vol {
			k = vol
		}
		counts := randomCounts(rng, vol, k)
		tests = append(tests, testCase{a: a, b: b, c: c, d: counts})
		totalVolume += vol
	}
	return tests
}

func randomCounts(rng *rand.Rand, sum, parts int) []int {
	if parts == 0 {
		return nil
	}
	breaks := make([]int, parts-1)
	for i := range breaks {
		breaks[i] = rng.Intn(sum-1) + 1
	}
	sortInts(breaks)
	res := make([]int, parts)
	prev := 0
	for i := 0; i < parts-1; i++ {
		res[i] = breaks[i] - prev
		prev = breaks[i]
	}
	res[parts-1] = sum - prev
	sortInts(res)
	return res
}

func sortInts(a []int) {
	if len(a) <= 1 {
		return
	}
	quickSort(a, 0, len(a)-1)
}

func quickSort(a []int, l, r int) {
	if l >= r {
		return
	}
	i, j := l, r
	pivot := a[l+(r-l)/2]
	for i <= j {
		for a[i] < pivot {
			i++
		}
		for a[j] > pivot {
			j--
		}
		if i <= j {
			a[i], a[j] = a[j], a[i]
			i++
			j--
		}
	}
	if l < j {
		quickSort(a, l, j)
	}
	if i < r {
		quickSort(a, i, r)
	}
}

func compareAnswers(expected, actual []int) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("answer count mismatch")
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("test %d mismatch: expected %d, got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	expectedOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	actualOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	expectedAns, err := parseOutput(expectedOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expectedOut)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actualOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actualOut)
		os.Exit(1)
	}

	if err := compareAnswers(expectedAns, actualAns); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}
