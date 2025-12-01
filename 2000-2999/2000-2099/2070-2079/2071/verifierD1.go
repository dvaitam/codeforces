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

const refSource = "./2071D1.go"

type testCase struct {
	input string
}

type testInstance struct {
	n   int
	l   uint64
	arr []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
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

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2071D1-ref-*")
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

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20712071))
	var tests []testCase

	tests = append(tests, sampleTest())

	tests = append(tests, buildInput([]testInstance{
		{n: 1, l: 1, arr: []int{0}},
		{n: 1, l: 2, arr: []int{1}},
		{n: 2, l: 1, arr: []int{0, 1}},
	}))

	tests = append(tests, buildInput([]testInstance{
		{n: 5, l: 9, arr: []int{1, 0, 1, 0, 1}},
		{n: 10, l: 1000000000000000000, arr: alternating(10)},
	}))

	for i := 0; i < 40; i++ {
		tests = append(tests, randomTestCase(rng, rng.Intn(5)+1, 200000))
	}

	tests = append(tests, randomTestCase(rng, 20, 200000))
	tests = append(tests, edgeCase())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "9\n" +
			"1 1 1\n1\n" +
			"2 3 3\n1 0\n" +
			"3 5 5\n1 1 1\n" +
			"5 1 1\n1 0 1 0 1\n" +
			"1 1000000000000000000 1000000000000000000\n1\n" +
			"6 87 87\n0 1 1 1 1 1\n" +
			"12 69 69\n1 0 0 0 0 1 0 1 0 1 1 0\n" +
			"13 46 46\n0 1 0 1 1 1 1 1 0 1 1 1 0\n" +
			"3 4 4\n1 1 1\n",
	}
}

func buildInput(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(instances))
	for _, inst := range instances {
		fmt.Fprintf(&b, "%d %d %d\n", inst.n, inst.l, inst.l)
		for i, v := range inst.arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}

func randomTestCase(rng *rand.Rand, maxCases, maxTotalN int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	remaining := maxTotalN
	for i := 0; i < t && remaining > 0; i++ {
		n := rng.Intn(min(remaining, 200000)) + 1
		l := randomL(rng)
		arr := randomBinaryArray(rng, n)
		instances = append(instances, testInstance{n: n, l: l, arr: arr})
		remaining -= n
	}
	if len(instances) == 0 {
		instances = append(instances, testInstance{n: 1, l: 1, arr: []int{0}})
	}
	return buildInput(instances)
}

func randomL(rng *rand.Rand) uint64 {
	exponent := rng.Intn(60)
	base := uint64(1) << uint(exponent)
	offset := rng.Uint64() % (base + 1)
	return base + offset
}

func randomBinaryArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(2)
	}
	return arr
}

func alternating(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i % 2
	}
	return arr
}

func edgeCase() testCase {
	n := 200000
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = 1
	}
	return buildInput([]testInstance{
		{n: n, l: 1000000000000000000, arr: arr},
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
