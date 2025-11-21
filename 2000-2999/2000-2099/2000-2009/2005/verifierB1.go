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

type caseInput struct {
	n        int64
	teachers []int64
	query    int64
}

type testCase struct {
	input string
	cases []caseInput
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB1.go /path/to/binary")
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
		refVals, err := parseOutputs(refOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := range tc.cases {
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
	tmp, err := os.CreateTemp("", "2005B1_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2005B1.go")
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
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
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
	tests = append(tests, randomTests(rng, 80)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	cases := []caseInput{
		{n: 10, teachers: []int64{1, 4}, query: 2},
		{n: 28, teachers: []int64{3, 6}, query: 1},
		{n: 18, teachers: []int64{3, 6}, query: 8},
	}
	return []testCase{makeTestCase(cases)}
}

func randomTests(rng *rand.Rand, batches int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCnt := rng.Intn(5) + 1
		cases := make([]caseInput, caseCnt)
		for i := 0; i < caseCnt; i++ {
			n := int64(rng.Intn(1_000_000_000-2) + 3)
			t1 := int64(rng.Intn(int(n)) + 1)
			t2 := int64(rng.Intn(int(n)) + 1)
			for t2 == t1 {
				t2 = int64(rng.Intn(int(n)) + 1)
			}
			if t1 > t2 {
				t1, t2 = t2, t1
			}
			query := int64(rng.Intn(int(n)) + 1)
			for query == t1 || query == t2 {
				query = int64(rng.Intn(int(n)) + 1)
			}
			cases[i] = caseInput{n: n, teachers: []int64{t1, t2}, query: query}
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func stressTests() []testCase {
	return []testCase{
		makeTestCase([]caseInput{
			{n: 1_000_000_000, teachers: []int64{1, 1_000_000_000}, query: 500_000_000},
			{n: 1_000_000_000, teachers: []int64{2, 999_999_999}, query: 1},
		}),
	}
}

func makeTestCase(cases []caseInput) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d 2 1\n", c.n))
		sb.WriteString(fmt.Sprintf("%d %d\n", c.teachers[0], c.teachers[1]))
		sb.WriteString(fmt.Sprintf("%d\n", c.query))
	}
	return testCase{
		input: sb.String(),
		cases: cases,
	}
}
