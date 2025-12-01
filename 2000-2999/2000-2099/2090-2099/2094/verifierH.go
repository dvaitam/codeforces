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

const refSource = "./2094H.go"

type testBatch struct {
	text    string
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/candidate")
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
		refOut, err := runProgram(exec.Command(refBin), tc.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.text)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candCmd := commandFor(candidate)
		candOut, err := runProgram(candCmd, tc.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.text, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.answers; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d line %d\ninput:\n%s\nexpected: %d\nfound: %d\n", idx+1, i+1, tc.text, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2094H-ref-*")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmpPath, filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmpPath, nil
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() []testBatch {
	tests := []testBatch{sampleTest()}
	tests = append(tests, fixedTests()...)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 {
		tests = append(tests, randomBatch(rng))
	}
	return tests
}

func sampleTest() testBatch {
	text := "" +
		"2\n" +
		"5 3\n" +
		"2 3 5 7 11\n" +
		"2 1 5\n" +
		"2 2 4\n" +
		"2310 1 5\n" +
		"4 3\n" +
		"18 12 8 9\n" +
		"216 1 2\n" +
		"48 2 4\n" +
		"82944 1 4\n"
	return testBatch{text: text, answers: 7}
}

func fixedTests() []testBatch {
	// Covers k=1 (immediate early break), arrays with repeated primes, and full removal of factors.
	text1 := "" +
		"2\n" +
		"1 3\n" +
		"2\n" +
		"1 1 1\n" +
		"2 1 1\n" +
		"8 1 1\n" +
		"4 4\n" +
		"6 6 6 6\n" +
		"12 1 4\n" +
		"72 2 3\n" +
		"7 4 4\n" +
		"1 2 2\n"

	// Includes mixed primes, wide ranges, and larger k values.
	text2 := "" +
		"3\n" +
		"3 4\n" +
		"2 4 8\n" +
		"64 1 3\n" +
		"3 1 1\n" +
		"5 2 2\n" +
		"7 3 3\n" +
		"6 5\n" +
		"11 13 17 19 23 29\n" +
		"100000 1 6\n" +
		"77 2 4\n" +
		"121 3 5\n" +
		"341 1 3\n" +
		"6 2\n" +
		"100000 100000 100000 100000 100000 100000\n" +
		"99991 1 6\n" +
		"100000 2 5\n"

	return []testBatch{
		{text: text1, answers: 7},
		{text: text2, answers: 11},
	}
}

func randomBatch(rng *rand.Rand) testBatch {
	t := rng.Intn(4) + 1

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	totalAnswers := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(80) + 1
		q := rng.Intn(40) + 1
		totalAnswers += q

		fmt.Fprintf(&sb, "%d %d\n", n, q)
		for j := 0; j < n; j++ {
			val := rng.Intn(100000-2) + 2
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')

		for j := 0; j < q; j++ {
			k := rng.Intn(100000) + 1
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			fmt.Fprintf(&sb, "%d %d %d\n", k, l, r)
		}
	}

	return testBatch{text: sb.String(), answers: totalAnswers}
}
