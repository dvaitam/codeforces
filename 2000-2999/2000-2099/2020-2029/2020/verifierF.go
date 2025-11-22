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

const refSource = "2000-2999/2000-2099/2020-2029/2020/2020F.go"

type testCase struct {
	n int64
	k int
	d int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := renderInput(tests)

	refOut, err := runWithInput(exec.Command(refBin), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candCmd := commandFor(candidate)
	candOut, err := runWithInput(candCmd, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	expected, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\ninput:\n%s", i+1, expected[i], got[i], formatSingleInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2020F-ref-*")
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
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
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

func parseOutputs(output string, expected int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) < expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(fields))
	}
	if len(fields) > expected {
		return nil, fmt.Errorf("extra output detected after %d numbers", expected)
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse output %d (%q): %v", i+1, f, err)
		}
		res[i] = val
	}
	return res, nil
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.k, tc.d))
	}
	return sb.String()
}

func formatSingleInput(tc testCase) string {
	return fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.k, tc.d)
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 6, k: 1, d: 1},           // sample 1
		{n: 1, k: 3, d: 3},           // sample 2
		{n: 10, k: 1, d: 2},          // sample 3
		{n: 1, k: 100000, d: 100000}, // max k,d with tiny n
		{n: 2, k: 100000, d: 1},
		{n: 30, k: 5, d: 5},
		{n: 1000, k: 2, d: 3},
	}

	totalN := int64(0)
	for _, tc := range tests {
		totalN += tc.n
	}

	const maxTotalN = 300000
	rng := rand.New(rand.NewSource(2020))
	for totalN < maxTotalN {
		n := int64(rng.Intn(5000) + 1)
		if totalN+n > maxTotalN {
			n = maxTotalN - totalN
		}
		k := rng.Intn(100000) + 1
		d := rng.Intn(100000) + 1
		tests = append(tests, testCase{n: n, k: k, d: d})
		totalN += n
		if len(tests) > 200 {
			break
		}
	}

	if totalN < maxTotalN {
		n := maxTotalN - totalN
		tests = append(tests, testCase{n: n, k: 3, d: 7})
	}

	return tests
}
