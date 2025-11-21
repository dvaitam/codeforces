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

const refSource = "2000-2999/2000-2099/2030-2039/2037/2037B.go"

type testCase struct {
	k    int
	nums []int
}

type testInput struct {
	cases []testCase
}

func (ti testInput) buildInput() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(ti.cases))
	for _, cs := range ti.cases {
		fmt.Fprintln(&b, cs.k)
		for i, v := range cs.nums {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.Itoa(v))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

type pair struct {
	n int64
	m int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		input := tc.buildInput()
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refPairs, err := parsePairs(refOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candPairs, err := parsePairs(candOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i, cs := range tc.cases {
			if err := validatePair(cs, refPairs[i], candPairs[i]); err != nil {
				fmt.Fprintf(os.Stderr, "test %d case %d failed: %v\ninput: k=%d nums=%v\nreference: %d %d\ncandidate: %d %d\n",
					idx+1, i+1, err, cs.k, cs.nums, refPairs[i].n, refPairs[i].m, candPairs[i].n, candPairs[i].m)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func validatePair(tc testCase, ref pair, cand pair) error {
	product := int64(tc.k - 2)
	if cand.n*cand.m != product {
		return fmt.Errorf("candidate pair multiplies to %d, expected %d", cand.n*cand.m, product)
	}
	// At least the numbers must appear in the list with required multiplicities.
	freq := make(map[int64]int)
	for _, v := range tc.nums {
		freq[int64(v)]++
	}
	if cand.n == cand.m {
		if freq[cand.n] < 2 {
			return fmt.Errorf("value %d does not appear twice", cand.n)
		}
	} else {
		if freq[cand.n] < 1 || freq[cand.m] < 1 {
			return fmt.Errorf("pair values %d or %d missing", cand.n, cand.m)
		}
	}
	return nil
}

func parsePairs(out string, expected int) ([]pair, error) {
	fields := strings.Fields(out)
	if len(fields) != expected*2 {
		return nil, fmt.Errorf("expected %d integers, got %d", expected*2, len(fields))
	}
	res := make([]pair, expected)
	for i := 0; i < expected; i++ {
		n, err1 := strconv.ParseInt(fields[2*i], 10, 64)
		m, err2 := strconv.ParseInt(fields[2*i+1], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid integers for case %d", i+1)
		}
		res[i] = pair{n: n, m: m}
	}
	return res, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2037B-ref-*")
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTests())
	tests = append(tests, edgeTests())

	rng := rand.New(rand.NewSource(2037))
	const limit = 200000
	used := 0
	var randomCases []testCase
	for used < limit {
		k := rng.Intn(200000) + 3
		if used+k > limit {
			k = limit - used
			if k < 3 {
				break
			}
		}
		arr := make([]int, k)
		for i := range arr {
			arr[i] = rng.Intn(k) + 1
		}
		randomCases = append(randomCases, testCase{k: k, nums: arr})
		used += k
		if len(randomCases) >= 5000 {
			break
		}
	}
	if len(randomCases) > 0 {
		tests = append(tests, testInput{cases: randomCases})
	}
	return tests
}

func sampleTests() testInput {
	return testInput{cases: []testCase{
		{k: 3, nums: []int{1, 1, 2}},
		{k: 11, nums: []int{3, 3, 4, 5, 6, 7, 8, 9, 9, 10, 11}},
		{k: 8, nums: []int{4, 8, 3, 8, 2, 8, 16, 2}},
		{k: 6, nums: []int{2, 1, 4, 5, 3, 3}},
		{k: 8, nums: []int{1, 2, 6, 3, 8, 5, 5, 3}},
	}}
}

func edgeTests() testInput {
	return testInput{cases: []testCase{
		{k: 3, nums: []int{1, 1, 1}},
		{k: 3, nums: []int{3, 3, 3}},
		{k: 5, nums: []int{1, 2, 3, 4, 5}},
		{k: 200000, nums: generateSequence(200000)},
	}}
}

func generateSequence(k int) []int {
	arr := make([]int, k)
	for i := 0; i < k; i++ {
		arr[i] = (i % k) + 1
	}
	return arr
}
