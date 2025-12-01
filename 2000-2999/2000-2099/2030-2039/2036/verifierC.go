package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource2036C = "./2036C.go"

type query struct {
	pos int
	val int
}

type caseInput struct {
	s       string
	queries []query
}

type testCase struct {
	name      string
	input     string
	ansCount  int
	caseCount int
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, tc.ansCount)
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
		candAns, err := parseAnswers(candOut, tc.ansCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		for i := 0; i < tc.ansCount; i++ {
			if refAns[i] != candAns[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed at answer %d: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refAns[i], candAns[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2036C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2036C.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2036C)
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

func parseAnswers(output string, expected int) ([]string, error) {
	tokens := strings.Fields(output)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d tokens", expected, len(tokens))
	}
	res := make([]string, expected)
	for i, tok := range tokens {
		normalized := normalizeYesNo(tok)
		if normalized != "YES" && normalized != "NO" {
			return nil, fmt.Errorf("invalid answer %q (must be YES/NO)", tok)
		}
		res[i] = normalized
	}
	return res, nil
}

func normalizeYesNo(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func buildTests() []testCase {
	tests := []testCase{
		makeManualTest("single_char", []caseInput{
			{s: "0", queries: []query{{1, 0}, {1, 1}}},
		}),
		makeManualTest("exact_match", []caseInput{
			{s: "1100", queries: []query{{2, 0}, {3, 0}, {3, 1}, {4, 1}}},
		}),
		makeManualTest("short_string", []caseInput{
			{s: "101", queries: []query{{2, 0}, {3, 1}}},
		}),
		makeManualTest("multi_case", []caseInput{
			{s: "111000", queries: []query{{3, 1}, {4, 0}, {5, 0}}},
			{s: "0000", queries: []query{{1, 1}, {2, 1}, {3, 0}}},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func makeManualTest(name string, cases []caseInput) testCase {
	input, answers := formatCases(cases)
	return testCase{
		name:     name,
		input:    input,
		ansCount: answers,
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCnt := rng.Intn(4) + 1
	cases := make([]caseInput, caseCnt)
	totalAns := 0
	for i := 0; i < caseCnt; i++ {
		n := rng.Intn(10) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		q := rng.Intn(10) + 1
		qs := make([]query, q)
		for j := 0; j < q; j++ {
			qs[j] = query{
				pos: rng.Intn(n) + 1,
				val: rng.Intn(2),
			}
		}
		cases[i] = caseInput{s: sb.String(), queries: qs}
		totalAns += q
	}
	input, answers := formatCases(cases)
	return testCase{
		name:     fmt.Sprintf("random_%d", idx+1),
		input:    input,
		ansCount: answers,
	}
}

func formatCases(cases []caseInput) (string, int) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	total := 0
	for _, cs := range cases {
		sb.WriteString(cs.s)
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", len(cs.queries)))
		for _, q := range cs.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", q.pos, q.val))
		}
		total += len(cs.queries)
	}
	return sb.String(), total
}
