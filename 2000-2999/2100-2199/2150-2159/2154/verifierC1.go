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

const refSource = "2000-2999/2100-2199/2150-2159/2154/2154C1.go"

type testCase struct {
	name  string
	input string
}

type caseData struct {
	a []int
	b []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
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
		refVals, err := parseOutputs(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d answers, got %d\ninput:\n%s\n", idx+1, tc.name, len(refVals), len(candVals), tc.input)
			os.Exit(1)
		}
		for caseIdx := range refVals {
			if refVals[caseIdx] != candVals[caseIdx] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) case %d: expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, caseIdx+1, refVals[caseIdx], candVals[caseIdx], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2154C1-ref-*")
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

func parseOutputs(input, out string) ([]int64, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n int
		fmt.Fscan(reader, &n)
		for i := 0; i < n; i++ {
			var tmp int
			fmt.Fscan(reader, &tmp)
		}
		for i := 0; i < n; i++ {
			var tmp int
			fmt.Fscan(reader, &tmp)
		}
	}

	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d answers, got %d tokens", t, len(tokens))
	}

	res := make([]int64, t)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", tok, err)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("manual-1", []caseData{
			{a: []int{1, 1}, b: []int{1, 1}},
			{a: []int{4, 8}, b: []int{1, 1}},
		}),
		buildCase("manual-2", []caseData{
			{a: []int{5, 7, 11}, b: []int{1, 1, 1}},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		t := rng.Intn(4) + 1
		c := make([]caseData, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(5) + 2
			arr := make([]int, n)
			for k := 0; k < n; k++ {
				arr[k] = rng.Intn(30) + 1
			}
			ones := make([]int, n)
			for k := range ones {
				ones[k] = 1
			}
			c[j] = caseData{a: arr, b: ones}
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), c))
	}

	return tests
}

func buildCase(name string, cases []caseData) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, cs := range cases {
		n := len(cs.a)
		if len(cs.b) != n {
			panic("length mismatch between a and b")
		}
		fmt.Fprintf(&sb, "%d\n", n)
		writeArray(&sb, cs.a)
		writeArray(&sb, cs.b)
	}
	return testCase{name: name, input: sb.String()}
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
}
