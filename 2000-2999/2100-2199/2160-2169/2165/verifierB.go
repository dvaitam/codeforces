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

type testCase struct {
	n   int
	arr []int
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2165B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2165B.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(target string) *exec.Cmd {
	switch filepath.Ext(target) {
	case ".go":
		return exec.Command("go", "run", target)
	case ".py":
		return exec.Command("python3", target)
	default:
		return exec.Command(target)
	}
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func exampleInput() string {
	cases := []testCase{
		{n: 3, arr: []int{1, 2, 3}},
		{n: 3, arr: []int{1, 1, 1}},
		{n: 3, arr: []int{1, 2, 2}},
		{n: 10, arr: []int{1, 1, 1, 1, 2, 2, 2, 3, 3, 4}},
		{n: 10, arr: []int{1, 1, 1, 2, 2, 2, 3, 3, 3, 4}},
	}
	return buildInput(cases)
}

func smallEdgeCases() string {
	cases := []testCase{
		{n: 1, arr: []int{1}},
		{n: 2, arr: []int{1, 1}},
		{n: 2, arr: []int{1, 2}},
		{n: 4, arr: []int{2, 2, 2, 2}},
		{n: 5, arr: []int{1, 1, 2, 2, 3}},
	}
	return buildInput(cases)
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	return arr
}

func randomTestInput(rng *rand.Rand, maxCases, maxN int) string {
	target := rng.Intn(maxCases) + 1
	var cases []testCase
	sumN := 0
	for len(cases) < target && sumN < 5000 {
		remain := 5000 - sumN
		if remain == 0 {
			break
		}
		nUpper := maxN
		if nUpper > remain {
			nUpper = remain
		}
		if nUpper < 1 {
			nUpper = 1
		}
		n := rng.Intn(nUpper) + 1
		cases = append(cases, testCase{n: n, arr: randomArray(rng, n)})
		sumN += n
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 1, arr: []int{1}})
	}
	return buildInput(cases)
}

func buildStressInput() string {
	n := 5000
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = (i%n)%n + 1
	}
	return buildInput([]testCase{{n: n, arr: arr}})
}

func compareOutputs(expected, got string) error {
	expVals := strings.Fields(expected)
	gotVals := strings.Fields(got)
	if len(expVals) != len(gotVals) {
		return fmt.Errorf("expected %d outputs, got %d", len(expVals), len(gotVals))
	}
	for i := range expVals {
		if expVals[i] != gotVals[i] {
			return fmt.Errorf("mismatch at line %d: expected %s, got %s", i+1, expVals[i], gotVals[i])
		}
	}
	return nil
}

func buildTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{exampleInput(), smallEdgeCases()}
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTestInput(rng, 5, 50))
	}
	for i := 0; i < 60; i++ {
		tests = append(tests, randomTestInput(rng, 5, 500))
	}
	tests = append(tests, randomTestInput(rng, 10, 2000))
	tests = append(tests, buildStressInput())
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, input := range tests {
		exp, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		if err := compareOutputs(exp, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, err, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
