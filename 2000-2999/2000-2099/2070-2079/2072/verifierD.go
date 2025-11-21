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

type testCase struct {
	arrays [][]int
	input  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, tc := range tests {
		expect, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		got, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, got)
			os.Exit(1)
		}
		if err := verifyOutputs(tc.arrays, expect, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n", idx+1, err, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2072D_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2072D.go")
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
		return "", fmt.Errorf("unable to locate verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".go"):
		cmd = exec.Command("go", "run", path)
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

func verifyOutputs(arrays [][]int, expect, got string) error {
	refPairs, err := parsePairs(expect, len(arrays))
	if err != nil {
		return fmt.Errorf("reference output invalid: %v", err)
	}
	partPairs, err := parsePairs(got, len(arrays))
	if err != nil {
		return fmt.Errorf("participant output invalid: %v", err)
	}

	for i := range arrays {
		n := len(arrays[i])
		lr := partPairs[i]
		if lr.l < 1 || lr.r < lr.l || lr.r > n {
			return fmt.Errorf("case %d: invalid pair %d %d", i+1, lr.l, lr.r)
		}
		refLR := refPairs[i]
		refCount, err := resultingInversions(arrays[i], refLR.l, refLR.r)
		if err != nil {
			return fmt.Errorf("case %d: %v", i+1, err)
		}
		partCount, err := resultingInversions(arrays[i], lr.l, lr.r)
		if err != nil {
			return fmt.Errorf("case %d: %v", i+1, err)
		}
		if partCount != refCount {
			return fmt.Errorf("case %d: inversion count mismatch, expected %d got %d", i+1, refCount, partCount)
		}
	}
	return nil
}

type pair struct {
	l int
	r int
}

func parsePairs(out string, t int) ([]pair, error) {
	fields := strings.Fields(out)
	if len(fields) < 2*t {
		return nil, fmt.Errorf("expected %d integers got %d", 2*t, len(fields))
	}
	pairs := make([]pair, t)
	idx := 0
	for i := 0; i < t; i++ {
		l, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[idx])
		}
		r, err := strconv.Atoi(fields[idx+1])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[idx+1])
		}
		pairs[i] = pair{l: l, r: r}
		idx += 2
	}
	return pairs, nil
}

func resultingInversions(arr []int, l, r int) (int64, error) {
	if l < 1 || l > len(arr) || r < 1 || r > len(arr) || l > r {
		return 0, fmt.Errorf("invalid l,r = %d,%d for n=%d", l, r, len(arr))
	}
	shifted := rotateOnce(arr, l-1, r-1)
	return inversionCount(shifted), nil
}

func rotateOnce(arr []int, l, r int) []int {
	res := append([]int(nil), arr...)
	if l == r {
		return res
	}
	first := res[l]
	copy(res[l:r], res[l+1:r+1])
	res[r] = first
	return res
}

func inversionCount(arr []int) int64 {
	maxVal := 2000
	bit := make([]int, maxVal+2)
	var sum func(int) int64
	sum = func(idx int) int64 {
		if idx > maxVal+1 {
			idx = maxVal + 1
		}
		var res int64
		for idx > 0 {
			res += int64(bit[idx])
			idx -= idx & -idx
		}
		return res
	}
	add := func(idx int) {
		for idx <= maxVal+1 {
			bit[idx]++
			idx += idx & -idx
		}
	}
	var inv int64
	for i := len(arr) - 1; i >= 0; i-- {
		val := arr[i]
		if val < 0 {
			val = 0
		}
		if val > maxVal {
			val = maxVal
		}
		inv += sum(val)
		add(val + 1)
	}
	return inv
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	tests = append(tests, smallExhaustiveTests()...)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomBatches(rng, 60, 5, 12)...)
	tests = append(tests, randomBatches(rng, 60, 200, 2000)...)
	tests = append(tests, largeStressTests(rng)...)
	return tests
}

