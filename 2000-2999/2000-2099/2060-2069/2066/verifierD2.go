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

const refSource = "./2066D2.go"

type testInput struct {
	text  string
	cases int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/candidate")
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
		refVals, err := parseOutputs(refOut, tc.cases)
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
		candVals, err := parseOutputs(candOut, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.cases; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%s\nexpected: %d\nfound: %d\n", idx+1, i+1, tc.text, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2066D2-ref-*")
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

func generateTests() []testInput {
	tests := []testInput{sampleTest()}
	tests = append(tests, smallFixedTests()...)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 {
		t := rng.Intn(8) + 1
		tests = append(tests, randomInput(rng, t))
	}
	return tests
}

func sampleTest() testInput {
	text := "" +
		"8\n" +
		"3 2 4\n" +
		"0 0 0 0\n" +
		"5 5 7\n" +
		"0 0 0 0 0 0 0\n" +
		"6 1 3\n" +
		"2 0 0\n" +
		"2 3 5\n" +
		"0 0 1 0 2\n" +
		"3 3 4\n" +
		"3 3 3 0\n" +
		"2 1 2\n" +
		"0 1\n" +
		"2 1 2\n" +
		"0 2\n" +
		"5 3 12\n" +
		"0 0 1 0 2 4 0 0 0 5 0 5\n"
	return testInput{text: text, cases: 8}
}

func smallFixedTests() []testInput {
	text1 := "" +
		"3\n" +
		"1 1 1\n" +
		"0\n" +
		"3 2 6\n" +
		"1 2 3 0 2 1\n" +
		"4 4 8\n" +
		"0 0 0 0 0 0 0 0\n"

	text2 := "" +
		"4\n" +
		"2 1 2\n" +
		"1 2\n" +
		"2 2 4\n" +
		"2 1 0 0\n" +
		"5 1 5\n" +
		"5 4 3 2 1\n" +
		"5 5 10\n" +
		"0 1 2 3 4 5 0 5 4 3\n"

	return []testInput{
		{text: text1, cases: 3},
		{text: text2, cases: 4},
	}
}

func randomInput(rng *rand.Rand, t int) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		c := rng.Intn(100) + 1
		maxM := n * c
		limit := maxM
		if limit > 300 {
			if rng.Intn(6) != 0 {
				limit = 100 + rng.Intn(201)
			}
		}
		if limit < c {
			limit = c
		}
		m := c
		if limit > c {
			m = c + rng.Intn(limit-c+1)
		}
		fmt.Fprintf(&sb, "%d %d %d\n", n, c, m)
		for j := 0; j < m; j++ {
			val := 0
			if rng.Intn(100) < 45 {
				val = rng.Intn(n) + 1
			}
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return testInput{text: sb.String(), cases: t}
}
