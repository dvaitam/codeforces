package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2020-2029/2020/2020E.go"

type testCase struct {
	name  string
	input string
	t     int
}

type testData struct {
	a []int
	p []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
		}
		exp, err := parseOutputs(refOut, tc.t)
		if err != nil {
			fail("failed to parse reference output on test %d (%s): %v\nOutput:\n%s", idx+1, tc.name, err, refOut)
		}

		candOut, err := runProgramCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
		}
		got, err := parseOutputs(candOut, tc.t)
		if err != nil {
			fail("failed to parse candidate output on test %d (%s): %v\nOutput:\n%s", idx+1, tc.name, err, candOut)
		}
		for i := 0; i < tc.t; i++ {
			if got[i] != exp[i] {
				fail("wrong answer on test %d (%s) case %d: expected %d got %d\nInput:\n%s",
					idx+1, tc.name, i+1, exp[i], got[i], tc.input)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2020E-ref-*")
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
		return "", fmt.Errorf("go build failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runProgramCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, dataset []testData) {
		input, t := formatDataset(dataset)
		tests = append(tests, testCase{name: name, input: input, t: t})
	}

	add("basic", []testData{
		{a: []int{1, 2}, p: []int{5000, 5000}},
		{a: []int{1, 1}, p: []int{1000, 2000}},
		{a: []int{3, 4, 3, 6, 2, 4}, p: []int{343, 624, 675, 451, 902, 820}},
		{a: []int{6, 5, 3, 6, 5, 3}, p: []int{5326, 7648, 2165, 9430, 5428, 1110}},
	})

	add("edge-small", []testData{
		{a: []int{1023}, p: []int{10000}},
		{a: []int{5, 5}, p: []int{10000, 10000}},
		{a: []int{7, 8, 9}, p: []int{1, 1, 1}},
		{a: []int{10, 20, 30, 40}, p: []int{9000, 8000, 7000, 6000}},
	})

	rng := rand.New(rand.NewSource(2020))
	for i := 0; i < 5; i++ {
		tests = append(tests, randomTest(fmt.Sprintf("random-%d", i+1), rng, 5, 50))
	}
	tests = append(tests, randomTest("large", rng, 1, 200000))

	return tests
}

func formatDataset(dataset []testData) (string, int) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(dataset))
	for _, tc := range dataset {
		n := len(tc.a)
		fmt.Fprintf(&sb, "%d\n", n)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), len(dataset)
}

func randomTest(name string, rng *rand.Rand, t, maxN int) testCase {
	dataset := make([]testData, t)
	totalN := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-1) + 1
		if totalN+n > 200000 {
			n = 1
		}
		totalN += n
		a := make([]int, n)
		p := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(1023) + 1
			p[j] = rng.Intn(10000) + 1
		}
		dataset[i] = testData{a: a, p: p}
	}
	input, _ := formatDataset(dataset)
	return testCase{name: name, input: input, t: t}
}
