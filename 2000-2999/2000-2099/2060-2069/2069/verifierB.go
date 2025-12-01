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

const refSource = "./2069B.go"

type testCase struct {
	input string
}

type testInstance struct {
	grid [][]int
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

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2069B-ref-*")
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
	rng := rand.New(rand.NewSource(20692069))
	var tests []testCase

	tests = append(tests, sampleTest())

	tests = append(tests, buildInput([]testInstance{
		{grid: [][]int{{1}}},
		{grid: [][]int{{1, 2}, {3, 4}}},
	}))

	tests = append(tests, buildInput([]testInstance{
		{grid: [][]int{
			{1, 1, 1},
			{1, 1, 1},
		}},
		{grid: [][]int{
			{1, 2, 3, 4},
			{4, 3, 2, 1},
			{5, 6, 7, 8},
		}},
	}))

	for i := 0; i < 40; i++ {
		tests = append(tests, randomTestCase(rng, rng.Intn(5)+1, 10000))
	}

	tests = append(tests, randomTestCase(rng, 20, 500000))
	tests = append(tests, extremeCase())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "4\n" +
			"1 1\n1\n" +
			"3 3\n1 2 1\n2 3 2\n1 3 1\n" +
			"2 3\n1 3 1\n1 6 5\n" +
			"3 4\n1 4 2 2\n1 4 3 5\n6 6 3 5\n",
	}
}

func buildInput(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(instances))
	for _, inst := range instances {
		n := len(inst.grid)
		m := len(inst.grid[0])
		fmt.Fprintf(&b, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					b.WriteByte(' ')
				}
				fmt.Fprintf(&b, "%d", inst.grid[i][j])
			}
			b.WriteByte('\n')
		}
	}
	return testCase{input: b.String()}
}

func randomTestCase(rng *rand.Rand, maxCases, maxCells int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	remaining := maxCells
	for i := 0; i < t && remaining > 0; i++ {
		n := rng.Intn(min(remaining, 700)) + 1
		m := rng.Intn(min(remaining/n, 700)) + 1
		if n > m {
			n, m = m, n
		}
		grid := randomGrid(rng, n, m)
		instances = append(instances, testInstance{grid: grid})
		remaining -= n * m
	}
	if len(instances) == 0 {
		instances = append(instances, testInstance{grid: [][]int{{1}}})
	}
	return buildInput(instances)
}

func randomGrid(rng *rand.Rand, n, m int) [][]int {
	grid := make([][]int, n)
	val := 1
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = val
			val++
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			swapI := rng.Intn(n)
			swapJ := rng.Intn(m)
			grid[i][j], grid[swapI][swapJ] = grid[swapI][swapJ], grid[i][j]
		}
	}
	return grid
}

func extremeCase() testCase {
	n, m := 700, 700
	grid := make([][]int, n)
	val := 1
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = val
			val++
		}
	}
	return buildInput([]testInstance{{grid: grid}})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
