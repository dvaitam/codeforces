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

const refSource = "./2072F.go"

type testCase struct {
	input string
}

type testInstance struct {
	n int
	k int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2072F-ref-*")
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
	rng := rand.New(rand.NewSource(20722072))
	var tests []testCase

	tests = append(tests, sampleTest())

	tests = append(tests, buildInput([]testInstance{
		{n: 1, k: 1},
		{n: 2, k: 1},
		{n: 3, k: 1},
		{n: 4, k: 1},
	}))

	tests = append(tests, buildInput([]testInstance{
		{n: 5, k: (1 << 30) - 1},
		{n: 6, k: 123456789},
	}))

	for i := 0; i < 40; i++ {
		tests = append(tests, randomTestCase(rng, rng.Intn(5)+1, 200000))
	}

	tests = append(tests, randomTestCase(rng, 50, 1000000))
	tests = append(tests, singleLargeCase())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "5\n" +
			"1 5\n" +
			"2 10\n" +
			"3 16\n" +
			"9 1\n" +
			"1 52\n",
	}
}

func buildInput(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(instances))
	for _, inst := range instances {
		fmt.Fprintf(&b, "%d %d\n", inst.n, inst.k)
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
		n := rng.Intn(min(remaining, 1000000)) + 1
		k := rng.Int63n(1<<31-1) + 1
		instances = append(instances, testInstance{n: n, k: k})
		remaining -= n
	}
	if len(instances) == 0 {
		instances = append(instances, testInstance{n: 1, k: 1})
	}
	return buildInput(instances)
}

func singleLargeCase() testCase {
	return buildInput([]testInstance{
		{n: 1000000, k: 1<<30 - 3},
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
