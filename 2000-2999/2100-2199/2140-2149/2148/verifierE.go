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
	input string
	t     int
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", out, "2148E.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(output))
	}
	return out, nil
}

func runProgram(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(lines))
	}
	return lines, nil
}

func deterministicTests() []testCase {
	return []testCase{
		formatTest([]testInput{
			{n: 3, k: 2, arr: []int{1, 1, 1}},
			{n: 4, k: 2, arr: []int{1, 2, 1, 2}},
		}),
		formatTest([]testInput{
			{n: 5, k: 5, arr: []int{1, 2, 3, 4, 5}},
		}),
		formatTest([]testInput{
			{n: 6, k: 3, arr: []int{1, 1, 1, 1, 1, 1}},
			{n: 6, k: 2, arr: []int{1, 2, 3, 1, 2, 3}},
		}),
	}
}

type testInput struct {
	n, k int
	arr  []int
}

func formatTest(cases []testInput) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), t: len(cases)}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	totalN := 0
	for len(tests) < count && totalN < 200000 {
		t := rnd.Intn(5) + 1
		cases := make([]testInput, t)
		for i := 0; i < t; i++ {
			if totalN >= 200000 {
				t = i
				break
			}
			n := rnd.Intn(2000) + 2
			if totalN+n > 200000 {
				n = 200000 - totalN
			}
			totalN += n
			k := rnd.Intn(n-1) + 2
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				arr[j] = rnd.Intn(n) + 1
			}
			cases[i] = testInput{n: n, k: k, arr: arr}
		}
		tests = append(tests, formatTest(cases))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(200)...)

	for idx, tc := range tests {
		expOut, err := runProgram(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expVals, err := parseOutput(expOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		for i := 0; i < tc.t; i++ {
			if gotVals[i] != expVals[i] {
				fmt.Fprintf(os.Stderr, "case %d test %d mismatch: expected %s got %s\n", idx+1, i+1, expVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
