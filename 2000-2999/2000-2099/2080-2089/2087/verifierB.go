package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2087B_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2087B.go")
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
		return "", fmt.Errorf("unable to determine verifier path")
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
	expTokens := normalizeTokens(expect)
	gotTokens := normalizeTokens(got)
	if len(expTokens) != len(gotTokens) {
		return fmt.Errorf("expected %d answers got %d", len(expTokens), len(gotTokens))
	}
	for i := range expTokens {
		if expTokens[i] != gotTokens[i] {
			return fmt.Errorf("mismatch at position %d: expected %s got %s", i+1, expTokens[i], gotTokens[i])
		}
	}
	return nil
}

func normalizeTokens(out string) []string {
	fields := strings.Fields(out)
	for i := range fields {
		fields[i] = strings.ToUpper(fields[i])
	}
	return fields
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	tests = append(tests, exhaustiveN2Tests()...)
	tests = append(tests, edgeTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 150)...)
	tests = append(tests, batchedRandomTests(rng, 60)...)
	return tests
}

func manualTests() []testCase {
	a1 := []int{3, 7, 5, 12}
	a2 := []int{3, 7, 5, 8}
	a3 := []int{1, 4, 7, 10, 13, 16}
	a4 := []int{1, 100, 2, 99, 3, 98, 4, 97}
	a5 := []int{5, 1, 9, 13, 2, 6}
	a6 := []int{2, 9, 4, 11, 6, 13, 8, 15}
	return []testCase{
		makeTestCase(a1),
		makeTestCase(a2),
		makeTestCase(a3),
		makeTestCase(a4),
		makeTestCase(a5),
		makeTestCase(a6),
		makeTestCase(a1, a2, a4),
	}
}

func exhaustiveN2Tests() []testCase {
	values := []int{1, 2, 3, 4, 5, 6}
	var tests []testCase
	for i := 0; i < len(values); i++ {
		for j := 0; j < len(values); j++ {
			if j == i {
				continue
			}
			for k := 0; k < len(values); k++ {
				if k == i || k == j {
					continue
				}
				for l := 0; l < len(values); l++ {
					if l == i || l == j || l == k {
						continue
					}
					arr := []int{values[i], values[j], values[k], values[l]}
					tests = append(tests, makeTestCase(arr))
				}
			}
		}
	}
	return tests
}

func edgeTests() []testCase {
	var tests []testCase
	arr1 := make([]int, 100)
	for i := range arr1 {
		arr1[i] = 2*i + 1
	}
	arr2 := make([]int, 100)
	for i := range arr2 {
		arr2[i] = 100000 - 2*i
	}
	arr3 := make([]int, 100)
	for i := 0; i < len(arr3); i++ {
		if i%2 == 0 {
			arr3[i] = 1000 + i
		} else {
			arr3[i] = 90000 - i
		}
	}
	arr4 := []int{1, 2, 100000, 99999}
	arr5 := []int{10, 30, 20, 40, 25, 35, 5, 15}
	tests = append(tests,
		makeTestCase(arr1),
		makeTestCase(arr2),
		makeTestCase(arr3),
		makeTestCase(arr4),
		makeTestCase(arr5),
		makeTestCase(arr1, arr4),
	)
	return tests
}

func randomTests(rng *rand.Rand, count int) []testCase {
	var tests []testCase
	for i := 0; i < count; i++ {
		caseCount := 1
		if rng.Intn(3) == 0 {
			caseCount = rng.Intn(3) + 2
		}
		var cases [][]int
		for j := 0; j < caseCount; j++ {
			n := rng.Intn(49) + 2
			if rng.Intn(7) == 0 {
				n = 50
			}
			cases = append(cases, randomRatings(rng, 2*n, rng.Intn(4) == 0))
		}
		tests = append(tests, makeTestCase(cases...))
	}
	return tests
}

func batchedRandomTests(rng *rand.Rand, count int) []testCase {
	var tests []testCase
	for i := 0; i < count; i++ {
		caseCount := rng.Intn(4) + 2
		var cases [][]int
		for j := 0; j < caseCount; j++ {
			n := rng.Intn(49) + 2
			cases = append(cases, randomRatings(rng, 2*n, rng.Intn(2) == 0))
		}
		tests = append(tests, makeTestCase(cases...))
	}
	return tests
}

func randomRatings(rng *rand.Rand, count int, structured bool) []int {
	arr := make([]int, count)
	if structured {
		step := rng.Intn(7) + 1
		maxStart := 100000 - step*(count-1)
		if maxStart < 1 {
			maxStart = 1
		}
		start := rng.Intn(maxStart) + 1
		for i := 0; i < count; i++ {
			arr[i] = start + i*step
		}
		if rng.Intn(3) == 0 {
			// interleave to create equal-distance neighbors
			tmp := make([]int, 0, count)
			for l, r := 0, count-1; l <= r; l, r = l+1, r-1 {
				tmp = append(tmp, arr[l])
				if l != r {
					tmp = append(tmp, arr[r])
				}
			}
			copy(arr, tmp)
		}
	} else {
		used := make(map[int]struct{}, count)
		idx := 0
		for idx < count {
			val := rng.Intn(100000) + 1
			if _, ok := used[val]; ok {
				continue
			}
			used[val] = struct{}{}
			arr[idx] = val
			idx++
		}
	}
	mode := rng.Intn(4)
	switch mode {
	case 0:
		sort.Ints(arr)
	case 1:
		sort.Ints(arr)
		for l, r := 0, len(arr)-1; l < r; l, r = l+1, r-1 {
			arr[l], arr[r] = arr[r], arr[l]
		}
	default:
		rng.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
	return arr
}

func makeTestCase(cases ...[]int) testCase {
	return testCase{input: buildInput(cases)}
}

func buildInput(cases [][]int) string {
	if len(cases) == 0 {
		panic("test case must contain at least one query")
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, arr := range cases {
		if len(arr)%2 != 0 {
			panic("ratings count must be even")
		}
		sb.WriteString(strconv.Itoa(len(arr) / 2))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
