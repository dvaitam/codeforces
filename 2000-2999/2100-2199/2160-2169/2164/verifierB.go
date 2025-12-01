package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2164B.go"

type testCase struct {
	input string
	cases []caseData
}

type caseData struct {
	arr []int64
	set map[int64]struct{}
}

type testInstance struct {
	arr []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for idx, tc := range tests {
		expectOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test #%d: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}

		expect, err := parseReference(expectOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test #%d: %v\noutput:\n%s", idx+1, err, expectOut)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test #%d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}

		if err := validateCandidate(gotOut, tc.cases, expect); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test #%d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2164B-ref-*")
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
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
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
	return strings.TrimSpace(out.String()), err
}

func parseReference(output string, cases int) ([]bool, error) {
	res := make([]bool, cases)
	reader := bufio.NewReader(strings.NewReader(output))
	for i := 0; i < cases; i++ {
		var first int64
		if _, err := fmt.Fscan(reader, &first); err != nil {
			return nil, fmt.Errorf("case %d: missing first value: %w", i+1, err)
		}
		if first == -1 {
			res[i] = false
			continue
		}
		res[i] = true
		var second int64
		if _, err := fmt.Fscan(reader, &second); err != nil {
			return nil, fmt.Errorf("case %d: missing second value: %w", i+1, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra token %q in reference output", extra)
	} else if err != io.EOF {
		return nil, fmt.Errorf("failed to read reference output: %w", err)
	}
	return res, nil
}

func validateCandidate(output string, cases []caseData, expect []bool) error {
	reader := bufio.NewReader(strings.NewReader(output))
	for i, c := range cases {
		var first int64
		if _, err := fmt.Fscan(reader, &first); err != nil {
			return fmt.Errorf("case %d: missing first value: %w", i+1, err)
		}
		if !expect[i] {
			if first != -1 {
				return fmt.Errorf("case %d: expected -1, got %d", i+1, first)
			}
			continue
		}
		x := first
		var y int64
		if _, err := fmt.Fscan(reader, &y); err != nil {
			return fmt.Errorf("case %d: expected two integers, missing y: %w", i+1, err)
		}
		if x >= y {
			return fmt.Errorf("case %d: require x < y, got %d %d", i+1, x, y)
		}
		if _, ok := c.set[x]; !ok {
			return fmt.Errorf("case %d: x=%d not in sequence", i+1, x)
		}
		if _, ok := c.set[y]; !ok {
			return fmt.Errorf("case %d: y=%d not in sequence", i+1, y)
		}
		if ((y % x) & 1) != 0 {
			return fmt.Errorf("case %d: y mod x = %d is odd", i+1, y%x)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("extra token %q in output", extra)
	} else if err != io.EOF {
		return fmt.Errorf("failed to read output: %w", err)
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(21642164))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest([]testInstance{
		{arr: []int64{2, 5}},
		{arr: []int64{4, 6, 8, 10}},
		{arr: []int64{5, 6}},
	}))

	tests = append(tests, makeTest([]testInstance{
		{arr: []int64{1, 2}},
		{arr: []int64{3, 7, 11, 15}},
		{arr: []int64{6, 13, 20, 27, 34}},
	}))

	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, rng.Intn(5)+1, 5000))
	}

	tests = append(tests, randomTestCase(rng, 12, 100000))
	tests = append(tests, largeCase())

	return tests
}

func sampleTest() testCase {
	return makeTest([]testInstance{
		{arr: []int64{1, 3, 4, 5, 6}},
		{arr: []int64{2, 3, 5, 7, 11, 13}},
		{arr: []int64{4, 7, 13}},
		{arr: []int64{17, 117, 1117}},
	})
}

func makeTest(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(instances))
	cases := make([]caseData, len(instances))
	for i, inst := range instances {
		fmt.Fprintf(&b, "%d\n", len(inst.arr))
		for j, v := range inst.arr {
			if j > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		set := make(map[int64]struct{}, len(inst.arr))
		arrCopy := make([]int64, len(inst.arr))
		copy(arrCopy, inst.arr)
		for _, v := range inst.arr {
			set[v] = struct{}{}
		}
		cases[i] = caseData{arr: arrCopy, set: set}
	}
	return testCase{input: b.String(), cases: cases}
}

func randomTestCase(rng *rand.Rand, maxCases, maxTotalN int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	remaining := maxTotalN
	for i := 0; i < t && remaining >= 2; i++ {
		capN := min(remaining, 2000)
		if capN < 2 {
			break
		}
		n := rng.Intn(capN-1) + 2
		instances = append(instances, testInstance{arr: randomSequence(rng, n)})
		remaining -= n
	}
	if len(instances) == 0 {
		instances = append(instances, testInstance{arr: randomSequence(rng, 2)})
	}
	return makeTest(instances)
}

func largeCase() testCase {
	n := 100000
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(i + 1)
	}
	return makeTest([]testInstance{{arr: arr}})
}

func randomSequence(rng *rand.Rand, n int) []int64 {
	arr := make([]int64, n)
	cur := int64(rng.Intn(10) + 1)
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(10) + 1)
		arr[i] = cur
	}
	return arr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
