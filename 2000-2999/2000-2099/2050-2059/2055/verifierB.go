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

const refSource = "2000-2999/2000-2099/2050-2059/2055/2055B.go"

type testCase struct {
	input string
}

type testInstance struct {
	a []int64
	b []int64
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
	tmp, err := os.CreateTemp("", "2055B-ref-*")
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
	rng := rand.New(rand.NewSource(20552055))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest([]testInstance{
		{a: []int64{0, 5, 5}, b: []int64{1, 4, 4}},
		{a: []int64{3, 1, 1}, b: []int64{3, 2, 1}},
	}))
	tests = append(tests, makeTest([]testInstance{
		{a: []int64{1, 10, 3}, b: []int64{3, 3, 3}},
		{a: []int64{0, 0, 0}, b: []int64{0, 0, 0}},
	}))

	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(4)+1, 40000))
	}

	tests = append(tests, randomCase(rng, 20, 200000))
	tests = append(tests, edgeCase())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "3\n" +
			"3\n0 5 5\n1 4 4\n" +
			"3\n1 1 3\n2 2 1\n" +
			"2\n1 10\n3 3\n",
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

func randomCase(rng *rand.Rand, maxCases, maxTotal int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	remaining := maxTotal
	for i := 0; i < t && remaining >= 2; i++ {
		capN := min(remaining, 200000)
		if capN < 2 {
			break
		}
		n := rng.Intn(capN-1) + 2
		a := make([]int64, n)
		b := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Int63n(1_000_000_000)
			b[j] = rng.Int63n(1_000_000_000)
		}
		instances = append(instances, testInstance{a: a, b: b})
		remaining -= n
	}
	if len(instances) == 0 {
		instances = append(instances, testInstance{
			a: []int64{1, 1},
			b: []int64{1, 1},
		})
	}
	return makeTest(instances)
}

func edgeCase() testCase {
	n := 200000
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = 0
		b[i] = int64(i % 2)
	}
	return makeTest([]testInstance{{a: a, b: b}})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
