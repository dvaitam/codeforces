package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

type caseData struct {
	n   int
	arr []int
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		cases, err := parseInputCases(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse generated test %q: %v\n", tc.name, err)
			os.Exit(1)
		}

		if _, err := runProgram(refBin, tc.input); err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := verifyCandidateOutput(cases, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1270G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1270G.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseInputCases(input string) ([]caseData, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("read t: %w", err)
	}
	cases := make([]caseData, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, fmt.Errorf("case %d: read n: %w", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &arr[j]); err != nil {
				return nil, fmt.Errorf("case %d: read a[%d]: %w", i+1, j+1, err)
			}
		}
		cases[i] = caseData{n: n, arr: arr}
	}
	return cases, nil
}

func verifyCandidateOutput(cases []caseData, output string) error {
	reader := bufio.NewReader(strings.NewReader(output))
	for idx, tc := range cases {
		var s int
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return fmt.Errorf("case %d: failed to read subset size: %v", idx+1, err)
		}
		if s < 1 || s > tc.n {
			return fmt.Errorf("case %d: subset size %d out of range [1, %d]", idx+1, s, tc.n)
		}
		used := make([]bool, tc.n)
		var sum int64
		for j := 0; j < s; j++ {
			var pos int
			if _, err := fmt.Fscan(reader, &pos); err != nil {
				return fmt.Errorf("case %d: failed to read index %d: %v", idx+1, j+1, err)
			}
			if pos < 1 || pos > tc.n {
				return fmt.Errorf("case %d: index %d out of range [1, %d]", idx+1, pos, tc.n)
			}
			if used[pos-1] {
				return fmt.Errorf("case %d: index %d repeated", idx+1, pos)
			}
			used[pos-1] = true
			sum += int64(tc.arr[pos-1])
		}
		if sum != 0 {
			return fmt.Errorf("case %d: subset sum %d != 0", idx+1, sum)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err != nil {
		if err != io.EOF {
			return fmt.Errorf("failed to parse trailing output: %v", err)
		}
	} else {
		return fmt.Errorf("unexpected extra token %q after processing all cases", extra)
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, buildManualTest("single_zero", [][]int{
		{0},
	}))

	tests = append(tests, buildManualTest("manual_mix", [][]int{
		{0, -1, 1, 0},
		{-4, 1, 0, 3, 0},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomManySmall("many_small", 50, 20, rng))
	tests = append(tests, randomWithNs("random_small", []int{5, 7, 11, 13, 17}, rng))
	tests = append(tests, randomWithNs("random_medium", []int{500, 800, 1200, 1600}, rng))
	tests = append(tests, randomWithNs("random_large", []int{50000, 70000, 90000}, rng))
	tests = append(tests, randomWithNs("random_huge", []int{200000, 200000, 200000, 200000}, rng))

	return tests
}

func buildManualTest(name string, arrays [][]int) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(arrays))
	for _, arr := range arrays {
		fmt.Fprintf(&b, "%d\n", len(arr))
		for i, val := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteByte('\n')
	}
	return testCase{name: name, input: b.String()}
}

func randomManySmall(name string, caseCount int, maxN int, rng *rand.Rand) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", caseCount)
	for i := 0; i < caseCount; i++ {
		n := rng.Intn(maxN) + 1
		fmt.Fprintf(&b, "%d\n", n)
		arr := randomArray(n, rng)
		for j, val := range arr {
			if j > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteByte('\n')
	}
	return testCase{name: name, input: b.String()}
}

func randomWithNs(name string, ns []int, rng *rand.Rand) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(ns))
	for _, n := range ns {
		fmt.Fprintf(&b, "%d\n", n)
		arr := randomArray(n, rng)
		for j, val := range arr {
			if j > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteByte('\n')
	}
	return testCase{name: name, input: b.String()}
}

func randomArray(n int, rng *rand.Rand) []int {
	arr := make([]int, n)
	for i := 1; i <= n; i++ {
		low := i - n
		high := i - 1
		arr[i-1] = low + rng.Intn(high-low+1)
	}
	return arr
}
