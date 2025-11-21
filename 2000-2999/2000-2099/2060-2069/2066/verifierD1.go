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
	input   string
	caseCnt int
	totalM  int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
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
		refVals, err := parseOutputs(refOut, tc.caseCnt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, tc.caseCnt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < tc.caseCnt; caseIdx++ {
			if refVals[caseIdx] != gotVals[caseIdx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					idx+1, caseIdx+1, refVals[caseIdx], gotVals[caseIdx], tc.input, refOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2066D1_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2066D1.go")
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

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers got %d", expected, len(fields))
	}
	vals := make([]int64, expected)
	for i, token := range fields {
		v, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		vals[i] = v
	}
	return vals, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 60)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	return []testCase{
		makeTestCase([][3]int{
			{3, 2, 4},
			{5, 5, 7},
		}),
		makeTestCase([][3]int{
			{1, 1, 1},
			{2, 1, 2},
			{3, 1, 3},
		}),
	}
}

func randomTests(rng *rand.Rand, batches int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCnt := rng.Intn(6) + 1
		entries := make([][3]int, 0, caseCnt)
		sumM := 0
		for len(entries) < caseCnt {
			n := rng.Intn(100) + 1
			c := rng.Intn(100) + 1
			if c > n {
				c = n
			}
			m := rng.Intn(n*c-c+1) + c
			if sumM+m > 10000 {
				break
			}
			entries = append(entries, [3]int{n, c, m})
			sumM += m
		}
		if len(entries) == 0 {
			entries = append(entries, [3]int{1, 1, 1})
		}
		tests = append(tests, makeTestCase(entries))
	}
	return tests
}

func stressTests() []testCase {
	return []testCase{
		makeTestCase([][3]int{
			{100, 100, 10000},
		}),
		makeTestCase([][3]int{
			{100, 50, 5000},
			{100, 75, 7500},
		}),
		makeTestCase([][3]int{
			{50, 1, 50},
			{50, 2, 100},
			{50, 3, 150},
			{50, 4, 200},
		}),
	}
}

func makeTestCase(entries [][3]int) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(entries)))
	sb.WriteByte('\n')
	totalM := 0
	for _, e := range entries {
		n, c, m := e[0], e[1], e[2]
		if c > n {
			c = n
		}
		if m < c {
			m = c
		}
		if m > n*c {
			m = n * c
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, c, m))
		totalM += m
		for i := 0; i < m; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteByte('0')
		}
		sb.WriteByte('\n')
	}
	return testCase{
		input:   sb.String(),
		caseCnt: len(entries),
		totalM:  totalM,
	}
}
