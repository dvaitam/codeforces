package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2000-2999/2100-2199/2160-2169/2163/2163A.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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

	for i, tc := range tests {
		expect, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}
		if !equalVerdicts(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2163A-ref-*")
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

func runCandidate(bin, input string) (string, error) {
	cmd := commandFor(bin)
	return runWithInput(cmd, input)
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func normalizeVerdicts(s string) []string {
	return strings.Fields(strings.ToLower(s))
}

func equalVerdicts(expected, got string) bool {
	expFields := normalizeVerdicts(expected)
	gotFields := normalizeVerdicts(got)
	if len(expFields) != len(gotFields) {
		return false
	}
	for i := range expFields {
		if expFields[i] != gotFields[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(21632163))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest([][]int{
		{1, 1, 1},
	}))
	tests = append(tests, makeTest([][]int{
		{3, 2, 1, 1},
		{1, 2, 3, 3},
	}))
	tests = append(tests, makeTest([][]int{
		buildRange(3, 1),
		buildRange(5, 5),
	}))

	for i := 0; i < 40; i++ {
		tests = append(tests, randomTestCase(rng, rng.Intn(5)+1, 3, 100))
	}

	tests = append(tests, randomTestCase(rng, 100, 3, 100))
	tests = append(tests, pathologicalCase())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "5\n" +
			"4\n4 2 2 1\n" +
			"4\n1 1 1 1\n" +
			"5\n1 5 1 5 1\n" +
			"3\n1 2 3\n" +
			"5\n1 3 2 3 5\n",
	}
}

func makeTest(arrays [][]int) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(arrays))
	for _, arr := range arrays {
		fmt.Fprintf(&b, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}

func randomTestCase(rng *rand.Rand, tCount int, minN, maxN int) testCase {
	if tCount < 1 {
		tCount = 1
	}
	var arrays [][]int
	for i := 0; i < tCount; i++ {
		n := rng.Intn(maxN-minN+1) + minN
		if n < minN {
			n = minN
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(n) + 1
		}
		arrays = append(arrays, arr)
	}
	return makeTest(arrays)
}

func buildRange(n, val int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = val
		if arr[i] > n {
			arr[i] = (val % n) + 1
		}
	}
	return arr
}

func pathologicalCase() testCase {
	var arrays [][]int
	arr := make([]int, 100)
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			arr[i] = 1
		} else {
			arr[i] = 100
		}
	}
	arrays = append(arrays, arr)
	return makeTest(arrays)
}
