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
)

const refSource = "./2162E.go"

type testCase struct {
	n, k int
	arr  []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	expectedOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\noutput:\n%s\n", err, expectedOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate crashed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refAns, err := parseAnswers(expectedOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		refVals := refAns[i]
		candVals := candAns[i]
		if len(refVals) != tc.k {
			fmt.Fprintf(os.Stderr, "reference output mismatch on test %d: expected %d values got %d\n", i+1, tc.k, len(refVals))
			os.Exit(1)
		}
		if len(candVals) != tc.k {
			fmt.Fprintf(os.Stderr, "candidate output mismatch on test %d: expected %d values got %d\n", i+1, tc.k, len(candVals))
			os.Exit(1)
		}
		for idx, val := range candVals {
			if val < 1 || val > tc.n {
				fmt.Fprintf(os.Stderr, "test %d: value %d at position %d is outside [1, %d]\n", i+1, val, idx+1, tc.n)
				os.Exit(1)
			}
		}

		refArr := mergeSequence(tc.arr, refVals)
		candArr := mergeSequence(tc.arr, candVals)
		refCnt := countPalindromes(refArr)
		candCnt := countPalindromes(candArr)

		if candCnt != refCnt {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d palindromic subarrays, got %d\n", i+1, refCnt, candCnt)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2162E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswers(output string, tests []testCase) ([][]int, error) {
	tokens := strings.Fields(output)
	idx := 0
	ans := make([][]int, len(tests))
	for i, tc := range tests {
		cur := make([]int, 0, tc.k)
		for j := 0; j < tc.k; j++ {
			if idx >= len(tokens) {
				return nil, fmt.Errorf("missing output for test %d", i+1)
			}
			val, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf("token %q is not an integer", tokens[idx])
			}
			cur = append(cur, val)
			idx++
		}
		ans[i] = cur
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[idx])
	}
	return ans, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, sampleTests()...)

	tests = append(tests, testCase{n: 3, k: 1, arr: []int{1, 1, 1}})
	tests = append(tests, testCase{n: 3, k: 3, arr: []int{1, 2, 1}})
	tests = append(tests, alternatingCase(12, 6))

	rng := rand.New(rand.NewSource(21622086))
	for i := 0; i < 25; i++ {
		tests = append(tests, randomCase(rng, 3+rng.Intn(7), 1+rng.Intn(3)))
	}
	for i := 0; i < 20; i++ {
		n := 20 + rng.Intn(180)
		k := 1 + rng.Intn(n)
		tests = append(tests, randomCase(rng, n, k))
	}
	tests = append(tests, randomCase(rng, 500, 250))
	tests = append(tests, randomCase(rng, 1000, 999))
	tests = append(tests, maxCase(rng))

	return tests
}

func sampleTests() []testCase {
	return []testCase{
		{n: 4, k: 1, arr: []int{1, 3, 3, 4}},
		{n: 4, k: 2, arr: []int{2, 2, 2, 2}},
		{n: 5, k: 1, arr: []int{4, 1, 5, 2, 2}},
		{n: 6, k: 3, arr: []int{1, 2, 3, 4, 5, 6}},
		{n: 5, k: 3, arr: []int{3, 2, 5, 2, 3}},
	}
}

func alternatingCase(n, k int) testCase {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			arr[i] = 1
		} else {
			arr[i] = 2
		}
	}
	return testCase{n: n, k: k, arr: arr}
}

func randomCase(rng *rand.Rand, n, k int) testCase {
	if k > n {
		k = n
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	return testCase{n: n, k: k, arr: arr}
}

func maxCase(rng *rand.Rand) testCase {
	n := 200000
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	return testCase{n: n, k: n, arr: arr}
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.k)
		writeArray(&b, tc.arr)
	}
	return b.String()
}

func writeArray(b *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(b, "%d", v)
	}
	b.WriteByte('\n')
}

func mergeSequence(base, extra []int) []int {
	res := make([]int, len(base)+len(extra))
	copy(res, base)
	copy(res[len(base):], extra)
	return res
}

func countPalindromes(arr []int) int64 {
	n := len(arr)
	if n == 0 {
		return 0
	}
	d1 := make([]int, n)
	l, r := 0, -1
	var total int64
	for i := 0; i < n; i++ {
		k := 1
		if i <= r {
			mirror := l + r - i
			if d1[mirror] < r-i+1 {
				k = d1[mirror]
			} else {
				k = r - i + 1
			}
		}
		for i-k >= 0 && i+k < n && arr[i-k] == arr[i+k] {
			k++
		}
		d1[i] = k
		total += int64(k)
		if i+k-1 > r {
			l = i - k + 1
			r = i + k - 1
		}
	}

	d2 := make([]int, n)
	l, r = 0, -1
	for i := 0; i < n; i++ {
		k := 0
		if i <= r {
			mirror := l + r - i + 1
			if d2[mirror] < r-i+1 {
				k = d2[mirror]
			} else {
				k = r - i + 1
			}
		}
		for i+k < n && i-k-1 >= 0 && arr[i+k] == arr[i-k-1] {
			k++
		}
		d2[i] = k
		total += int64(k)
		if i+k-1 > r {
			l = i - k
			r = i + k - 1
		}
	}

	return total
}
