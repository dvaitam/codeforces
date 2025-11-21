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

const refSource2044G2 = "2000-2999/2000-2099/2040-2049/2044/2044G2.go"

type caseData struct {
	n int
	r []int
}

type testCase struct {
	name     string
	input    string
	ansCount int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG2.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(tc.ansCount, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(tc.ansCount, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.ansCount; i++ {
			if refAns[i] != candAns[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refAns[i], candAns[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2044G2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2044G2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2044G2)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(expected int, output string) ([]int64, error) {
	tokens := strings.Fields(strings.TrimSpace(output))
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d tokens", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		if val < 2 {
			return nil, fmt.Errorf("answer must be at least 2, got %d", val)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManual("two_cycle", []caseData{
			{n: 2, r: []int{2, 1}},
		}),
		makeManual("sample_small", []caseData{
			{n: 5, r: []int{2, 3, 4, 5, 1}},
			{n: 5, r: []int{2, 1, 4, 2, 3}},
		}),
		makeManual("multi_component", []caseData{
			{n: 4, r: []int{2, 3, 4, 1}},
			{n: 6, r: []int{2, 3, 1, 6, 4, 5}},
			{n: 3, r: []int{2, 3, 1}},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func makeManual(name string, cases []caseData) testCase {
	input, ansCount := formatInput(cases)
	return testCase{
		name:     name,
		input:    input,
		ansCount: ansCount,
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCnt := rng.Intn(4) + 1
	cases := make([]caseData, caseCnt)
	for i := 0; i < caseCnt; i++ {
		n := rng.Intn(20) + 2
		r := make([]int, n)
		for j := 0; j < n; j++ {
			to := rng.Intn(n)
			if to == j {
				to = (to + 1) % n
			}
			r[j] = to + 1
		}
		cases[i] = caseData{n: n, r: r}
	}
	input, ansCount := formatInput(cases)
	return testCase{
		name:     fmt.Sprintf("random_%d", idx+1),
		input:    input,
		ansCount: ansCount,
	}
}

func formatInput(cases []caseData) (string, int) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		for i, v := range cs.r {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), len(cases)
}
