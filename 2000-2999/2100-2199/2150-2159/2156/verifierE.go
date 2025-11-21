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

const maxValue = int64(1_000_000_000)

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		expect, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}

		got, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if err := compareOutputs(expect, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n", i+1, err, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}

	tmp, err := os.CreateTemp("", "2156E_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2156E.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".go"):
		cmd = exec.Command("go", "run", path)
	case strings.HasSuffix(path, ".py"):
		cmd = exec.Command("python3", path)
	default:
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func compareOutputs(expect, got string) error {
	expTokens := strings.Fields(expect)
	gotTokens := strings.Fields(got)
	if len(expTokens) != len(gotTokens) {
		return fmt.Errorf("expected %d tokens got %d", len(expTokens), len(gotTokens))
	}
	for i := range expTokens {
		if expTokens[i] != gotTokens[i] {
			return fmt.Errorf("mismatch at position %d: expected %s got %s", i+1, expTokens[i], gotTokens[i])
		}
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	tests = append(tests, structuredEdgeTests()...)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomBatchTests(rng, 60, 60)...)
	tests = append(tests, randomBatchTests(rng, 60, 2000)...)
	tests = append(tests, randomBatchTests(rng, 40, 40000)...)
	tests = append(tests, largeStressTests(rng)...)
	return tests
}

func manualTests() []testCase {
	cases1 := [][]int64{
		{5, 1, 2, 3, 4},
		{4, 3, 1, 2, 1},
		{6, 100, 97, 95, 92, 90, 88},
	}
	cases2 := [][]int64{
		{8, 2, 2, 2, 2, 2, 2, 2, 2},
		{7, 1, 3, 5, 7, 9, 11, 13, 15},
	}
	cases3 := [][]int64{
		{10, 7, 1, 3, 5, 8, 2, 8, 3, 5, 16},
		{4, 9, 9, 8, 2},
		{5, 4, 4, 3, 5, 34},
	}
	return []testCase{
		makeTestCase(cases1...),
		makeTestCase(cases2...),
		makeTestCase(cases3...),
	}
}

func structuredEdgeTests() []testCase {
	arr1 := increasingSequence(1000, 1)
	arr2 := decreasingSequence(1000, 2)
	arr3 := alternatingSequence(4000)
	arr4 := plateauSequence(6000, 5)
	arr5 := alternatingSequence(10000)
	return []testCase{
		makeTestCase(arr1),
		makeTestCase(arr2),
		makeTestCase(arr3),
		makeTestCase(arr4),
		makeTestCase(arr5),
	}
}

func randomBatchTests(rng *rand.Rand, batches, maxN int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		var cases [][]int64
		sumN := 0
		t := rng.Intn(4) + 1
		for len(cases) < t {
			if sumN+4 > 100000 {
				break
			}
			limit := maxN
			if limit > 100000-sumN {
				limit = 100000 - sumN
			}
			if limit < 4 {
				break
			}
			n := rng.Intn(limit-3) + 4
			pattern := rng.Intn(5)
			cases = append(cases, randomArray(rng, n, pattern))
			sumN += n
		}
		if len(cases) == 0 {
			cases = append(cases, randomArray(rng, 4, 0))
		}
		tests = append(tests, makeTestCase(cases...))
	}
	return tests
}

func largeStressTests(rng *rand.Rand) []testCase {
	return []testCase{
		makeTestCase(increasingSequence(100000, 3)),
		makeTestCase(decreasingSequence(100000, 2)),
		makeTestCase(alternatingSequence(90000)),
		makeTestCase(randomArray(rng, 100000, rng.Intn(5))),
		makeTestCase(plateauSequence(95000, 30)),
	}
}

func randomArray(rng *rand.Rand, n int, pattern int) []int64 {
	arr := make([]int64, n)
	switch pattern {
	case 0:
		for i := range arr {
			arr[i] = rng.Int63n(maxValue) + 1
		}
	case 1:
		start := rng.Int63n(maxValue/2) + 1
		cur := start
		for i := 0; i < n; i++ {
			cur += int64(rng.Intn(5) + 1)
			if cur > maxValue {
				cur = maxValue
			}
			arr[i] = cur
		}
	case 2:
		inc := randomArray(rng, n, 1)
		for i := 0; i < n; i++ {
			arr[i] = inc[n-1-i]
		}
	case 3:
		k := rng.Intn(6) + 2
		pool := make([]int64, k)
		for i := range pool {
			pool[i] = rng.Int63n(maxValue) + 1
		}
		for i := range arr {
			arr[i] = pool[rng.Intn(k)]
		}
	case 4:
		low := rng.Int63n(maxValue/4) + 1
		high := low + int64(rng.Intn(int(maxValue-low))) + 1
		if high > maxValue {
			high = maxValue
		}
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				arr[i] = low + int64(rng.Intn(7))
				if arr[i] > maxValue {
					arr[i] = maxValue
				}
			} else {
				val := high - int64(rng.Intn(7))
				if val < 1 {
					val = 1
				}
				arr[i] = val
			}
		}
	default:
		for i := range arr {
			arr[i] = rng.Int63n(maxValue) + 1
		}
	}
	if pattern == 0 || pattern == 3 || rng.Intn(2) == 0 {
		rng.Shuffle(n, func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
	return arr
}

func increasingSequence(n int, step int64) []int64 {
	arr := make([]int64, n)
	cur := int64(1)
	for i := 0; i < n; i++ {
		arr[i] = cur
		cur += step
		if cur > maxValue {
			cur = maxValue
		}
	}
	return arr
}

func decreasingSequence(n int, step int64) []int64 {
	arr := make([]int64, n)
	cur := maxValue
	for i := 0; i < n; i++ {
		arr[i] = cur
		if cur > step {
			cur -= step
		} else {
			cur = 1
		}
	}
	return arr
}

func alternatingSequence(n int) []int64 {
	arr := make([]int64, n)
	low := int64(1)
	high := maxValue
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			arr[i] = low + int64(i%7)
		} else {
			val := high - int64(i%11)
			if val < 1 {
				val = 1
			}
			arr[i] = val
		}
	}
	return arr
}

func plateauSequence(n int, width int) []int64 {
	arr := make([]int64, n)
	base := rand.New(rand.NewSource(int64(n))).Int63n(maxValue/2) + 1
	cur := base
	for i := 0; i < n; i++ {
		if i%width == 0 {
			cur += int64(i%13 + 1)
			if cur > maxValue {
				cur = maxValue
			}
		}
		arr[i] = cur
	}
	return arr
}

func makeTestCase(cases ...[]int64) testCase {
	if len(cases) == 0 {
		panic("test case must contain at least one array")
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, arr := range cases {
		if len(arr) < 4 {
			panic("each array must have length at least 4")
		}
		sb.WriteString(strconv.Itoa(len(arr)))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String()}
}
