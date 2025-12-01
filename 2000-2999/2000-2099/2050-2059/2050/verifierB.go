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
	"time"
)

const refSource = "./2050B.go"

type testCase struct {
	name  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
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
		exp, err := parseAnswers(refOut, tc.t)
		if err != nil {
			fail("failed to parse reference output on test %d (%s): %v\nOutput:\n%s", idx+1, tc.name, err, refOut)
		}

		candOut, err := runProgramCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
		}
		got, err := parseAnswers(candOut, tc.t)
		if err != nil {
			fail("failed to parse candidate output on test %d (%s): %v\nOutput:\n%s", idx+1, tc.name, err, candOut)
		}

		for i := range exp {
			if got[i] != exp[i] {
				fail("wrong answer on test %d (%s) case %d: expected %s got %s\nInput:\n%s",
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
	tmp, err := os.CreateTemp("", "2050B-ref-*")
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

func parseAnswers(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(tokens))
	}
	res := make([]string, expected)
	for i, tok := range tokens {
		ans := strings.ToUpper(tok)
		if ans != "YES" && ans != "NO" {
			return nil, fmt.Errorf("invalid answer %q", tok)
		}
		res[i] = ans
	}
	return res, nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, cases [][]int64) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(cases))
		for _, arr := range cases {
			fmt.Fprintf(&sb, "%d\n", len(arr))
			for i, v := range arr {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.FormatInt(v, 10))
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{name: name, input: sb.String(), t: len(cases)})
	}

	add("sample", [][]int64{
		{3, 2, 13},
		{1, 1, 3},
		{1, 2, 5, 4},
		{1, 6, 6, 15},
		{6, 2, 1, 4, 2},
		{1, 4, 2, 15},
		{3, 1, 2, 1, 3},
		{2, 4, 2},
	})

	add("single-element", [][]int64{
		{5, 5, 5},
		{1, 1, 1},
	})

	add("two-values", [][]int64{
		{1, 1, 2},
		{1000000000, 999999999, 1000000000},
	})

	tests = append(tests, randomCase("random-small", 20, 5))
	tests = append(tests, randomCase("random-large", 5, 200000))
	return tests
}

func randomCase(name string, cases, maxN int) testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	arrs := make([][]int64, cases)
	for i := 0; i < cases; i++ {
		n := rng.Intn(maxN-2) + 3
		arr := make([]int64, n)
		base := int64(rng.Intn(10) + 1)
		if rng.Intn(2) == 0 {
			for j := 0; j < n; j++ {
				arr[j] = base
			}
		} else {
			for j := 0; j < n; j++ {
				arr[j] = base + int64(rng.Intn(5))
			}
		}
		arrs[i] = arr
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arrs))
	for _, arr := range arrs {
		fmt.Fprintf(&sb, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String(), t: len(arrs)}
}
