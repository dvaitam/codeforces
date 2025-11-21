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

const referenceSource = "2000-2999/2000-2099/2020-2029/2028/2028C.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d answers got %d\ninput:\n%sreference:\n%s\ncandidate:\n%s",
				idx+1, tc.name, len(refVals), len(candVals), tc.input, refOut, candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2028C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2028C.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	manual := []string{
		"7\n6 2 1\n1 1 10 1 1 10\n6 2 2\n1 1 10 1 1 10\n6 2 3\n1 1 10 1 1 10\n6 2 10\n1 1 10 1 1 10\n6 2 11\n1 1 10 1 1 10\n6 2 12\n1 1 10 1 1 10\n6 2 12\n1 1 1 1 10 10\n",
	}
	for idx, input := range manual {
		tests = append(tests, testCase{name: fmt.Sprintf("manual-%d", idx+1), input: input})
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, testCase{name: fmt.Sprintf("random-%d", i+1), input: randomTest(rng)})
	}
	return tests
}

func randomTest(rng *rand.Rand) string {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(n) + 1
		v := rng.Int63n(50) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, v))
		for i := 0; i < n; i++ {
			val := rng.Int63n(50) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(val, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutputs(out string, input string) ([]int64, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, fmt.Errorf("empty output")
	}
	lines := strings.Split(out, "\n")
	expected, err := countCases(input)
	if err != nil {
		return nil, err
	}
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d lines, got %d", expected, len(lines))
	}
	res := make([]int64, expected)
	for i, line := range lines {
		line = strings.TrimSpace(line)
		val, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q on line %d", line, i+1)
		}
		res[i] = val
	}
	return res, nil
}

func countCases(input string) (int, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, fmt.Errorf("failed to read test count: %v", err)
	}
	return t, nil
}
