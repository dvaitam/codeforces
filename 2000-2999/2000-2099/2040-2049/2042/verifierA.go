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
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "2042A_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2042A.go")
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
	tests = append(tests, randomTests(rng, 80, 10)...)
	tests = append(tests, randomTests(rng, 60, 30)...)
	tests = append(tests, randomTests(rng, 40, 50)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	cases := []struct {
		n int
		k int64
		a []int64
	}{
		{5, 4, []int64{4, 1, 2, 3, 2}},
		{5, 10, []int64{4, 1, 2, 3, 2}},
		{2, 10, []int64{1, 3}},
		{3, 3, []int64{3, 3, 3}},
		{3, 3, []int64{1, 1, 1}},
	}
	return []testCase{makeTestCase(cases)}
}

func randomTests(rng *rand.Rand, batches int, maxN int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCnt := rng.Intn(5) + 1
		cases := make([]struct {
			n int
			k int64
			a []int64
		}, caseCnt)
		for i := 0; i < caseCnt; i++ {
			n := rng.Intn(maxN) + 1
			k := int64(rng.Intn(1_000_000) + 1)
			if k > 10_000_000 {
				k = 10_000_000
			}
			arr := make([]int64, n)
			sum := int64(0)
			for j := 0; j < n; j++ {
				val := int64(rng.Intn(int(min64(k, 1_000_000))) + 1)
				arr[j] = val
				sum += val
			}
			if sum >= k {
				k = sum
				k -= int64(rng.Intn(int(min64(sum, 1_000_000))))
				if k < 1 {
					k = 1
				}
			}
			cases[i] = struct {
				n int
				k int64
				a []int64
			}{n: n, k: k, a: arr}
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func stressTests() []testCase {
	return []testCase{
		makeTestCase([]struct {
			n int
			k int64
			a []int64
		}{
			{n: 50, k: 10_000_000, a: constantArray(50, 200_000)},
		}),
		makeTestCase([]struct {
			n int
			k int64
			a []int64
		}{
			{n: 50, k: 10_000_000, a: increasingArray(50)},
			{n: 1, k: 10_000_000, a: []int64{10_000_000}},
		}),
	}
}

func constantArray(n int, val int64) []int64 {
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = val
	}
	return arr
}

func increasingArray(n int) []int64 {
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(i + 1)
	}
	return arr
}

func makeTestCase(cases []struct {
	n int
	k int64
	a []int64
}) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", c.n, c.k))
		for i := 0; i < c.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(c.a[i], 10))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		input:   sb.String(),
		caseCnt: len(cases),
	}
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
