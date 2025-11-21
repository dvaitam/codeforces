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

const refSource = "2000-2999/2000-2099/2000-2009/2006/2006A.go"

type testCase struct {
	input string
}

type testInstance struct {
	n     int
	edges [][2]int
	s     string
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

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2006A-ref-*")
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
	rng := rand.New(rand.NewSource(20062006))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest([]testInstance{
		starCase(3, "0?1"),
		pathCase(4, "????"),
	}))

	for i := 0; i < 30; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(6)+1))
	}

	tests = append(tests, limitCase())

	return tests
}

func sampleTest() testCase {
	return makeTest([]testInstance{
		{
			n: 4,
			edges: [][2]int{
				{1, 2}, {1, 3}, {1, 4},
			},
			s: "0101",
		},
		{
			n: 4,
			edges: [][2]int{
				{1, 2}, {2, 3}, {3, 4},
			},
			s: "0??1",
		},
	})
}

func starCase(n int, s string) testInstance {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	return testInstance{n: n, edges: edges, s: s}
}

func pathCase(n int, s string) testInstance {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return testInstance{n: n, edges: edges, s: s}
}

func makeTest(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(instances))
	for _, inst := range instances {
		fmt.Fprintln(&b, inst.n)
		for _, e := range inst.edges {
			fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
		}
		fmt.Fprintln(&b, inst.s)
	}
	return testCase{input: b.String()}
}

func randomCase(rng *rand.Rand, maxCases int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	inst := make([]testInstance, t)
	totalN := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(50) + 2
		totalN += n
		if totalN > 200000 {
			inst = inst[:i]
			break
		}
		edges := randomTree(rng, n)
		s := randomString(rng, n)
		inst[i] = testInstance{n: n, edges: edges, s: s}
	}
	if len(inst) == 0 {
		inst = append(inst, starCase(3, "???"))
	}
	return makeTest(inst)
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func randomString(rng *rand.Rand, n int) string {
	bytes := make([]byte, n)
	alphabet := []byte{'0', '1', '?'}
	for i := 0; i < n; i++ {
		bytes[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return string(bytes)
}

func limitCase() testCase {
	n := 100000
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	sBytes := make([]byte, n)
	for i := range sBytes {
		if i%2 == 0 {
			sBytes[i] = '0'
		} else {
			sBytes[i] = '?'
		}
	}
	inst := []testInstance{
		{n: n, edges: edges, s: string(sBytes)},
	}
	return makeTest(inst)
}
