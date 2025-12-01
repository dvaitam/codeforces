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

const refSource = "./2085F1.go"

type testCase struct {
	input string
}

type testInstance struct {
	n, k int
	arr  []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2085F1-ref-*")
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
	rng := rand.New(rand.NewSource(20852085))
	var tests []testCase

	tests = append(tests, sampleTest())

	tests = append(tests, buildInput([]testInstance{
		{n: 2, k: 2, arr: []int{2, 1}},
	}))

	tests = append(tests, buildInput([]testInstance{
		{n: 5, k: 3, arr: []int{1, 1, 1, 2, 3}},
		{n: 6, k: 4, arr: []int{4, 4, 3, 2, 1, 1}},
	}))

	tests = append(tests, buildInput([]testInstance{
		{n: 8, k: 5, arr: []int{5, 5, 5, 1, 2, 3, 4, 4}},
	}))

	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, rng.Intn(4)+1, 50))
	}

	tests = append(tests, randomTestCase(rng, 5, 400))
	tests = append(tests, worstCase())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "6\n" +
			"3 2\n1 2 1\n" +
			"7 3\n2 1 1 3 1 1 2\n" +
			"6 3\n1 1 2 2 2 3\n" +
			"6 3\n1 2 2 2 2 3\n" +
			"10 5\n5 1 3 1 1 2 2 4 1 3\n" +
			"9 4\n1 2 3 3 3 3 3 2 4\n",
	}
}

func buildInput(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(instances))
	for _, inst := range instances {
		fmt.Fprintf(&b, "%d %d\n", inst.n, inst.k)
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

func randomTestCase(rng *rand.Rand, maxCases, maxN int) testCase {
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-1) + 2
		k := rng.Intn(n-1) + 2
		arr := randomArray(rng, n, k)
		instances = append(instances, testInstance{n: n, k: k, arr: arr})
	}
	return buildInput(instances)
}

func randomArray(rng *rand.Rand, n, k int) []int {
	arr := make([]int, n)
	for i := 0; i < k; i++ {
		arr[i] = i + 1
	}
	for i := k; i < n; i++ {
		arr[i] = rng.Intn(k) + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func worstCase() testCase {
	n := 3000
	k := 50
	arr := make([]int, n)
	idx := 0
	for val := 1; val <= k; val++ {
		arr[idx] = val
		idx++
	}
	for idx < n {
		arr[idx] = (idx % k) + 1
		idx++
	}
	return buildInput([]testInstance{{n: n, k: k, arr: arr}})
}
