package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type testCase struct {
	input  string
	people int
}

const refSource = "./2172B.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.people)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		userVals, err := parseOutputs(userOut, tc.people)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, tc.input, userOut)
			os.Exit(1)
		}

		if err := compareAnswers(refVals, userVals); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2172B-ref-*")
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

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]float64, error) {
	reader := strings.NewReader(out)
	vals := make([]float64, expected)
	for i := 0; i < expected; i++ {
		if _, err := fmt.Fscan(reader, &vals[i]); err != nil {
			return nil, fmt.Errorf("expected %d values, got %d (%v)", expected, i, err)
		}
	}
	var extra float64
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %.6f)", extra)
	} else if err != io.EOF {
		return nil, fmt.Errorf("error while checking extra output: %v", err)
	}
	return vals, nil
}

func compareAnswers(expected, actual []float64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("expected %d answers, got %d", len(expected), len(actual))
	}
	for i := range expected {
		diff := math.Abs(expected[i] - actual[i])
		limit := 1e-6 * math.Max(1.0, math.Abs(expected[i]))
		if diff > limit+1e-9 {
			return fmt.Errorf("answer %d mismatch: expected %.10f, got %.10f", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, simpleEdgeTest())

	rng := rand.New(rand.NewSource(2172))
	tests = append(tests, randomCase(rng, 5, 5, 100))
	tests = append(tests, randomCase(rng, 20, 25, 1000))
	tests = append(tests, randomCase(rng, 200, 200, 100000))
	tests = append(tests, randomCase(rng, 1000, 1000, 1000000000))
	tests = append(tests, randomCase(rng, 200000, 200000, 1000000000))

	return tests
}

func sampleTest() testCase {
	input := strings.TrimSpace(`3 3 10 4 1
0 5
2 4
7 9
3
8
5
`) + "\n"
	return testCase{input: input, people: 3}
}

func simpleEdgeTest() testCase {
	n, m := 2, 5
	L := int64(100)
	x := int64(10)
	y := int64(1)
	buses := [][2]int64{{0, 100}, {50, 100}}
	people := []int64{0, 25, 50, 75, 100}
	return buildCase(n, m, L, x, y, buses, people)
}

func randomCase(rng *rand.Rand, n, m int, maxL int64) testCase {
	if maxL < 1 {
		maxL = 1
	}
	L := rng.Int63n(maxL) + 1
	y := rng.Int63n(1000000-1) + 1
	x := y + 1 + rng.Int63n(1000000-y)

	buses := make([][2]int64, n)
	for i := 0; i < n; i++ {
		var s int64
		if L > 1 {
			s = rng.Int63n(L)
		}
		t := s + 1 + rng.Int63n(L-s)
		buses[i] = [2]int64{s, t}
	}

	people := make([]int64, m)
	for i := 0; i < m; i++ {
		people[i] = rng.Int63n(L + 1)
	}

	return buildCase(n, m, L, x, y, buses, people)
}

func buildCase(n, m int, L, x, y int64, buses [][2]int64, people []int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d %d\n", n, m, L, x, y)
	for _, b := range buses {
		fmt.Fprintf(&sb, "%d %d\n", b[0], b[1])
	}
	for _, p := range people {
		fmt.Fprintf(&sb, "%d\n", p)
	}
	return testCase{input: sb.String(), people: m}
}