func manualTests() []testCase {
	return []testCase{
		makeTestCase([][]int{{1}}),
		makeTestCase([][]int{{1, 2}}),
		makeTestCase([][]int{{2, 1}}),
		makeTestCase([][]int{{1, 1, 1, 1}}),
		makeTestCase([][]int{{3, 1, 4, 1, 5}}),
		makeTestCase([][]int{{5, 1, 2, 3, 4}, {4, 3, 1, 2, 1}}),
	}
}

func smallExhaustiveTests() []testCase {
	var tests []testCase
	values := []int{1, 2, 3}
	for n := 3; n <= 5; n++ {
		arr := make([]int, n)
		var dfs func(int)
		dfs = func(pos int) {
			if pos == n {
				tmp := make([]int, n)
				copy(tmp, arr)
				tests = append(tests, makeTestCase([][]int{tmp}))
				return
			}
			for _, v := range values {
				arr[pos] = v
				dfs(pos + 1)
			}
		}
		dfs(0)
	}
	return tests
}

func randomBatches(rng *rand.Rand, batches int, minN, maxN int) []testCase {
	var tests []testCase
	for i := 0; i < batches; i++ {
		t := rng.Intn(3) + 1
		var arrays [][]int
		sumSq := 0
		for len(arrays) < t {
			n := rng.Intn(maxN-minN+1) + minN
			if n <= 0 {
				n = minN
			}
			if sumSq+n*n > 3_500_000 {
				break
			}
			arrays = append(arrays, randomArray(rng, n))
			sumSq += n * n
		}
		if len(arrays) == 0 {
			arrays = append(arrays, randomArray(rng, minN))
		}
		tests = append(tests, makeTestCase(arrays))
	}
	return tests
}

func largeStressTests(rng *rand.Rand) []testCase {
	return []testCase{
		makeTestCase([][]int{increasingArray(2000)}),
		makeTestCase([][]int{decreasingArray(2000)}),
		makeTestCase([][]int{randomArray(rng, 2000)}),
		makeTestCase([][]int{alternatingArray(2000)}),
		makeTestCase([][]int{plateauArray(2000)}),
	}
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	pattern := rng.Intn(4)
	switch pattern {
	case 0:
		for i := range arr {
			arr[i] = rng.Intn(2000) + 1
		}
	case 1:
		base := rng.Intn(1000) + 1
		for i := range arr {
			arr[i] = base + (i % 5)
			if arr[i] > 2000 {
				arr[i] = 2000
			}
		}
	case 2:
		for i := range arr {
			arr[i] = 2000 - (i % 7)
		}
	case 3:
		values := []int{1, 1000, 2000}
		for i := range arr {
			arr[i] = values[rng.Intn(len(values))]
		}
	}
	if rng.Intn(2) == 0 {
		rng.Shuffle(n, func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
	return arr
}

func increasingArray(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = 1 + i%2000
	}
	return arr
}

func decreasingArray(n int) []int {
	arr := make([]int, n)
	val := 2000
	for i := 0; i < n; i++ {
		arr[i] = val
		val--
		if val <= 0 {
			val = 2000
		}
	}
	return arr
}

func alternatingArray(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			arr[i] = 1 + (i % 100)
		} else {
			arr[i] = 2000 - (i % 100)
		}
	}
	return arr
}

func plateauArray(n int) []int {
	arr := make([]int, n)
	val := rand.New(rand.NewSource(int64(n))).Intn(2000) + 1
	for i := 0; i < n; i++ {
		if i%50 == 0 {
			val = rand.New(rand.NewSource(int64(n+i))).Intn(2000) + 1
		}
		arr[i] = val
	}
	return arr
}

func makeTestCase(arrays [][]int) testCase {
	if len(arrays) == 0 {
		panic("need at least one array")
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(arrays)))
	sb.WriteByte('\n')
	for _, arr := range arrays {
		if len(arr) == 0 {
			panic("array must be non-empty")
		}
		sb.WriteString(strconv.Itoa(len(arr)))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	copied := make([][]int, len(arrays))
	for i := range arrays {
		copied[i] = append([]int(nil), arrays[i]...)
	}
	return testCase{arrays: copied, input: sb.String()}
}
