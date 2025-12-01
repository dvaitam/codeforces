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

const refSource = "./2022E1.go"

type testCase struct {
	input string
}

type testInstance struct {
	n, m, k int
	assign  []assignment
}

type assignment struct {
	r, c int
	v    int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2022E1-ref-*")
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
	rng := rand.New(rand.NewSource(20222022))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest([]testInstance{
		{n: 2, m: 2, k: 0},
		{n: 3, m: 3, k: 1, assign: []assignment{{r: 1, c: 1, v: 0}}},
	}))

	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(5)+1))
	}

	tests = append(tests, maxCase())

	return tests
}

func sampleTest() testCase {
	return makeTest([]testInstance{
		{
			n: 3, m: 3, k: 8,
			assign: []assignment{
				{1, 1, 0}, {1, 2, 6}, {1, 3, 10},
				{2, 1, 6}, {2, 2, 0}, {2, 3, 12},
				{3, 1, 10}, {3, 2, 12},
			},
		},
		{
			n: 2, m: 5, k: 2,
			assign: []assignment{
				{1, 1, 10},
				{1, 2, 30},
			},
		},
	})
}

func makeTest(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(instances))
	for _, inst := range instances {
		fmt.Fprintf(&b, "%d %d %d %d\n", inst.n, inst.m, inst.k, 0)
		for _, asg := range inst.assign {
			fmt.Fprintf(&b, "%d %d %d\n", asg.r, asg.c, asg.v)
		}
	}
	return testCase{input: b.String()}
}

func randomCase(rng *rand.Rand, maxCases int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	inst := make([]testInstance, t)
	var totalN, totalM, totalK int
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 2
		m := rng.Intn(10) + 2
		k := rng.Intn(min(20, n*m))
		totalN += n
		totalM += m
		if totalN > 1e5 || totalM > 1e5 {
			inst = inst[:i]
			break
		}
		seen := make(map[int]struct{})
		assignments := make([]assignment, 0, k)
		for len(assignments) < k {
			r := rng.Intn(n) + 1
			c := rng.Intn(m) + 1
			key := (r << 20) ^ c
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			v := rng.Intn(1 << 30)
			assignments = append(assignments, assignment{r: r, c: c, v: v})
		}
		inst[i] = testInstance{n: n, m: m, k: len(assignments), assign: assignments}
		totalK += len(assignments)
		if totalK > 1e5 {
			inst = inst[:i+1]
			break
		}
	}
	if len(inst) == 0 {
		inst = append(inst, testInstance{n: 2, m: 2, k: 0})
	}
	return makeTest(inst)
}

func maxCase() testCase {
	n := 100000
	m := 100000
	k := 100000
	assignments := make([]assignment, k)
	for i := 0; i < k; i++ {
		assignments[i] = assignment{r: i + 1, c: 1, v: i % (1 << 30)}
	}
	return makeTest([]testInstance{{n: n, m: m, k: k, assign: assignments}})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
