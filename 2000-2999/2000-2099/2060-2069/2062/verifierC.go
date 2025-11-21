package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		expected, err := parseOutputs(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < tc.t; caseIdx++ {
			if expected[caseIdx] != gotVals[caseIdx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					idx+1, caseIdx+1, expected[caseIdx], gotVals[caseIdx], tc.input, refOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2062C_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2062C.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".go"):
		cmd = exec.Command("go", "run", path)
	default:
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, t int) ([]int64, error) {
	lines := strings.Fields(out)
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d integers got %d", t, len(lines))
	}
	values := make([]int64, t)
	for i, token := range lines {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		values[i] = val
	}
	return values, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 80, 5)...)
	tests = append(tests, randomTests(rng, 60, 20)...)
	tests = append(tests, randomTests(rng, 40, 50)...)
	return tests
}

func manualTests() []testCase {
	return []testCase{
		makeTestCase([][]int{
			{1, -1000},
			{2, 5, -3},
			{3, 1001, -1000, 0},
			{5, 9, -9, 9, -8, 7},
			{10, 11678, 201, 340, 444, 453, 922, 128, 987, 127, 752, 0},
		}),
		makeTestCase([][]int{
			{1, 7},
			{1, -7},
			{2, 100, 50},
			{3, -2, -2, -2},
		}),
	}
}

func randomTests(rng *rand.Rand, batches int, maxN int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCount := rng.Intn(5) + 1
		cases := make([][]int, caseCount)
		for i := 0; i < caseCount; i++ {
			n := rng.Intn(maxN) + 1
			arr := make([]int, n+1)
			arr[0] = n
			for j := 1; j <= n; j++ {
				arr[j] = rng.Intn(2001) - 1000
			}
			cases[i] = arr
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func makeTestCase(cases [][]int) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		n := c[0]
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 1; j <= n; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(c[j]))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		input: sb.String(),
		t:     len(cases),
	}
}
