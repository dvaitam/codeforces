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

const refSource = "2000-2999/2100-2199/2150-2159/2159/2159D1.go"

type testCase struct {
	input string
}

type testInstance struct {
	arr []uint64
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
	tmp, err := os.CreateTemp("", "2159D1-ref-*")
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
	rng := rand.New(rand.NewSource(21592159))
	var tests []testCase

	tests = append(tests, sampleTest())

	tests = append(tests, makeTest([]testInstance{
		{arr: []uint64{1}},
		{arr: []uint64{1_000_000_000_000_000_000}},
	}))

	tests = append(tests, makeTest([]testInstance{
		{arr: []uint64{5, 5, 5, 5, 5}},
		{arr: []uint64{9, 7, 5, 3, 1}},
		{arr: []uint64{1, 2, 3, 4, 5, 6}},
	}))

	tests = append(tests, makeTest([]testInstance{
		{arr: alternating(60, 1, 1_000_000_000_000_000_000)},
	}))

	for i := 0; i < 20; i++ {
		tests = append(tests, randomTestCase(rng, 8, 2000, 200000, 1_000_000_000_000))
	}

	tests = append(tests, randomTestCase(rng, 20, 20000, 400000, 1_000_000_000_000_000_000))
	tests = append(tests, stressCase())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "4\n" +
			"5\n3 1 4 1 5\n" +
			"10\n9 2 6 5 3 5 8 9 7 9\n" +
			"8\n1 2 3 4 5 6 7 8\n" +
			"2\n1 1000000000000000000\n",
	}
}

func makeTest(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(instances))
	for _, inst := range instances {
		fmt.Fprintf(&b, "%d\n", len(inst.arr))
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

func randomTestCase(rng *rand.Rand, maxCases, maxN, limit int, maxVal uint64) testCase {
	if limit < 1 {
		limit = 1
	}
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	remaining := limit
	for i := 0; i < t && remaining > 0; i++ {
		capN := maxN
		if capN > remaining {
			capN = remaining
		}
		n := rng.Intn(capN) + 1
		arr := make([]uint64, n)
		for j := 0; j < n; j++ {
			arr[j] = randomValue(rng, maxVal)
		}
		instances = append(instances, testInstance{arr: arr})
		remaining -= n
	}
	if len(instances) == 0 {
		instances = append(instances, testInstance{arr: []uint64{randomValue(rng, maxVal)}})
	}
	return makeTest(instances)
}

func randomValue(rng *rand.Rand, maxVal uint64) uint64 {
	if maxVal == 0 {
		return 1
	}
	return uint64(rng.Int63n(int64(maxVal))) + 1
}

func alternating(n int, low, high uint64) []uint64 {
	if n <= 0 {
		return []uint64{low}
	}
	arr := make([]uint64, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			arr[i] = low
		} else {
			arr[i] = high
		}
	}
	return arr
}

func stressCase() testCase {
	n := 400000
	arr := make([]uint64, n)
	for i := 0; i < n; i++ {
		if i%3 == 0 {
			arr[i] = 1
		} else if i%3 == 1 {
			arr[i] = uint64(i+1) * 3
		} else {
			arr[i] = 1_000_000_000_000_000_000
		}
	}
	return makeTest([]testInstance{{arr: arr}})
}
