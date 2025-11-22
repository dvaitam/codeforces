package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2100-2109/2107/2107B.go"

type testBatch struct {
	text    string
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
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
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%s\nexpected: %s\nfound: %s\n", idx+1, i+1, tc.text, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2107B-ref-*")
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

func parseOutputs(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(lines))
	}
	for i, v := range lines {
		if v != "Tom" && v != "Jerry" {
			return nil, fmt.Errorf("invalid token at %d: %s", i+1, v)
		}
	}
	return lines, nil
}

func generateTests() []testBatch {
	tests := []testBatch{sampleTest()}
	tests = append(tests, fixedTests()...)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 150 {
		tests = append(tests, randomBatch(rng))
	}
	return tests
}

func sampleTest() testBatch {
	text := "" +
		"3\n" +
		"3 1\n" +
		"2 1 2\n" +
		"3 1\n" +
		"1 1 3\n" +
		"2 1\n" +
		"1 4\n"
	return testBatch{text: text, answers: 3}
}

func fixedTests() []testBatch {
	// Cases cover diff far above k, exact k+1 with multiple maximums, parity-based decisions, and large k.
	text1 := "" +
		"4\n" +
		"2 0\n" +
		"5 5\n" +
		"4 1\n" +
		"10 9 9 9\n" +
		"5 3\n" +
		"1 1 1 1 1\n" +
		"3 4\n" +
		"8 3 8\n" +
		"6 1000000000\n" +
		"1000000000 1 1 1 1 1\n"

	text2 := "" +
		"2\n" +
		"5 2\n" +
		"10 10 10 10 9\n" +
		"5 2\n" +
		"10 10 10 10 8\n"

	return []testBatch{
		{text: text1, answers: 5},
		{text: text2, answers: 2},
	}
}

func randomBatch(rng *rand.Rand) testBatch {
	t := rng.Intn(8) + 1

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	total := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 2
		k := int64(rng.Intn(1_000_000_000) + 1)
		total++

		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			val := rng.Int63n(1_000_000_000) + 1
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(val))
		}
		sb.WriteByte('\n')
	}

	return testBatch{text: sb.String(), answers: total}
}
