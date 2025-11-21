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

const refSource = "2000-2999/2000-2099/2080-2089/2087/2087E.go"

type testCase struct {
	input string
}

type testInstance struct {
	dir string
	val []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2087E-ref-*")
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

func equalTokens(expected, got string) bool {
	ta := strings.Fields(expected)
	tb := strings.Fields(got)
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
	rng := rand.New(rand.NewSource(20872087))
	var tests []testCase

	tests = append(tests, sampleTest())

	tests = append(tests, makeTest([]testInstance{
		{dir: ">", val: []int64{5}},
		{dir: "<", val: []int64{-5}},
	}))

	tests = append(tests, makeTest([]testInstance{
		{dir: "<>", val: []int64{1, -2}},
		{dir: ">>><<<", val: []int64{5, 4, -3, 2, 0, -1}},
	}))

	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, rng.Intn(5)+1, 2000))
	}

	tests = append(tests, randomTestCase(rng, 20, 300000))
	tests = append(tests, alternatingTest())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "5\n" +
			"3\n<>>\n5 4 6\n" +
			"5\n<><>>\n5 -2 4 -3 7\n" +
			"2\n>>\n-1 -2\n" +
			"8\n>>>><<<<\n1 -1 1 -1 1 -1 1 -1\n" +
			"5\n><<<>\n-1 100 100 100 100\n",
	}
}

func makeTest(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(instances))
	for _, inst := range instances {
		fmt.Fprintf(&b, "%d\n", len(inst.dir))
		b.WriteString(inst.dir)
		b.WriteByte('\n')
		for i, v := range inst.val {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}

func randomTestCase(rng *rand.Rand, maxCases int, maxTotal int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	remaining := maxTotal
	for i := 0; i < t && remaining > 0; i++ {
		capN := remaining
		if capN > 50000 {
			capN = 50000
		}
		n := rng.Intn(capN) + 1
		dir := randomDir(rng, n)
		vals := randomValues(rng, n)
		instances = append(instances, testInstance{dir: dir, val: vals})
		remaining -= n
	}
	if len(instances) == 0 {
		instances = append(instances, testInstance{dir: ">", val: []int64{0}})
	}
	return makeTest(instances)
}

func randomDir(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '<'
		} else {
			b[i] = '>'
		}
	}
	return string(b)
}

func randomValues(rng *rand.Rand, n int) []int64 {
	vals := make([]int64, n)
	for i := 0; i < n; i++ {
		vals[i] = rng.Int63n(2_000_000_001) - 1_000_000_000
	}
	return vals
}

func alternatingTest() testCase {
	n := 300000
	dir := make([]byte, n)
	vals := make([]int64, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			dir[i] = '>'
			vals[i] = int64(i + 1)
		} else {
			dir[i] = '<'
			vals[i] = -int64(i + 1)
		}
	}
	return makeTest([]testInstance{{dir: string(dir), val: vals}})
}
