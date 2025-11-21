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
	out := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", out, "1551C.go")
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

func parseOutput(out string, expected int) ([]int, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(lines))
	}
	res := make([]int, expected)
	for i, line := range lines {
		val, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", line)
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	tests := []testCase{
		formatTest([][]string{
			{"a", "b", "c"},
		}),
		formatTest([][]string{
			{"aa", "bb", "cc", "dd", "ee"},
		}),
		formatTest([][]string{
			{"aaaaa", "bb", "cc"},
		}),
		formatTest([][]string{
			{"abcde", "e", "a"},
		}),
	}
	return tests
}

func formatTest(cases [][]string) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, words := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", len(words)))
		for _, w := range words {
			sb.WriteString(w)
			sb.WriteByte('\n')
		}
	}
	return testCase{input: sb.String(), t: len(cases)}
}

func randomWord(rnd *rand.Rand, maxLen int) string {
	length := rnd.Intn(maxLen) + 1
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(byte('a' + rnd.Intn(5)))
	}
	return sb.String()
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	totalWords := 0
	for len(tests) < count && totalWords < 200000 {
		t := rnd.Intn(5) + 1
		cases := make([][]string, t)
		for i := 0; i < t; i++ {
			n := rnd.Intn(100) + 1
			if totalWords+n > 200000 {
				n = 200000 - totalWords
			}
			totalWords += n
			words := make([]string, n)
			for j := 0; j < n; j++ {
				words[j] = randomWord(rnd, 100)
			}
			cases[i] = words
		}
		tests = append(tests, formatTest(cases))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
				fmt.Fprintf(os.Stderr, "case %d test %d mismatch: expected %d got %d\n", idx+1, i+1, expVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
