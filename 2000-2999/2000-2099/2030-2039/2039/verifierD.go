package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch.\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					idx+1, caseIdx+1, tc.input, refVals[caseIdx], gotVals[caseIdx])
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
	tmp, err := os.CreateTemp("", "2039D_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2039D.go")
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

func parseOutputs(out string, expected int) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	compacted := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		compacted = append(compacted, line)
	}
	if len(compacted) != expected {
		return nil, fmt.Errorf("expected %d lines got %d", expected, len(compacted))
	}
	return compacted, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 60, 2000)...)
	tests = append(tests, randomTests(rng, 40, 20000)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	cases := []struct {
		n int
		s []int
	}{
		{n: 6, s: []int{3, 4, 6}},
		{n: 1, s: []int{1}},
		{n: 2, s: []int{2}},
		{n: 5, s: []int{1, 2, 3, 4, 5}},
	}
	return []testCase{
		makeTestCase(cases),
	}
}

func randomTests(rng *rand.Rand, batches int, maxN int) []testCase {
	const globalLimit = 300000
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCnt := rng.Intn(4) + 1
		sumN := 0
		var cases []struct {
			n int
			s []int
		}
		for len(cases) < caseCnt {
			n := rng.Intn(maxN) + 1
			if sumN+n > globalLimit {
				break
			}
			m := rng.Intn(n) + 1
			s := randomSet(rng, n, m)
			cases = append(cases, struct {
				n int
				s []int
			}{n: n, s: s})
			sumN += n
		}
		if len(cases) == 0 {
			cases = append(cases, struct {
				n int
				s []int
			}{n: 1, s: []int{1}})
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func stressTests() []testCase {
	return []testCase{
		makeTestCase([]struct {
			n int
			s []int
		}{
			{n: 100000, s: sequenceSet(100000)},
		}),
		makeTestCase([]struct {
			n int
			s []int
		}{
			{n: 100000, s: []int{1}},
		}),
	}
}

func randomSet(rng *rand.Rand, n, m int) []int {
	set := make(map[int]struct{}, m)
	for len(set) < m {
		val := rng.Intn(n) + 1
		set[val] = struct{}{}
	}
	arr := make([]int, 0, m)
	for v := range set {
		arr = append(arr, v)
	}
	sort.Ints(arr)
	return arr
}

func sequenceSet(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	return arr
}

func makeTestCase(cases []struct {
	n int
	s []int
}) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", c.n, len(c.s)))
		for i := 0; i < len(c.s); i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(c.s[i]))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		input:   sb.String(),
		caseCnt: len(cases),
	}
}
