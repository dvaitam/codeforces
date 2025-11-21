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

const refSource = "1000-1999/1700-1799/1770-1779/1774/1774A.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseOutputs(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutputs(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(candAns) != len(refAns) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d answers, got %d\ninput:\n%s\n", idx+1, tc.name, len(refAns), len(candAns), tc.input)
			os.Exit(1)
		}

		reader := strings.NewReader(tc.input)
		var t int
		fmt.Fscan(reader, &t)
		for caseIdx := 0; caseIdx < t; caseIdx++ {
			var n int
			var s string
			fmt.Fscan(reader, &n)
			fmt.Fscan(reader, &s)
			refSign := refAns[caseIdx]
			candSign := candAns[caseIdx]
			if len(candSign) != n-1 {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected sign string of length %d, got %d\n", idx+1, caseIdx+1, n-1, len(candSign))
				os.Exit(1)
			}
			refVal := evalExpression(s, refSign)
			candVal := evalExpression(s, candSign)
			if candVal != refVal {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected value %d, got %d\ninput:\n%s\n", idx+1, caseIdx+1, refVal, candVal, tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1774A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseOutputs(input, out string) ([]string, error) {
	reader := strings.NewReader(input)
	var t int
	fmt.Fscan(reader, &t)
	lines := filterNonEmpty(strings.Split(out, "\n"))
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d lines of output, got %d", t, len(lines))
	}
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return lines, nil
}

func evalExpression(s, ops string) int {
	result := int(s[0] - '0')
	for i := 1; i < len(s); i++ {
		val := int(s[i] - '0')
		if ops[i-1] == '+' {
			result += val
		} else {
			result -= val
		}
	}
	return abs(result)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func filterNonEmpty(lines []string) []string {
	var res []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("basic", [][]string{
			{"2", "11"},
			{"5", "01011"},
			{"5", "10001"},
		}),
	}

	rng := rand.New(rand.NewSource(1774))
	for i := 0; i < 40; i++ {
		t := rng.Intn(5) + 1
		var cases [][]string
		for j := 0; j < t; j++ {
			n := rng.Intn(99) + 2
			var sb strings.Builder
			for k := 0; k < n; k++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			cases = append(cases, []string{strconv.Itoa(n), sb.String()})
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), cases))
	}

	return tests
}

func buildCase(name string, cases [][]string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%s\n%s\n", c[0], c[1])
	}
	return testCase{name: name, input: sb.String()}
}
