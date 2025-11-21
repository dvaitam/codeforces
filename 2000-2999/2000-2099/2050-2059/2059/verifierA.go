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

const refSource = "2000-2999/2000-2099/2050-2059/2059/2059A.go"

type testCase struct {
	input string
}

type testInstance struct {
	a []int64
	b []int64
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

		if !equalYesNo(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2059A-ref-*")
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

func equalYesNo(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if strings.ToUpper(ta[i]) != strings.ToUpper(tb[i]) {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20592059))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest([]testInstance{
		{a: []int64{1, 1, 2}, b: []int64{2, 2, 3}},
		{a: []int64{5, 5, 6, 6}, b: []int64{1, 1, 2, 2}},
	}))
	tests = append(tests, makeTest([]testInstance{
		{a: []int64{1, 2, 1, 2}, b: []int64{3, 3, 4, 4}},
		{a: []int64{10, 10, 10, 10}, b: []int64{1, 1, 1, 1}},
	}))

	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, rng.Intn(4)+1, 80))
	}

	tests = append(tests, randomTestCase(rng, 10, 500))
	tests = append(tests, maxCase())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "5\n" +
			"4\n1 2 1 2\n1 2 1 2\n" +
			"6\n1 2 3 3 2 1\n1 1 1 1 1 1\n" +
			"3\n1 1 1\n1 1 1\n" +
			"6\n1 5 2 5 2 3\n1 35 9 4 3 5\n" +
			"5\n100 1 100 1 100\n2 2 2 2 2\n",
	}
}

func makeTest(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(instances))
	for _, inst := range instances {
		fmt.Fprintln(&b, len(inst.a))
		writeArray(&b, inst.a)
		writeArray(&b, inst.b)
	}
	return testCase{input: b.String()}
}

func writeArray(b *strings.Builder, arr []int64) {
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(b, "%d", v)
	}
	b.WriteByte('\n')
}

func randomTestCase(rng *rand.Rand, maxCases, maxTotal int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	remaining := maxTotal
	for i := 0; i < t && remaining >= 3; i++ {
		n := rng.Intn(min(remaining, 50)-2) + 3
		a := randomGoodArray(rng, n)
		b := randomGoodArray(rng, n)
		instances = append(instances, testInstance{a: a, b: b})
		remaining -= n
	}
	if len(instances) == 0 {
		instances = append(instances, testInstance{a: []int64{1, 1, 2}, b: []int64{2, 2, 1}})
	}
	return makeTest(instances)
}

func randomGoodArray(rng *rand.Rand, n int) []int64 {
	arr := make([]int64, n)
	vals := rng.Intn(n/2) + 1
	valueSet := make([]int64, vals)
	for i := 0; i < vals; i++ {
		valueSet[i] = rng.Int63n(1_000_000_000) + 1
	}
	for i := 0; i < n; i += 2 {
		v := valueSet[rng.Intn(vals)]
		arr[i] = v
		if i+1 < n {
			arr[i+1] = v
		} else {
			arr[i] = valueSet[0]
		}
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func maxCase() testCase {
	n := 50
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i += 2 {
		a[i], a[i+1] = int64(i+1), int64(i+1)
		b[i], b[i+1] = int64(1000-i), int64(1000-i)
	}
	return makeTest([]testInstance{{a: a, b: b}})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
